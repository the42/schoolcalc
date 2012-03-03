// Copyright 2011, 2012  Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// This command reads input from stdin in the form of dividend:dvisor
// and displays the intermediate division steps. It is useful to check
// pen and paper calculations from pupils who learn to divide.
package main

import (
	"bufio"
	"fmt"
	"github.com/the42/schoolcalc"
	"io"
	"os"
	"strings"
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

func printdivresult(sd *schoolcalc.SDivide, err error) {
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

	reader := bufio.NewReaderSize(os.Stdin, 1000)

	for data, prefix, err := reader.ReadLine(); err != io.EOF; data, prefix, err = reader.ReadLine() {

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

		printdivresult(schoolcalc.SchoolDivide(splitstrings[0], splitstrings[1], schoolcalc.SDivPrecReached|2))
	}
}
