package opencv

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"gocv.io/x/gocv"
)

type Classification struct {
	Stats gocv.Mat
	Img   gocv.Mat
}

func NewCounting() Classification {
	return Classification{
		Stats: gocv.Mat{},
		Img:   gocv.Mat{},
	}
}

// maskMat применяет маску к изображению: dst[i,j] = src[i,j] если mask[i,j] != 0, иначе 0
func maskMat(src, mask gocv.Mat) gocv.Mat {
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

func BuildCleanDotMaskWithMask(src gocv.Mat, debugDir string) (gocv.Mat, gocv.Mat) {
	hsv := gocv.NewMat()
	defer hsv.Close()
	gocv.CvtColor(src, &hsv, gocv.ColorBGRToHSV)

	// Define HSV range for substrate (yellowish)
	lowerMat := gocv.NewMatWithSize(1, 3, gocv.MatTypeCV8UC3)
	upperMat := gocv.NewMatWithSize(1, 3, gocv.MatTypeCV8UC3)
	defer lowerMat.Close()
	defer upperMat.Close()

	// Set values for lower bound (H:15, S:50, V:100)
	lowerMat.SetUCharAt(0, 0, 15)  // H
	lowerMat.SetUCharAt(0, 1, 50)  // S
	lowerMat.SetUCharAt(0, 2, 100) // V

	// Set values for upper bound (H:35, S:255, V:255)
	upperMat.SetUCharAt(0, 0, 35)  // H
	upperMat.SetUCharAt(0, 1, 255) // S
	upperMat.SetUCharAt(0, 2, 255) // V

	mask := gocv.NewMat()
	defer mask.Close()
	gocv.InRange(hsv, lowerMat, upperMat, &mask)

	// Remove noise with morphological opening
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(5, 5))
	defer kernel.Close()
	gocv.MorphologyEx(mask, &mask, gocv.MorphOpen, kernel)

	gocv.IMWrite(filepath.Join(debugDir, "01_mask.png"), mask)

	// Find contours
	contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	defer contours.Close()

	// Find the largest contour
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

	// Get minimum area rectangle (rotated bounding box)
	rect := gocv.MinAreaRect(largestContour)
	angle := rect.Angle

	// Calculate rotation matrix to align the substrate horizontally
	center := image.Point{X: int(rect.Center.X), Y: int(rect.Center.Y)}
	rotationMatrix := gocv.GetRotationMatrix2D(center, angle, 1.0)
	defer rotationMatrix.Close()

	// Rotate the original image
	rotated := gocv.NewMat()
	defer rotated.Close()
	size := image.Pt(src.Cols(), src.Rows())
	gocv.WarpAffine(src, &rotated, rotationMatrix, size)

	gocv.IMWrite(filepath.Join(debugDir, "02_rotated.png"), rotated)

	// Also rotate the mask to align with the rotated image
	rotatedMask := gocv.NewMat()
	defer rotatedMask.Close()
	gocv.WarpAffine(mask, &rotatedMask, rotationMatrix, size)

	// === Create solid mask (fill holes inside substrate) ===
	solidMask := gocv.NewMatWithSize(rotatedMask.Rows(), rotatedMask.Cols(), gocv.MatTypeCV8UC1)
	defer solidMask.Close()

	// Fill solidMask with zeros manually
	for r := 0; r < solidMask.Rows(); r++ {
		for c := 0; c < solidMask.Cols(); c++ {
			solidMask.SetUCharAt(r, c, 0)
		}
	}

	// Find contours in rotated mask again
	contoursRotated := gocv.FindContours(rotatedMask, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	defer contoursRotated.Close()

	var largestRotatedContour gocv.PointVector
	maxRotatedArea := float64(0)
	for i := 0; i < contoursRotated.Size(); i++ {
		area := gocv.ContourArea(contoursRotated.At(i))
		if area > maxRotatedArea {
			maxRotatedArea = area
			largestRotatedContour = contoursRotated.At(i)
		}
	}

	if largestRotatedContour.Size() == 0 {
		return gocv.Mat{}, gocv.Mat{}
	}

	// Simply draw the largest contour filled to create solid mask
	contourList := gocv.NewPointsVectorFromPoints([][]image.Point{largestRotatedContour.ToPoints()})
	gocv.DrawContours(&solidMask, contourList, -1, color.RGBA{255, 255, 255, 0}, -1)

	gocv.IMWrite(filepath.Join(debugDir, "03_solid_mask.png"), solidMask)

	// === Create white background ===
	background := gocv.NewMatWithSize(rotated.Rows(), rotated.Cols(), gocv.MatTypeCV8UC3)
	defer background.Close()

	// Fill with white (B=255, G=255, R=255)
	for r := 0; r < rotated.Rows(); r++ {
		for c := 0; c < rotated.Cols(); c++ {
			background.SetUCharAt(r, c*3+0, 255) // B
			background.SetUCharAt(r, c*3+1, 255) // G
			background.SetUCharAt(r, c*3+2, 255) // R
		}
	}

	// === Combine: object from rotated where solidMask != 0, else white ===
	finalResult := gocv.NewMatWithSize(rotated.Rows(), rotated.Cols(), gocv.MatTypeCV8UC3)
	defer finalResult.Close()

	for r := 0; r < rotated.Rows(); r++ {
		for c := 0; c < rotated.Cols(); c++ {
			if solidMask.GetUCharAt(r, c) > 0 {
				// Copy pixel from rotated
				finalResult.SetUCharAt(r, c*3+0, rotated.GetUCharAt(r, c*3+0))
				finalResult.SetUCharAt(r, c*3+1, rotated.GetUCharAt(r, c*3+1))
				finalResult.SetUCharAt(r, c*3+2, rotated.GetUCharAt(r, c*3+2))
			} else {
				// Keep white
				finalResult.SetUCharAt(r, c*3+0, 255)
				finalResult.SetUCharAt(r, c*3+1, 255)
				finalResult.SetUCharAt(r, c*3+2, 255)
			}
		}
	}

	gocv.IMWrite(filepath.Join(debugDir, "05_bg_replaced.png"), finalResult)

	return finalResult.Clone(), solidMask.Clone()
}

func (cl *Classification) Counter(imagePath string, outDir string) int {
	// img := gocv.IMRead(imagePath, gocv.IMReadColor)
	// if img.Empty() {
	// 	panic("cannot read image")
	// }
	// defer img.Close()
	data, _ := os.ReadFile(imagePath)
	// origin := img.Clone()
	// defer origin.Close()

	// processed, solidMask := BuildCleanDotMaskWithMask(img, outDir)
	// defer processed.Close()
	// defer solidMask.Close()

	// gocv.IMWrite(filepath.Join(outDir, "processed.png"), processed)
	// gocv.IMWrite(filepath.Join(outDir, "solidMask.png"), solidMask)

	cmd := exec.Command("python3", "test.py")

	writer, _ := cmd.StdinPipe()
	reader, _ := cmd.StdoutPipe()
	stdErr, _ := cmd.StderrPipe()

	cmd.Start()

	go func() {
		data, err := io.ReadAll(stdErr)
		if err != nil {
			return
		}

		fmt.Println(string(data))
	}()

	writer.Write(data)
	writer.Close()

	buf, err := io.ReadAll(reader)
	if err != nil {
		return -1
	}

	if err := cmd.Wait(); err != nil {
		return -1
	}

	fmt.Printf("%s", string(buf))
	// result, _ := gocv.IMDecode(buf, gocv.IMReadColor)
	// defer result.Close()
	// gocv.IMWrite(filepath.Join(outDir, "result.jpg"), result)

	return 0
}
