package application

import (
	"bytes"
	"errors"
	"testing"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

//
// 🔹 FAKE CAMERA (чтобы не тянуть реальные зависимости)
//

type fakeCamera struct {
	takePhotoBuf *bytes.Buffer
	takePhotoErr error

	getBytesBuf *bytes.Buffer
	getBytesErr error

	saveErr error
	deleteErr error
	runErr error
}

func (f *fakeCamera) TakePhoto() (*bytes.Buffer, error) {
	if f.takePhotoBuf == nil {
		return bytes.NewBuffer([]byte("fake-image")), f.takePhotoErr
	}
	return f.takePhotoBuf, f.takePhotoErr
}

func (f *fakeCamera) SavePhoto(path, id string, buf *bytes.Buffer) error {
	return f.saveErr
}

func (f *fakeCamera) GetBytesFromPhoto(path string) (*bytes.Buffer, error) {
	if f.getBytesBuf == nil {
		return bytes.NewBuffer([]byte("fake-bytes")), f.getBytesErr
	}
	return f.getBytesBuf, f.getBytesErr
}

func (f *fakeCamera) DeletePhoto(id, name string) error {
	return f.deleteErr
}

func (f *fakeCamera) Run() error {
	return f.runErr
}


type fakeSQLite struct {
	getCalibrationFunc func(string) (models.Calibration, error)
}

func (f *fakeSQLite) GetCalibration(id string) (models.Calibration, error) {
	return f.getCalibrationFunc(id)
}

// заглушки чтобы не ругался интерфейс
func (f *fakeSQLite) AddCalibration(string, time.Time) (bool, error) {
	return false, nil
}
func (f *fakeSQLite) UpdateCalibration(models.Calibration, string) (bool, error) {
	return false, nil
}

type fakeOpenCV struct {
	finderFunc func([]byte, []byte) (float64, float64, error)
}

func (f *fakeOpenCV) Finder(a, b []byte) (float64, float64, error) {
	if f.finderFunc != nil {
		return f.finderFunc(a, b)
	}
	return 0, 0, nil
}
//
// 🔹 TESTS
//

func TestClear(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		app := Calibration{
			calibrate: map[string]models.Calibration{
				"session-1": {},
			},
		}

		err := app.Clear("session-1")

		assert.NoError(t, err)

		_, ok := app.calibrate["session-1"]
		assert.False(t, ok)
	})

	t.Run("EmptySession", func(t *testing.T) {
		app := &Calibration{}

		err := app.Clear("")

		assert.Error(t, err)
	})
}

func TestBytesFromPhoto(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		app := &Calibration{
			camera: &fakeCamera{
				getBytesBuf: bytes.NewBuffer([]byte("image")),
			},
		}

		data, err := app.BytesFromPhoto("path")

		assert.NoError(t, err)
		assert.Equal(t, []byte("image"), data)
	})

	t.Run("CameraError", func(t *testing.T) {
		app := &Calibration{
			camera: &fakeCamera{
				getBytesErr: errors.New("fail"),
			},
		}

		_, err := app.BytesFromPhoto("path")

		assert.Error(t, err)
	})

	t.Run("EmptyBuffer", func(t *testing.T) {
		app := &Calibration{
			camera: &fakeCamera{
				getBytesBuf: bytes.NewBuffer([]byte{}),
			},
		}

		_, err := app.BytesFromPhoto("path")

		assert.Error(t, err)
	})
}

func TestCalculateResult_ValidationFail(t *testing.T) {
	app := &Calibration{
		validate: validator.New(),
	}

	_, err := app.CalculateResult(models.Calibration{})

	assert.Error(t, err)
}

// func TestCalculateResult_PhotosNotReady(t *testing.T) {
// 	sqlite := &fakeSQLite{
// 		getCalibrationFunc: func(id string) (models.Calibration, error) {
// 			return models.Calibration{
// 				SessionId:      id,
// 				FirstPhotoPath: nil,
// 				SecondPhotoPath: nil,
// 			}, nil
// 		},
// 	}

// 	app := &Calibration{
// 		repo: &repository.Repository{
// 			SQLite: sqlite,
// 		},
// 		camera:    &fakeCamera{},
// 		opencv:    &fakeOpenCV{},
// 		validate:  validator.New(),
// 		calibrate: map[string]models.Calibration{},
// 	}

// 	steps := float64(10)

// 	input := models.Calibration{
// 		SessionId: "session-1",
// 		Steps:     &steps,
// 	}

// 	_, err := app.CalculateResult(input)

// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "Photos not ready")
// }

// func TestSave_RepoNil(t *testing.T) {
// 	app := &Calibration{
// 		repo: nil,
// 	}

// 	err := app.Save("session-1")

// 	assert.Error(t, err)
// }