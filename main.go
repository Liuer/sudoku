package main

import (
	"fmt"
	"flag"
	"net/http"
	"encoding/json"
)

var fAddr = flag.String("addr", "", "http service address")
var fGrid = flag.String("g", "", "grid")

type resultResp struct{
	Result string `json:"result,omitempty"`
	FinishTime string `json:"time,omitempty"`
	Err string `json:"err,omitempty"`
}
type solveReq struct{
	Grid string `json:"grid"`
}

func init(){
	http.HandleFunc("/solve", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		w.Header().Set("Content-Type","application/json")

		var resp resultResp
		var req solveReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			resp.Err = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		val, err := NewSudoku(req.Grid)
		if err != nil {
			resp.Err = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}
		res, _, _, ft := val.Solve()

		resp.Result = res.Line()
		resp.FinishTime = fmt.Sprintf("%.4f's", ft.Seconds())

		json.NewEncoder(w).Encode(resp)
	})
}

func main() {

	flag.Parse()

	if (*fGrid) != "" {
		val, err := NewSudoku(*fGrid)
		if err != nil {
			fmt.Println("error: ", err)
			return
		}

		unParsedSudoku, _ := gridValues(*fGrid)
		fmt.Printf("the sudoku puzzle: %+v \n", unParsedSudoku.Display())

		res, _, fd, ft := val.Solve()
		fmt.Printf("finished cost %.4f's finishDeep: %v, result: %+v \n", ft.Seconds(), fd, res.Display())
	}

	if (*fAddr) != "" {
		fmt.Println("listen addr",*fAddr)
		err := http.ListenAndServe(*fAddr, nil)
		if err != nil {
			fmt.Println("listen error: ", err)
			return
		}
	}
}