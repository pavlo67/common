package identity

import (
	"regexp"

	"strings"

	"github.com/pavlo67/workshop/common/libraries/strlib"
)

type Domain string

type Item struct {
	Domain Domain `bson:"domain,omitempty"  json:"domain,omitempty"`
	Path   string `bson:"path,omitempty"    json:"path,omitempty"`
	ID     string `bson:"id,omitempty"      json:"id,omitempty"`
}

const PathDelim = `/`
const IDDelim = `##`

var reProto = regexp.MustCompile(`^https?://`)

var reDomainDelim = regexp.MustCompile(PathDelim + `.*`)

var rePathDelimFirst = regexp.MustCompile(`^(` + PathDelim + `)+`)
var rePathDelim = regexp.MustCompile(IDDelim + `.*`)

var reIDDelimFirst = regexp.MustCompile(`^(` + IDDelim + `)+`)

func FromURLRaw(urlRaw string) Item {
	urlWithoutProto := reProto.ReplaceAllString(strings.TrimSpace(urlRaw), "")
	domain := reDomainDelim.ReplaceAllString(urlWithoutProto, "")

	// TODO!!! clean more

	return Item{
		Domain: Domain(domain),
		Path:   urlWithoutProto[len(domain):],
	}
}

// Key is a string representation of Item.
type Key string

//func (key Key) NotEmpty() bool {
//
//}

func (item *Item) IsValid() bool {
	if item == nil {
		return false
	}
	return strlib.ReSpaces.ReplaceAllString(string(item.Domain), "") != "" &&
		strlib.ReSpaces.ReplaceAllString(item.Path, "") != "" &&
		strlib.ReSpaces.ReplaceAllString(item.ID, "") != ""
}

func (item *Item) Key() Key {
	if item == nil {
		return Key("")
	}

	domain := strlib.ReSpaces.ReplaceAllString(string(item.Domain), "")
	path := strlib.ReSpaces.ReplaceAllString(item.Path, "")
	id := strlib.ReSpaces.ReplaceAllString(item.ID, "")

	if len(id) > 0 {
		return Key(domain + PathDelim + path + IDDelim + id)
	} else if len(path) > 0 {
		return Key(domain + PathDelim + path)
	} else if len(domain) > 0 {
		return Key(domain)
	} else {
		return Key("")
	}
}

func (item *Item) String() string {
	return string(item.Key())
}

func (key Key) Normalize() Key {
	return Key(strings.TrimSpace(string(key)))
}

func (key Key) Identity() *Item {
	keyTrimmed := strings.TrimSpace(string(key))
	if len(keyTrimmed) < 1 {
		return nil
	}

	domain := reDomainDelim.ReplaceAllString(keyTrimmed, "")
	pathid := rePathDelimFirst.ReplaceAllString(strings.TrimSpace(keyTrimmed[len(domain):]), "")

	path := rePathDelim.ReplaceAllString(pathid, "")
	id := reIDDelimFirst.ReplaceAllString(strings.TrimSpace(pathid[len(path):]), "")

	return &Item{
		Domain: Domain(domain),
		Path:   path,
		ID:     id,
	}
}

func (key Key) Short(domain string) Key {
	if len(key) > len(domain) && string(key[:len(domain)]) == domain && key[len(domain):len(domain)+1] == PathDelim {
		return Key(key[len(domain):])
	}
	return key
}

func (key Key) Full(domain string) Key {
	if len(key) > 0 && key[:1] == PathDelim {
		return Key(domain + string(key))
	}
	return key
}

func IsEqual(identity *Item, is Key, domain string) bool {
	return identity != nil && is == identity.Key()
}
