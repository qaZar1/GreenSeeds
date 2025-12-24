package models

type Calibration struct {
	SessionId       string   `json:"session_id" validate:"required" db:"session_id"`
	FirstPhotoPath  *string  `json:"first_photo" db:"first_photo_path"`
	SecondPhotoPath *string  `json:"second_photo" db:"second_photo_path"`
	Dx              *float64 `json:"dx" db:"dx"`
	Dy              *float64 `json:"dy" db:"dy"`
	Cir             *float64 `json:"cir" db:"cir"`
	DPerStep        *float64 `json:"d_per_step" db:"d_per_step"`
	CreatedAt       *string  `json:"created_at" db:"created_at"`
} //@name calibration

type Calculation struct {
	SessionId     string `json:"session_id"`
	NumberOfPhoto int64  `json:"number_of_photo"`
}

type GetPhoto struct {
	SessionId string `json:"session_id"`
	Photo     []byte `json:"photo"`
} //@name get-photo
