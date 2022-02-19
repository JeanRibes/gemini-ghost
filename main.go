package main

import (
	"flag"
	"github.com/pitr/gig"
	"log"
)

var (
	crtFilename = flag.String("crt", "./certs/crt.pem", "cert filename")
	keyFilename = flag.String("key", "./certs/key.pem", "key filename")
	addr        = flag.String("addr", ":1965", "address to listen to")

	ghostUrl = flag.String("ghost-url", "http://localhost:2368/ghost/api/v4/content", "Ghost Content API Url")
	ghostKey = flag.String("ghost-key", "a513a3dc949855fb654a545bd7", "Ghost Content API Key")
)

func main() {
	flag.Parse()
	go contentFetcher(*ghostUrl, *ghostKey)

	g := gig.Default()

	g.Handle("/search", searchPost)
	g.Handle("/:slug", ghostContent)
	g.Handle("/", ghostIndex)

	gig.Debug = false

	if err := g.Run(*addr, *crtFilename, *keyFilename); err != nil {
		log.Print(err)
	}
}

func input_helper(ctx gig.Context, help string) string {
	input, err := ctx.QueryString()
	if err != nil {
		ctx.Error(err)
	}
	if input == "" {
		if err := ctx.NoContent(gig.StatusInput, help); err != nil {
			ctx.Error(err)
		}
	}
	return input
}
