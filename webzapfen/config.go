// Copyright 2012  Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.
//
// configuration for webzapfen
//
// +build !appengine
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var configFileName = flag.String("config", "config.json", "location of JSON configuration file")

type config struct {
	RootDomain      string
	Binding         string
	RootTemplateDir string
	Languages       map[string]string
}

var conf = &config{RootDomain: "webzapfen.hoechtl.at", Binding: ":1112", Languages: map[string]string{"de": "Deutsch", "en": "Englisch"}, RootTemplateDir: "./templates/"}

func readConfig(filename string, conf *config) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if conf == nil {
		conf = &config{}
	}

	err = json.Unmarshal(b, &conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		panic("Unable to parse json configuration file")
	}
	return
}

func conf_rootdomain() string {
	flag.Parse()
	readConfig(*configFileName, conf)
	return conf.RootDomain
}

func conf_ISOlanguages() []string {
	flag.Parse()
	readConfig(*configFileName, conf)
	languages := []string{}
	for key, _ := range conf.Languages {
		languages = append(languages, key)
	}
	return languages
}

func conf_languages() map[string]string {
	flag.Parse()
	readConfig(*configFileName, conf)
	return conf.Languages
}

func conf_binding() string {
	flag.Parse()
	readConfig(*configFileName, conf)
	return conf.Binding
}

func conf_roottemplatedir() string {
	flag.Parse()
	readConfig(*configFileName, conf)
	return conf.RootTemplateDir
}
