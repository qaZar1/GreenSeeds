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

func (app *App) CalibrationHandshake() (string, error) {
	// if ok := app.ws.Serial.InitialHandshake(); !ok {
	// 	return "", errors.New("Failed to initial handshake")
	// }

	sessionId := uuid.NewString()
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

	// msg := app.ws.Serial.RunGcode(toStart.Value)
	// if msg.Error != nil {
	// 	return "", errors.New(*msg.Error)
	// }

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

	buf, err := app.camera.TakePhoto()
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
		path = fmt.Sprintf("./tmp/%s/%app.jpg", sessionId, "1")
		cal.FirstPhotoPath = &path

		if err := app.camera.SavePhoto(path, sessionId, buf); err != nil {
			return nil, err
		}

		// toEnd, err := app.repo.DevSet.GetSettingsByKey("toEnd")
		// if err != nil {
		// 	return nil, err
		// }

		// if toEnd == (models.DeviceSettings{}) {
		// 	return nil, errors.New("toEnd is empty")
		// }

		// msg := app.ws.Serial.RunGcode(toEnd.Value)
		// if msg.Error != nil {
		// 	return nil, errors.New(*msg.Error)
		// }

	case 2:
		if cal.FirstPhotoPath == nil {
			return nil, errors.New("first photo is nil")
		}

		path = fmt.Sprintf("./tmp/%s/%app.jpg", sessionId, "2")
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

	dx, dy, err := app.calib.Finder(firstPhoto, secondPhoto)
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

	delete(app.calibration, sessionId)

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
