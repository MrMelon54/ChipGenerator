package main

import (
	"fmt"
	"github.com/Kikithecat12345/ChipGenerator"
	"math/big"
	"math/rand"
	"time"
)

var bi10 = big.NewInt(1000)

func main() {
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
		_ = ChipGenerator.GenerateIllion(z)
		fmt.Printf("%s, ", time.Since(t2))

		// total time
		fmt.Printf("%s\n", time.Since(t))
		n.Mul(n, bi10)
	}
}

func genNumber(n *big.Int) string {
        var s big.Int
        s.Rand(rand.New(rand.NewSource(time.Now().UnixNano())), n)
	return s.String()
}
