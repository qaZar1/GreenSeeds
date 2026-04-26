package models

import "time"

type Calibration struct {
	SessionId       string     `db:"session_id" validate:"required"`
	FirstPhotoPath  *string    `json:"first_photo" db:"first_photo_path"`
	SecondPhotoPath *string    `json:"second_photo" db:"second_photo_path"`
	Dx              *float64   `json:"dx" db:"dx"`
	Dy              *float64   `json:"dy" db:"dy"`
	Steps           *float64   `json:"steps" db:"steps"`
	DPerStep        *float64   `json:"d_per_step" db:"d_per_step"`
	CreatedAt       *time.Time `json:"created_at" db:"created_at"`
} //@name calibration
