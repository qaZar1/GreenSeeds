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

type Default interface {
	BinaryWatercress(gocv.Mat, string) (gocv.Mat, gocv.Mat, gocv.Mat)
	ClassifyWatercressSeeds(gocv.Mat, gocv.Mat, gocv.Mat, gocv.Mat, string) int
}

func (cl *Classification) BinaryWatercress(img gocv.Mat, outDir string) (gocv.Mat, gocv.Mat, gocv.Mat) {
	// -----------------------------
	// 1. R-G цветовой индекс
	// -----------------------------
	ch := gocv.Split(img)

	g := ch[1]
	r := ch[2]

	rg := gocv.NewMat()
	defer rg.Close()

	gocv.Subtract(r, g, &rg)

	gocv.GaussianBlur(rg, &rg, image.Pt(7, 7), 0, 0, gocv.BorderDefault)

	maskRG := gocv.NewMat()
	defer maskRG.Close()

	gocv.Threshold(
		rg,
		&maskRG,
		0,
		255,
		gocv.ThresholdBinary|gocv.ThresholdOtsu,
	)

	// -----------------------------
	// 2. LAB яркость
	// -----------------------------
	lab := gocv.NewMat()
	defer lab.Close()

	gocv.CvtColor(img, &lab, gocv.ColorBGRToLab)

	labCh := gocv.Split(lab)
	L := labCh[0]

	maskL := gocv.NewMat()
	defer maskL.Close()

	gocv.Threshold(
		L,
		&maskL,
		140,
		255,
		gocv.ThresholdBinaryInv,
	)

	// -----------------------------
	// 3. HSV насыщенность
	// -----------------------------
	hsv := gocv.NewMat()
	defer hsv.Close()

	gocv.CvtColor(img, &hsv, gocv.ColorBGRToHSV)

	hsvCh := gocv.Split(hsv)
	S := hsvCh[1]

	maskS := gocv.NewMat()
	defer maskS.Close()

	gocv.Threshold(
		S,
		&maskS,
		60,
		255,
		gocv.ThresholdBinary,
	)

	// -----------------------------
	// 4. объединяем маски
	// -----------------------------
	tmp := gocv.NewMat()
	defer tmp.Close()

	gocv.BitwiseAnd(maskRG, maskS, &tmp)

	clean := gocv.NewMat()
	defer clean.Close()

	gocv.BitwiseAnd(tmp, maskL, &clean)

	// -----------------------------
	// 5. морфология
	// -----------------------------
	kernel := gocv.GetStructuringElement(
		gocv.MorphEllipse,
		image.Pt(3, 3),
	)
	defer kernel.Close()

	gocv.MorphologyEx(clean, &clean, gocv.MorphOpen, kernel)
	gocv.MorphologyEx(clean, &clean, gocv.MorphClose, kernel)

	// -----------------------------
	// 6. компоненты
	// -----------------------------
	labels := gocv.NewMat()
	stats := gocv.NewMat()
	centroids := gocv.NewMat()

	defer labels.Close()
	defer centroids.Close()

	gocv.ConnectedComponentsWithStats(
		clean,
		&labels,
		&stats,
		&centroids,
	)

	return stats.Clone(), img.Clone(), clean.Clone()
}

func (cl *Classification) ClassifyWatercressSeeds(
	solidMask gocv.Mat,
	img gocv.Mat,
	stats gocv.Mat,
	binary gocv.Mat,
	outDir string,
) int {
	mask := gocv.NewMat()
	defer mask.Close()
	gocv.BitwiseAnd(binary, solidMask, &mask)

	// лёгкая очистка шума (не уничтожает мелкие семена)
	openKernel := gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(3, 3))
	defer openKernel.Close()
	gocv.MorphologyEx(mask, &mask, gocv.MorphOpen, openKernel)

	// === МЕДИАНА ПЛОЩАДИ (из твоей первой функции) ===
	var areas []float64
	for i := 1; i < stats.Rows(); i++ {
		a := float64(stats.GetIntAt(i, int(gocv.CC_STAT_AREA)))
		if a >= 18 {
			areas = append(areas, a)
		}
	}
	if len(areas) == 0 {
		return 0
	}
	sort.Float64s(areas)
	medianArea := areas[len(areas)/2]

	// === КОНТУРЫ + ЛОКАЛЬНЫЙ ОБРАБОТКА ===
	contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	defer contours.Close()

	result := img.Clone()
	defer result.Close()

	dilKernel := gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(2, 2))
	defer dilKernel.Close()

	totalSeeds := 0

	for i := 0; i < contours.Size(); i++ {
		cnt := contours.At(i)
		area := gocv.ContourArea(cnt)
		if area < 9 {
			continue
		}

		rect := gocv.BoundingRect(cnt)

		// фильтр формы (из первой функции)
		aspect := math.Min(float64(rect.Min.X), float64(rect.Min.Y)) / math.Max(float64(rect.Max.X), float64(rect.Max.Y))
		if aspect < 0.30 {
			continue
		}

		cluster := mask.Region(rect)

		label := gocv.NewMat()
		defer label.Close()

		// локальный DistanceTransform
		dist := gocv.NewMat()
		gocv.DistanceTransform(cluster, &dist, &label, gocv.DistL2, gocv.DistanceMask5, gocv.DistanceLabelCComp)
		gocv.GaussianBlur(dist, &dist, image.Pt(3, 3), 0, 0, gocv.BorderDefault)
		gocv.Normalize(dist, &dist, 0, 1.0, gocv.NormMinMax)

		peaks := gocv.NewMat()
		gocv.Threshold(dist, &peaks, 0.22, 1.0, gocv.ThresholdBinary) // ниже порог
		gocv.Dilate(peaks, &peaks, dilKernel)                         // усиливаем пики

		peaks8 := gocv.NewMat()
		peaks.ConvertTo(&peaks8, gocv.MatTypeCV8U)

		labels := gocv.NewMat()
		stats2 := gocv.NewMat()
		centroids := gocv.NewMat()

		n := gocv.ConnectedComponentsWithStats(peaks8, &labels, &stats2, &centroids)

		peaksCount := 0
		for j := 1; j < n; j++ {
			if stats2.GetIntAt(j, int(gocv.CC_STAT_AREA)) >= 1 {
				cx := int(centroids.GetDoubleAt(j, 0)) + rect.Min.X
				cy := int(centroids.GetDoubleAt(j, 1)) + rect.Min.Y
				gocv.Circle(&result, image.Pt(cx, cy), 3, color.RGBA{0, 255, 255, 0}, -1) // голубые центры
				peaksCount++
			}
		}

		// === ГИБРИДНЫЙ ПОДСЧЁТ (главное улучшение) ===
		areaBased := int(math.Round(area / medianArea))
		if areaBased < 1 {
			areaBased = 1
		}
		if areaBased > 8 {
			areaBased = 8
		}

		finalCount := peaksCount
		if peaksCount <= 1 && areaBased >= 3 {
			finalCount = areaBased // fallback для очень слипшихся кластеров
		} else if peaksCount == 1 && areaBased == 2 {
			finalCount = 2
		}

		if finalCount < 1 {
			finalCount = 1
		}

		totalSeeds += finalCount

		// визуализация
		col := color.RGBA{0, 255, 0, 0}
		if finalCount > 1 {
			col = color.RGBA{255, 140, 0, 0}
		}
		gocv.Rectangle(&result, rect, col, 2)
		gocv.PutText(&result, fmt.Sprintf("%d", finalCount),
			image.Pt(rect.Min.X+2, rect.Min.Y-5),
			gocv.FontHersheySimplex, 0.55, col, 2)

		// очистка
		cluster.Close()
		dist.Close()
		peaks.Close()
		peaks8.Close()
		labels.Close()
		stats2.Close()
		centroids.Close()
	}

	gocv.IMWrite(filepath.Join(outDir, "watercress_seeds_hybrid.png"), result)
	return totalSeeds
}
