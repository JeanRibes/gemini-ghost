package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

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
	contentDir  = flag.String("dir", "./gemini", "content directory")
	port        = flag.Int("port", 1965, "port number")
)

func main() {
	flag.Parse()

	// Load TLS certificate - crt.pem and key.pem are the public and private key
	// parts of a TLS certificate.
	cert, err := tls.LoadX509KeyPair(*crtFilename, *keyFilename)
	if err != nil {
		log.Fatalf("Unable to load TLS certficate: %s", err)
	}

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
		log.Println("Accept connection")

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
	if err != nil {
		sendResponseHeader(conn, statusPermanentFailure, "URL incorrectly formatted")
		return
	}

	// If the URL ends with a '/' character, assume that the user wants the index.gmi
	// file in the corresponding directory.
	var reqPath string
	if strings.HasSuffix(reqURL.Path, "/") || reqURL.Path == "" {
		reqPath = filepath.Join(reqURL.Path, "index.gmi")
	} else {
		reqPath = reqURL.Path
	}
	cleanPath := filepath.Clean(reqPath)

	// If the content directory is not specified as an absolute path, make it absolute.
	var workDir string
	var rootDir http.Dir
	if !strings.HasPrefix(*contentDir, "/") {
		workDir, _ = os.Getwd()
		// Use this function to avoid directory traversal type attacks.
		rootDir = http.Dir(workDir + strings.Replace(*contentDir, ".", "", -1))
	} else {
		rootDir = http.Dir(strings.Replace(*contentDir, ".", "", -1))
	}

	// Open the requested resource.
	log.Printf("Path: %s", cleanPath)
	f, err := rootDir.Open(cleanPath)
	if err != nil {
		sendResponseHeader(conn, statusPermanentFailure, "Resource not found")
		return
	}
	defer f.Close()

	// Read the contents of the file.
	content, err := ioutil.ReadAll(f)
	if err != nil {
		sendResponseHeader(conn, statusPermanentFailure, "Resource could not be read")
		return
	}

	// Determine MIME type.
	meta := http.DetectContentType(content)
	if strings.HasSuffix(cleanPath, ".gmi") {
		meta = "text/gemini; lang=en; charset=utf-8"
	}

	log.Println("Write response header")
	sendResponseHeader(conn, statusSuccess, meta)

	log.Println("Write content")
	sendResponseContent(conn, content)

	log.Println("Close connection")
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
