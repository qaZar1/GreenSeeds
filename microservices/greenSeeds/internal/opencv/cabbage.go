package opencv

// // import (
// // 	"image"
// // 	"image/color"
// // 	"path/filepath"

// // 	"gocv.io/x/gocv"
// // )

// // type Cabbage interface {
// // 	BinaryCabbage(gocv.Mat, string) (gocv.Mat, gocv.Mat, gocv.Mat)
// // 	ClassifyCabbageSeeds(gocv.Mat, gocv.Mat, gocv.Mat, gocv.Mat, string) int
// // }

// // func (cl *Classification) BinaryCabbage(img gocv.Mat, outDir string) (gocv.Mat, gocv.Mat, gocv.Mat) {
// // 	// grayscale
// // 	gray := gocv.NewMat()
// // 	defer gray.Close()
// // 	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

// // 	// blur
// // 	blur := gocv.NewMat()
// // 	defer blur.Close()
// // 	gocv.GaussianBlur(gray, &blur, image.Pt(5, 5), 0, 0, gocv.BorderDefault)

// // 	// adaptive threshold
// // 	binary := gocv.NewMat()
// // 	defer binary.Close()

// // 	gocv.AdaptiveThreshold(
// // 		blur,
// // 		&binary,
// // 		255,
// // 		gocv.AdaptiveThresholdGaussian,
// // 		gocv.ThresholdBinaryInv, // семена станут белыми
// // 		21,
// // 		3,
// // 	)

// // 	// морфология (чистим шум)
// // 	kernel := gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(3, 3))
// // 	defer kernel.Close()

// // 	clean := gocv.NewMat()
// // 	defer clean.Close()

// // 	gocv.MorphologyEx(binary, &clean, gocv.MorphOpen, kernel)

// // 	// компоненты
// // 	labels := gocv.NewMat()
// // 	stats := gocv.NewMat()
// // 	centroids := gocv.NewMat()
// // 	defer labels.Close()
// // 	defer stats.Close()
// // 	defer centroids.Close()

// // 	gocv.ConnectedComponentsWithStats(clean, &labels, &stats, &centroids)

// // 	return stats.Clone(), img.Clone(), clean.Clone()
// // }

// // func (cl *Classification) ClassifyCabbageSeeds(
// // 	solidMask gocv.Mat,
// // 	img gocv.Mat,
// // 	stats gocv.Mat,
// // 	binary gocv.Mat,
// // 	outDir string,
// // ) int {
// // 	// --- HSV цветовой фильтр для темных семян ---
// // 	hsv := gocv.NewMat()
// // 	defer hsv.Close()

// // 	gocv.CvtColor(img, &hsv, gocv.ColorBGRToHSV)

// // 	colorMask := gocv.NewMat()
// // 	defer colorMask.Close()

// // 	lower := gocv.NewScalar(0, 40, 0, 0)
// // 	upper := gocv.NewScalar(180, 255, 120, 0)

// // 	gocv.InRangeWithScalar(hsv, lower, upper, &colorMask)

// // 	// --- объединяем binary + colorMask ---
// // 	tmp := gocv.NewMat()
// // 	defer tmp.Close()

// // 	gocv.BitwiseAnd(binary, colorMask, &tmp)

// // 	// --- применяем solidMask ---
// // 	mask := gocv.NewMat()
// // 	defer mask.Close()

// // 	gocv.BitwiseAnd(tmp, solidMask, &mask)

// // 	// --- чистка шума ---
// // 	kernel := gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(3, 3))
// // 	gocv.MorphologyEx(mask, &mask, gocv.MorphOpen, kernel)

// // 	contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)

// // 	result := img.Clone()
// // 	defer result.Close()

// // 	count := 0

// // 	for i := 0; i < contours.Size(); i++ {

// // 		cnt := contours.At(i)

// // 		area := gocv.ContourArea(cnt)
// // 		rect := gocv.BoundingRect(cnt)

// // 		// --- фильтр по яркости (убирает хлеб) ---
// // 		roi := img.Region(rect)
// // 		meanMat := gocv.NewMat()
// // 		stdMat := gocv.NewMat()

// // 		gocv.MeanStdDevWithMask(roi, &meanMat, &stdMat, gocv.NewMat())
// // 		roi.Close()

// // 		if meanMat.Mean().Val1 > 150 {
// 			continue
// 		}

// 		// --- одиночное семя ---
// 		if area > 50 && area < 200 {

// 			count++

// 			gocv.Rectangle(
// 				&result,
// 				rect,
// 				color.RGBA{255, 0, 0, 0},
// 				1,
// 			)

// 			continue
// 		}

// 		// слишком маленькие и слишком большие пропускаем
// 		if area < 200 || area > 5000 {
// 			continue
// 		}

// 		cluster := mask.Region(rect)
// 		defer cluster.Close()

// 		dist := gocv.NewMat()
// 		defer dist.Close()

// 		label := gocv.NewMat()
// 		defer label.Close()

// 		gocv.DistanceTransform(
// 			cluster,
// 			&dist,
// 			&label,
// 			gocv.DistL2,
// 			gocv.DistanceMask5,
// 			gocv.DistanceLabelCComp,
// 		)

// 		gocv.GaussianBlur(dist, &dist, image.Pt(7, 7), 0, 0, gocv.BorderDefault)

// 		gocv.Normalize(dist, &dist, 0, 1.0, gocv.NormMinMax)

// 		centers := gocv.NewMat()
// 		defer centers.Close()

// 		gocv.Threshold(dist, &centers, 0.6, 1.0, gocv.ThresholdBinary)

// 		centers8 := gocv.NewMat()
// 		defer centers8.Close()

// 		centers.ConvertTo(&centers8, gocv.MatTypeCV8U)

// 		labels := gocv.NewMat()
// 		stats2 := gocv.NewMat()
// 		centroids := gocv.NewMat()

// 		n := gocv.ConnectedComponentsWithStats(
// 			centers8,
// 			&labels,
// 			&stats2,
// 			&centroids,
// 		)

// 		for j := 1; j < n; j++ {

// 			a := stats2.GetIntAt(j, int(gocv.CC_STAT_AREA))

// 			if a < 3 {
// 				continue
// 			}

// 			cx := int(centroids.GetDoubleAt(j, 0)) + rect.Min.X
// 			cy := int(centroids.GetDoubleAt(j, 1)) + rect.Min.Y

// 			count++

// 			gocv.Circle(
// 				&result,
// 				image.Pt(cx, cy),
// 				4,
// 				color.RGBA{0, 255, 0, 0},
// 				2,
// 			)
// 		}

// 		labels.Close()
// 		stats2.Close()
// 		centroids.Close()
// 	}

// 	gocv.IMWrite(filepath.Join(outDir, "seeds_detected.png"), result)

// 	return count
// }
