package ghost

type ContentAPI struct {
	url string
	key string
}

func New(apiUrl string, key string) *ContentAPI {
	return &ContentAPI{
		url: apiUrl,
		key: key,
	}
}

func (g *ContentAPI) AllPosts() ([]Post, error) {
	return fetchPosts(g.url, g.key)
}

func (g *ContentAPI) AllPages() ([]Page, error) {
	return fetchPages(g.url, g.key)
}

func (g *ContentAPI) Settings() (*Settings, error) {
	return fetchSettings(g.url, g.key)
}
