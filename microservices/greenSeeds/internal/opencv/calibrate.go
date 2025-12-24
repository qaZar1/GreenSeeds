package opencv

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"gocv.io/x/gocv"
)

type Calibration struct {
	Center image.Point
	Side   float64
}

func NewCalibration() *Calibration {
	return &Calibration{}
}

func (calib *Calibration) Calibrate(buf []byte) (image.Point, float64, bool) {
	img, err := gocv.IMDecode(buf, gocv.IMReadColor)
	if err != nil {
		log.Println("cannot read image")
		return image.Point{}, 0, false
	}
	defer img.Close()

	log.Println("img size:", img.Cols(), img.Rows())

	if img.Rows() == 0 || img.Cols() == 0 || img.Empty() {
		log.Println("decoded image has zero size")
		return image.Point{}, 0, false
	}

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	blur := gocv.NewMat()
	defer blur.Close()
	gocv.GaussianBlur(gray, &blur, image.Pt(5, 5), 0, 0, gocv.BorderDefault)

	edges := gocv.NewMat()
	defer edges.Close()
	gocv.Canny(blur, &edges, 50, 150)

	contours := gocv.FindContours(edges, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	var bestRect gocv.RotatedRect
	found := false
	maxArea := 0.0
	var bestSide float64

	for i := 0; i < contours.Size(); i++ {
		c := contours.At(i)

		perim := gocv.ArcLength(c, true)
		approx := gocv.ApproxPolyDP(c, 0.02*perim, true)
		defer approx.Close()

		pts := approx.ToPoints()

		if approx.Size() != 4 || !calib.isRectangle(approx) || !calib.isSquare(pts) {
			continue
		}

		area := math.Abs(gocv.ContourArea(approx))
		if area < 1000 {
			continue
		}

		rect := gocv.MinAreaRect(approx)

		if rect.Width <= 0 || rect.Height <= 0 {
			continue
		}

		if !found || area > maxArea {
			bestRect = rect
			bestSide = calib.squareSidePx(pts)
			maxArea = area
			found = true
		}
	}

	if !found {
		log.Println("square not found")
		return image.Point{}, 0, false
	}

	center := bestRect.Center

	if bestSide <= 0 {
		return image.Point{}, 0, false
	}

	// визуализация
	pts := bestRect.Points
	rectPts := gocv.NewPointVectorFromPoints(pts)
	defer rectPts.Close()

	pv := gocv.NewPointsVector()
	defer pv.Close()

	pv.Append(rectPts)

	gocv.Polylines(&img, pv, true, color.RGBA{0, 255, 0, 0}, 3)
	gocv.Circle(&img, center, 5, color.RGBA{255, 0, 0, 0}, 3)

	gocv.PutText(
		&img,
		fmt.Sprintf("%d, %d", center.X, center.Y),
		image.Pt(center.X+10, center.Y),
		gocv.FontHersheySimplex,
		0.6,
		color.RGBA{255, 0, 0, 0},
		2,
	)

	gocv.IMWrite("output.png", img)

	return center, bestSide, true
}

func (calib *Calibration) isRectangle(pts gocv.PointVector) bool {
	if pts.Size() != 4 {
		return false
	}

	const angleThresh = 0.3 // допуск

	for i := 0; i < 4; i++ {
		p0 := pts.At(i)
		p1 := pts.At((i + 1) % 4)
		p2 := pts.At((i + 2) % 4)

		v1x := float64(p0.X - p1.X)
		v1y := float64(p0.Y - p1.Y)
		v2x := float64(p2.X - p1.X)
		v2y := float64(p2.Y - p1.Y)

		dot := v1x*v2x + v1y*v2y
		n1 := v1x*v1x + v1y*v1y
		n2 := v2x*v2x + v2y*v2y

		cos := dot / (math.Sqrt(n1*n2) + 1e-10)
		if calib.abs(cos) > angleThresh {
			return false
		}
	}
	return true
}

func (calib *Calibration) isSquare(pts []image.Point) bool {
	if len(pts) != 4 {
		return false
	}

	var sides [4]float64
	for i := 0; i < 4; i++ {
		sides[i] = calib.distance(pts[i], pts[(i+1)%4])
	}

	avg := (sides[0] + sides[1] + sides[2] + sides[3]) / 4.0

	for i := 0; i < 4; i++ {
		if math.Abs(sides[i]-avg) > avg*0.15 {
			return false
		}
	}

	return true
}

func (calib *Calibration) abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func (c *Calibration) distance(a, b image.Point) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

func (c *Calibration) Finder(buf1, buf2 []byte) (float64, float64, error) {
	c1, s1, ok1 := c.Calibrate(buf1)
	c2, s2, ok2 := c.Calibrate(buf2)

	if s1 <= 0 || s2 <= 0 {
		return 0, 0, errors.New("invalid calibration result")
	}

	if !ok1 || !ok2 {
		return 0, 0, errors.New("calibration failed")
	}

	// реальный размер квадрата (см)
	const squareSizeCm = 5.0

	pxPerCm := ((s1 + s2) / 2.0) / squareSizeCm

	dxPx := float64(c2.X - c1.X)
	dyPx := float64(c2.Y - c1.Y)

	log.Printf(
		"dx = %.2f cm, dy = %.2f cm",
		dxPx/pxPerCm,
		dyPx/pxPerCm,
	)
	return dxPx / pxPerCm, dyPx / pxPerCm, nil
}

func (c *Calibration) squareSidePx(pts []image.Point) float64 {
	var sum float64
	for i := 0; i < 4; i++ {
		sum += c.distance(pts[i], pts[(i+1)%4])
	}
	return sum / 4.0
}
