package main

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

// image segmentation with wathershed algorithm
// following this tutorial https://docs.opencv.org/master/d3/db4/tutorial_py_watershed.html

func main() {
	img := gocv.IMRead("water_coins.jpg", gocv.IMReadUnchanged)
	// color palette to Gray
	gray := gocv.NewMat()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)
	gocv.IMWrite("1.png", gray)

	// Apply Otsu threshold
	otsu := gocv.NewMat()
	gocv.Threshold(gray, &otsu, 0, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)
	gocv.IMWrite("2.png", otsu)

	// Remove noise
	noNoise := gocv.NewMat()
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	defer kernel.Close()
	gocv.MorphologyExWithParams(otsu, &noNoise, gocv.MorphOpen, kernel, 2, gocv.BorderConstant)
	gocv.IMWrite("3.png", noNoise)

	// sure background area
	sureBG := gocv.NewMat()
	gocv.Dilate(noNoise, &sureBG, kernel)
	gocv.Dilate(sureBG, &sureBG, kernel)
	gocv.Dilate(sureBG, &sureBG, kernel)
	gocv.IMWrite("4.png", sureBG)

	// sure foreground area
	distanceTransform := gocv.NewMat()
	labels := gocv.NewMat()
	gocv.DistanceTransform(noNoise, &distanceTransform, &labels, gocv.DistL2, gocv.DistanceMask5, 0)
	gocv.IMWrite("5.png", distanceTransform)
	sureFG := gocv.NewMat()
	_, max, _, _ := gocv.MinMaxIdx(distanceTransform)
	gocv.Threshold(distanceTransform, &sureFG, .7*max, 255, gocv.ThresholdBinary)
	gocv.IMWrite("6.png", sureFG)

	// substract fg and bg
	sureFG8U := gocv.NewMat()
	sureFG.ConvertTo(&sureFG8U, gocv.MatTypeCV8U)
	substract := gocv.NewMat()
	gocv.Subtract(sureBG, sureFG8U, &substract)
	gocv.IMWrite("7.png", substract)

	// Marker labelling
	markers := gocv.NewMat()
	gocv.ConnectedComponents(sureFG8U, &markers)
	gocv.IMWrite("8.png", markers)
	// Add one to all labels so that sure background is not 0, but 1
	markers.AddUChar(1)
	gocv.IMWrite("9.png", markers)
	// Now, mark the region of unknown with zero
	for row := 0; row < substract.Rows(); row++ {
		for col := 0; col < substract.Cols(); col++ {
			if substract.GetUCharAt(row, col) == 255 {
				markers.SetIntAt(row, col, 0)
			}
		}
	}
	gocv.IMWrite("10.png", markers)
	convert := gocv.NewMat()
	markers.ConvertTo(&convert, gocv.MatTypeCV8U)
	colored := gocv.NewMat()
	gocv.ApplyColorMap(convert, &colored, gocv.ColormapJet)
	gocv.IMWrite("12.png", colored)

	// apply watershed
	gocv.Watershed(img, &markers)

	gocv.IMWrite("13.png", markers)

	// update img from markers
	var m int
	for row := 0; row < markers.Rows(); row++ {
		for col := 0; col < markers.Cols(); col++ {
			if markers.GetIntAt(row, col) == -1 {
				m++
				img.SetUCharAt(row, col*3, 0)
				img.SetUCharAt(row, col*3+1, 0)
				img.SetUCharAt(row, col*3+2, 255)
			}
		}
	}
	fmt.Println("updated", m)
	gocv.IMWrite("final.png", img)
}
