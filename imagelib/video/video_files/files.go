package video_files

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib"
	"path/filepath"
	"sync"
	"time"

	"gocv.io/x/gocv"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/serialization"

	"github.com/pavlo67/imagelib/opencvlib"
	"github.com/pavlo67/imagelib/video"
)

var _ video.Operator = &videoFiles{}

type videoFiles struct {
	video.Info

	colorConversionCode gocv.ColorConversionCode

	lastFrame      *video.Processing
	lastFrameMutex *sync.Mutex

	nFilesI int
	nFiles  []filelib.NumberedFile
}

const onNew = "on videoFiles.New()"

func New(path string, colorConversionCode gocv.ColorConversionCode) (video.Operator, error) {

	nFiles, err := filelib.NumberedFilesSequence(path, video.RePNGFilesStr, false)
	if err != nil {
		return nil, errors.Wrap(err, onNew)
	} else if len(nFiles) < 1 {
		return nil, fmt.Errorf(`%s / "%s": len(nFiles) == 0 / `+onNew, path, video.RePNGFilesStr)
	}

	var info video.Info
	if err = serialization.Read(filepath.Join(path, video.InfoFilename), serialization.MarshalerJSON, &info); err != nil {
		return nil, errors.Wrap(err, onNew)
	}

	return &videoFiles{
		Info:                info,
		colorConversionCode: colorConversionCode,
		lastFrameMutex:      &sync.Mutex{},
		nFiles:              nFiles,
		nFilesI:             0,
	}, nil
}

func (op *videoFiles) GetInfo() (*video.Info, error) {
	if op == nil {
		return nil, errors.New("op == nil / on videoOpenCV.GetInfo()")
	}

	return &op.Info, nil
}

func (op videoFiles) IsFinished() bool {
	return op.nFilesI >= len(op.nFiles)
}

const onNextFrane = "on videoFiles.NextFrame()"

func (op *videoFiles) NextFrame(scale float64, delayMax time.Duration) (*video.Processing, error) {

	N := op.nFilesI + op.nFiles[0].I

	if op.IsFinished() {
		return nil, fmt.Errorf("video sequence is finished / " + onNextFrane)
	}

	imgPath := op.nFiles[op.nFilesI].Path
	imgOriginal, err := imagelib.Read(imgPath)
	if err != nil {
		return nil, errors.Wrap(err, onNextFrane)
	} else if imgOriginal == nil {
		return nil, fmt.Errorf("imgOriginal == nil (%s) / "+onNextFrane, imgPath)
	}

	imgGray, err := imagelib.ImageToGray(imgOriginal)
	if err != nil {
		return nil, errors.Wrapf(err, "on imagelib.ImageToGray(%s) / "+onNextFrane, imgPath)
	} else if imgOriginal == nil {
		return nil, fmt.Errorf("imgGray == nil (%s) / "+onNextFrane, imgPath)
	}

	mat, err := gocv.ImageGrayToMatGray(imgGray)
	if err != nil {
		return nil, errors.Wrapf(err, "on gocv.ImageGrayToMatGray(%s) / "+onNextFrane, imgPath)
	} else if imgOriginal == nil {
		return nil, fmt.Errorf("mat == nil (%s) / "+onNextFrane, imgPath)
	}
	defer mat.Close()

	img, err := opencvlib.Prepare(mat, op.colorConversionCode, scale)
	if err != nil {
		return nil, fmt.Errorf("frame#%d: %s / "+onNextFrane, N, err)
	} else if img == nil {
		return nil, fmt.Errorf("frame#%d: img == nil / "+onNextFrane, N)
	}

	op.lastFrameMutex.Lock()
	defer op.lastFrameMutex.Unlock()

	op.lastFrame = &video.Processing{
		FrameInfo: video.FrameInfo{
			N:    N,
			Time: op.Info.StartedAt.Add(time.Duration(op.nFilesI) * time.Second / time.Duration(op.Info.FPS)),
		},
		Image: img,
	}

	op.nFilesI++

	return op.lastFrame, nil
}

func (op videoFiles) LastFrame() (*video.Processing, error) {
	op.lastFrameMutex.Lock()
	defer op.lastFrameMutex.Unlock()

	return op.lastFrame, nil
}
