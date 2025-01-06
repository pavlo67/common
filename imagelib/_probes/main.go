package main

import (
	"fmt"
	"log"
	"time"

	"gocv.io/x/gocv"
)

func main() {

	// gocv.GetgetBuildInformation

	mat := gocv.NewMat()

	//webcam, err := gocv.OpenVideoCaptureWithAPI("sudo gst-launch-1.0 -vc udpsrc port=5600"+
	//	" close-socket=false auto-multicast=true ! application/x-rtp, payload=96 ! rtph264depay ! decodebin3 "+
	//	" ! fpsdisplaysink sync=false ! appsink drop=1", gocv.VideoCaptureGstreamer)

	webcam, err := gocv.OpenVideoCaptureWithAPI("gst-launch-1.0 -vc udpsrc port=5600 ! application/x-rtp,payload=96,encoding-name=H264 ! rtpjitterbuffer mode=1 ! rtph264depay ! h264parse ! decodebin ! videoconvert ! appsink", gocv.VideoCaptureGstreamer)

	// webcam, err := gocv.OpenVideoCapture("udp://127.0.0.1:5600")
	// webcam, err := gocv.VideoCaptureFile("_in/_landmarks/anat_tiahur/trim10.mp4")

	if err != nil {
		log.Fatal(err)
	}

	//webcam.Set(gocv.VideoCaptureFOURCC, webcam.ToCodec("MJPG"))
	//webcam.Set(gocv.VideoCaptureAutoExposure, 0)
	//webcam.Set(gocv.VideoCaptureFrameWidth, 640)
	//webcam.Set(gocv.VideoCaptureFrameHeight, 480)
	//webcam.Set(gocv.VideoCaptureFPS, 30)
	//
	//fmt.Println(webcam.Get(gocv.VideoCaptureFPS))

	// webcam.Set(gocv.VideoCaptureBufferSize, 1)
	//webcam.Set(gocv.VideoCaptureFrameWidth, 1280)
	//webcam.Set(gocv.VideoCaptureFrameHeight, 720)

	//fmt.Println(webcam.List(gocv.VideoCaptureBufferSize))

	window := gocv.NewWindow("Hello1")
	i := 0

	for {
		i++
		// webcam.Grab(50)
		//webcam.Set(gocv.VideoCapturePosMsec, float64(i*1000))
		// webcam.Grab(30)

		ok := webcam.Read(&mat)
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
