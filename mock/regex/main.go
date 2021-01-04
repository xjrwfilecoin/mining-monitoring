package main

import (
	"fmt"
	"regexp"
)

var reg = regexp.MustCompile(` t[\d]* `)

func main() {
	src:=`Chain: [sync ok] [basefee 100 aFIL]
Miner: t01000 (2 KiB sectors)
Power: 42 Ki / 42 Ki (100.0000%)
	Raw: 6 KiB / 6 KiB (100.0000%)
	Committed: 8 KiB
	Proving: 6 KiB
Expected block win rate: 21600.0000/day (every 4s)

Deals: 0, 0 B
	Active: 0, 0 B (Verified: 0, 0 B)

Miner Balance:    45735.905 FIL
      PreCommit:  0
      Pledge:     1.311 μFIL
      Vesting:    34443.367 FIL
      Available:  11292.538 FIL
Market Balance:   0
       Locked:    0
       Available: 0
Worker Balance:   49899999 FIL
       Control:   50000000 FIL
Total Spendable:  99911291.538 FIL

Sectors:
	Total: 4
	Proving: 4`

	result := reg.FindString(src)
	fmt.Println(result)
}


