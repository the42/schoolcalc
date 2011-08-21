// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

//target:github.com/the42/sdivision

// This package provides functionality to calculate intermediate
// steps for division "the pen and paper method". Thus it may support
// controlling intermediate steps when pupils start learning to divide.
package sdivision

import (
	"strings"
	"fmt"
	"os"
	"big"
)

// Format specifiers
const (
	DisplayControlDefault    int8 = 10 // per default, calculate up to 10 digits precision
	DisplayControlDefaultMax int8 = -1 // 127 digits maximum precision
)

type SDivide struct {
	Dividend, Divisor, Result, Remainder string
	ISteps                               []string
	Prec                                 int8
	Exact                                bool
	Negative                             bool
}

// dividend and divisor are strings, both may contain fractions denoted by '.'
// prec set's the number of digits fraction to calculate, if the remainder is non-zero.
// prec set's the highest precision; If the remainder reaches zero before, calculation
// will stop.
func SchoolDivide(dividend, divisor string, prec int8) (sd *SDivide, err os.Error) {

	mydividend := dividend
	mydivisor := divisor
	var dividendsuffixlen, divisorsuffixlen int
	var endresult, intermediatedividend string
	var runningprec int8
	var dividendep int

	if len(mydividend) <= 0 {
		return nil, fmt.Errorf("Dividend must not be null")
	}

	if len(mydivisor) <= 0 {
		return nil, fmt.Errorf("Divisor must not be null")
	}

	negative := false

	if mydividend[0] == '-' {
		negative = !negative
		mydividend = mydividend[1:]
	}

	if mydivisor[0] == '-' {
		negative = !negative
		mydivisor = mydivisor[1:]
	}

	if len(mydividend) <= 0 {
		return nil, fmt.Errorf("Dividend must not be null")
	}

	if len(mydivisor) <= 0 {
		return nil, fmt.Errorf("Divisor must not be null")
	}

	splitstrings := strings.Split(mydivisor, ".")
	if slen := len(splitstrings); slen > 2 {
		return nil, fmt.Errorf("Not a valid divisor: \"%s\"", divisor)
	} else if slen == 2 {
		divisorsuffixlen = len(splitstrings[1])
		mydivisor = splitstrings[0] + splitstrings[1]
	}

	splitstrings = strings.Split(mydividend, ".")
	if slen := len(splitstrings); slen > 2 {
		return nil, fmt.Errorf("Not a valid dividend: \"%s\"", dividend)
	} else if slen == 2 {
		dividendsuffixlen = len(splitstrings[1])
		mydividend = splitstrings[0] + splitstrings[1]
	}

	padlen := dividendsuffixlen - divisorsuffixlen

	if padlen < 0 {
		mydividend += strings.Repeat("0", padlen*-1)
	} else if padlen > 0 {
		mydivisor += strings.Repeat("0", padlen)
	}

	bigdivisor, _ := big.NewInt(0).SetString(mydivisor, 0)
	bigdividend, _ := big.NewInt(0).SetString(mydividend, 0)
	bigintermediateremainder := big.NewInt(0)
	bigintermediatedividend := big.NewInt(0)

	for {
		intermediatedividend = mydividend[0:dividendep]
		bigintermediatedividend, _ = big.NewInt(0).SetString(intermediatedividend, 0)
		if !(bigintermediatedividend.Cmp(bigdivisor) < 0 && dividendep < len(mydividend)) {
			break
		}
		dividendep++
	}

	for {

		intermediateresult := big.NewInt(0).Div(bigintermediatedividend, bigdivisor)
		bigintermediateremainder = big.NewInt(0).Rem(bigintermediatedividend, bigdivisor)

		// endresult erst später dazufügen
		endresult += intermediateresult.String()

		if dividendep < len(mydividend) {

			onebig, _ := big.NewInt(0).SetString(string(mydividend[dividendep]), 0)
			dividendep++

			bigintermediatedividend = big.NewInt(0).Mul(big.NewInt(10), bigintermediateremainder)
			bigintermediatedividend = big.NewInt(0).Add(onebig, bigintermediatedividend)

		} else if prec-runningprec > 0 {
			bigintermediatedividend = big.NewInt(0).Mul(bigintermediateremainder, big.NewInt(10))
			runningprec++
		} else {
			break
		}
	}

	// EVTL gar nicht nötig?
	if bigdividend.Cmp(bigdivisor) < 0 {
		endresult = "0." + endresult
	}

	if negative {
		endresult = "-" + endresult
	}

	return &SDivide{Dividend: dividend, Divisor: divisor}, nil
}
