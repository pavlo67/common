package preparation

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib"
	"image"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gocv.io/x/gocv"

	"github.com/pavlo67/common/common/serialization"

	"github.com/pavlo67/imagelib/sources"
	"github.com/pavlo67/imagelib/video"
)

func TestVideoToImages(t *testing.T) {

	if os.Getenv("T") != "" {
		t.Skip()
	}

	device := "/home/pavlo/0/_apacer/_data/landmarks/at/trim.mp4"
	imagesPath := "/home/pavlo/0/partner/_/"

	momentMin := time.Minute*10 + time.Second*15
	momentMax := time.Minute*11 + time.Second*52

	nGap := 4

	capture, err := gocv.OpenVideoCapture(device)
	require.NoError(t, err)
	require.NotNil(t, capture)

	fps := capture.Get(gocv.VideoCaptureFPS)
	xWidth := int(capture.Get(gocv.VideoCaptureFrameWidth))
	yHeight := int(capture.Get(gocv.VideoCaptureFrameHeight))

	t.Log(fps)

	nMin := int(math.Ceil(fps * float64(momentMin/time.Second)))
	nMax := int(math.Ceil(fps * float64(momentMax/time.Second)))

	fpsDivider := nGap + 1

	err = serialization.Save(video.Info{
		Device:     device,
		NFrom:      nMin,
		FPS:        fps,
		FPSDivider: &fpsDivider,
		Rectangle:  image.Rectangle{Max: image.Point{xWidth, yHeight}},
	}, serialization.MarshalerJSON, filepath.Join(imagesPath, sources.VideoInfoFilename))
	require.NoError(t, err)

	fmt.Println(nMin, nMax)

	n := nMin
	for n > 0 {
		nToGrab := min(100, n)
		n -= nToGrab

		capture.Grab(nToGrab)
		fmt.Printf("GRABBED: %d / %d\n", nMin-n, nMin)
	}

	mat := gocv.NewMat()
	defer mat.Close()

	i := 1

	for n = nMin; n < nMax; {
		if i%10 == 0 {
			fmt.Printf("PROCESSED: %d / %d\n", n-nMin, nMax-nMin)
		}
		ok := capture.Read(&mat)
		require.Truef(t, ok, "n: %d", n-nMin)

		img, err := mat.ToImage()
		require.NoError(t, err)
		require.NotNil(t, img)

		resFilename := filepath.Join(imagesPath, fmt.Sprintf("%04d.pgm", i))
		t.Log(resFilename)

		err = imagelib.SavePGM(img, resFilename)
		require.NoError(t, err)

		capture.Grab(nGap)
		n += nGap + 1

		i++
	}

}
