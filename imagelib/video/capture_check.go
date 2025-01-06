package video

import (
	"fmt"
	"gocv.io/x/gocv"
	"time"
)

func CaptureCheck(capture *gocv.VideoCapture) {

	mat := gocv.NewMat()
	window := gocv.NewWindow("Hello")
	i := 0

	for {
		i++
		ok := capture.Read(&mat)
		if i%10 == 0 {
			fmt.Printf("%d %t %s\n", i, ok, time.Now()) // int(webcam.List(gocv.VideoCapturePosFrames)),
		}

		////img, err := mat.ToImage()
		////if err != nil {
		////	log.Fatal(err)
		////}
		////
		////if err = imagelib.SavePNG(img, strconv.Itoa(i)+".png"); err != nil {
		////	log.Fatal(err)
		////}

		if ok {
			window.IMShow(mat)
		}
		window.WaitKey(1)
	}
}
