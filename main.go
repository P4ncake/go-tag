package main

import (
	"fmt"
	"net/http"
    "labix.org/v2/mgo"
	"time"
)
/* ********************************************
 *                   TYPES                    *
 * ********************************************/
type Handler struct{}
/* ********************************************
 *            VARIABLES GLOBALES              *
 * ********************************************/
var collection string
var database string
/* ********************************************
 *                FUNCTIONS                   *
 * ********************************************/

func main() {
	fmt.Println("Listening 9090")
	http.ListenAndServe(":9090", new(Handler))
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := make( map [string][]string)
	m = r.URL.Query()
	go Insertdata(m)
	return
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
