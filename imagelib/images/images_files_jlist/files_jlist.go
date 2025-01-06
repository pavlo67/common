package images_files_jlist

import (
	"fmt"
	"image"
	"path/filepath"
	"strconv"

	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/serialization"

	"github.com/pavlo67/imagelib/images"
	"github.com/pavlo67/imagelib/sources"
)

var _ images.Operator = &imagesFilesJList{}

type imagesFilesJList struct {
	basePath string
	descrs   []sources.Description
	colored  bool
}

const onNew = "on images_files_jlist.New()"

func New(basePath string, colored bool) (images.Operator, db.Cleaner, error) {

	var err error
	if basePath, err = filelib.Dir(basePath); err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	imagesOp := imagesFilesJList{
		basePath: basePath,
		colored:  colored,
	}

	if err = serialization.ReadAllPartsJSON(filepath.Join(basePath, sources.FramesAllDescriptionsFilename), &imagesOp.descrs); err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	return &imagesOp, &imagesOp, nil
}

const onSave = "on imagesFilesJList.Save()"

func (op *imagesFilesJList) Save(img image.Image, descr sources.Description, key images.Key) (imgPath string, err error) {
	if op == nil {
		return "", errors.New("op == nil / " + onSave)
	}

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
		imgGray, err := imagelib.ImageToGray(img)
		if err != nil {
			return "", fmt.Errorf("%s: %s / "+onSave, imgPath, err)
		} else if imgGray == nil {
			return "", fmt.Errorf("imgGray == nil: %s / "+onSave, imgPath)
		}

		if err = imagelib.SavePGM(imgGray, imgPath); err != nil {
			return "", fmt.Errorf("%s: %s / "+onSave, imgPath, err)
		}
	}

	descr.ImagePath = imgPath
	if len(op.descrs) < 1 {
		op.descrs = []sources.Description{descr}
	} else if descr.N < op.descrs[0].N {
		op.descrs = append([]sources.Description{descr}, op.descrs...)
	} else {
		var added bool
		for i, d := range op.descrs {
			if descr.N == d.N {
				op.descrs[i], added = descr, true
				break
			} else if descr.N < d.N {
				op.descrs, added = append(op.descrs[:i], append([]sources.Description{descr}, op.descrs[i:]...)...), true
				break
			}
		}
		if !added {
			op.descrs = append(op.descrs, descr)
		}
	}
	if err = serialization.SaveAllPartsJSON(op.descrs, filepath.Join(filepath.Dir(imgPath), sources.FramesAllDescriptionsFilename)); err != nil {
		return imgPath, errors.Wrap(err, onSave)
	}

	return imgPath, nil
}

const onCheck = "on imagesFilesJList.Check()"

func (op imagesFilesJList) Check(key images.Key) (bool, error) {

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

	if op.descrI(imgPath) < 0 {
		return false, fmt.Errorf("%s: descr is absent / "+onCheck, imgPath)
	}

	return true, nil
}

const onGet = "on imagesFilesJList.Get()"

func (op imagesFilesJList) Get(key images.Key) (image.Image, *sources.Description, error) {
	keyPath := images.KeyPath(key, op.colored)
	if keyPath == "" {
		return nil, nil, fmt.Errorf("empty keyPath for key '%s' / "+onGet, key)
	}
	imgPath := filepath.Join(op.basePath, keyPath)
	imgExists, err := filelib.FileExists(imgPath, false)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %s / "+onGet, imgPath, err)
	} else if !imgExists {
		return nil, nil, nil
	}

	descrI := op.descrI(imgPath)
	if descrI < 0 {
		return nil, nil, fmt.Errorf("%s: descr is absent / "+onGet, imgPath)
	}

	if op.colored {
		img, err := imagelib.Read(imgPath)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %s / "+onGet, imgPath, err)
		}
		return img, &op.descrs[descrI], nil
	}

	imgGray, err := imagelib.ReadPGMSpecial(imgPath)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %s / "+onGet, imgPath, err)
	}

	descr := op.descrs[descrI]
	descr.ImagePath = imgPath

	return imgGray, &descr, nil
}

const onListPaths = "on imagesFiles.ListPaths()"

func (op imagesFilesJList) ListPaths(keyRegexStr string) ([]string, error) {
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
		descrI := op.descrI(imgPath)
		if descrI < 0 {
			l.Errorf("description doesn't exist for %s / "+onListPaths, imgPath)
			continue
		}
		imgPathsOk = append(imgPathsOk, imgPath)
	}

	return imgPathsOk, nil
}

func (op imagesFilesJList) descrI(imgPath string) int {
	imgBase := filepath.Base(imgPath)

	//fmt.Printf("op.descrs: %+v, imgBase: %s", op.descrs, imgBase)

	if len(imgBase) == 8 {
		if n, err := strconv.Atoi(imgBase[:4]); err == nil {
			for i := range op.descrs {
				if op.descrs[i].N == n {
					return i
				}
			}
		}
	}

	return -1
}
