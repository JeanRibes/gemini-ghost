package main

import (
	"html/template"
	"io"
	"time"
)

func ghostResponse(conn io.ReadWriteCloser, path string) bool {
	post, exists := db.Posts[path]
	if !exists {
		post, exists = db.Pages[path]
	}
	if !exists {
		return false
	}
	sendResponseHeader(conn, statusSuccess, "text/gemini; lang=en; charset=utf-8")
	sendResponseContent(conn, []byte(post.Content))
	conn.Close()
	return true
}

func ghostIndex(conn io.ReadWriteCloser) bool {
	tmpl, err := template.ParseFiles("index.tpl")
	if err != nil {
		println(err.Error())
		return false
	}
	sendResponseHeader(conn, statusSuccess, "text/gemini; lang=en; charset=utf-8")
	err = tmpl.Execute(conn, map[string]interface{}{
		"Posts":    db.Posts,
		"Pages":    db.Pages,
		"Settings": db.Settings,
		"DateF": func(date time.Time) string {
			return date.Format(time.ANSIC)
		},
	})
	if err != nil {
		println(err.Error())
		return false
	}
	return true
}
