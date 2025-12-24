package service

import (
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (s *Service) AddSetting(setting models.DeviceSettings) (models.DeviceSettings, error) {
	if err := s.validate.Struct(setting); err != nil {
		return models.DeviceSettings{}, err
	}

	return s.repo.DevSet.AddSetting(setting)
}

func (s *Service) GetSettings() ([]models.DeviceSettings, error) {
	return s.repo.DevSet.GetSettings()
}

func (s *Service) GetSettingsByKey(key string) (models.DeviceSettings, error) {
	return s.repo.DevSet.GetSettingsByKey(key)
}

func (s *Service) UpdateSetting(setting models.DeviceSettings) (models.DeviceSettings, error) {
	if err := s.validate.Struct(setting); err != nil {
		return models.DeviceSettings{}, err
	}

	return s.repo.DevSet.UpdateSettings(setting)
}

func (s *Service) DeleteSetting(key string) (bool, error) {
	return s.repo.DevSet.DeleteSettings(key)
}
