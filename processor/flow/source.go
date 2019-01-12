package flow

type Source struct {
	URL      string `bson:"url"                json:"url"`
	Key      string `bson:"key,omitempty"      json:"key,omitempty"`
	Original string `bson:"original,omitempty" json:"original,omitempty"`
}
