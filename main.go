package main

import (
	"github.com/blevesearch/bleve"
	"github.com/dustin/go-humanize"
	"golang.org/x/crypto/acme/autocert"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	err          error
	posts        []Post
	index        [5][]uint64
	user_counter uint64
	templates    map[string]*template.Template
	text_index   bleve.Index
)

func main() {
	text_index, err = bleve.Open("database/search.bleve")
	if err != nil {
		log.Fatal(err)
	}
	defer text_index.Close()

	templates = make(map[string]*template.Template)
	templates["base.html"] = template.Must(template.ParseFiles("templates/base.html")).Funcs(template.FuncMap{"human_time": humanize.Time})
	templates["home.html"] = template.Must(template.Must(templates["base.html"].Clone()).ParseFiles("templates/home.html"))
	templates["post.html"] = template.Must(template.Must(templates["base.html"].Clone()).ParseFiles("templates/post.html"))
	templates["search.html"] = template.Must(template.Must(templates["base.html"].Clone()).ParseFiles("templates/search.html"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", hsts(fresh_cookie(home)))
	mux.HandleFunc("/post/", hsts(fresh_cookie(post)))
	mux.HandleFunc("/search/", hsts(fresh_cookie(search)))
	mux.HandleFunc("/submit_post", hsts(user_cookie(submit_post)))
	mux.HandleFunc("/submit_comment", hsts(user_cookie(submit_comment)))
	mux.Handle("/static/", http.FileServer(http.Dir("./")))

    log.Println("server ready")
	go http.ListenAndServe(":80", http.HandlerFunc(https_redirect))
	go http.Serve(autocert.NewListener("commune.is"), mux)
	close := make(chan os.Signal, 2)
	signal.Notify(close, os.Interrupt, syscall.SIGTERM)
	<-close
}
