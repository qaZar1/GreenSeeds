package opencv

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"path/filepath"
	"sort"

	"gocv.io/x/gocv"
)

type Amarant interface {
	Binary(gocv.Mat, string) (gocv.Mat, gocv.Mat, gocv.Mat)
	ClassifyComponents(gocv.Mat, gocv.Mat, gocv.Mat, gocv.Mat, string) int
}

func (cl *Classification) Binary(img gocv.Mat, outDir string) (gocv.Mat, gocv.Mat, gocv.Mat) {
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	gocv.GaussianBlur(gray, &gray, image.Pt(5, 5), 0, 0, gocv.BorderDefault)

	kernel := gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(11, 11))
	defer kernel.Close()

	blackhat := gocv.NewMat()
	defer blackhat.Close()
	gocv.MorphologyEx(gray, &blackhat, gocv.MorphBlackhat, kernel)

	th := gocv.NewMat()
	defer th.Close()
	gocv.Threshold(blackhat, &th, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)

	kernelSmall := gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(3, 3))
	defer kernelSmall.Close()

	clean := gocv.NewMat()
	defer clean.Close()
	gocv.MorphologyEx(th, &clean, gocv.MorphOpen, kernelSmall)

	labels := gocv.NewMat()
	stats := gocv.NewMat()
	centroids := gocv.NewMat()
	defer labels.Close()
	defer stats.Close()
	defer centroids.Close()

	gocv.ConnectedComponentsWithStats(clean, &labels, &stats, &centroids)

	return stats.Clone(), img.Clone(), clean.Clone()
}

func (cl *Classification) ClassifyComponents(
	solidMask gocv.Mat,
	img gocv.Mat,
	stats gocv.Mat,
	binary gocv.Mat,
	outDir string,
) int {
	result := img.Clone()
	defer result.Close()

	dist := gocv.NewMat()
	defer dist.Close()

	labels := gocv.NewMat()
	defer labels.Close()

	gocv.DistanceTransform(
		solidMask,
		&dist,
		&labels,
		gocv.DistL2,
		gocv.DistanceMask3,
		gocv.DistanceLabelCComp,
	)

	// Площадь
	var areas []float64
	for i := 1; i < stats.Rows(); i++ {
		a := float64(stats.GetIntAt(i, int(gocv.CC_STAT_AREA)))
		if a >= 20 {
			areas = append(areas, a)
		}
	}

	if len(areas) == 0 {
		return 0
	}

	sort.Float64s(areas)
	median := areas[len(areas)/2]

	// Классификация
	minDist := float32(8)
	totalSeeds := 0

	for i := 1; i < stats.Rows(); i++ {
		x := stats.GetIntAt(i, int(gocv.CC_STAT_LEFT))
		y := stats.GetIntAt(i, int(gocv.CC_STAT_TOP))
		w := stats.GetIntAt(i, int(gocv.CC_STAT_WIDTH))
		h := stats.GetIntAt(i, int(gocv.CC_STAT_HEIGHT))
		area := float64(stats.GetIntAt(i, int(gocv.CC_STAT_AREA)))

		if area < 15 {
			continue
		}

		// Центр компонента
		cx := x + w/2
		cy := y + h/2

		// Расстояние до края
		d := dist.GetFloatAt(int(cy), int(cx))
		if d < minDist {
			continue
		}

		// Форма
		aspect := math.Min(float64(w), float64(h)) / math.Max(float64(w), float64(h))
		if aspect < 0.3 {
			continue
		}

		// Оценка количества семян
		count := int(math.Round(area / median))
		if count < 1 {
			count = 1
		}
		if count > 6 {
			count = 6
		}

		totalSeeds += count

		col := color.RGBA{0, 255, 0, 0}
		if count > 1 {
			col = color.RGBA{255, 0, 0, 0}
		}

		gocv.Rectangle(
			&result,
			image.Rect(int(x), int(y), int(x+w), int(y+h)),
			col,
			1,
		)

		gocv.PutText(
			&result,
			fmt.Sprintf("%d", count),
			image.Pt(int(x), int(y-3)),
			gocv.FontHersheySimplex,
			0.4,
			col,
			1,
		)
	}

	gocv.IMWrite(filepath.Join(outDir, "6_counted.png"), result)
	return totalSeeds
}
