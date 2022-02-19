package main

import (
	"github.com/JeanRibes/gemini-ghost/ghost"
	"github.com/blevesearch/bleve/v2"
	"github.com/pitr/gig"
	"html/template"
	"log"
	"time"
)

func ghostContent(c gig.Context) error {
	slug := c.Param("slug")
	post, exists := db.Posts[slug]
	if !exists {
		post, exists = db.Pages[slug]
	}
	if !exists {
		return c.NoContent(gig.StatusTemporaryFailure, "no post or page here :(")
	}
	return c.Gemini(post.Content)
}

func ghostIndex(c gig.Context) error {
	tmpl, err := template.ParseFiles("index.tpl")
	if err != nil {
		return err
	}
	if err := c.Response().WriteHeader(gig.StatusSuccess, gig.MIMETextGemini); err != nil {
		return err
	}
	return tmpl.Execute(c.Response().Writer, map[string]interface{}{
		"Posts":    db.Posts,
		"Pages":    db.Pages,
		"Settings": db.Settings,
		"DateF": func(date time.Time) string {
			return date.Format(time.ANSIC)
		},
	})
}

func searchPost(c gig.Context) error {
	query := input_helper(c, "Your query ?")
	results, err := index.Search(bleve.NewSearchRequest(bleve.NewMatchQuery(query)))
	if err != nil {
		println(err.Error())
		return err
	}
	log.Println(results.String())

	posts := []StoredPost{}

	for _, result := range results.Hits {
		println("id:", result.ID)
		if post, exists := db.Posts[result.ID]; exists {
			posts = append(posts, post)
		}
	}

	tmpl, err := template.ParseFiles("index.tpl")
	if err != nil {
		return err
	}
	if err := c.Response().WriteHeader(gig.StatusSuccess, gig.MIMETextGemini); err != nil {
		return err
	}
	return tmpl.Execute(c.Response().Writer, map[string]interface{}{
		"Posts":    posts,
		"Pages":    []ghost.Page{},
		"Settings": db.Settings,
		"DateF": func(date time.Time) string {
			return date.Format(time.ANSIC)
		},
	})

}
