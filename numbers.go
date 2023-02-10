package ChipGenerator

/*
This code (and it's accompanying Lua script) calculates, generates, and imports custom poker chips of arbitrary denominations and colors into Tabletop Simulator.
The code is written in Go, except for the import function which is written in Lua, because the API for TTS uses Lua exclusively.
The code is written in a way that allows for easy modification of the chip denominations, colors, and other parameters via a config file.
To support arbitrary chip denominations, the code uses math/big to handle arbitrary precision numbers, and then converts them to strings for the generation of the chip images.
*/

import (
	"math/big"
	"strconv"
	"strings"
	"sync"
)

// == variables ==
var prefixes = map[int]string{
	0:   "nilli",
	1:   "un",
	2:   "duo",
	3:   "tre",
	4:   "quattuor",
	5:   "quinqua",
	6:   "se",
	7:   "septe",
	8:   "octo",
	9:   "nove",
	10:  "deci",
	20:  "viginti",
	30:  "triginta",
	40:  "quadraginta",
	50:  "quinquaginta",
	60:  "sexaginta",
	70:  "septuaginta",
	80:  "octoginta",
	90:  "nonaginta",
	100: "centi",
	200: "ducenti",
	300: "trecenti",
	400: "quadringenti",
	500: "quingenti",
	600: "sescenti",
	700: "septingenti",
	800: "octingenti",
	900: "nongenti",
}
var littlePrefixes = map[int]string{
	1:  "milli", // 1
	2:  "billi",
	3:  "trilli",
	4:  "quadrilli",
	5:  "quintilli",
	6:  "sextilli",
	7:  "septilli",
	8:  "octilli",
	9:  "nonilli",
	10: "decilli", // 10
}

var cacheTripleString = make(map[string]string, 1000)

func init() {
	// Cache triples
	for i := 0; i < 1000; i++ {
		padded := padToMultipleOf3(strconv.Itoa(i))
		cacheTripleString[padded] = generateTriplePrefix(padded)
	}
}

// GenerateIllionBigInt is a wrapper for GenerateIllion which converts `math/big.Int` -> string
func GenerateIllionBigInt(n *big.Int, multi bool) string {
	return GenerateIllion(n.String(), multi)
}

// GenerateIllion takes a number and returns a string with the number in illion form, where the number is the illion in the sequence of illions.
// takes in a big.Int and returns a string.
// examples: 1 -> "million", 10 -> "decillion", 24 -> "quattorvigintillion" etc.
func GenerateIllion(str string, multi bool) string {
	// ignore an empty string
	if len(str) == 0 {
		return ""
	}

	// pad the start of the string with 0s so that it's divisible by 3.
	str = padToMultipleOf3(str)

	l := len(str)
	var illionGen strings.Builder

	if multi {
		var splitLen int
		switch {
		case l > 99999999:
			splitLen = 999999
		case l > 999999:
			splitLen = 9999
		case l > 9999:
			splitLen = 99
		default:
			splitLen = l
		}

		maxN := l / splitLen
		if l%splitLen > 0 {
			maxN++
		}

		c := make(chan illionPair, maxN)
		wg := new(sync.WaitGroup)
		wg.Add(maxN)
		n := 0
		for i := 0; i < l; i += splitLen {
			if n+1 == maxN {
				go pipeIllionRange(wg, c, n, l-i, str[i:l])
			} else {
				go pipeIllionRange(wg, c, n, splitLen, str[i:i+splitLen])
			}
			n++
		}

		c2 := make(chan []string, 1)
		go (func() {
			// collect channel outputs into array
			a := make([]string, maxN)
			for v := range c {
				a[v.n] = v.word
			}
			c2 <- a
		})()
		wg.Wait()
		close(c)

		// write to final string, waits for result on c2
		for _, i := range <-c2 {
			illionGen.WriteString(i)
		}
		close(c2)
	} else {
		// iterate in reverse order
		// i = hundreds digit for a set of 3
		for i := 0; i < len(str); i += 3 {
			illionGen.WriteString(cacheTripleString[str[i:i+3]])
		}
	}

	illionWord := illionGen.String()

	// add the "illion" suffix.
	// however, if it ends in "illi" we only add "on", for example "milli" -> "million"
	// and if it only ends in a vowel we remove it, then add "illion". for example "quadraginta" -> "quadragintillion"
	if strings.HasSuffix(illionWord, "illi") {
		illionWord += "on"
	} else {
		lastIndex := len(illionWord) - 1
		if isVowel(illionWord[lastIndex]) {
			illionWord = illionWord[:lastIndex]
		}
		illionWord += "illion"
	}
	return illionWord
}

type illionPair struct {
	n    int
	word string
}

func pipeIllionRange(wg *sync.WaitGroup, c chan illionPair, n, l int, str string) {
	defer wg.Done()
	var illionGen strings.Builder
	for i := 0; i < l; i += 3 {
		illionGen.WriteString(cacheTripleString[str[i:i+3]])
	}
	c <- illionPair{n, illionGen.String()}
}

func generateTriplePrefix(triple string) string {
	// are all the digits in this group 0
	if triple == "000" {
		return "nilli"
	}

	var lastPrefix int

	// if tens and hundreds digits are 0, we use the littlePrefixes map instead of the prefixes map.
	if triple[:2] == "00" {
		lastPrefix = int(triple[2] - '0')
		return littlePrefixes[lastPrefix]
	}

	var prefix string

	// hundreds digit
	if triple[0] != '0' {
		lastPrefix = int(triple[0]-'0') * 100
		prefix = prefixes[lastPrefix] + prefix
	}

	// tens digit
	if triple[1] != '0' {
		lastPrefix = int(triple[1]-'0') * 10
		prefix = prefixes[lastPrefix] + prefix
	}

	// we use the prefixes map, but account for english grammar rules.
	prefix = onesDigitPrefix(int(triple[2]-'0'), lastPrefix) + prefix
	return prefix
}

func onesDigitPrefix(digit int, lastPrefix int) string {
	switch {
	case digit == 0:
		return ""
	case (digit == 3 || digit == 6) && (lastPrefix == 20 || lastPrefix == 30 || lastPrefix == 40 || lastPrefix == 50 || lastPrefix == 300 || lastPrefix == 400 || lastPrefix == 500):
		return prefixes[digit] + "s"
	case digit == 6 && (lastPrefix == 80 || lastPrefix == 100 || lastPrefix == 800):
		return "sex"
	case (digit == 7 || digit == 9) && (lastPrefix == 20 || lastPrefix == 80 || lastPrefix == 800):
		return prefixes[digit] + "m"
	case (digit == 7 || digit == 9) && (lastPrefix == 10 || lastPrefix == 30 || lastPrefix == 40 || lastPrefix == 50 || lastPrefix == 60 || lastPrefix == 70 || lastPrefix == 100 || lastPrefix == 200 || lastPrefix == 300 || lastPrefix == 400 || lastPrefix == 500 || lastPrefix == 600 || lastPrefix == 700):
		return prefixes[digit] + "n"
	default:
		return prefixes[digit]
	}
}

func padToMultipleOf3(a string) string {
	switch len(a) % 3 {
	case 1:
		return "00" + a
	case 2:
		return "0" + a
	default:
		return a
	}
}

func isVowel(a uint8) bool {
	switch a {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	default:
		return false
	}
}
