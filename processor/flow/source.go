package flow

import (
	"strings"
)

type SourceType string

type Source struct {
	URL      string `bson:"url"                 json:"url"`
	SourceID string `bson:"source_id,omitempty" json:"source_id,omitempty"`
	Original string `bson:"original,omitempty"  json:"original,omitempty"`
}

func (src *Source) Key(keyAdd string) string {
	if src == nil {
		return ""
	}

	url := strings.TrimSpace(src.URL)

	pos := strings.Index(url, "#")
	if pos >= 0 {
		url = url[:pos]
	}

	if url == "" {
		return ""
	}

	if len(keyAdd) > 0 {
		url += "#" + keyAdd
	}

	sourceID := strings.TrimSpace(src.SourceID)
	if sourceID == "" {
		return url
	}

	return url + "#" + sourceID
}
