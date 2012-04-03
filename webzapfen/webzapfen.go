// Copyright 2012  Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// This command implements a web server which serves Zapfenrechnung
// and displays intermediate division steps 
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/the42/schoolcalc"
	"html/template"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

const (
	roottplfilename  = "index.tpl"
	divisionfilename = "division.tpl"
	zapfenfilename   = "zapfen.tpl"
	defaultlang      = "en"
	cookie_lang      = "uilang"
	setlanguageparam = "setlanguage"
	redirectparam    = "redirect"
)

var (
	roottemplatedir = conf_roottemplatedir()
	rootdomain      = conf_rootdomain()
	applicationport = conf_binding()
	languages       = conf_languages()
)

type webhandler struct {
	handler  func(w io.Writer, req *http.Request, lang string) error
	compress bool
}

type divisionPage struct {
	Dividend, Divisor, Precision string
	*schoolcalc.SDivide
	Error           []string
	StopRemz, Boxed bool
}

func tplfuncdivdisplay(sd *schoolcalc.SDivide, boxed bool) (htmlResult template.HTML) {
	if sd != nil {
		dividendivisorresult := fmt.Sprintf("%s:%s=%s", sd.NormalizedDividend, sd.NormalizedDivisor, sd.Result)
		var column template.HTML
		var runlen int

		if len(dividendivisorresult) > int(sd.ActualPrec) {
			runlen = len(dividendivisorresult)
		} else {
			runlen = int(sd.ActualPrec)
		}

		lastIntermediate := sd.DivisionSteps[len(sd.DivisionSteps)-1]
		lastdivColumn := lastIntermediate.Indent + len(lastIntermediate.Iremainder)

		for i := 0; i < runlen; i++ {

			if i < lastdivColumn {
				column = template.HTML(fmt.Sprintf(`<div class="divisionColumn" data-division="true" data-boxed="%t">`, boxed))
			} else {
				column = template.HTML(`<div class="divisionColumn" data-result="true">`)
			}

			if i < len(dividendivisorresult) {
				column += template.HTML(dividendivisorresult[i]) + "<br />"
			} else {
				column += "<br />"
			}

			for _, elm := range sd.DivisionSteps {
				if i >= elm.Indent && i < elm.Indent+len(elm.Iremainder) {
					column += template.HTML(elm.Iremainder[i-elm.Indent]) + "<br />"
				} else {
					column += "<br />"
				}
			}
			htmlResult += column + "</div>\n"
		}
	}
	return
}

var templdivfuncMap = template.FuncMap{
	"tplfuncdivdisplay": tplfuncdivdisplay,
}

func divisionHandler(w io.Writer, req *http.Request, lang string) error {

	var tpl *template.Template
	var err error

	dividend := strings.TrimSpace(req.URL.Query().Get("dividend"))
	divisor := strings.TrimSpace(req.URL.Query().Get("divisor"))
	page := &divisionPage{Dividend: dividend, Divisor: divisor, Boxed: true, StopRemz: true}

	prec := 0
	page.Precision = strings.TrimSpace(req.URL.Query().Get("prec"))

	stopremzs := req.URL.Query().Get("stopremz")
	if len(stopremzs) > 0 {
		stopremz, err := strconv.ParseBool(stopremzs)
		if err != nil {
			page.Error = append(page.Error, fmt.Sprintf("Parameter 'stopremz' tainted: %s", err))
		} else {
			page.StopRemz = stopremz
		}
	}

	boxes := req.URL.Query().Get("boxed")
	if len(boxes) > 0 {
		boxed, err := strconv.ParseBool(boxes)
		if err != nil {
			page.Error = append(page.Error, fmt.Sprintf("Parameter 'boxed' tainted: %s", err))
		} else {
			page.Boxed = boxed
		}
	}

	for retry := 0; retry <= 1; retry++ {
		tpl, err = template.New("DivisionTemplate").Funcs(templdivfuncMap).ParseFiles(roottemplatedir + lang + "." + divisionfilename)
		if err != nil {
			if _, ok := err.(*os.PathError); ok && retry < 1 {
				lang = defaultlang
			} else {
				panic(err)
			}
		}
	}

	func() {
		defer func() { // want to handle division by zero
			if err := recover(); err != nil {
				log.Printf("%s", debug.Stack())
				page.Error = append(page.Error, fmt.Sprint(err))
			}
		}()

		if len(page.Precision) > 0 {
			if prec, err = strconv.Atoi(page.Precision); err != nil {
				page.Error = append(page.Error, fmt.Sprint(err))
			}
		}
		if len(dividend) > 0 || len(divisor) > 0 {
			if page.StopRemz {
				prec = int(schoolcalc.SDivPrecReached | uint8(prec))
			}
			result, err := schoolcalc.SchoolDivide(dividend, divisor, uint8(prec))
			if err != nil {
				page.Error = append(page.Error, fmt.Sprint(err))
			} else {
				page.SDivide = result
			}
		}
	}()

	return tpl.Execute(w, page)
}

type zapfenPage struct {
	Error     []string
	Number    string
	ZapfenStr string
	*schoolcalc.Zapfen
}

// Uses a cell for every digit to display the zapfen
/*
func tplfunczapfendisplay(zapfen *schoolcalc.Zapfen) template.HTML {
	var retstring string
	if zapfen != nil {

		input := zapfen.Zapfenzahl
		inputzahlstr := input.String()

		for mults := 0; mults < 8; mults++ {
			retstring += "<tr>\n" + "  " + strings.Repeat("<td></td>", zapfen.Longest-len(inputzahlstr))
			for _, multrune := range inputzahlstr {
				retstring += fmt.Sprintf("<td>%c</td>", multrune)
			}
			retstring += fmt.Sprintf("<td>x</td><td>%d</td><td>=</td>", mults+2)

			input = zapfen.Multzapfen[mults]
			inputzahlstr = input.String()

			for _, multrune := range inputzahlstr {
				retstring += fmt.Sprintf("<td>%c</td>", multrune)
			}
			retstring += "\n</tr>"
		}

		for divs := 0; divs < 8; divs++ {
			retstring += "\n<tr>\n" + "  " + strings.Repeat("<td></td>", zapfen.Longest-len(inputzahlstr))
			for _, divrune := range inputzahlstr {
				retstring += fmt.Sprintf("<td>%c</td>", divrune)
			}
			retstring += fmt.Sprintf("<td>/</td><td>%d</td><td>=</td>", divs+2)

			input = zapfen.Divzapfen[divs]
			inputzahlstr = input.String()

			for _, divrune := range inputzahlstr {
				retstring += fmt.Sprintf("<td>%c</td>", divrune)
			}
			retstring += "\n</tr>"
		}
	}
	return template.HTML(retstring)
}
*/

func tplfunczapfendisplay(zapfen *schoolcalc.Zapfen) template.HTML {
	var retstring string
	if zapfen != nil {
		input := zapfen.Zapfenzahl

		for i := 0; i < 8; i++ {
			retstring += fmt.Sprintf("<tr>\n<td class='zapfenmultiplier'>%s</td><td>x</td><td>%d</td><td>=</td><td>%s</td>\n</tr>\n", input, i+2, zapfen.Multzapfen[i])
			input = zapfen.Multzapfen[i]
		}

		for i := 0; i < 8; i++ {
			retstring += fmt.Sprintf("\n<tr id='zapfendividend%d'>\n<td class='zapfendividend'>%s</td><td>/</td><td>%d</td><td>=</td><td>%s</td>\n</tr>\n", i, input, i+2, zapfen.Divzapfen[i])
			input = zapfen.Divzapfen[i]
		}
	}
	return template.HTML(retstring)
}

var templzapfenfuncMap = template.FuncMap{
	"tplfunczapfendisplay": tplfunczapfendisplay,
}

func zapfenHandler(w io.Writer, req *http.Request, lang string) error {

	var tpl *template.Template
	var err error

	page := &zapfenPage{Number: strings.TrimSpace(req.URL.Query().Get("number"))}

	for retry := 0; retry <= 1; retry++ {
		tpl, err = template.New("ZapfenTemplate").Funcs(templzapfenfuncMap).ParseFiles(roottemplatedir + lang + "." + zapfenfilename)
		if err != nil {
			if _, ok := err.(*os.PathError); ok && retry < 1 {
				lang = defaultlang
			} else {
				panic(err)
			}
		}
	}

	if len(page.Number) >= 1 {
		if num, ok := big.NewInt(0).SetString(page.Number, 10); ok {
			result := schoolcalc.ZapfenRechnung(num)
			page.ZapfenStr = fmt.Sprint(result)
			page.Zapfen = result
		} else {
			page.Error = append(page.Error, fmt.Sprintf("Not a valid integer: '%s'", page.Number))
		}
	}

	return tpl.Execute(w, page)
}

func rootHandler(w io.Writer, req *http.Request, lang string) (err error) {

	_, err = fmt.Fprintf(w, "Got language: %s", lang)

	tpl, err := template.ParseFiles(lang + "." + roottplfilename)
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(w, nil)
	return
}

func validlanguage(language string) bool {
	for _, lang := range languages {
		if lang == language {
			return true
		}
	}
	return false
}

// ServeHTTP installs a catch-all error recovery for the specific handler functions
// It may gzip-compress the output depending on webhandler.compress 
func (wh webhandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// error handler: Recover from panic by setting http error 500 and letting the user know the reason
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%s", debug.Stack())
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		}
	}()

	var language string
	var setcookie bool
	var mustredirect bool

	// set the language according to the request parameter.
	if language = req.URL.Query().Get(setlanguageparam); language != "" {
		if validlanguage(language) {
			setcookie = true
		} else {
			panic(fmt.Sprintf("Not a valid language '%s'. Valid languages are: %v", language, languages))
		}
	} else {
		language = defaultlang
	}

	// read language cookie. If not set, initialize and set it with the language we determined so far,
	// either the default language or the language given as a request parameter
	langCookie, err := req.Cookie(cookie_lang)
	if err == http.ErrNoCookie || setcookie {
		langCookie = &http.Cookie{Name: cookie_lang, Value: language, Path: "/", Domain: "." + rootdomain}
		http.SetCookie(w, langCookie)
		// We set the cookie and redirect to the language subdomain
		mustredirect = true
	} else if !validlanguage(language) {
		panic(fmt.Sprintf("Not a valid language '%s'. Valid languages are: %v", language, languages))
	} else if err != nil {
		panic(err)
	}

	// read the language out of the cookie
	language = langCookie.Value

	// if we have a cookie, but we are not on the correct language-specific subdomain, redirect
	if strings.Index(req.Host, language) != 0 || mustredirect {

		// BEGIN: functionality required? If yes, should redirect relatively or absolutely?
		var redirect, scheme string
		if redirect = req.URL.Query().Get(redirectparam); redirect == "" {
			redirect = req.URL.Path
		}
		// END: functionality required? If yes, should redirect relatively or absolutely?

		if req.TLS == nil {
			scheme = "http://"
		} else {
			scheme = "https://"
		}
		http.Redirect(w, req, scheme+language+"."+rootdomain+applicationport+redirect, http.StatusSeeOther)
		return
	}

	buf := new(bytes.Buffer)

	if err := wh.handler(buf, req, language); err != nil {
		panic(err)

	}

	// gzip-compression of the output of the given webhandler
	if wh.compress && strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		gzbuf := new(bytes.Buffer)
		gzwriter := gzip.NewWriter(gzbuf)

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", http.DetectContentType(buf.Bytes())) // We have to set the content type, otherwise the ResponseWriter will guess it's application/x-gzip
		gzwriter.Write(buf.Bytes())
		gzwriter.Close() // Otherwise the content won't get flushed to the output stream

		buf = gzbuf
	}
	// Prevent chunking of result
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	buf.WriteTo(w)
}

func main() {
	http.Handle("/", &webhandler{rootHandler, true})
	http.Handle("/division/", &webhandler{divisionHandler, true})
	http.Handle("/zapfen/", &webhandler{zapfenHandler, true})
	http.ListenAndServe(applicationport, nil)
}
