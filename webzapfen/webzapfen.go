package main

import (
	"fmt"
	"html/template"
	"net/http"
)

const rootfiletpl = "index.tpl"
const defaultlang = "en"
const cookie_lang = "webzapfen.lang"

func zapfenHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "In Zapfen")
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	// read language cookie. If not set, initialize and set it with the default language,
	// then redirect to the language subdomain

	langCookie, err := req.Cookie(cookie_lang)
	if err == http.ErrNoCookie {
		langCookie = &http.Cookie{Name: cookie_lang, Value: defaultlang, Path: "/", Domain: "localhost"}
		http.SetCookie(w, langCookie)
	} else if err != nil {
		panic(err)
	}

	// if we have a cookie defining a language, but no corresponding language domain, we redirect
	mustredirect := true
	if langCookie != nil && mustredirect {
		http.Redirect(w, req, defaultlang, http.StatusSeeOther)
		return
	}
	fmt.Fprint(w, "Got cookie: %v", langCookie)

	tpl, err := template.ParseFiles(rootfiletpl + "_" + defaultlang + ".tpl")
	if err != nil {
		panic(err)
	}

	err = tpl.Execute(w, nil)

}

func main() {

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/zapfen", zapfenHandler)
	http.ListenAndServe(":1112", nil)
}
