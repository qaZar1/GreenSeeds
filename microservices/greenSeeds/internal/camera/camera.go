package camera

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type Camera struct {
	cameraName  string
	inputDevice string
	framerate   string
	videoSize   string
}

func NewCamera(cameraName string, inputDevice string, framerate string, videoSize string) *Camera {
	return &Camera{
		cameraName:  cameraName,
		inputDevice: inputDevice,
		framerate:   framerate,
		videoSize:   videoSize,
	}
}

func (cam *Camera) TakePhoto() (*bytes.Buffer, error) {
	log.Info().Msg("Starting FFmpeg take photo")

	cmd := exec.Command(
		"ffmpeg",
		"-f", cam.inputDevice,
		"-framerate", cam.framerate,
		"-video_size", cam.videoSize,
		"-i", cam.cameraName,
		"-frames:v", "1",
		"-f", "image2",
		"-update", "1",
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error().Err(err).Msg("Cannot create stdout pipe")
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		log.Error().Msgf("Failed to run ffmpeg: %v", err)
		return nil, err
	}

	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(stdout)
	if err != nil {
		log.Error().Err(err).Msg("Cannot read photo")
		return nil, err
	}

	if n == 0 {
		log.Error().Msg("Photo is empty")
		return nil, err
	}

	if err = cmd.Wait(); err != nil {
		log.Error().Msgf("Failed to run ffmpeg: %v", err)
		return nil, err
	}

	return buf, nil
}

func (cam *Camera) SavePhoto(path, id string, buf *bytes.Buffer) error {
	if err := os.MkdirAll(fmt.Sprintf("./tmp/%s/", id), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), 0644)
}

func (cam *Camera) GetBytesFromPhoto(path string) (*bytes.Buffer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	return buf, nil
}

func (cam *Camera) DeletePhoto(id, name string) error {
	path := fmt.Sprintf("./tmp/%s/%s.jpg", id, name)

	return os.Remove(path)
}
