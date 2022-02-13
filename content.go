package main

import (
	"encoding/json"
	"fmt"
	"github.com/LukeEmmet/html2gemini"
	"log"
	"net/http"
	"os"
	"time"
)

const API_KEY = "API_KEY"
const URL = "URL"

type GhostPost struct {
	ID                   string      `json:"id"`
	UUID                 string      `json:"uuid"`
	Title                string      `json:"title"`
	Slug                 string      `json:"slug"`
	HTML                 string      `json:"html"`
	CommentID            string      `json:"comment_id"`
	FeatureImage         string      `json:"feature_image"`
	Featured             bool        `json:"featured"`
	Visibility           string      `json:"visibility"`
	EmailRecipientFilter string      `json:"email_recipient_filter"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
	PublishedAt          time.Time   `json:"published_at"`
	CustomExcerpt        string      `json:"custom_excerpt"`
	CodeinjectionHead    interface{} `json:"codeinjection_head"`
	CodeinjectionFoot    interface{} `json:"codeinjection_foot"`
	CustomTemplate       interface{} `json:"custom_template"`
	CanonicalURL         interface{} `json:"canonical_url"`
	URL                  string      `json:"url"`
	Excerpt              string      `json:"excerpt"`
	ReadingTime          int         `json:"reading_time"`
	Access               bool        `json:"access"`
	OgImage              interface{} `json:"og_image"`
	OgTitle              interface{} `json:"og_title"`
	OgDescription        interface{} `json:"og_description"`
	TwitterImage         interface{} `json:"twitter_image"`
	TwitterTitle         interface{} `json:"twitter_title"`
	TwitterDescription   interface{} `json:"twitter_description"`
	MetaTitle            interface{} `json:"meta_title"`
	MetaDescription      interface{} `json:"meta_description"`
	EmailSubject         interface{} `json:"email_subject"`
	Frontmatter          interface{} `json:"frontmatter"`
	FeatureImageAlt      interface{} `json:"feature_image_alt"`
	FeatureImageCaption  interface{} `json:"feature_image_caption"`
	Plaintext            string      `json:"plaintext,omitempty"`
}

type GhostContent struct {
	Posts []GhostPost `json:"posts"`
	Meta  struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
			Next  int `json:"next"`
			Prev  int `json:"prev"`
		} `json:"pagination"`
	} `json:"meta"`
}

type StoredPost struct {
	Slug        string
	Title       string
	PublishedAt time.Time
	Content     string

	Excerpt string
}

func convertpost(post *GhostPost) *StoredPost {
	stored := StoredPost{
		Slug:        post.Slug,
		Title:       post.Title,
		PublishedAt: post.PublishedAt,
		Excerpt:     post.Excerpt,
	}
	gemtext, err := html2gemini.FromString(post.HTML, *html2gemini.NewTraverseContext(*html2gemini.NewOptions()))
	if err != nil {
		panic(err)
	}
	stored.Content = gemtext
	return &stored
}
func fetchcontent() map[string]*StoredPost {
	key := os.Getenv(API_KEY)
	localDb := map[string]*StoredPost{}
	url := "http://localhost:2368/ghost/api/v4/content/posts/"
	if _url := os.Getenv(URL); _url != "" {
		url = _url
	}
	nextPage := 1
	for {
		res, err := http.Get(fmt.Sprintf("%s?limit=3&key=%s&page=%d", url, key, nextPage))
		if err != nil {
			panic(err)
		}
		var data GhostContent
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			panic(nil)
		}
		log.Printf("fetching posts, page %d/%d", data.Meta.Pagination.Page, data.Meta.Pagination.Total)
		for _, post := range data.Posts {
			localDb[post.Slug] = convertpost(&post)
		}
		nextPage = data.Meta.Pagination.Next
		if nextPage == 0 {
			break
		}
	}

	log.Printf("fetched %d posts", len(localDb))

	//log.Println(localDb["welcome"].Content)
	return localDb
}

func contentFetcher() {
	for {
		db = fetchcontent()
		time.Sleep(1 * time.Hour)
	}
}
