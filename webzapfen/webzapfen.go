package main

import (
	"fmt"
	//	"html/template"
	"net/http"
	"strings"
)

const roottplfilename = "index.tpl"
const defaultlang = "en"
const cookie_lang = "uilang"
const rootdomain = "webzapfen.hoechtl.at"
const applicationport = ":1112"
const setlanguageparam = "setlanguage"
const redirectparam = "redirect"

func zapfenHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "In Zapfen")
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	// read language cookie. If not set, initialize and set it with the default language,
	// then redirect to the language subdomain

	var language string = "en"
	var setcookie bool
	var mustredirect bool

	if languages, paramset := req.URL.Query()[setlanguageparam]; paramset {
		language = languages[0]
		setcookie = true
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
		if redirects, paramset := req.URL.Query()[redirectparam]; paramset {
			redirect = redirects[0]
		}
		http.Redirect(w, req, "http://"+language+"."+rootdomain+applicationport+redirect, http.StatusSeeOther)
		return
	}

	fmt.Fprintf(w, "Got cookie: %v", langCookie)

	/*tpl, err := template.ParseFiles(language + "." + roottplfilename)
	if err != nil {
		panic(err)
	}

	err = tpl.Execute(w, nil)
	*/

}

func main() {

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/zapfen/", zapfenHandler)
	http.ListenAndServe(applicationport, nil)
}
