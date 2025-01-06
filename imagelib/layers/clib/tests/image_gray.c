#include <stdio.h>
#include "image_gray.h"

ImageGray img_check_connect(ImageGray img) {
    printf("x_width: %i, y_height: %i, x_min: %i, y_min: %i, stride: %i, pix[0]: %i, pix[1]: %i, pix[2]: %i\n",
        img.x_width, img.y_height, img.x_min, img.y_min, img.stride, img.pixels[0], img.pixels[1], img.pixels[2]
    );

    return img;
}

void img_br_classes(ImageGray src, Layer* dst) {
//
//	var offset, whiteCnt, blackCnt int
//
//	for y := 0; y < yHeight; y++ {
//		for x := 0; x < xWidth; x++ {
//			v := mask.Calculate(rect.Min.X+x*scale, rect.Min.Y+y*scale)
//			lyrConvolved.Pix[offset+x] = v
//			if v == pix.ValueMax {
//				whiteCnt++
//				maxValue = v
//			} else if v > maxValue {
//				maxValue = v
//			} else if v == 0 {
//				blackCnt++
//				minValue = v
//			} else if v < minValue {
//				minValue = v
//			}
//		}
//		offset += xWidth
//	}
//
//	classNum := mask.imgGray.Pix[offset] / mask.classRange
//	mask.classes[classNum]++
//
//	return mask.classRange * classNum
//

}
