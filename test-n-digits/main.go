package main

import (
	"flag"
	"fmt"
	"github.com/Kikithecat12345/ChipGenerator"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

var bi10 = big.NewInt(10)
var bi18 = big.NewInt(18)

func main() {
	var speed bool
	flag.BoolVar(&speed, "speed", false, "Multithreading speed enhancement")
	flag.Parse()

	n := big.NewInt(1)
	fmt.Println("time, digits, gen random, gen word, total time")
	for {
		fmt.Printf("%s, %s, ", time.Now().Format(time.DateTime), n)
		t := time.Now()

		// gen random
		t2 := time.Now()
		z := genNumber(n)
		fmt.Printf("%s, ", time.Since(t2))

		// gen word
		t2 = time.Now()
		_ = ChipGenerator.GenerateIllion(z, speed)
		fmt.Printf("%s, ", time.Since(t2))

		// total time
		fmt.Printf("%s\n", time.Since(t))
		n.Mul(n, bi10)
	}
}

func genNumber(n *big.Int) string {
	end := new(big.Int).Set(n)
	var b strings.Builder
	for i := big.NewInt(0); i.Cmp(end) < 0; i.Add(i, bi18) {
		b.WriteString(fmt.Sprintf("%018d", rand.Int63())[:18])
	}
	return b.String()
}
