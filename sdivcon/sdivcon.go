package main

import (
	"github.com/the42/sdivision"
	"fmt"
	"os"
	"strings"
)

func printsdresult(sd *sdivision.SDivide, err os.Error) {
	var blank string
	if err == nil {
		fmt.Printf("%s : %s = %s\n", sd.Dividend, sd.Divisor, sd.Result)
		fmt.Printf("%s : %s = %s\n", sd.NormalizedDividend, sd.NormalizedDivisor, sd.Result)
		for _, elm := range sd.DivisionSteps {
			blank = strings.Repeat(" ", elm.Indent)
			fmt.Printf("%s%s\n", blank, elm.Iremainder)
		}
		fmt.Printf("\n%#v\n\n", sd)
	} else {
		fmt.Printf("%s", err)
	}
	fmt.Printf("\n")
}

func main() {

	printsdresult(sdivision.SchoolDivide("100", "5", sdivision.SDivPrecReached|8))
	printsdresult(sdivision.SchoolDivide("100.5", "5", sdivision.SDivPrecReached|8))
	printsdresult(sdivision.SchoolDivide("100.5", "5.5", sdivision.SDivPrecReached|8))
	printsdresult(sdivision.SchoolDivide("-100.5", "5.56", sdivision.SDivPrecReached|2))
	printsdresult(sdivision.SchoolDivide("5", "100", sdivision.SDivPrecReached|2))
	printsdresult(sdivision.SchoolDivide("-5", "100", sdivision.SDivPrecReached|2))
	printsdresult(sdivision.SchoolDivide("2", "0.5", sdivision.SDivPrecReached|2))
	printsdresult(sdivision.SchoolDivide("0.5", "0.5", sdivision.SDivPrecReached|2))
	printsdresult(sdivision.SchoolDivide("10065767", "55.7", sdivision.SDivPrecReached|2))
	printsdresult(sdivision.SchoolDivide("0", "1", sdivision.SDivPrecReached|2))

	printsdresult(sdivision.SchoolDivide("100X65767", "55.7", sdivision.SDivPrecReached|2))
	printsdresult(sdivision.SchoolDivide("100X65767", "Y.7", sdivision.SDivPrecReached|2))
	printsdresult(sdivision.SchoolDivide("", "Y.7", sdivision.SDivPrecReached|2))
	printsdresult(sdivision.SchoolDivide("10065767", "0", sdivision.SDivPrecReached|2))
}
