package main

import (
	"fmt"
	//	"html/template"
	"bytes"
	"io"
	"net/http"
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

type webhandler func(w io.Writer, req *http.Request, lang string) error

func zapfenHandler(w io.Writer, req *http.Request, lang string) (err error) {
	_, err = fmt.Fprint(w, "In Zapfen")
	return
}

func rootHandler(w io.Writer, req *http.Request, lang string) (err error) {

	_, err = fmt.Fprintf(w, "Got language: %s", lang)

	/*tpl, err := template.ParseFiles(language + "." + roottplfilename)
	if err != nil {
		panic(err)
	}

	err = tpl.Execute(w, nil)
	*/
	return
}

func frHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Bonjour1")
	return
}

// ServeHTTP installs a catch-all error recovery for the specific handler functions
func (fn webhandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// error handler 
	// Recover from panic by setting http error 500 and letting the user know the reason
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		}
	}()

	// read language cookie. If not set, initialize and set it with the default language,
	// then redirect to the language subdomain
	var language string
	var setcookie bool
	var mustredirect bool

	if language = req.URL.Query().Get(setlanguageparam); language != "" {
		setcookie = true
	} else {
		language = defaultlang
	}

	langCookie, err := req.Cookie(cookie_lang)
	if err == http.ErrNoCookie || setcookie {
		langCookie = &http.Cookie{Name: cookie_lang, Value: language, Path: "/", Domain: "." + rootdomain}
		http.SetCookie(w, langCookie)
		mustredirect = true
	} else if err != nil {
		panic(err)
	}

	language = langCookie.Value

	if strings.Index(req.Host, language) != 0 || mustredirect {
		var redirect string
		if redirect = req.URL.Query().Get(redirectparam); redirect == "" {
			redirect = req.URL.Path
		}

		scheme := req.URL.Scheme
		if scheme == "" {
			scheme = "http://"
		}
		// functionality required? If yes, should redirect relatively or absolutely?
		http.Redirect(w, req, scheme+language+"."+rootdomain+applicationport+redirect, http.StatusSeeOther)
		return
	}

	buf := new(bytes.Buffer)
	if err := fn(buf, req, language); err != nil {
		// might as well panic(err) but we add some more info
		// we  serialize the error here in the chosen encoding
		fmt.Fprint(buf, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	buf.WriteTo(w)
}

func main() {

	http.Handle("/", webhandler(rootHandler))
	http.Handle("/zapfen/", webhandler(zapfenHandler))
	http.HandleFunc("fr.webzapfen.hoechtl.at:1112/", frHandler)
	http.ListenAndServe(applicationport, nil)

}
