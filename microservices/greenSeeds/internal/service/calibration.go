package service

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (s *Service) CalibrationHandshake() (models.Calibration, error) {
	if ok := s.ws.Serial.InitialHandshake(); !ok {
		return models.Calibration{}, errors.New("Failed to initial handshake")
	}

	sessionId := uuid.NewString()
	handshake := models.Calibration{
		SessionId: sessionId,
	}

	now := time.Now().String()
	// добавляем в базу
	s.repo.SQLite.AddCalibration(models.Calibration{
		SessionId: sessionId,
		CreatedAt: &now,
	})

	toStart, err := s.repo.DevSet.GetSettingsByKey("toStart")
	if err != nil {
		return models.Calibration{}, err
	}

	if toStart == (models.DeviceSettings{}) {
		return models.Calibration{}, errors.New("toStart is empty")
	}

	msg := s.ws.Serial.RunGcode(toStart.Value)
	if msg.Error != nil {
		return models.Calibration{}, errors.New(*msg.Error)
	}
	return handshake, nil
}

func (s *Service) GetPhoto(calculation models.Calculation) (models.GetPhoto, error) {
	if err := s.validate.Struct(calculation); err != nil {
		return models.GetPhoto{}, err
	}

	cal, err := s.repo.SQLite.GetCalibration(calculation.SessionId)
	if err != nil {
		return models.GetPhoto{}, err
	}

	getPhoto := models.GetPhoto{
		SessionId: calculation.SessionId,
	}

	buf, err := s.camera.TakePhoto()
	if err != nil {
		return models.GetPhoto{}, err
	}

	if buf == nil || buf.Len() == 0 {
		return models.GetPhoto{}, errors.New("photo is nil")
	}

	photo := buf.Bytes()

	var path string
	switch calculation.NumberOfPhoto {
	case 1:
		path = fmt.Sprintf("./tmp/%s/%s.jpg", calculation.SessionId, "1")
		cal.FirstPhotoPath = &path

		if err := s.camera.SavePhoto(path, calculation.SessionId, buf); err != nil {
			return models.GetPhoto{}, err
		}

		toEnd, err := s.repo.DevSet.GetSettingsByKey("toEnd")
		if err != nil {
			return models.GetPhoto{}, err
		}

		if toEnd == (models.DeviceSettings{}) {
			return models.GetPhoto{}, errors.New("toEnd is empty")
		}

		msg := s.ws.Serial.RunGcode(toEnd.Value)
		if msg.Error != nil {
			return models.GetPhoto{}, errors.New(*msg.Error)
		}

	case 2:
		if cal.FirstPhotoPath == nil {
			return models.GetPhoto{}, errors.New("first photo is nil")
		}

		path = fmt.Sprintf("./tmp/%s/%s.jpg", calculation.SessionId, "2")
		cal.SecondPhotoPath = &path

		if err := s.camera.SavePhoto(path, calculation.SessionId, buf); err != nil {
			return models.GetPhoto{}, err
		}

	default:
		return models.GetPhoto{}, errors.New("number of photo is invalid")
	}

	ok, err := s.repo.SQLite.UpdateCalibration(cal)
	if err != nil {
		return models.GetPhoto{}, err
	}
	if !ok {
		return models.GetPhoto{}, errors.New("failed to update calibration")
	}

	getPhoto.Photo = photo

	return getPhoto, nil
}

func (s *Service) CalculateResult(calibration models.Calibration) (models.Calibration, error) {
	if err := s.validate.Struct(calibration); err != nil {
		return models.Calibration{}, err
	}

	if calibration.SessionId == "" {
		return models.Calibration{}, errors.New("session id is empty")
	}

	cal, err := s.repo.SQLite.GetCalibration(calibration.SessionId)
	if err != nil {
		return models.Calibration{}, err
	}

	firstPhoto, err := s.BytesFromPhoto(*cal.FirstPhotoPath)
	if err != nil {
		return models.Calibration{}, err
	}

	secondPhoto, err := s.BytesFromPhoto(*cal.SecondPhotoPath)
	if err != nil {
		return models.Calibration{}, err
	}

	dx, dy, err := s.calib.Finder(firstPhoto, secondPhoto)
	if err != nil {
		return models.Calibration{}, err
	}
	if dx == 0 && dy == 0 {
		return models.Calibration{}, errors.New("Distance not changed")
	}

	max := math.Max(math.Abs(dx), math.Abs(dy))
	dPerStep := max / *calibration.Cir

	cal = models.Calibration{
		SessionId: calibration.SessionId,
		Dx:        &dx,
		Dy:        &dy,
		DPerStep:  &dPerStep,
	}

	ok, err := s.repo.SQLite.UpdateCalibration(cal)
	if err != nil {
		return models.Calibration{}, err
	}
	if !ok {
		return models.Calibration{}, errors.New("failed to update calibration")
	}

	return cal, nil
}

func (s *Service) Clear(sessionId string) error {
	if sessionId == "" {
		return errors.New("session id is empty")
	}

	delete(s.calibration, sessionId)

	return nil
}

func (s *Service) Save(calibration models.Calibration) error {
	if err := s.validate.Struct(calibration); err != nil {
		return err
	}

	step := int(math.Round(*calibration.DPerStep * 10))
	stepStr := strconv.Itoa(step)

	return s.repo.CalRepo.TxUpsert(stepStr)
}

func (s *Service) BytesFromPhoto(path string) ([]byte, error) {
	buf, err := s.camera.GetBytesFromPhoto(path)
	if err != nil {
		return nil, err
	}

	if buf == nil || buf.Len() == 0 {
		return nil, errors.New("photo is nil")
	}

	return buf.Bytes(), nil
}
