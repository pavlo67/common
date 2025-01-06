package video

import (
	"encoding/json"
	"fmt"
	"github.com/pavlo67/common/common/logger"
	"image"
	"path/filepath"
	"time"

	"github.com/pavlo67/common/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "video"

const LogOptionsPrefix = "video_options_"

type Info struct {
	FPS        float64
	Rectangle  image.Rectangle
	Monochrome bool `json:",omitempty"`
	Realtime   bool

	Source  string `json:",omitempty"`
	SrcNMin int    `json:",omitempty"`
	SrcNMax int    `json:",omitempty"`

	Device       interface{} `json:",omitempty"`
	StartedAt    *time.Time  `json:",omitempty"`
	NFrom        int         `json:",omitempty"`
	FPSDivider   *int        `json:",omitempty"`
	DPMConverted *float64    `json:",omitempty"`
}

type FrameInfo struct {
	N    int
	Time time.Time `json:",omitempty"`
}

type Processing struct {
	FrameInfo   `json:",inline"`
	image.Image `json:",omitempty"`
}

//type DescriptionVideo struct {
//	Description `json:",inline"`
//	FPS         int
//	StartHMS    string
//	TimeHMS     string
//
//	// Grayscale   bool
//	// Marks       []Mark
//	// Targets     []Mark
//}

type Operator interface {
	GetInfo() (*Info, error)
	NextFrame(scale float64, delayMax time.Duration) (*Processing, error)
	LastFrame() (*Processing, error)
	IsFinished() bool
}

const onSaveDebug = "on videoFiles.SaveDebug()"

func SaveDebug(op Operator, l logger.Operator, path string) {
	if op == nil {
		return
	}

	info, err := op.GetInfo()
	if err != nil {
		l.Error(err, " / ", onSaveDebug)
	}

	infoBytes, err := json.Marshal(info)
	if err != nil {
		l.Error(err, " / ", onSaveDebug)
	}

	l.File(filepath.Join(path, LogOptionsPrefix+fmt.Sprintf("%s.log", l.Key())), false, infoBytes)
}

//const onGetRGBA = "on mappingLandmarks.ToRGBA()"
//
//func (pr Processing) ToRGBA(dpmRange [2]float64) (*image.RGBA, float64, error) {
//
//	imgRGB, err := imagelib.ImageToRGBA(pr.Image)
//	if err != nil {
//		return nil, 0, errors.Wrap(err, onGetRGBA)
//	} else if imgRGB == nil {
//		return nil, 0, errors.New("imgRGB == nil / " + onGetRGBA)
//	}
//
//	imgRGB.Rect = imagelib.Normalize(imgRGB.Rect)
//
//	return opencvlib.ResizeToRange(*imgRGB, pr.Settings.DPM, dpmRange)
//}
