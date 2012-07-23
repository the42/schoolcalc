// Copyright 2012  Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// This program implements a web server which serves Zapfenrechnung and
// divisions, displays intermediate division steps and can present the
// user some rendered divisons for self-test
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
	"path"
	"runtime/debug"
	"strconv"
	"strings"
)

const (
	roottplfilename   = "root.tpl"
	indextplfilename  = "index.tpl"
	divisionfilename  = "division.tpl"
	zapfenfilename    = "zapfen.tpl"
	excersisefilename = "excercise.tpl"
	defaultlang       = "en"
	cookie_lang       = "uilang"
	setlanguageparam  = "setlanguage"
	redirectparam     = "redirect"
)

type menulayout struct {
	URL           string
	assocmenuitem map[string]string
}

var menustructure = []menulayout{
	{"/", map[string]string{"de": "Home", "en": "Home"}},
	{"/excercise", map[string]string{"de": "Übungen", "en": "Excercises"}},
	{"/division", map[string]string{"de": "Division", "en": "Division"}},
	{"/zapfen", map[string]string{"de": "Webzapfen", "en": "zopfen"}},
}

type webhandler struct {
	handler      func(w io.Writer, req *http.Request, lang string) error
	compress     bool
	seccachetime int
}

func langselector(currlang string) template.HTML {
	var htmlString string

	for _, key := range conf_ISOlanguages() {
		if key != currlang {
			htmlString += fmt.Sprintf(`<a href="?setlanguage=%s">%s</a>`, key, conf_languages()[key])
		} else {
			htmlString += conf_languages()[key]
		}
	}

	return template.HTML(htmlString)
}

func rendermenu(currlang, url string) template.HTML {
	var htmlString string

	for _, menuitem := range menustructure {
		htmlString += fmt.Sprintf(`<li><a href="%s"`, menuitem.URL)
		if url == menuitem.URL {
			htmlString += ` class="current"`
		}
		htmlString += fmt.Sprintf(`>%s</a></li>`, menuitem.assocmenuitem[currlang])
	}
	return template.HTML(htmlString)
}

type rootPage struct {
	Error    []string
	CurrLang string
	URL      string
}

var templrootfuncMap = template.FuncMap{
	"langselector": langselector,
	"rendermenu":   rendermenu,
}

func rootHandler(w io.Writer, req *http.Request, lang string) error {

	var tpl *template.Template
	var err error

	page := &rootPage{CurrLang: lang, URL: path.Dir(req.URL.Path)}

	for retry := 0; retry <= 1; retry++ {
		tpl, err = template.New("DivisionSetup").Funcs(templrootfuncMap).ParseFiles(conf_roottemplatedir()+roottplfilename, conf_roottemplatedir()+lang+"."+indextplfilename)
		if err != nil {
			if os.IsNotExist(err) && retry < 1 {
				errstring := fmt.Sprintf("Template file for language '%s' not found, resorting to default language '%s'", lang, defaultlang)
				log.Printf(errstring)
				page.Error = append(page.Error, errstring)
				lang = defaultlang
			} else {
				panic(err)
			}
		}
	}

	return tpl.Execute(w, page)
}

type divisionPage struct {
	rootPage
	Dividend, Divisor, Precision string
	*schoolcalc.SDivide
	StopRemz, Boxed bool
}

func tplfuncdivdisplay(sd *schoolcalc.SDivide, boxed bool) (htmlResult template.HTML) {
	if sd != nil {
		dividendivisorresult := fmt.Sprintf("%s:%s=%s", sd.NormalizedDividend, sd.NormalizedDivisor, sd.Result)
		var column string
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
				column = fmt.Sprintf(`<div class="divisionColumn" data-division="true" data-boxed="%t">`, boxed)
			} else {
				column = `<div class="divisionColumn" data-result="true">`
			}

			if i < len(dividendivisorresult) {
				column += string(dividendivisorresult[i]) + "<br />"
			} else {
				column += "<br />"
			}

			for _, elm := range sd.DivisionSteps {
				if i >= elm.Indent && i < elm.Indent+len(elm.Iremainder) {
					column += string(elm.Iremainder[i-elm.Indent]) + "<br />"
				} else {
					column += "<br />"
				}
			}
			htmlResult += template.HTML(column) + "</div>"
		}
	}
	return
}

var templdivfuncMap = template.FuncMap{
	"rendermenu":        rendermenu,
	"langselector":      langselector,
	"tplfuncdivdisplay": tplfuncdivdisplay,
}

func divisionHandler(w io.Writer, req *http.Request, lang string) error {

	var tpl *template.Template
	var err error

	dividend := strings.TrimSpace(req.URL.Query().Get("dividend"))
	divisor := strings.TrimSpace(req.URL.Query().Get("divisor"))

	page := &divisionPage{Dividend: dividend, Divisor: divisor, Boxed: true, StopRemz: true, rootPage: rootPage{CurrLang: lang, URL: path.Dir(req.URL.Path)}}

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
		tpl, err = template.New("DivisionSetup").Funcs(templdivfuncMap).ParseFiles(conf_roottemplatedir()+roottplfilename, conf_roottemplatedir()+lang+"."+divisionfilename)
		if err != nil {
			if os.IsNotExist(err) && retry < 1 {
				errstring := fmt.Sprintf("Template file for language '%s' not found, resorting to default language '%s'", lang, defaultlang)
				log.Printf(errstring)
				page.Error = append(page.Error, errstring)
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
	rootPage
	Number string
	*schoolcalc.Zapfen
	IntermedZapfen [8]*schoolcalc.SDivide
	boxed          bool
}

func tplfunczapfendisplay(zapfen *schoolcalc.Zapfen, steps [8]*schoolcalc.SDivide) template.HTML {
	var retstring string
	if zapfen != nil {
		input := zapfen.Zapfenzahl

		for i := 0; i < 8; i++ {
			retstring += fmt.Sprintf("<tr>\n<td class='zapfenmultiplier'>%s</td><td>x</td><td>%d</td><td>=</td><td>%s</td>\n</tr>\n", input, i+2, zapfen.Multzapfen[i])
			input = zapfen.Multzapfen[i]
		}

		for i := 0; i < 8; i++ {
			retstring += fmt.Sprintf("\n<tr>\n  <td class='zapfendividend'><a href='javascript:void(0);' class='zapfendividenditem emptylink' id='zapfendividenditem%d'>%s</a></td><td>:</td><td>%d</td><td>=</td><td>%s</td>\n</tr>\n", i, input, i+2, zapfen.Divzapfen[i])

			var divisorintermed string = ""

			for dividendcolums := 0; dividendcolums < len(steps[i].NormalizedDividend); dividendcolums++ {

				column := `<div class="divisionColumn" data-division="true" data-boxed="true">`
				column += string(steps[i].NormalizedDividend[dividendcolums]) + "<br />"

				for _, elm := range steps[i].DivisionSteps {
					if dividendcolums >= elm.Indent && dividendcolums < elm.Indent+len(elm.Iremainder) {
						column += string(elm.Iremainder[dividendcolums-elm.Indent]) + "<br />"
					} else {
						column += "<br />"
					}
				}
				divisorintermed += column + "</div>"
			}

			retstring += fmt.Sprintf("\n<tr class='zapfenintermeddivisionrow' id='zapfenintermeddivisionrow%d'>\n<td class='zapfendividendintermed'>\n%s\n</td><td>:</td><td>%d</td><td>=</td><td>%s</td>\n</tr>", i, divisorintermed, i+2, zapfen.Divzapfen[i])
			input = zapfen.Divzapfen[i]
		}
	}
	return template.HTML(retstring)
}

var templzapfenfuncMap = template.FuncMap{
	"rendermenu":           rendermenu,
	"langselector":         langselector,
	"tplfunczapfendisplay": tplfunczapfendisplay,
}

func zapfenHandler(w io.Writer, req *http.Request, lang string) error {

	var tpl *template.Template
	var err error

	page := &zapfenPage{Number: strings.TrimSpace(req.URL.Query().Get("number")), rootPage: rootPage{CurrLang: lang, URL: path.Dir(req.URL.Path)}}

	for retry := 0; retry <= 1; retry++ {
		tpl, err = template.New("DivisionSetup").Funcs(templzapfenfuncMap).ParseFiles(conf_roottemplatedir()+roottplfilename, conf_roottemplatedir()+lang+"."+zapfenfilename)
		if err != nil {
			if os.IsNotExist(err) && retry < 1 {
				errstring := fmt.Sprintf("Template file for language '%s' not found, resorting to default language '%s'", lang, defaultlang)
				log.Printf(errstring)
				page.Error = append(page.Error, errstring)
				lang = defaultlang
			} else {
				panic(err)
			}
		}
	}

	if len(page.Number) >= 1 {
		if num, ok := big.NewInt(0).SetString(page.Number, 10); ok {
			result := schoolcalc.ZapfenRechnung(num)
			page.Zapfen = result

			dividend := result.Multzapfen[7].String()
			for i := 0; i < 8; i++ {
				divresult, _ := schoolcalc.SchoolDivide(dividend, strconv.Itoa(i+2), 0)
				page.IntermedZapfen[i] = divresult
				dividend = result.Divzapfen[i].String()
			}
		} else {
			page.Error = append(page.Error, fmt.Sprintf("Not a valid integer: '%s'", page.Number))
		}
	}

	return tpl.Execute(w, page)
}

type ExcersiseLevel byte
const (
	ExcersiseLevel_Beginner ExcersiseLevel = 1
	ExcersiseLevel_Apprentice
	ExcersiseLevel_Sophomore
	ExcersiseLevel_Advanced
	ExcersiseLevel_Master
)

type excersisePage struct {
	rootPage
	MinDividend, MaxDividend,
	MinDivisor, MaxDivisor string
	MaxDigitisPastPointUntilZero string
	NumberofExcersises           string
	Level                        ExcersiseLevel
}

func excersiseHandler(w io.Writer, req *http.Request, lang string) error {
	var tpl *template.Template
	var err error

	mindividend := strings.TrimSpace(req.URL.Query().Get("mindividend"))
	maxdividend := strings.TrimSpace(req.URL.Query().Get("maxdividend"))

	mindivisor := strings.TrimSpace(req.URL.Query().Get("mindivisor"))
	maxdivisor := strings.TrimSpace(req.URL.Query().Get("maxdivisor"))

	numberofexcersises := strings.TrimSpace(req.URL.Query().Get("n"))
	maxdigitispastpointuntilzero := strings.TrimSpace(req.URL.Query().Get("numremz"))

	page := &excersisePage{rootPage: rootPage{CurrLang: lang, URL: path.Dir(req.URL.Path)},
		MinDividend:                  mindividend,
		MaxDividend:                  maxdividend,
		MinDivisor:                   mindivisor,
		MaxDivisor:                   maxdivisor,
		NumberofExcersises:           numberofexcersises,
		MaxDigitisPastPointUntilZero: maxdigitispastpointuntilzero,
	}

	for retry := 0; retry <= 1; retry++ {
		tpl, err = template.New("DivisionSetup").Funcs(templzapfenfuncMap).ParseFiles(conf_roottemplatedir()+roottplfilename, conf_roottemplatedir()+lang+"."+excersisefilename)
		if err != nil {
			if os.IsNotExist(err) && retry < 1 {
				errstring := fmt.Sprintf("Template file for language '%s' not found, resorting to default language '%s'", lang, defaultlang)
				log.Printf(errstring)
				page.Error = append(page.Error, errstring)
				lang = defaultlang
			} else {
				panic(err)
			}
		}
	}
	return tpl.Execute(w, page)
}

func staticHTMLPageHandler(w io.Writer, req *http.Request, lang string) error {

	var tpl *template.Template
	var err error

	page := &rootPage{CurrLang: lang, URL: path.Dir(req.URL.Path)}

	directory := path.Dir(req.URL.Path)
	file := path.Base(req.URL.Path)

	for retry := 0; retry <= 1; retry++ {

		detailedfile := "." + directory + "/" + lang + "." + file

		tpl, err = template.New("DivisionSetup").Funcs(templrootfuncMap).ParseFiles(conf_roottemplatedir()+roottplfilename, detailedfile)
		if err != nil {
			if os.IsNotExist(err) && retry < 1 {
				errstring := fmt.Sprintf("Template file for language '%s' not found, resorting to default language '%s'", lang, defaultlang)
				log.Printf(errstring)
				page.Error = append(page.Error, errstring)
				lang = defaultlang
			} else {
				panic(err)
			}
		}
	}

	return tpl.Execute(w, page)
}

func validlanguage(language string) bool {
	for _, lang := range conf_ISOlanguages() {
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
			panic(fmt.Sprintf("Not a valid language '%s'. Valid languages are: %v", language, conf_ISOlanguages()))
		}
	} else {
		language = defaultlang
	}

	// read language cookie. If not set, initialize and set it with the language we determined so far,
	// either the default language or the language given as a request parameter
	langCookie, err := req.Cookie(cookie_lang)
	if err == http.ErrNoCookie || setcookie {
		langCookie = &http.Cookie{Name: cookie_lang, Value: language, Path: "/", Domain: "." + conf_rootdomain()}
		http.SetCookie(w, langCookie)
		// We set the cookie and redirect to the language subdomain
		mustredirect = true
	} else if !validlanguage(language) {
		panic(fmt.Sprintf("Not a valid language '%s'. Valid languages are: %v", language, conf_ISOlanguages()))
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
		http.Redirect(w, req, scheme+language+"."+conf_rootdomain()+conf_binding()+redirect, http.StatusSeeOther)
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
		gzwriter.Close() // otherwise the content won't get flushed to the output stream

		buf = gzbuf
	}
	// Prevent chunking of result
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))

	// set cache timeout
	if wh.seccachetime > 0 {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", wh.seccachetime))
	}
	buf.WriteTo(w)
}

// modelled after https://groups.google.com/d/msg/golang-nuts/n-GjwsDlRco/2l8kpJbAHHwJ
// another good read: http://www.mnot.net/cache_docs/
func maxAgeHandler(seconds int, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", seconds))
		h.ServeHTTP(w, r)
	})
}

func init() {
	http.Handle("/", &webhandler{rootHandler, true, 0})
	http.Handle("/static/", maxAgeHandler(conf_statictimeout(), http.StripPrefix("/static/", http.FileServer(http.Dir("static")))))
	http.Handle("/pages/", &webhandler{staticHTMLPageHandler, true, conf_statictimeout()})
	http.Handle("/zapfen/", &webhandler{zapfenHandler, true, 0})
	http.Handle("/division/", &webhandler{divisionHandler, true, 0})
	http.Handle("/excercise/", &webhandler{excersiseHandler, true, 0})
}
