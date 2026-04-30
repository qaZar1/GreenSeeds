package application

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/camera"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/opencv"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
)

//go:generate mockgen -source=calibration.go -destination=./../mocks/mock_calibration.go -package=mocks
type ICalibrationApp interface {
	CalibrationHandshake() (string, error)
	GetPhoto(sessionId, number string) ([]byte, error)
	CalculateResult(models.Calibration) (models.Calibration, error)
	Clear(sessionId string) error
	Save(sessionId string) error
	BytesFromPhoto(path string) ([]byte, error)
}

type Calibration struct {
	repo 		*repository.Repository
	camera 		camera.ICamera
	calibrate 	map[string]models.Calibration
	opencv      *opencv.Calibration
	validate    *validator.Validate
}

func NewCalibration(
	repo *repository.Repository,
	camera camera.ICamera,
) ICalibrationApp {
	validate := validator.New()

	opencv := opencv.NewCalibration()
	return &Calibration{
		repo: repo,
		camera: camera,
		calibrate: make(map[string]models.Calibration),
		opencv: opencv,
		validate: validate,
	}
}

func (c *Calibration) CalibrationHandshake() (string, error) {
	// if ok := c.ws.Serial.InitialHandshake(); !ok {
	// 	return "", errors.New("Failed to initial handshake")
	// }

	sessionId := uuid.NewString()
	now := time.Now()

	// добавляем в базу
	ok, err := c.repo.SQLite.AddCalibration(sessionId, now)
	if err != nil || !ok {
		return "", err
	}

	toStart, err := c.repo.DevSet.GetSettingsByKey("toStart")
	if err != nil {
		return "", err
	}

	if toStart == (models.DeviceSettings{}) {
		return "", errors.New("toStart is empty")
	}

	// msg := c.ws.Serial.RunGcode(toStart.Value)
	// if msg.Error != nil {
	// 	return "", errors.New(*msg.Error)
	// }

	return sessionId, nil
}

func (c *Calibration) GetPhoto(sessionId string, numberPhotoStr string) ([]byte, error) {
	numberPhoto, err := strconv.Atoi(numberPhotoStr)
	if err != nil {
		return nil, err
	}

	cal, err := c.repo.SQLite.GetCalibration(sessionId)
	if err != nil {
		return nil, err
	}

	buf, err := c.camera.TakePhoto()
	if err != nil {
		return nil, err
	}

	if buf == nil || buf.Len() == 0 {
		return nil, errors.New("photo is nil")
	}

	photo := buf.Bytes()

	var path string
	switch numberPhoto {
	case 1:
		path = fmt.Sprintf("./tmp/%s/%s.jpg", sessionId, "1")
		cal.FirstPhotoPath = &path

		if err := c.camera.SavePhoto(path, sessionId, buf); err != nil {
			return nil, err
		}

		// toEnd, err := c.repo.DevSet.GetSettingsByKey("toEnd")
		// if err != nil {
		// 	return nil, err
		// }

		// if toEnd == (models.DeviceSettings{}) {
		// 	return nil, errors.New("toEnd is empty")
		// }

		// msg := c.ws.Serial.RunGcode(toEnd.Value)
		// if msg.Error != nil {
		// 	return nil, errors.New(*msg.Error)
		// }

	case 2:
		if cal.FirstPhotoPath == nil {
			return nil, errors.New("first photo is nil")
		}

		path = fmt.Sprintf("./tmp/%s/%s.jpg", sessionId, "2")
		cal.SecondPhotoPath = &path

		if err := c.camera.SavePhoto(path, sessionId, buf); err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("number of photo is invalid")
	}

	ok, err := c.repo.SQLite.UpdateCalibration(cal, sessionId)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("failed to update calibration")
	}

	return photo, nil
}

func (c *Calibration) CalculateResult(calibration models.Calibration) (models.Calibration, error) {
	if err := c.validate.Struct(calibration); err != nil {
		return models.Calibration{}, err
	}

	cal, err := c.repo.SQLite.GetCalibration(calibration.SessionId)
	if err != nil {
		return models.Calibration{}, err
	}

	firstPath := cal.FirstPhotoPath
	secondPath := cal.SecondPhotoPath
	if firstPath == nil || secondPath == nil {
		return models.Calibration{}, errors.New("Photos not ready")
	}

	firstPhoto, err := c.BytesFromPhoto(*firstPath)
	if err != nil {
		return models.Calibration{}, err
	}

	secondPhoto, err := c.BytesFromPhoto(*secondPath)
	if err != nil {
		return models.Calibration{}, err
	}

	dx, dy, err := c.opencv.Finder(firstPhoto, secondPhoto)
	if err != nil {
		return models.Calibration{}, err
	}
	if dx == 0 && dy == 0 {
		return models.Calibration{}, errors.New("Distance not changed")
	}

	max := math.Max(math.Abs(dx), math.Abs(dy))
	dPerStep := max / *calibration.Steps

	cal = models.Calibration{
		SessionId: calibration.SessionId,
		Dx:        &dx,
		Dy:        &dy,
		Steps:     cal.Steps,
		DPerStep:  &dPerStep,
	}

	ok, err := c.repo.SQLite.UpdateCalibration(cal, cal.SessionId)
	if err != nil {
		return models.Calibration{}, err
	}
	if !ok {
		return models.Calibration{}, errors.New("failed to update calibration")
	}

	return cal, nil
}

func (c *Calibration) Clear(sessionId string) error {
	if sessionId == "" {
		return errors.New("session id is empty")
	}

	delete(c.calibrate, sessionId)

	return nil
}

func (c *Calibration) Save(sessionId string) error {
	calibrate, err := c.repo.SQLite.GetCalibration(sessionId)
	if err != nil {
		return err
	}

	return c.repo.CalRepo.TxUpsert(*calibrate.DPerStep)
}

func (c *Calibration) BytesFromPhoto(path string) ([]byte, error) {
	buf, err := c.camera.GetBytesFromPhoto(path)
	if err != nil {
		return nil, err
	}

	if buf == nil || buf.Len() == 0 {
		return nil, errors.New("photo is nil")
	}

	return buf.Bytes(), nil
}

// func (s *Service) Stream() {
// 	if err := c.camera.Run(); err != nil {
// 		logrus.Printf("%s", err)
// 	}
// }
