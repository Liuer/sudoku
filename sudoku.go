package main

import (
	"fmt"
	"errors"
	"strings"
	"time"
)

// origin from https://norvig.com/sudoku.html
// sudoku the go reimplemention
var digits   = "123456789"
var rows     = "ABCDEFGHI"
var cols     = digits
var squares  = cross(rows, cols)
var unitlist = makeUnitList(rows, cols)
var units    = makeUnits(squares, unitlist)
var peers    = makePeers(squares, units)

var grid1  = "003020600900305001001806400008102900700000008006708200002609500800203009005010300"
var grid2  = "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"
var grid3  = ".....6....59.....82....8....45........3........6..3.54...325..6.................."
var grid4  = "8..........36......7..9.2...5...7.......457.....1...3...1....68..85...1..9....4.."
var grid5  = "..53.....8......2..7..1.5..4....53...1..7...6..32...8..6.5....9..4....3......97.."

type sudoku map[string]string

func NewSudoku(grid string) (sudoku, error) {
	val := sudoku{}
	for _, v := range squares {
		val[v] = digits
	}

	gv, err := gridValues(grid)
	if err != nil {
		return nil, err
	}
	for s, d := range gv {
		// fmt.Println("gridValues: ", s, d)
		if strings.Contains(digits, d) {
			_, ok := assign(val, s, d)
			if !ok {
				return nil, fmt.Errorf("not assign: s=%v, d=%v", s, d)
			}
		}
	}

	return val, nil
}

func gridValues(grid string) (sudoku, error) {
	if len(grid) != 81 {
		return nil, errors.New("parse error grid length error")
	}

	val := make(sudoku)

	for i, v := range squares {
		v2 := grid[i:i+1]
		val[v] = v2
	}

	return val, nil
}

func (sudo sudoku) Solve() (result sudoku, isFinish bool, finishDeep int, finishTime time.Duration) {
	return solve(sudo)
}

func (sudo *sudoku) Line() string {
	var line string
	for _, s := range squares {
		line += (*sudo)[s]
	}
	return line
}

func (sudo *sudoku) Display() string {

	line := "------+------+------\n"
	var res string = "\n"

	for _, s := range rows {
		for _, s2 := range cols {
			res += " " + (*sudo)[string(s) + string(s2)]
			if strings.Contains("36", string(s2)) {
				res += "|"
			}
		}
		res += "\n"

		if strings.Contains("CF", string(s)) {
			res += line
		}
	}

	return res
}

func (sudo *sudoku) Copy() sudoku {
	ret := sudoku{}

	for s, v := range (*sudo) {
		ret[s] = v
	}

	return ret
}

func solve(vals sudoku) (result sudoku, isFinish bool, finishDeep int, finishTime time.Duration) {
	t := time.Now()

	result, isFinish, finishDeep = search(vals, true, 0)
	if !isFinish {
		fmt.Printf("search err not finish \n")
	}
	finishTime = time.Now().Sub(t)

	return
}

func search(values sudoku, isOk bool, deep int) (result sudoku, isFinish bool, finishDeep int) {

	if !isOk {
		result = values
		isFinish = false
		return
	}

	isFinish = true
	result = values
	for _, s := range squares {
		if len(result[s]) != 1 {
			isFinish =  false
			break
		}
	}
	if isFinish {
		finishDeep = deep
		return
	}

	// 找到字符串最短的那个
	n, vs := 9, ""
	for _, s := range squares {
		n2 := len(result[s])
		if n2 > 1 && n > n2 {
			n = n2
			vs = s
		}
	}

	if vs != "" {
		for _, s := range result[vs] {
			retVal, ok := assign(values.Copy(), vs, string(s))
			nextDeep := deep + 1
			result, isFinish, finishDeep = search(retVal, ok, nextDeep)
			if isFinish {
				return
			}
		}
	}

	isFinish = false
	return
}

func assign(values sudoku, s, d string) (sudoku, bool) {
	otherValue := strings.Replace(values[s], d, "", 1)
	for _, v := range otherValue {
		_, ok := eliminate(values, s, string(v))
		if !ok {
			return values, false
		}
	}

	return values, true
}

func eliminate(values sudoku, s, d string) (sudoku, bool) {

	// Already eliminated
	if !strings.Contains(values[s], d) {
		return values, true
	}

	// eliminate
	values[s] = strings.Replace(values[s], d, "", 1)

	// (1) If a square s is reduced to one value d2, then eliminate d2 from the peers.
	if len(values[s]) == 0 {
		// 出错 排除后没有数字了
		return values, false
	} else if len(values[s]) == 1 {
		// 只剩最后一个数字，那么这最后一位数字就是确定的，即可排除peers[s]中的数字d2 
		d2 := values[s]
		for _, s2 := range peers[s] {
			if _, ok := eliminate(values, s2, d2); !ok {
				return values, false
			}
		}
	}

	// 被排除的 d 在其他 units[s] 中一定会存在
	for _, u := range units[s] {
		dplaces := make([]string,0)
		for _, su := range u {
			if strings.Contains(values[su], d) {
				dplaces = append(dplaces, su)
			}
		}

		if len(dplaces) == 0 {
			// 出错 不在其他 units[s] 中
			return values, false
		} else if len(dplaces) == 1 {
			// 如果只找出一位那么这个 d 的位置就是确定的
			if _, ok := assign(values, dplaces[0], d); !ok {
				return values, false
			}
		}
	}

	return values, true
}

func cross(A, B string) []string {
	res := make([]string, 0)
	
	for i := 0; i < len(A); i++ {
		for j:= 0; j < len(B); j++ {
			res = append(res, A[i:i+1] + B[j:j+1])
		}
	}

	return res
}

func makeUnitList(rows, cols string) [][]string {
	res := make([][]string, 0)

	for i := 0; i < len(cols); i++ {
		c := cols[i:i+1]
		e := cross(rows, c)
		res = append(res, e)
	}

	for i := 0; i < len(rows); i++ {
		c := rows[i:i+1]
		e := cross(c, cols)
		res = append(res, e)
	}

	rs := []string{"ABC","DEF","GHI"}
	ns := []string{"123","456","789"}
	for i := 0; i < len(rs); i++ {
		for j := 0; j < len(ns); j++ {
			is := rs[i]
			js := ns[j]
			res = append(res, cross(is, js))
		}
	}

	return res
}

func makeUnits(squares []string, unitlist [][]string) map[string][][]string {
	us := make(map[string][][]string)

	for _, v := range squares {
		us[v] = make([][]string, 0)
		for _, v2 := range unitlist {
			for _, v3 := range v2 {
				if v3 == v {
					us[v] = append(us[v], v2)
					break
				}
			}
		}
	}

	return us
}

func makePeers(squares []string, units map[string][][]string) map[string][]string {
	ps := make(map[string][]string)

	for _, v := range squares {
		ps[v] = make([]string, 0)

		us := units[v]
		for _, v2 := range us {
			for _, v3 := range v2 {
				if v3 != v {
					ps[v] = append(ps[v], v3)
				}
			}
		}
	}

	return ps
}

func test(){
	fmt.Println(len(squares) == 81)
	fmt.Println(len(unitlist) == 27)

	for _, s := range squares {
		fmt.Println(len(units[s]) == 3)
	}

	for _, s := range squares {
		fmt.Println(len(peers[s]) == 20)
	}

	// units['C2'] == [['A2', 'B2', 'C2', 'D2', 'E2', 'F2', 'G2', 'H2', 'I2'],
    //                        ['C1', 'C2', 'C3', 'C4', 'C5', 'C6', 'C7', 'C8', 'C9'],
	//                        ['A1', 'A2', 'A3', 'B1', 'B2', 'B3', 'C1', 'C2', 'C3']]
	
	// peers['C2'] == set(['A2', 'B2', 'D2', 'E2', 'F2', 'G2', 'H2', 'I2',
    //                            'C1', 'C3', 'C4', 'C5', 'C6', 'C7', 'C8', 'C9',
    //                            'A1', 'A3', 'B1', 'B3'])

}