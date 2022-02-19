package ghost

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Settings struct {
	Title             string      `json:"title"`
	Description       string      `json:"description"`
	Logo              interface{} `json:"logo"`
	Icon              interface{} `json:"icon"`
	AccentColor       string      `json:"accent_color"`
	CoverImage        string      `json:"cover_image"`
	Facebook          string      `json:"facebook"`
	Twitter           string      `json:"twitter"`
	Lang              string      `json:"lang"`
	Locale            string      `json:"locale"`
	Timezone          string      `json:"timezone"`
	CodeinjectionHead interface{} `json:"codeinjection_head"`
	CodeinjectionFoot string      `json:"codeinjection_foot"`
	Navigation        []struct {
		Label string `json:"label"`
		URL   string `json:"url"`
	} `json:"navigation"`
	SecondaryNavigation []struct {
		Label string `json:"label"`
		URL   string `json:"url"`
	} `json:"secondary_navigation"`
	MetaTitle             interface{} `json:"meta_title"`
	MetaDescription       interface{} `json:"meta_description"`
	OgImage               interface{} `json:"og_image"`
	OgTitle               interface{} `json:"og_title"`
	OgDescription         interface{} `json:"og_description"`
	TwitterImage          interface{} `json:"twitter_image"`
	TwitterTitle          interface{} `json:"twitter_title"`
	TwitterDescription    interface{} `json:"twitter_description"`
	MembersSupportAddress string      `json:"members_support_address"`
	URL                   string      `json:"url"`
}

type ghostsettings struct {
	Settings Settings `json:"settings"`
	Meta     struct {
	} `json:"meta"`
}

func fetchSettings(url string, key string) (*Settings, error) {
	res, err := http.Get(fmt.Sprintf("%s/settings/?key=%s", url, key))
	if err != nil {
		return nil, err
	}
	var data ghostsettings
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data.Settings, nil
}
