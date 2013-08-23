package main

import (
	"fmt"
	"net/http"
	"labix.org/v2/mgo"
	"time"
	"flag"
	"text/template"
)
/* ********************************************
 *            VARIABLES GLOBALES              *
 * ********************************************/
var collection string
var database string

var port = flag.String("p","9090","Listening port")
var js = flag.String("js","both","Using or not JS Loader / HTTP Listener . ( no-js / no-go / both )")
/* ********************************************
 *                FUNCTIONS                   *
 * ********************************************/

func main() {
	flag.Parse()
	listenport := fmt.Sprint(":",*port)
	fmt.Println("Listening ", listenport)

	if *js != "no-go" {
		http.HandleFunc("/tag", TagHandler)
	}
	if *js != "no-js" {
		http.HandleFunc("/loader.js",LoaderHandler)
	}

	http.ListenAndServe(listenport,nil)
}

func TagHandler(w http.ResponseWriter, r *http.Request) {
	m := make( map [string][]string)
	m = r.URL.Query()
	go Insertdata(m)
	return
}

func LoaderHandler (w http.ResponseWriter, r *http.Request) {
	loader, _ :=  template.ParseFiles("templates/loader.tmpl")
	data := make(map[string]interface{})
	data["domain"] = "jonathanschmidt.fr:9090"
	data["name"] = "visite"
	data["prefix"] = "gotag_"
	loader.Execute(w,data)
}

func Insertdata(m map[string][]string) {
	database := "gotest"
	collection := "visite"
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(database).C(collection)
	if len(m)!=0 {
		m["insert_date"] = []string {time.Now().Format("2006-01-02 15:04:05")}
		m["source"] = []string {"GoTag"}
		err = c.Insert(m)
		if err != nil {
			panic(err)
		}
	}
}
