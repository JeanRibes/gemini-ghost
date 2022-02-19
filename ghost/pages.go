package ghost

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Page struct {
	ID                  string      `json:"id"`
	UUID                string      `json:"uuid"`
	Title               string      `json:"title"`
	Slug                string      `json:"slug"`
	HTML                string      `json:"html"`
	CommentID           string      `json:"comment_id"`
	FeatureImage        interface{} `json:"feature_image"`
	Featured            bool        `json:"featured"`
	Visibility          string      `json:"visibility"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
	PublishedAt         time.Time   `json:"published_at"`
	CustomExcerpt       interface{} `json:"custom_excerpt"`
	CodeinjectionHead   interface{} `json:"codeinjection_head"`
	CodeinjectionFoot   interface{} `json:"codeinjection_foot"`
	CustomTemplate      interface{} `json:"custom_template"`
	CanonicalURL        interface{} `json:"canonical_url"`
	URL                 string      `json:"url"`
	Excerpt             string      `json:"excerpt"`
	ReadingTime         int         `json:"reading_time"`
	Page                bool        `json:"page"`
	Access              bool        `json:"access"`
	OgImage             interface{} `json:"og_image"`
	OgTitle             interface{} `json:"og_title"`
	OgDescription       interface{} `json:"og_description"`
	TwitterImage        interface{} `json:"twitter_image"`
	TwitterTitle        interface{} `json:"twitter_title"`
	TwitterDescription  interface{} `json:"twitter_description"`
	MetaTitle           interface{} `json:"meta_title"`
	MetaDescription     interface{} `json:"meta_description"`
	Frontmatter         interface{} `json:"frontmatter"`
	FeatureImageAlt     interface{} `json:"feature_image_alt"`
	FeatureImageCaption interface{} `json:"feature_image_caption"`
}

type ghostpages struct {
	Pages []Page `json:"pages"`
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

func fetchPages(url string, key string) (pages []Page, err error) {
	nextPage := 1
	for {
		res, err := http.Get(fmt.Sprintf("%s/pages/?limit=3&key=%s&page=%d", url, key, nextPage))
		if err != nil {
			return nil, err
		}
		var data ghostpages
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			return nil, err
		}
		pages = append(pages, data.Pages...)
		nextPage = data.Meta.Pagination.Next
		if nextPage == 0 {
			break
		}
	}
	return pages, nil
}
