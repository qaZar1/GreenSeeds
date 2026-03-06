package opencv

import (
	"image"
	"image/color"
	"path/filepath"

	"gocv.io/x/gocv"
)

type Counting interface {
	Counter(string, string) int
	RemoveBackground(gocv.Mat, string) (gocv.Mat, gocv.Mat)
	Amarant
	Cabbage
	Default
}

type Classification struct {
	Total int
}

func NewCounting() Counting {
	return &Classification{
		Total: 0,
	}
}

func (cl *Classification) Counter(imagePath string, outDir string) int {
	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		panic("cannot read image")
	}
	defer img.Close()

	origin := img.Clone()
	defer origin.Close()

	processed, solidMask := cl.RemoveBackground(img, outDir)
	defer processed.Close()
	defer solidMask.Close()

	gocv.IMWrite(filepath.Join(outDir, "result.png"), processed)

	stats, img1, binary := cl.Binary(processed, outDir)

	// Маска
	maskedBinary := cl.maskMat(binary, solidMask)
	defer maskedBinary.Close()
	gocv.IMWrite(filepath.Join(outDir, "4_masked_binary.png"), maskedBinary)

	return cl.ClassifyWatercressSeeds(solidMask, img1, stats, maskedBinary, outDir)
}

// maskMat применяет маску к изображению
func (cl *Classification) maskMat(src, mask gocv.Mat) gocv.Mat {
	dst := gocv.NewMatWithSize(src.Rows(), src.Cols(), src.Type())
	defer dst.Close()

	for r := 0; r < src.Rows(); r++ {
		for c := 0; c < src.Cols(); c++ {
			if mask.GetUCharAt(r, c) > 0 {
				dst.SetUCharAt(r, c, src.GetUCharAt(r, c))
			} else {
				dst.SetUCharAt(r, c, 0)
			}
		}
	}
	return dst.Clone()
}

func (cl *Classification) RemoveBackground(src gocv.Mat, debugDir string) (gocv.Mat, gocv.Mat) {
	// --- 1. BGR -> LAB
	lab := gocv.NewMat()
	defer lab.Close()
	gocv.CvtColor(src, &lab, gocv.ColorBGRToLab)

	channels := gocv.Split(lab)
	defer func() {
		for _, ch := range channels {
			ch.Close()
		}
	}()

	L := channels[0]
	// A := channels[1]
	// B := channels[2]

	// --- 2. Усиливаем яркость
	clahe := gocv.NewCLAHEWithParams(3.0, image.Pt(8, 8))
	defer clahe.Close()

	Lboost := gocv.NewMat()
	defer Lboost.Close()

	clahe.Apply(L, &Lboost)

	gocv.Normalize(Lboost, &Lboost, 0, 255, gocv.NormMinMax)

	gocv.IMWrite(filepath.Join(debugDir, "01_L_boost.png"), Lboost)

	// --- 3. Детектируем желтый
	// yellow: высокий B + средний A

	maskYellow := gocv.NewMat()
	defer maskYellow.Close()

	lower := gocv.NewScalar(0, 120, 140, 0)
	upper := gocv.NewScalar(255, 160, 255, 0)

	gocv.InRangeWithScalar(lab, lower, upper, &maskYellow)

	gocv.IMWrite(filepath.Join(debugDir, "02_yellow_mask.png"), maskYellow)

	// --- 4. Морфология

	closeKernel := gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(11, 11))
	defer closeKernel.Close()

	gocv.MorphologyEx(maskYellow, &maskYellow, gocv.MorphClose, closeKernel)

	openKernel := gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(5, 5))
	defer openKernel.Close()

	gocv.MorphologyEx(maskYellow, &maskYellow, gocv.MorphOpen, openKernel)

	gocv.Dilate(maskYellow, &maskYellow, openKernel)

	gocv.IMWrite(filepath.Join(debugDir, "03_mask_clean.png"), maskYellow)

	// --- 5. Контуры

	contours := gocv.FindContours(maskYellow, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	defer contours.Close()

	var largestContour gocv.PointVector
	maxArea := float64(0)

	for i := 0; i < contours.Size(); i++ {
		area := gocv.ContourArea(contours.At(i))
		if area > maxArea {
			maxArea = area
			largestContour = contours.At(i)
		}
	}

	if largestContour.Size() == 0 {
		return gocv.Mat{}, gocv.Mat{}
	}

	// --- 6. Поворот по MinAreaRect

	rect := gocv.MinAreaRect(largestContour)

	center := image.Point{
		X: int(rect.Center.X),
		Y: int(rect.Center.Y),
	}

	rotationMatrix := gocv.GetRotationMatrix2D(center, rect.Angle, 1.0)
	defer rotationMatrix.Close()

	size := image.Pt(src.Cols(), src.Rows())

	rotated := gocv.NewMat()
	defer rotated.Close()

	gocv.WarpAffine(src, &rotated, rotationMatrix, size)

	gocv.IMWrite(filepath.Join(debugDir, "04_rotated.png"), rotated)

	// --- 7. Поворот маски

	rotatedMask := gocv.NewMat()
	defer rotatedMask.Close()

	gocv.WarpAffine(maskYellow, &rotatedMask, rotationMatrix, size)

	// --- 8. Заполняем контур

	solidMask := gocv.NewMatWithSize(rotatedMask.Rows(), rotatedMask.Cols(), gocv.MatTypeCV8UC1)
	defer solidMask.Close()

	contoursRotated := gocv.FindContours(rotatedMask, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	defer contoursRotated.Close()

	var largestRotated gocv.PointVector
	maxRotatedArea := float64(0)

	for i := 0; i < contoursRotated.Size(); i++ {
		area := gocv.ContourArea(contoursRotated.At(i))
		if area > maxRotatedArea {
			maxRotatedArea = area
			largestRotated = contoursRotated.At(i)
		}
	}

	if largestRotated.Size() == 0 {
		return gocv.Mat{}, gocv.Mat{}
	}

	contourList := gocv.NewPointsVectorFromPoints([][]image.Point{
		largestRotated.ToPoints(),
	})

	gocv.DrawContours(&solidMask, contourList, -1, color.RGBA{255, 255, 255, 0}, -1)

	gocv.IMWrite(filepath.Join(debugDir, "05_solid_mask.png"), solidMask)

	// --- 9. Вырезаем объект

	finalResult := gocv.NewMat()
	defer finalResult.Close()

	white := gocv.NewMatWithSizeFromScalar(
		gocv.NewScalar(255, 255, 255, 0),
		rotated.Rows(),
		rotated.Cols(),
		gocv.MatTypeCV8UC3,
	)
	defer white.Close()

	rotated.CopyToWithMask(&finalResult, solidMask)

	bgMask := gocv.NewMat()
	defer bgMask.Close()

	gocv.BitwiseNot(solidMask, &bgMask)

	white.CopyToWithMask(&finalResult, bgMask)

	gocv.IMWrite(filepath.Join(debugDir, "06_final.png"), finalResult)

	return finalResult.Clone(), solidMask.Clone()
}

// func (cl *Classification) RemoveBackground(src gocv.Mat, debugDir string) (gocv.Mat, gocv.Mat) {
// 	// 1. Преобразуем в HSV
// 	hsv := gocv.NewMat()
// 	defer hsv.Close()
// 	gocv.CvtColor(src, &hsv, gocv.ColorBGRToHSV)

// 	channels := gocv.Split(hsv)
// 	defer func() {
// 		for _, ch := range channels {
// 			ch.Close()
// 		}
// 	}()

// 	v := channels[2] // V-канал — яркость

// 	// Выравниваем освещение
// 	clahe := gocv.NewCLAHEWithParams(2.5, image.Pt(8, 8))
// 	defer clahe.Close()
// 	equalizedV := gocv.NewMat()
// 	defer equalizedV.Close()
// 	clahe.Apply(v, &equalizedV)

// 	// Блюр
// 	blurred := gocv.NewMat()
// 	defer blurred.Close()
// 	gocv.GaussianBlur(equalizedV, &blurred, image.Pt(15, 15), 0, 0, gocv.BorderDefault)

// 	// Маска
// 	maskV := gocv.NewMat()
// 	defer maskV.Close()
// 	gocv.Threshold(blurred, &maskV, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)

// 	// Маска по С каналу
// 	s := channels[1]
// 	sBlur := gocv.NewMat()
// 	defer sBlur.Close()
// 	gocv.GaussianBlur(s, &sBlur, image.Pt(9, 9), 0, 0, gocv.BorderDefault)

// 	maskS := gocv.NewMat()
// 	defer maskS.Close()
// 	gocv.Threshold(sBlur, &maskS, 40, 255, gocv.ThresholdBinary)

// 	// Соединяем маски
// 	mask := gocv.NewMat()
// 	defer mask.Close()
// 	gocv.BitwiseOr(maskV, maskS, &mask)

// 	// Заполняем дырки
// 	closeKernel := gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(9, 9))
// 	defer closeKernel.Close()
// 	gocv.MorphologyEx(mask, &mask, gocv.MorphClose, closeKernel)

// 	dilateKernel := gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(5, 5))
// 	defer dilateKernel.Close()
// 	gocv.Dilate(mask, &mask, dilateKernel)

// 	gocv.IMWrite(filepath.Join(debugDir, "01_mask.png"), mask)

// 	contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)
// 	defer contours.Close()

// 	var largestContour gocv.PointVector
// 	maxArea := float64(0)
// 	for i := 0; i < contours.Size(); i++ {
// 		area := gocv.ContourArea(contours.At(i))
// 		if area > maxArea {
// 			maxArea = area
// 			largestContour = contours.At(i)
// 		}
// 	}

// 	if largestContour.Size() == 0 {
// 		return gocv.Mat{}, gocv.Mat{}
// 	}

// 	// MinAreaRect + поворот
// 	rect := gocv.MinAreaRect(largestContour)
// 	angle := rect.Angle
// 	center := image.Point{X: int(rect.Center.X), Y: int(rect.Center.Y)}
// 	rotationMatrix := gocv.GetRotationMatrix2D(center, angle, 1.0)
// 	defer rotationMatrix.Close()

// 	size := image.Pt(src.Cols(), src.Rows())

// 	rotated := gocv.NewMat()
// 	defer rotated.Close()
// 	gocv.WarpAffine(src, &rotated, rotationMatrix, size)

// 	gocv.IMWrite(filepath.Join(debugDir, "02_rotated.png"), rotated)

// 	// Поворот
// 	rotatedMask := gocv.NewMat()
// 	defer rotatedMask.Close()
// 	gocv.WarpAffine(mask, &rotatedMask, rotationMatrix, size)

// 	// Заполняем контур
// 	solidMask := gocv.NewMatWithSize(rotatedMask.Rows(), rotatedMask.Cols(), gocv.MatTypeCV8UC1)
// 	defer solidMask.Close()

// 	// Заливаем белым
// 	for r := 0; r < solidMask.Rows(); r++ {
// 		for c := 0; c < solidMask.Cols(); c++ {
// 			solidMask.SetUCharAt(r, c, 0)
// 		}
// 	}

// 	contoursRotated := gocv.FindContours(rotatedMask, gocv.RetrievalExternal, gocv.ChainApproxSimple)
// 	defer contoursRotated.Close()

// 	var largestRotated gocv.PointVector
// 	maxRotatedArea := float64(0)
// 	for i := 0; i < contoursRotated.Size(); i++ {
// 		area := gocv.ContourArea(contoursRotated.At(i))
// 		if area > maxRotatedArea {
// 			maxRotatedArea = area
// 			largestRotated = contoursRotated.At(i)
// 		}
// 	}

// 	if largestRotated.Size() == 0 {
// 		return gocv.Mat{}, gocv.Mat{}
// 	}

// 	contourList := gocv.NewPointsVectorFromPoints([][]image.Point{largestRotated.ToPoints()})
// 	gocv.DrawContours(&solidMask, contourList, -1, color.RGBA{255, 255, 255, 0}, -1)

// 	gocv.IMWrite(filepath.Join(debugDir, "03_solid_mask.png"), solidMask)

// 	// Белый фон + комбинация
// 	background := gocv.NewMatWithSize(rotated.Rows(), rotated.Cols(), gocv.MatTypeCV8UC3)
// 	defer background.Close()
// 	for r := 0; r < rotated.Rows(); r++ {
// 		for c := 0; c < rotated.Cols(); c++ {
// 			background.SetUCharAt(r, c*3+0, 255)
// 			background.SetUCharAt(r, c*3+1, 255)
// 			background.SetUCharAt(r, c*3+2, 255)
// 		}
// 	}

// 	finalResult := gocv.NewMatWithSize(rotated.Rows(), rotated.Cols(), gocv.MatTypeCV8UC3)
// 	defer finalResult.Close()

// 	for r := 0; r < rotated.Rows(); r++ {
// 		for c := 0; c < rotated.Cols(); c++ {
// 			if solidMask.GetUCharAt(r, c) > 0 {
// 				finalResult.SetUCharAt(r, c*3+0, rotated.GetUCharAt(r, c*3+0))
// 				finalResult.SetUCharAt(r, c*3+1, rotated.GetUCharAt(r, c*3+1))
// 				finalResult.SetUCharAt(r, c*3+2, rotated.GetUCharAt(r, c*3+2))
// 			} else {
// 				finalResult.SetUCharAt(r, c*3+0, 255)
// 				finalResult.SetUCharAt(r, c*3+1, 255)
// 				finalResult.SetUCharAt(r, c*3+2, 255)
// 			}
// 		}
// 	}

// 	gocv.IMWrite(filepath.Join(debugDir, "05_bg_replaced.png"), finalResult)

// 	return finalResult.Clone(), solidMask.Clone()
// }
