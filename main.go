package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	_ "modernc.org/sqlite"
	"net"
	"net/url"
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
	port        = flag.Int("port", 1965, "port number")

	ghostUrl = flag.String("ghost-url", "http://localhost:2368/ghost/api/v4/content", "Ghost Content API Url")
	ghostKey = flag.String("ghost-key", "a513a3dc949855fb654a545bd7", "Ghost Content API Key")
)

func main() {
	flag.Parse()
	go contentFetcher(*ghostUrl, *ghostKey)

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
