# sudoku

## how to build
 
```
cd <project-dir>
go build
```
in windows you will see the sudoku.exe at current dir

## cmd line
```
./sudoku.exe -g 8..........36......7..9.2...5...7.......457.....1...3...1....68..85...1..9....4..
```
will output 
```
the sudoku puzzle: 
 8 . .| . . .| . . .
 . . 3| 6 . .| . . .
 . 7 .| . 9 .| 2 . .
------+------+------
 . 5 .| . . 7| . . .
 . . .| . 4 5| 7 . .
 . . .| 1 . .| . 3 .
------+------+------
 . . 1| . . .| . 6 8
 . . 8| 5 . .| . 1 .
 . 9 .| . . .| 4 . .

finished cost 0.0150's finishDeep: 10, result: 
 8 1 2| 7 5 3| 6 4 9
 9 4 3| 6 8 2| 1 7 5
 6 7 5| 4 9 1| 2 8 3
------+------+------
 1 5 4| 2 3 7| 8 9 6
 3 6 9| 8 4 5| 7 2 1
 2 8 7| 1 6 9| 5 3 4
------+------+------
 5 2 1| 9 7 4| 3 6 8
 4 3 8| 5 2 6| 9 1 7
 7 9 6| 3 1 8| 4 5 2
```

## start as http serve
```
./sudoku.exe -addr :8080
```
it will listen 8080 port 
```
// use curl 
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "grid": "..53.....8......2..7..1.5..4....53...1..7...6..32...8..6.5....9..4....3......97.."
}' \
 'http://127.0.0.1:8080/solve'
```
will output 
```
{"result":"145327698839654127672918543496185372218473956753296481367542819984761235521839764","time":"0.0030's"}
```