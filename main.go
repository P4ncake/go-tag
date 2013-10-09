package main

import (
	"flag"
	"fmt"
	"github.com/sbinet/go-config/config"
	"labix.org/v2/mgo"
	"net/http"
	"text/template"
	"time"
)

/* ********************************************
 *            VARIABLES GLOBALES              *
 * ********************************************/
var collection string
var database string
var baseUrl string
var loaderName string
var loaderUrl string
var tagName string
var tagUrl string
var port string
var dbhost string

var db *mgo.Collection

/* ********************************************
 *                FUNCTIONS                   *
 * ********************************************/

func main() {

	Getconf()

	// DB Connection
	session, err := mgo.Dial(dbhost)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	db = session.DB(database).C(collection)

	// Here goes the magic
	flag.Parse()
	listenport := fmt.Sprint(":", port)
	fmt.Println("Listening ", listenport)

	if tagUrl != "" {
		http.HandleFunc(tagUrl, TagHandler)
	}
	if loaderUrl != "" {
		http.HandleFunc(loaderUrl, LoaderHandler)
	}

	http.ListenAndServe(listenport, nil)
}

func Getconf() {

	// Reading configuration file
	c, _ := config.ReadDefault("go-tag.cfg")

	// Get configuration variables
	baseUrl, _ = c.String("default", "base-url")
	port, _ = c.String("default", "port")
	loaderName, _ = c.String("default", "loader-name")
	loaderUrl, _ = c.String("default", "loader-url")
	tagName, _ = c.String("default", "tag-name")
	tagUrl, _ = c.String("default", "tag-url")
	database, _ = c.String("default", "database")
	collection, _ = c.String("default", "collection")
	dbhost, _ = c.String("default", "dbhost")
}

func TagHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string][]string)
	m = r.URL.Query()
	go Insertdata(m)
	return
}

func LoaderHandler(w http.ResponseWriter, r *http.Request) {
	loader, _ := template.ParseFiles("templates/loader.tmpl")
	data := make(map[string]interface{})
	data["domain"] = "jonathanschmidt.fr:9090"
	data["name"] = "visite"
	data["prefix"] = "gotag_"
	loader.Execute(w, data)
}

func Insertdata(m map[string][]string) {
	if len(m) != 0 {
		m["insert_date"] = []string{time.Now().Format("2006-01-02 15:04:05")}
		m["source"] = []string{"GoTag"}
		err := db.Insert(m)
		if err != nil {
			panic(err)
		}
	}
}
