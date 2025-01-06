package sources

import (
	"fmt"
	"github.com/pavlo67/common/common/geolib"
	"image"
	"math"
	"path/filepath"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/serialization"

	"github.com/pavlo67/imagelib/video"
)

type Series struct {
	Path         string
	Descriptions []Description
	Info         video.Info
	DPMRequired  *float64
	Grayscaled   bool
}

type SeriesConfig struct {
	VideoInfo video.Info
	// DEPRECATED!!!
	Grayscaled bool
	DPM        float64
	Trajectory geolib.Trajectory
	Velocity   float64
}

const onAdd = "on preparation.Series.Add()"

func (series *Series) Add(img image.Image, descr Description) error {
	if series == nil {
		return errors.New("series == nil / " + onAdd)
		//} else if img == nil {
		//	return errors.New("img == nil / " + onAdd)
		//} else if !(descr.DPM > 0) || math.IsInf(descr.DPM, 1) {
		//	return fmt.Errorf("wrong descr.DPM: %+v / "+onAdd, descr)
	} else if series.DPMRequired != nil && (!(*series.DPMRequired > 0) || math.IsInf(*series.DPMRequired, 1)) {
		return fmt.Errorf("wrong series.DPMRequired: %f / "+onAdd, *series.DPMRequired)
	} else if len(series.Descriptions) > 0 && series.Descriptions[len(series.Descriptions)-1].N+1 != descr.N {
		return fmt.Errorf("wrong descr.N: %+v (previous N = %d) / "+onAdd, descr, series.Descriptions[len(series.Descriptions)-1].N)
	}

	//imgGray, err := imagelib.ImageToGray(img)
	//if err != nil {
	//	return errors.Wrap(err, onAdd)
	//}
	//
	//if series.DPMRequired != nil {
	//	imgGrayResized, _, err := opencvlib.ResizeGray(*imgGray, *series.DPMRequired/descr.DPM)
	//	if err != nil {
	//		return errors.Wrap(err, onAdd)
	//	} else if imgGrayResized == nil {
	//		return errors.New("imgGrayResized == nil / " + onAdd)
	//	}
	//	imgGray, descr.DPM = imgGrayResized, *series.DPMRequired
	//}

	// TODO!!! remove it
	descr.DPM = *series.DPMRequired

	if series.Grayscaled {
		descr.ImagePath = filepath.Join(series.Path, fmt.Sprintf("%04d.pgm", descr.N))
		// err = imagelib.SavePGM(imgGray, descr.ImagePath)
	} else {
		descr.ImagePath = filepath.Join(series.Path, fmt.Sprintf("%04d.png", descr.N))
		// err = imagelib.SavePGM(imgGray, descr.ImagePath)
	}
	//if err != nil {
	//	return errors.Wrap(err, onAdd)
	//}

	series.Descriptions = append(series.Descriptions, descr)

	return nil
}

const onSave = "on preparation.Series.Save()"

func (series Series) Save() error {
	if err := serialization.SaveAllPartsJSON(series.Descriptions, filepath.Join(series.Path, FramesAllDescriptionsFilename)); err != nil {
		return errors.Wrap(err, onSave)
	}

	if series.DPMRequired != nil {
		series.Info.DPMConverted = series.DPMRequired
		if err := serialization.Save(series.Info, serialization.MarshalerJSON, filepath.Join(series.Path, VideoInfoFilename)); err != nil {
			return errors.Wrap(err, onSave)
		}
	}

	return nil
}
