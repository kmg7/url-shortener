package main

import (
	"flag"
	"fmt"
	"net/http"

	"slices"
)

var paths string
var mode string
var port string
var address string
var raw bool

func setFlags() {
	flag.StringVar(&paths, "u", "", "Paths to redirect their associated urls")
	flag.BoolVar(&raw, "r", true, "Raw data mode")
	flag.StringVar(&mode, "m", "json", "File mod to work with json, yaml, csv")
	flag.StringVar(&port, "p", "8080", "Port to serve")
	flag.StringVar(&address, "ip", "0.0.0.0", "IP Address to serve")
	flag.Parse()

}

func validateFlags() (bool, string) {
	var availableMods = []string{"json", "yaml", "csv"}

	if paths == "" {
		return false, "\nRedirect urls must be set, \nto set urls use flag '-u' \n"
	}

	if !slices.Contains(availableMods, mode) {
		return false, "Mode can only be json, yaml or csv. To set mod use '-m' option"
	}
	if port == "" {
		return false, "Port must be set with '-p' option"
	}

	if address == "" {
		return false, "A valid ip address must be set with '-a' option"
	}
	return true, ""

}

func main() {
	setFlags()
	ok, msg := validateFlags()
	if !ok {
		panic(msg)
	}
	mux := defaultMux()

	fmt.Printf("Mode: %v\nStarting the server\nAddress: %v:%v\n", mode, address, port)
	err := http.ListenAndServe(address+":"+port, RedirectHandler(mode, paths, raw, mux))
	if err != nil {
		panic(err)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Not found :(<h1/>")
}
