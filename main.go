package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"github.com/LukeEmmet/html2gemini"
	"html/template"
	"io"
	"log"
	_ "modernc.org/sqlite"
	"net"
	"net/url"
	"strings"
	"time"
)

type Post struct {
	title        string
	slug         string
	published_at time.Time
}

const (
	statusInput            = 10
	statusSuccess          = 20
	statusRedirectTemp     = 30
	statusTemporaryFailure = 40
	statusPermanentFailure = 50
)

var (
	hostname    = flag.String("hostname", "localhost", "hostname")
	crtFilename = flag.String("crt", "./certs/crt.pem", "cert filename")
	keyFilename = flag.String("key", "./certs/key.pem", "key filename")
	port        = flag.Int("port", 1965, "port number")
	dbfile      = flag.String("dbfile", "ghost.db", "SQLITE file to use")
)

var h2gCtx *html2gemini.TextifyTraverseContext

func main() {
	flag.Parse()

	// Load TLS certificate - crt.pem and key.pem are the public and private key
	// parts of a TLS certificate.
	cert, err := tls.LoadX509KeyPair(*crtFilename, *keyFilename)
	if err != nil {
		log.Fatalf("Unable to load TLS certficate: %s", err)
	}

	h2gCtx = html2gemini.NewTraverseContext(*html2gemini.NewOptions())

	db, dconn := openDb()
	dconn.Close()
	db.Close()

	// Create TSL over TCP session.
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}, ServerName: *hostname, MinVersion: tls.VersionTLS12}
	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", *port), cfg)
	if err != nil {
		log.Fatalf("Unable to listen: %s", err)
	}
	log.Printf("Listening for connections on port: %d", *port)

	serveGemini(listener)
}

func serveGemini(listener net.Listener) {
	for {
		// Accept incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn io.ReadWriteCloser) {
	defer conn.Close()

	// Check the size of the request buffer.
	s := bufio.NewScanner(conn)
	if len(s.Bytes()) > 1024 {
		sendResponseHeader(conn, statusPermanentFailure, "Request exceeds maximum permitted length")
		return
	}

	// Sanity check incoming request URL content.
	if ok := s.Scan(); !ok {
		sendResponseHeader(conn, statusPermanentFailure, "Request not valid")
		return
	}

	// Parse incoming request URL.
	reqURL, err := url.Parse(s.Text())
	log.Printf("request path: %s\n", reqURL.Path)
	if err != nil {
		sendResponseHeader(conn, statusPermanentFailure, "URL incorrectly formatted")
		return
	}

	if reqURL.Path == "/" || reqURL.Path == "/index.gmi" || reqURL.Path == "" {
		ghostIndex(conn)
		conn.Close()
		return
	}

	if ghostResponse(conn, strings.Trim(reqURL.Path, "/")) {
		return
	} else {
		sendResponseHeader(conn, statusPermanentFailure, "no content found at this address")
		conn.Close()
	}
}

func sendResponseHeader(conn io.ReadWriteCloser, statusCode int, meta string) {
	header := fmt.Sprintf("%d %s\r\n", statusCode, meta)
	_, err := conn.Write([]byte(header))
	if err != nil {
		log.Printf("There was an error writing to the connection: %s", err)
	}
}

func sendResponseContent(conn io.ReadWriteCloser, content []byte) {
	_, err := conn.Write(content)
	if err != nil {
		log.Printf("There was an error writing to the connection: %s", err)
	}
}

func openDb() (*sql.DB, *sql.Conn) {
	db, err := sql.Open("sqlite", *dbfile)
	if err != nil {
		log.Panicf("error opening database %s: %s\n", *dbfile, err)
	}

	conn, err := db.Conn(context.TODO())
	if err != nil {
		log.Panicf("error opening database %s: %s\n", *dbfile, err)
	}
	return db, conn
}

func ghostResponse(conn io.ReadWriteCloser, path string) bool {
	db, dconn := openDb()
	defer dconn.Close()
	defer db.Close()

	rows, err := dconn.QueryContext(context.TODO(),
		"SELECT html FROM posts WHERE status='published' AND type ='post' AND slug=?", path)
	if err != nil {
		log.Fatal("error", err)
	}

	rows.Next()
	var html string
	err = rows.Scan(&html)
	if err != nil {
		println(err.Error())
		return false
	}
	gmi, err := html2gemini.FromReader(strings.NewReader(html), *h2gCtx)
	if err != nil {
		return false
	}
	sendResponseHeader(conn, statusSuccess, "text/gemini; lang=en; charset=utf-8")
	sendResponseContent(conn, []byte(gmi))
	conn.Close()
	return true
}

type IndexData struct {
	Posts []IndexPost
}

type IndexPost struct {
	Slug  string
	Title string
	Date  string
}

func ghostIndex(conn io.ReadWriteCloser) bool {
	db, dconn := openDb()
	defer dconn.Close()
	defer db.Close()

	rows, err := dconn.QueryContext(context.TODO(),
		"SELECT slug,title,published_at FROM posts WHERE status='published' AND type ='post' ORDER BY published_at DESC ")
	if err != nil {
		log.Fatal("error", err)
	}
	posts := []IndexPost{}
	for rows.Next() {
		var slug string
		var title string
		var published_at string
		err = rows.Scan(&slug, &title, &published_at)
		if err != nil {
			println(err.Error())
			return false
		}
		date, _ := time.Parse("2006-01-02T15:04:05Z", published_at)
		if err != nil {
			return false
		}
		year, month, day := date.Date()
		post := IndexPost{
			Slug:  slug,
			Title: title,
			Date:  fmt.Sprintf("%d-%d-%d", year, month, day),
		}
		posts = append(posts, post)
	}

	tmpl, err := template.ParseFiles("index.tpl")
	if err != nil {
		println(err.Error())
		return false
	}
	sendResponseHeader(conn, statusSuccess, "text/gemini; lang=en; charset=utf-8")
	err = tmpl.Execute(conn, IndexData{Posts: posts})
	if err != nil {
		println(err.Error())
		return false
	}
	return true
}
