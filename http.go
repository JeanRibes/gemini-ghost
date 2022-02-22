package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

var cors = "*"

func serveHttp(ghostUrl string){
	gurl, err := url.Parse(ghostUrl)
	if err != nil {
		log.Println(err)
	}
	cors = gurl.Host

	http.Handle("/static/",http.StripPrefix("/static/",http.FileServer(http.Dir("static"))))
	http.HandleFunc("/webhook/rebuild", func(w http.ResponseWriter, req *http.Request) {
		log.Println("webhook: fetching ghost content")
		fetchContent(api)
		log.Println("webhook: finished")
		w.WriteHeader(200)
	})

	http.HandleFunc("/search", httpSearch)

	http.HandleFunc("/", httpIndex)

	if err := http.ListenAndServe(":1980",nil); err != nil {
		log.Println(err)
	}
}


func httpSearch(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Add("Access-Control-Allow-Origin", cors)
	if err := request.ParseForm(); err != nil {
		log.Println(err)
		responseWriter.Write([]byte("bad input"))
	}
	query := request.Form.Get("query")
	if query=="" {
		responseWriter.Write([]byte("empty query"))
		return
	}
	log.Println("HTTP query: ", query)

	//------------

	posts,_ := SearchPost(query)
	responseWriter.Header().Add("Content-Type", "application/json")
	responseWriter.WriteHeader(200)


	json.NewEncoder(responseWriter).Encode(map[string]interface{}{
		"query": query,
		"posts": posts,
	})
}
func httpIndex(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	writer.WriteHeader(200)
	writer.Write([]byte(`<html>
<head><title>Search ghost</title></head>
<body>
<h1>Search ghost posts</h1>
<form action="/search" method="get">
	<label>Query : <input name="query"/></label>
	<button>Seach</button>
</form>
</body>
</html>`))
}
