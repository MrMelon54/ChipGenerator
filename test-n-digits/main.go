package main

import (
	"fmt"
	"github.com/Kikithecat12345/ChipGenerator"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

var bi1 = big.NewInt(1)
var bi10 = big.NewInt(10)

func main() {
	n := big.NewInt(1)
	fmt.Println("time, digits, gen random, gen big int, gen word, total time")
	for {
		fmt.Printf("%s, %s, ", time.Now().Format(time.UnixDate), n)
		t := time.Now()

		// gen random
		t2 := time.Now()
		z := genNumber(n)
		fmt.Printf("%s, ", time.Since(t2))

		// gen big int
		t2 = time.Now()
		var a big.Int
		a.SetString(z, 10)
		fmt.Printf("%s, ", time.Since(t2))

		// gen word
		t2 = time.Now()
		_ = ChipGenerator.GenerateIllion(&a)
		fmt.Printf("%s, ", time.Since(t2))

		// total time
		fmt.Printf("%s\n", time.Since(t))
		n.Mul(n, bi10)
	}
}

func genNumber(n *big.Int) string {
	end := new(big.Int).Set(n)
	b := ""
	for i := big.NewInt(0); i.Cmp(end) < 0; i.Add(i, bi1) {
		b += strconv.Itoa(rand.Intn(10))
	}
	return b
}
