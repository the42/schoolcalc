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
	"math/big"
	"net/http"
	"os"
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
	Intermediate                 *schoolcalc.SDivide
	Error                        []string
	StopRemz, Boxed              bool
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

	dividend := req.URL.Query().Get("dividend")
	divisor := req.URL.Query().Get("divisor")
	page := &divisionPage{Dividend: dividend, Divisor: divisor, Boxed: true, StopRemz: true}

	prec := 0
	page.Precision = req.URL.Query().Get("prec")

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

	retry := false
retry:
	tpl, err := template.New("DivisionTemplate").Funcs(templdivfuncMap).ParseFiles(roottemplatedir + lang + "." + divisionfilename)
	if err != nil {
		if _, ok := err.(*os.PathError); ok && !retry {
			retry = true
			lang = defaultlang
			goto retry
		} else {
			panic(err)
		}
	}

	func() {
		defer func() { // want to handle division by zero
			if err := recover(); err != nil {
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
				page.Intermediate = result
			}
		}
	}()

	return tpl.Execute(w, page)
}

type zapfenPage struct {
	Error        []string
	Number       string
	Intermediate string
}

func zapfenHandler(w io.Writer, req *http.Request, lang string) error {

	number := req.URL.Query().Get("number")
	page := &zapfenPage{Number: number}
	retry := false

retry:
	tpl, err := template.ParseFiles(roottemplatedir + lang + "." + zapfenfilename)
	if err != nil {
		if _, ok := err.(*os.PathError); ok && !retry {
			retry = true
			lang = defaultlang
			goto retry
		} else {
			panic(err)
		}
	}

	if len(number) >= 1 {
		if num, ok := big.NewInt(0).SetString(number, 10); ok {
			result := schoolcalc.ZapfenRechnung(num)
			page.Intermediate = fmt.Sprint(result)
		} else {
			page.Error = append(page.Error, fmt.Sprintf("Not a valid integer: '%s'", number))
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
