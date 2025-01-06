package video_opencv

import (
	"fmt"
	"image"
	"sync"
	"time"

	"gocv.io/x/gocv"

	"github.com/pavlo67/common/common/errors"

	"github.com/pavlo67/imagelib/opencvlib"
	"github.com/pavlo67/imagelib/video"
)

var _ video.Operator = &videoOpenCV{}

const onNew = "on videoOpenCV.New()"

func New(info video.Info, colorConversionCode gocv.ColorConversionCode) (video.Operator, error) {
	capture, err := gocv.OpenVideoCapture(info.Device)
	if err != nil {
		return nil, errors.Wrapf(err, "can't open device: %#v", info.Device)
	} else if capture == nil {
		return nil, fmt.Errorf("capture == nil (from device: %#v) / "+onNew, info.Device)
	}

	//capture.Set(gocv.VideoCaptureFOURCC, capture.ToCodec("MJPG"))
	//capture.Set(gocv.VideoCaptureAutoExposure, 0)
	//capture.Set(gocv.VideoCaptureFrameWidth, 640)
	//capture.Set(gocv.VideoCaptureFrameHeight, 480)
	//capture.Set(gocv.VideoCaptureFPS, 3)

	// TODO!!! is it correct???
	capture.Set(gocv.VideoCaptureBufferSize, 0)

	if bufferSize := int(capture.Get(gocv.VideoCaptureBufferSize)); bufferSize > 1 {
		return nil, fmt.Errorf("bufferSize (%d) > 1", bufferSize)
	}

	return &videoOpenCV{
		Info: info,

		colorConversionCode: colorConversionCode,

		capture:        capture,
		lastFrameMutex: &sync.Mutex{},
	}, nil
}

type videoOpenCV struct {
	video.Info

	colorConversionCode gocv.ColorConversionCode

	capture        *gocv.VideoCapture
	lastFrame      *video.Processing
	lastFrameMutex *sync.Mutex
	isFinished     bool
}

func (op *videoOpenCV) GetInfo() (*video.Info, error) {
	if op == nil {
		return nil, errors.New("op == nil / on videoOpenCV.GetInfo()")
	}

	op.Info.FPS = op.capture.Get(gocv.VideoCaptureFPS)
	op.Info.Rectangle = image.Rectangle{
		Max: image.Point{int(op.capture.Get(gocv.VideoCaptureFrameWidth)), int(op.capture.Get(gocv.VideoCaptureFrameHeight))},
	}

	return &op.Info, nil
}

func (op videoOpenCV) IsFinished() bool {
	return op.isFinished
}

const onNextFrane = "on videoOpenCV.NextFrame()"

func (op *videoOpenCV) NextFrame(scale float64, delayMax time.Duration) (*video.Processing, error) {

	waitingTo := time.Now().Add(delayMax)

	var N int

	mat := gocv.NewMat()
	defer mat.Close()

	for {
		ok := op.capture.Read(&mat)
		if ok {
			break
		} else if delayMax >= 0 && !waitingTo.After(time.Now()) {
			op.isFinished = true
			return nil, fmt.Errorf("delayMax: %d (%d sec) is exceeded", delayMax, delayMax/time.Second)
		}
	}

	img, err := opencvlib.Prepare(mat, op.colorConversionCode, scale)
	if err != nil {
		return nil, fmt.Errorf("frame#%d: %s / "+onNextFrane, N, err)
	} else if img == nil {
		return nil, fmt.Errorf("frame#%d: img == nil / "+onNextFrane, N)
	}

	op.lastFrameMutex.Lock()
	defer op.lastFrameMutex.Unlock()

	// log.Print("video frame processed!")

	op.lastFrame = &video.Processing{
		FrameInfo: video.FrameInfo{
			N:    N,
			Time: op.Info.StartedAt.Add(time.Duration(N) * time.Second / time.Duration(op.Info.FPS)),
		},
		Image: img,
	}

	return op.lastFrame, nil
}

func (op videoOpenCV) LastFrame() (*video.Processing, error) {
	op.lastFrameMutex.Lock()
	defer op.lastFrameMutex.Unlock()

	return op.lastFrame, nil
}
