package main

import (
	"net/http"
	"path"
	"strconv"
)

const page_length uint64 = 50

func home(w http.ResponseWriter, r *http.Request, freshness uint64) {
	start, err := strconv.ParseUint(r.FormValue("start"), 10, 64)
	if err != nil {
		start = uint64(0)
	}
	content := struct {
		Posts []Post
		Start uint64
		Next  uint64
		Prev  uint64
	}{
		Start: start,
	}
	if start > page_length {
		content.Prev = start - page_length
	} else if start <= page_length && start > 0 {
		content.Prev = 0
	} else {
		content.Prev = start
	}
	posts, more := get_posts(freshness, start, page_length)
	content.Posts = posts
	if more {
		content.Next = start + page_length
	} else {
		content.Next = start
	}
	err = template_home.Execute(w, Page{
		Title:     "commune",
		Content:   content,
		Freshness: freshness,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func post(w http.ResponseWriter, r *http.Request, freshness uint64) {
	_, str := path.Split(r.URL.Path)
	post_id, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	post, err := view_post(post_id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	err = template_post.Execute(w, Page{
		Title:     post.Title,
		Content:   post,
		Freshness: freshness,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func search(w http.ResponseWriter, r *http.Request, freshness uint64) {
	start, err := strconv.ParseUint(r.FormValue("start"), 10, 64)
	if err != nil {
		start = uint64(0)
	}
	search_query := r.FormValue("query")
	if search_query == "" {
		home(w, r, freshness)
		return
	}
	content := struct {
		Query   string
		Results []Post
		Start   uint64
		Next    uint64
		Prev    uint64
	}{
		Query: search_query,
		Start: start,
	}
	if start > page_length {
		content.Prev = start - page_length
	} else if start <= page_length && start > 0 {
		content.Prev = 0
	} else {
		content.Prev = start
	}
	results, more := text_search(search_query, freshness, start, page_length)
	content.Results = results
	if more {
		content.Next = start + page_length
	} else {
		content.Next = start
	}
	err = template_search.Execute(w, Page{
		Title:     search_query,
		Content:   content,
		Freshness: freshness,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
