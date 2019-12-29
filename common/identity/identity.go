package identity

import (
	"regexp"

	"github.com/pavlo67/workshop/common/libraries/strlib"
)

type Item struct {
	Domain string `bson:"domain,omitempty"  json:"domain,omitempty"`
	Path   string `bson:"path,omitempty"    json:"path,omitempty"`
	ID     string `bson:"id,omitempty"      json:"id,omitempty"`
}

// Key is a string representation of Item.
type Key string

func (item *Item) IsValid() bool {
	if item == nil {
		return false
	}
	return strlib.ReSpaces.ReplaceAllString(item.Domain, "") != "" &&
		strlib.ReSpaces.ReplaceAllString(item.Path, "") != "" &&
		strlib.ReSpaces.ReplaceAllString(item.ID, "") != ""
}

func (item *Item) Key() Key {
	if item == nil {
		return Key("")
	}

	domain := strlib.ReSpaces.ReplaceAllString(item.Domain, "")
	path := strlib.ReSpaces.ReplaceAllString(item.Path, "")
	id := strlib.ReSpaces.ReplaceAllString(item.ID, "")

	if len(id) > 0 {
		return Key(domain + "/" + path + "/" + id)
	} else if len(path) > 0 {
		return Key(domain + "/" + path)
	} else if len(domain) > 0 {
		return Key(domain)
	} else {
		return Key("")
	}
}

func (item *Item) String() string {
	return string(item.Key())
}

var reKeyDelim = regexp.MustCompile(`/`)

func (is Key) Identity() Item {
	is0 := strlib.ReSpacesFin.ReplaceAllString(string(is), "")

	indexes := reKeyDelim.FindAllStringIndex(is0, -1)
	if len(indexes) < 1 {
		return Item{Domain: is0}
	}

	if len(indexes) < 2 {
		return Item{
			Domain: is0[:indexes[0][0]],
			Path:   is0[indexes[0][1]:],
		}
	}

	return Item{
		Domain: is0[:indexes[0][0]],
		Path:   is0[indexes[0][1]:indexes[len(indexes)-1][0]],
		ID:     is0[indexes[len(indexes)-1][1]:],
	}
}

func (is Key) Short(domain string) Key {
	if len(is) > len(domain) && string(is[:len(domain)]) == domain && is[len(domain):len(domain)+1] == "/" {
		return Key(is[len(domain):])
	}
	return is
}

func (is Key) Full(domain string) Key {
	if len(is) > 0 && is[:1] == "/" {
		return Key(domain + string(is))
	}
	return is
}

func IsEqual(identity *Item, is Key, domain string) bool {
	return identity != nil && is == identity.Key()
}
