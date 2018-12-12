package server

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/mnordsletten/lotto_dashboard/lottotest"
)

var resultStore = &lottotest.ResultStore{}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("main handler", r.RequestURI)
	rs := &lottotest.ResultStore{}
	*rs = *resultStore
	// Queries
	filter := lottotest.Identifier{}
	if mv, ok := r.URL.Query()["mothershipVersion"]; ok && len(mv) > 0 {
		if mv[0] != "all" {
			filter.MothershipVersion = mv[0]
		}
	}
	if mv, ok := r.URL.Query()["includeosVersion"]; ok && len(mv) > 0 {
		if mv[0] != "all" {
			filter.IncludeOSVersion = mv[0]
		}
	}
	if mv, ok := r.URL.Query()["environment"]; ok && len(mv) > 0 {
		if mv[0] != "all" {
			filter.Environment = mv[0]
		}
	}
	if mv, ok := r.URL.Query()["testName"]; ok && len(mv) > 0 {
		if mv[0] != "all" {
			filter.TestName = mv[0]
		}
	}
	rs.ResultIDs = rs.FilterIDs(filter)

	viewTemplate := path.Join("templates", "view.html")
	t, _ := template.ParseFiles(viewTemplate)
	t.Execute(w, rs)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	result := lottotest.TestResult{}

	if err := decoder.Decode(&result); err != nil {
		fmt.Printf("Could not decode json: %v\n", err)
		return
	}
	fmt.Println("received: ", result)
	resultStore.AddResult(result)
}

func Serve(port int) {
	resultStore = lottotest.NewResultStore()
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "There is no icon, %q", html.EscapeString(r.URL.Path))
	})
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
