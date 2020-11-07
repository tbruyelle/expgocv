package main

import (
	"fmt"
	"testing"

	"gocv.io/x/gocv"
)

func TestMat(t *testing.T) {
	m := gocv.NewMatWithSize(2, 2, gocv.MatTypeCV32S)
	defer m.Close()
	n := gocv.NewMatWithSize(2, 2, gocv.MatTypeCV8U)
	defer n.Close()

	m.SetIntAt(0, 0, 10)
	m.SetIntAt(1, 0, 20)
	m.SetIntAt(0, 1, 30)
	m.SetIntAt(1, 1, 40)

	n.SetTo(gocv.NewScalar(1, 0, 0, 0))
	fmt.Println(n.DataPtrInt8())

	fmt.Println(m.DataPtrInt8())
	fmt.Println(m.GetIntAt(0, 0))
	fmt.Println(m.GetIntAt(1, 0))

	o := gocv.NewMat()
	m.CopyToWithMask(&o, n)
	fmt.Println(o.DataPtrInt8())
	fmt.Println(o.GetIntAt(0, 0))
	fmt.Println(o.GetIntAt(1, 0))

}
