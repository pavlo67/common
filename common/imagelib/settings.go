package imagelib

import (
	"maps"

	"github.com/pavlo67/common/common"
)

type Settings struct {
	DPM     float64    `json:",omitempty"`
	Options common.Map `json:",omitempty"`
}

func (settings *Settings) SetOptions(key string, value interface{}) {
	if settings == nil {
		return
	}

	if settings.Options == nil {
		settings.Options = common.Map{key: value}
	} else {
		settings.Options[key] = value
	}
}

func (settings Settings) Copy() Settings {

	settings.Options = maps.Clone(settings.Options)

	return settings
}

//lyr.Options["color_range"]
