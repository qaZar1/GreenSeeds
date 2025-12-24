package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IDeviceSettingsRepository interface {
	AddSetting(settings models.DeviceSettings) (models.DeviceSettings, error)
	GetSettings() ([]models.DeviceSettings, error)
	GetSettingsByKey(key string) (models.DeviceSettings, error)
	UpdateSettings(settings models.DeviceSettings) (models.DeviceSettings, error)
	DeleteSettings(key string) (bool, error)
}

type deviceSettingsRepository struct {
	db *sqlx.DB
}

func NewDeviceSettingsRepository(db *sqlx.DB) *deviceSettingsRepository {
	return &deviceSettingsRepository{
		db: db,
	}
}

func (devSet *deviceSettingsRepository) AddSetting(settings models.DeviceSettings) (models.DeviceSettings, error) {
	const query = `
INSERT INTO green_seeds.device_settings (
	key,
	value
)
VALUES (
	:key,
	:value
)
RETURNING key, value`

	rows, err := devSet.db.NamedQuery(query, settings)
	if err != nil {
		return models.DeviceSettings{}, err
	}

	defer rows.Close()

	var inserted models.DeviceSettings
	if rows.Next() {
		err = rows.StructScan(&inserted)
		if err != nil {
			return models.DeviceSettings{}, err
		}
	}

	return inserted, nil
}

func (devSet *deviceSettingsRepository) GetSettings() ([]models.DeviceSettings, error) {
	const query = `
SELECT key, value
FROM green_seeds.device_settings
ORDER BY key ASC`

	var settings []models.DeviceSettings
	if err := devSet.db.Select(&settings, query); err != nil {
		return nil, err
	}

	return settings, nil
}

func (devSet *deviceSettingsRepository) GetSettingsByKey(key string) (models.DeviceSettings, error) {
	const query = `
SELECT key, value
FROM green_seeds.device_settings
WHERE key = $1`

	var setting models.DeviceSettings
	if err := devSet.db.Get(&setting, query, key); err != nil {
		return models.DeviceSettings{}, err
	}

	return setting, nil
}

func (devSet *deviceSettingsRepository) UpdateSettings(settings models.DeviceSettings) (models.DeviceSettings, error) {
	const query = `
UPDATE green_seeds.device_settings
SET value = :value
WHERE key = :key
RETURNING key, value`

	rows, err := devSet.db.NamedQuery(query, settings)
	if err != nil {
		return models.DeviceSettings{}, err
	}

	defer rows.Close()

	var updated models.DeviceSettings
	if rows.Next() {
		err = rows.StructScan(&updated)
		if err != nil {
			return models.DeviceSettings{}, err
		}
	}

	return updated, nil
}

func (devSet *deviceSettingsRepository) DeleteSettings(key string) (bool, error) {
	const query = `
DELETE FROM green_seeds.device_settings
WHERE key = $1`

	result, err := devSet.db.Exec(query, key)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}
