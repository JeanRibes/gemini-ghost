# {{.Settings.Title}}
{{.Settings.Description}}

=> /search Search

{{range .Settings.Navigation }}=> {{.URL}} {{.Label}}
{{end}}

# Posts
{{$df := .DateF}}
{{range .Posts}}=> /{{.Slug}} {{.Title}}
{{.Excerpt}}
Published at {{call $df .PublishedAt }}
{{end}}


{{range .Settings.SecondaryNavigation }}=> {{.URL}} {{.Label}}
{{end}}

Happy reading !

```
Generated from the Ghost API
```
