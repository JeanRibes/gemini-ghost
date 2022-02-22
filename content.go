package main

import (
	"github.com/JeanRibes/gemini-ghost/ghost"
	"github.com/LukeEmmet/html2gemini"
	"github.com/blevesearch/bleve/v2"
	"log"
	"time"
)

type LocalDB struct {
	Posts    map[string]StoredPost
	Pages    map[string]StoredPost
	Settings ghost.Settings
}

var db LocalDB
var index bleve.Index

func init() {
	index, _ = bleve.NewMemOnly(bleve.NewIndexMapping())
}
var api *ghost.ContentAPI
func contentFetcher(baseurl string, key string, timeout_minutes int) {
	api = ghost.New(baseurl, key)
	for {
		fetchContent(api)
		if timeout_minutes < 0 {
			break
		}
		time.Sleep(time.Minute* time.Duration(timeout_minutes))
	}
}

func fetchContent(api *ghost.ContentAPI) {
	log.Println("fetching content from Ghost ...")
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

	indexposts(db.Posts)
	log.Println("fetched content from Ghost")
}

type StoredPost struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	PublishedAt time.Time `json:"published_at"`
	Content     string `json:"-"`

	Excerpt string `json:"excerpt"`
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

func indexposts(posts map[string]StoredPost) {
	for slug, post := range posts {
		if err := index.Index(slug, post); err != nil {
			panic(err)
		}
	}
}

func SearchPost(query string) (posts []StoredPost, err error){
	results, err := index.Search(bleve.NewSearchRequest(bleve.NewMatchQuery(query)))
	if err != nil {
		return nil,err
	}
	log.Println(results.String())


	for _, result := range results.Hits {
		println("id:", result.ID)
		if post, exists := db.Posts[result.ID]; exists {
			posts = append(posts, post)
		}
	}
	return posts,nil
}
