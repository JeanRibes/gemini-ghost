Ghost blog on Gemini
# Posts
{{$df := .DateF}}
{{range .Posts}}=> /{{.Slug}} {{.Title}}
{{.Excerpt}}
Published at {{call $df .PublishedAt }}
{{end}}

Happy reading !

```
Generated from the Ghost API
```
