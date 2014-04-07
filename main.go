package main

import (
	"code.google.com/p/gcfg"
	"flag"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Server struct {
		Active      bool
		Host        string
		Port        string
		DbHost      string
		Database    string
		Collections []string
		UrlPrefix   string
	}

	Tag struct {
		Active bool
		Host   string
		Port   string
		Name   string
		Url    string
	}
}

/*
 * Global vars
 */

var db *mgo.Collection
var cfg Config

var server = flag.Bool("s", false, "Server Bool")
var tag = flag.Bool("t", false, "Tag Bool")

/*
 * Init config
 */

func Init() {
	err := gcfg.ReadFileInto(&cfg, "go-tag.ini")
	if err != nil {
		log.Print(err)
	}
	//	&cfg.Tag.Active = flag.Bool("t",false,"Tag Bool")
	log.Print(cfg)
}

func main() {
	Init()
	m := martini.Classic()
	m.Use(martini.Recovery())
	flag.Parse()

	if *server {
		m.Get(cfg.Server.UrlPrefix+":database/:collection", Server)
	}

	m.Use(render.Renderer(render.Options{
		HTMLContentType: "application/javascript",
	}))

	if *tag {
		m.Get(cfg.Tag.Url, func(r render.Render) {
			r.HTML(200, cfg.Tag.Name, &cfg.Tag)
		})
	}
	m.Run()
}

func Server(params martini.Params, r *http.Request) (int, string) {
	m := r.URL.Query()
	var v = false
	if params["database"] != cfg.Server.Database {
		return http.StatusNotFound, "404 page not found"
	}
	for _, c := range cfg.Server.Collections {
		if c == params["collection"] {
			v = true
			go Insertdata(params["client"], params["collection"], m)
		}
	}
	if !v {
		return http.StatusNotFound, "404 page not found"
	}
	return http.StatusOK, ""
}

func Insertdata(client string, collection string, m map[string][]string) {

	session, err := mgo.Dial(cfg.Server.DbHost)
	if err != nil {
		log.Print(err)
	}
	defer session.Close()
	db = session.DB(client).C(collection)
	if len(m) != 0 {
		m["insert_date"] = []string{time.Now().Format("2006-01-02 15:04:05")}
		m["source"] = []string{"GoTag"}
		err := db.Insert(m)
		if err != nil {
			log.Print(err)
		}
	}
}
