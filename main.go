package main

import (
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

	// Marker labelling (?)
	markers := gocv.NewMat()
	gocv.ConnectedComponents(sureFG8U, &markers)
	markers.AddUChar(1)
	gocv.IMWrite("8.png", markers)
}
