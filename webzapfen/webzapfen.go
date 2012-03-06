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
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	roottplfilename  = "index.tpl"
	defaultlang      = "en"
	cookie_lang      = "uilang"
	rootdomain       = "webzapfen.hoechtl.at"
	applicationport  = ":1112"
	setlanguageparam = "setlanguage"
	redirectparam    = "redirect"
)

type webhandler struct {
	handler  func(w io.Writer, req *http.Request, lang string) error
	compress bool
}

func zapfenHandler(w io.Writer, req *http.Request, lang string) (err error) {
	_, err = fmt.Fprint(w, "In Zapfen")
	return
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
		// if the language was specified as a paramter, we have to set a cookie	
		setcookie = true
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
		fmt.Fprint(buf, err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	// gzip-compression of the output of the given webhandler
	if wh.compress && strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		gzbuf := new(bytes.Buffer)
		gzwriter := gzip.NewWriter(gzbuf)
		defer gzwriter.Close() // Otherwise the content won't get flushed to the output stream

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", http.DetectContentType(buf.Bytes())) // We have to set the content type, otherwise the ResponseWriter will guess it's application/x-gzip
		gzwriter.Write(buf.Bytes())

		buf = gzbuf
	}
	// Prevent chunking of result
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	buf.WriteTo(w)
}

func main() {
	http.Handle("/", &webhandler{rootHandler, true})
	http.Handle("/zapfen/", &webhandler{zapfenHandler, true})
	http.ListenAndServe(applicationport, nil)
}
