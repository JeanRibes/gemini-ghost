package main

import (
	"github.com/JeanRibes/gemini-ghost/ghost"
	"github.com/LukeEmmet/html2gemini"
	"log"
	"time"
)

type LocalDB struct {
	Posts    map[string]StoredPost
	Pages    map[string]StoredPost
	Settings ghost.Settings
}

var db LocalDB

func contentFetcher(baseurl string, key string) {
	api := ghost.New(baseurl, key)
	for {
		posts, err := api.AllPosts()
		if err != nil {
			log.Println(err)
		}
		pages, err := api.AllPages()
		if err != nil {
			log.Println(err)
		}
		settings, err := api.Settings()
		if err != nil {
			log.Println(err)
		}
		if err == nil {
			db.Posts = convertposts(posts)
			db.Pages = convertpages(pages)
			db.Settings = *settings
		}

		time.Sleep(1 * time.Hour)
	}
}

type StoredPost struct {
	Slug        string
	Title       string
	PublishedAt time.Time
	Content     string

	Excerpt string
}

func convertposts(posts []ghost.Post) map[string]StoredPost {
	out := map[string]StoredPost{}
	for _, post := range posts {
		out[post.Slug] = convertpost(post)
	}
	return out
}

func convertpages(pages []ghost.Page) map[string]StoredPost {
	out := map[string]StoredPost{}
	for _, page := range pages {
		out[page.Slug] = convertpage(page)
	}
	return out
}

func convertpost(post ghost.Post) StoredPost {
	stored := StoredPost{
		Slug:        post.Slug,
		Title:       post.Title,
		PublishedAt: post.PublishedAt,
		Excerpt:     post.Excerpt,
	}
	gemtext, err := html2gemini.FromString(post.HTML, *html2gemini.NewTraverseContext(*html2gemini.NewOptions()))
	if err != nil {
		log.Println(err)
	}
	stored.Content = gemtext
	return stored
}

func convertpage(post ghost.Page) StoredPost {
	stored := StoredPost{
		Slug:        post.Slug,
		Title:       post.Title,
		PublishedAt: post.PublishedAt,
		Excerpt:     post.Excerpt,
	}
	gemtext, err := html2gemini.FromString(post.HTML, *html2gemini.NewTraverseContext(*html2gemini.NewOptions()))
	if err != nil {
		log.Println(err)
	}
	stored.Content = gemtext
	return stored
}
