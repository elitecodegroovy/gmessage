package main

import (
	"net/http"
	"fmt"
	"log"
	"mux"
	//"html/template"
	"encoding/json"
	"encoding/xml"
	"time"
	"html/template"
)

type Todo struct {
	Task string
	Done bool
}

type User struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
}

type Profile struct {
	Name    string
	Hobbies []string
}

func HandleSimpleJsonXml(w http.ResponseWriter, req *http.Request){
	profile := Profile{"John", []string{"Play football", "Running"}}
	renderType := req.URL.Query().Get("respType")
	if renderType == "xml" {
		x, err := xml.MarshalIndent(profile, "", "  ")
		if err != nil {
			log.Fatal(" fail to parse xml object instance", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		w.Write(x)
	} else {
		js , err := json.Marshal(&profile)
		if err != nil {
			log.Fatal(" fail to parse json object instance", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

	w.Header().Set("Server", "A Intelligent Robot Server")
}

func doHttp(){
	log.Println("start server ....")
	asset, _ := Asset("index.html")
	tmpl, _ := template.New("index").Parse(string(asset))


	todos := []Todo{
		{"Mastering.Go.Web.Services", true},
		{"Golagn In Action", false},
		{"The Go Programming Language", true},
	}

	//routine
	userAges := map[string]int {
		"Bob": 23,
		"John":30,
		"Lucy":25,
	}
	//http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request){
	//	user := r.URL.Path[len("/user/"):]
	//	age := userAges[user]
	//	fmt.Fprintf(w, "%s is %d years old!", user, age)
	//})
	//http.ListenAndServe(":9990", nil)

	r:= mux.NewRouter()
	r.HandleFunc("/user/{name}", func(w http.ResponseWriter, req *http.Request){
		vars := mux.Vars(req)
		userName := vars["name"]
		age := userAges[userName]
		fmt.Fprintf(w, "%s is %d years old!", userName, age)
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, struct{ Todos[]Todo }{todos})
	})

	r.HandleFunc("/decode", func(w http.ResponseWriter, r *http.Request){
		var user User
		json.NewDecoder(r.Body).Decode(&user)

		fmt.Fprintf(w, "%s %s is %d years old!", user.Firstname, user.Lastname, user.Age)
	})

	r.HandleFunc("/encode", func(w http.ResponseWriter, r *http.Request){
		john := User{
			Firstname:"John",
			Lastname:"Lau",
			Age: 30,
		}

		json.NewEncoder(w).Encode(john)
	})

	//json rendering
	r.HandleFunc("/simpleReq", HandleSimpleJsonXml)
	srv := &http.Server{
		Handler:      r,
		Addr:         ":9990",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}


//func main() {
//	doHttp()
//}
