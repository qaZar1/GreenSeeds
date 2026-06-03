package application

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

//go:generate mockgen -source=calibration.go -destination=./../mocks/mock_calibration.go -package=mocks
type ICalibrationApp interface {
	CalibrationHandshake() (string, error)
	GetPhoto(sessionId string, numberPhotoStr string) ([]byte, error)
	CalculateResult(calibration models.Calibration) (models.Calibration, error)
	Clear(sessionId string) error
	Save(sessionId string) error
}

func (app *App) CalibrationHandshake() (string, error) {
	sessionId := uuid.NewString()

	if !app.device.Manager.TryAcquireSession(sessionId) {
		return "", errors.New("device busy")
	}
	defer app.device.Manager.ReleaseSession(sessionId)

	// 3. Останавливаем фоновый пинг на время сессии
	app.device.PausePolling()

	if err := app.device.Boot(sessionId, false); err != nil {
		return "", err
	}

	now := time.Now()

	// добавляем в базу
	ok, err := app.repo.SQLite.AddCalibration(sessionId, now)
	if err != nil || !ok {
		return "", err
	}

	toStart, err := app.repo.DevSet.GetSettingsByKey("toStart")
	if err != nil {
		return "", err
	}

	if toStart == (models.DeviceSettings{}) {
		return "", errors.New("toStart is empty")
	}

	if err := app.device.RunGcode(toStart.Value, sessionId); err != nil {
		return "", err
	}

	return sessionId, nil
}

func (app *App) GetPhoto(sessionId string, numberPhotoStr string) ([]byte, error) {
	numberPhoto, err := strconv.Atoi(numberPhotoStr)
	if err != nil {
		return nil, err
	}

	cal, err := app.repo.SQLite.GetCalibration(sessionId)
	if err != nil {
		return nil, err
	}

	photoPath := ""
	if numberPhoto == 1 {
		photoPath = "./1.jpg"
	} else {
		photoPath = "./2.jpg"
	}

	buf, err := app.camera.GetBytesFromPhoto(photoPath)
	if err != nil {
		return nil, err
	}

	// buf, err := app.camera.TakePhoto()
	// if err != nil {
	// 	return nil, err
	// }

	if buf == nil || buf.Len() == 0 {
		return nil, errors.New("photo is nil")
	}

	photo := buf.Bytes()

	var path string
	switch numberPhoto {
	case 1:
		path = fmt.Sprintf("./tmp/%s/%s.jpg", sessionId, "1")
		cal.FirstPhotoPath = &path

		if err := app.camera.SavePhoto(path, sessionId, buf); err != nil {
			return nil, err
		}

		toEnd, err := app.repo.DevSet.GetSettingsByKey("toEnd")
		if err != nil {
			return nil, err
		}

		if toEnd == (models.DeviceSettings{}) {
			return nil, errors.New("toEnd is empty")
		}

		err = app.device.RunGcode(toEnd.Value, sessionId)
		if err != nil {
			return nil, err
		}

	case 2:
		time.Sleep(2 * time.Second)
		if cal.FirstPhotoPath == nil {
			return nil, errors.New("first photo is nil")
		}

		path = fmt.Sprintf("./tmp/%s/%s.jpg", sessionId, "2")
		cal.SecondPhotoPath = &path

		if err := app.camera.SavePhoto(path, sessionId, buf); err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("number of photo is invalid")
	}

	ok, err := app.repo.SQLite.UpdateCalibration(cal, sessionId)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("failed to update calibration")
	}

	return photo, nil
}

func (app *App) CalculateResult(calibration models.Calibration) (models.Calibration, error) {
	if err := app.validate.Struct(calibration); err != nil {
		return models.Calibration{}, err
	}

	cal, err := app.repo.SQLite.GetCalibration(calibration.SessionId)
	if err != nil {
		return models.Calibration{}, err
	}

	firstPath := cal.FirstPhotoPath
	secondPath := cal.SecondPhotoPath
	if firstPath == nil || secondPath == nil {
		return models.Calibration{}, errors.New("Photos not ready")
	}

	firstPhoto, err := app.BytesFromPhoto(*firstPath)
	if err != nil {
		return models.Calibration{}, err
	}

	secondPhoto, err := app.BytesFromPhoto(*secondPath)
	if err != nil {
		return models.Calibration{}, err
	}

	dx, dy, err := app.opencv.Finder(firstPhoto, secondPhoto)
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

	ok, err := app.repo.SQLite.UpdateCalibration(cal, cal.SessionId)
	if err != nil {
		return models.Calibration{}, err
	}
	if !ok {
		return models.Calibration{}, errors.New("failed to update calibration")
	}

	return cal, nil
}

func (app *App) Clear(sessionId string) error {
	if sessionId == "" {
		return errors.New("session id is empty")
	}

	delete(app.calibrate, sessionId)

	return nil
}

func (app *App) Save(sessionId string) error {
	calibrate, err := app.repo.SQLite.GetCalibration(sessionId)
	if err != nil {
		return err
	}

	return app.repo.CalRepo.TxUpsert(*calibrate.DPerStep)
}

func (app *App) BytesFromPhoto(path string) ([]byte, error) {
	buf, err := app.camera.GetBytesFromPhoto(path)
	if err != nil {
		return nil, err
	}

	if buf == nil || buf.Len() == 0 {
		return nil, errors.New("photo is nil")
	}

	return buf.Bytes(), nil
}

// func (s *Service) Stream() {
// 	if err := app.camera.Run(); err != nil {
// 		logrus.Printf("%s", err)
// 	}
// }
