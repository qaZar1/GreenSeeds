package camera

import (
	"bytes"
	"fmt"
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
		return nil, fmt.Errorf("photo is empty")
	}

	if err = cmd.Wait(); err != nil {
		log.Error().Msgf("Failed to run ffmpeg: %v", err)
		return nil, err
	}

	return buf, nil
}
