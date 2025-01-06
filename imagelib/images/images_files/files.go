package images_files

import (
	"fmt"
	"image"
	"path/filepath"

	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/serialization"

	"github.com/pavlo67/imagelib/images"
	"github.com/pavlo67/imagelib/sources"
)

var _ images.Operator = &imagesFiles{}

type imagesFiles struct {
	basePath string
	colored  bool
}

const onNew = "on images_files.New()"

func New(basePath string, colored bool) (images.Operator, db.Cleaner, error) {

	var err error
	if basePath, err = filelib.Dir(basePath); err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	imagesOp := imagesFiles{
		basePath: basePath,
		colored:  colored,
	}

	return &imagesOp, &imagesOp, nil
}

const onSave = "on imagesFiles.Save()"

func (op imagesFiles) Save(img image.Image, descr sources.Description, key images.Key) (imgPath string, err error) {
	keyPath := images.KeyPath(key, op.colored)
	if keyPath == "" {
		return "", fmt.Errorf("empty path for key '%s' / "+onSave, key)
	}
	imgPath = filepath.Join(op.basePath, keyPath)

	if op.colored {
		if err = imagelib.SavePNG(img, imgPath); err != nil {
			return "", fmt.Errorf("%s: %s / "+onSave, imgPath, err)
		}
	} else {
		var imgGray *image.Gray
		if imgGray, err = imagelib.ImageToGray(img); err != nil {
			return "", fmt.Errorf("%s: %s / "+onSave, imgPath, err)
		} else if imgGray == nil {
			return "", fmt.Errorf("imgGray == nil: %s / "+onSave, imgPath)
		} else if err = imagelib.SavePGM(imgGray, imgPath); err != nil {
			return "", fmt.Errorf("%s: %s / "+onSave, keyPath, err)
		}
	}

	descr.ImagePath = imgPath
	if err = serialization.Save(descr, serialization.MarshalerJSON, imgPath+".json"); err != nil {
		return imgPath, errors.Wrap(err, onSave)
	}

	return imgPath, nil
}

const onCheck = "on imagesFiles.Check()"

func (op imagesFiles) Check(key images.Key) (bool, error) {

	keyPath := images.KeyPath(key, op.colored)
	if keyPath == "" {
		return false, fmt.Errorf("empty keyPath for key '%s' / "+onCheck, key)
	}
	imgPath := filepath.Join(op.basePath, keyPath)

	imgExists, err := filelib.FileExists(imgPath, false)
	if err != nil {
		return false, fmt.Errorf("%s: %s / "+onCheck, imgPath, err)
	}

	if !imgExists {
		return false, nil
	}

	descrPath := imgPath + ".json"
	descrExists, err := filelib.FileExists(descrPath, false)
	if err != nil {
		return false, fmt.Errorf("%s: %s / "+onCheck, descrPath, err)
	} else if !descrExists {
		return false, fmt.Errorf("%s: descr is absent / "+onCheck, descrPath)
	}

	return true, nil
}

const onGet = "on imagesFiles.Get()"

func (op imagesFiles) Get(key images.Key) (image.Image, *sources.Description, error) {
	keyPath := images.KeyPath(key, op.colored)
	if keyPath == "" {
		return nil, nil, fmt.Errorf("empty keyPath for key '%s' / "+onGet, key)
	}
	imgPath := filepath.Join(op.basePath, keyPath)
	descrPath := imgPath + ".json"

	descrExists, err := filelib.FileExists(descrPath, false)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %s / "+onGet, descrPath, err)
	} else if !descrExists {
		return nil, nil, nil
	}

	var descr sources.Description
	if err = serialization.Read(descrPath, serialization.MarshalerJSON, &descr); err != nil {
		return nil, nil, errors.Wrap(err, onGet)
	}

	imgExists, err := filelib.FileExists(imgPath, false)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %s / "+onGet, imgPath, err)
	} else if !imgExists {
		return nil, nil, nil
	}

	if op.colored {
		img, err := imagelib.Read(imgPath)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %s / "+onGet, imgPath, err)
		}
		return img, &descr, nil
	}

	imgGray, err := imagelib.ReadPGMSpecial(imgPath)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %s / "+onGet, imgPath, err)
	}

	descr.ImagePath = imgPath
	return imgGray, &descr, nil
}

const onListPaths = "on imagesFiles.ListPaths()"

func (op imagesFiles) ListPaths(keyRegexStr string) ([]string, error) {
	reKeyPath := images.KeyPathRegex(keyRegexStr, op.colored)
	if reKeyPath == nil {
		return nil, fmt.Errorf("empty keyPath regex for keyRegexStr '%s' / "+onListPaths, keyRegexStr)
	}

	imgPaths, err := filelib.List(op.basePath, reKeyPath, false, true)
	if err != nil {
		return nil, fmt.Errorf("%s / %s --> %s / "+onListPaths, op.basePath, keyRegexStr, err)
	}

	var imgPathsOk []string
	for _, imgPath := range imgPaths {
		descrPath := imgPath + ".json"
		descrExists, err := filelib.FileExists(descrPath, false)
		if err != nil {
			l.Error(err, onListPaths)
			continue
		} else if !descrExists {
			l.Errorf("description doesn't exist: %s / "+onListPaths, descrPath)
			continue
		}
		imgPathsOk = append(imgPathsOk, imgPath)
	}

	return imgPathsOk, nil
}
