package main

import (
	"github.com/the42/sdivision"
	"fmt"
	"os"
	"strings"
	"bufio"
)

func normalizeSDivInput(s string) (out string) {

	for _, ch := range s {
		if ch != ' ' {
			out += string(ch)
		}
	}

	return
}

func inputSDivisequaltoResultSDiv(inputdivisor, inputdividend, outputdivisor, outputdivided string) bool {
	return inputdividend == outputdivided && inputdivisor == outputdivisor
}

func printdivresult(sd *sdivision.SDivide, err os.Error) {
	var blank string
	if err == nil {
		fmt.Printf("%s : %s = %s\n", sd.Dividend, sd.Divisor, sd.Result)
		if !inputSDivisequaltoResultSDiv(sd.Dividend, sd.Divisor, sd.NormalizedDividend, sd.NormalizedDivisor) {
			fmt.Printf("%s : %s = %s\n", sd.NormalizedDividend, sd.NormalizedDivisor, sd.Result)
		}
		for _, elm := range sd.DivisionSteps {
			blank = strings.Repeat(" ", elm.Indent)
			fmt.Printf("%s%s\n", blank, elm.Iremainder)
		}
	} else {
		fmt.Printf("%s", err)
	}
	fmt.Printf("\n")
}

func main() {

	var lines int
	var instring string

	liner, err := bufio.NewReaderSize(os.Stdin, 1000)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sdivcon: %s (exiting)\n", err)
		os.Exit(3)
	}

	for data, prefix, err := liner.ReadLine(); err != os.EOF; data, prefix, err = liner.ReadLine() {

		if err != nil {
			fmt.Fprintf(os.Stderr, "sdivcon %d: %s\n", lines, err)
			continue
		}

		if prefix {
			fmt.Fprintf(os.Stderr, "sdivcon %d: too long, exceeding 1000 characters (ignoring)\n", lines)
			continue
		}

		lines++

		instring = normalizeSDivInput(string(data))

		splitstrings := strings.Split(instring, ":")
		if slen := len(splitstrings); slen != 2 {
			fmt.Fprintf(os.Stderr, "sdivcon %d: not a valid divisor:dividend (ignoring)\n", lines)
			continue
		}

		printdivresult(sdivision.SchoolDivide(splitstrings[0], splitstrings[1], sdivision.SDivPrecReached|2))
	}
}
