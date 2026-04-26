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
		"-f", "image2",
		"-frames:v", "1",
		"-q:v", "2",
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

func (cam *Camera) Run() error {
	log.Info().Msg("Starting FFmpeg stream via ffmpeg-go")

	stderr := &bytes.Buffer{}
	cmd := exec.Command(
		"ffmpeg",
		"-f", cam.inputDevice,
		"-framerate", cam.framerate,
		"-video_size", cam.videoSize,
		"-i", cam.cameraName,
		"-vcodec", "libx264",
		"-preset", "ultrafast",
		"-tune", "zerolatency",
		"-f", "mpegts",
		"pipe:1",
	)

	reader, err := cmd.StdoutPipe()

	log.Info().Msgf("FFmpeg stream info: %s", reader)

	// Закрываем writer по завершении контекста
	// go func() {
	// 	<-ctx.Done()
	// 	log.Info().Msg("Context canceled, closing ffmpeg writer")
	// 	app.writer.Close()
	// }()

	// Подключаем stderr для логов ffmpeg

	log.Info().Msgf("FFmpeg command: %s", cmd.String())
	log.Info().Msgf("FFmpeg stderr: %s", stderr.String())

	err = cmd.Run()

	// Логируем stderr
	if stderr.Len() > 0 {
		log.Error().Msgf("FFmpeg stderr: %s", stderr.String())
	}

	if err != nil {
		log.Error().Msgf("Failed to run ffmpeg-go: %v", err)
		return err
	}
	log.Info().Msg("FFmpeg finished")

	return nil
}

// // Рассылает видеопоток всем подписчикам
// func DistributeStream(log zerolog.Logger, ctx context.Context) error {
// 	buf := make([]byte, 1024)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			log.Error().Msg("DistributeStream stopped")
// 			return ctx.Err()
// 		default:
// 			n, err := app.reader.Read(buf)
// 			if err != nil {
// 				log.Error().Msgf("Pipe read error: %s", err)
// 				return err
// 			}
// 			chunk := make([]byte, n)
// 			copy(chunk, buf[:n])

// 			app.mu.RLock()
// 			for _, ch := range app.subscribers {
// 				select {
// 				case ch <- chunk:
// 				default:
// 					log.Info().Msg("Dropping frame for slow client")
// 				}
// 			}
// 			app.mu.RUnlock()
// 		}
// 	}
// }

// // Отдаёт поток подключённому клиенту
// func StreamHandler(log zerolog.Logger, w http.ResponseWriter, r *http.Request) {
// 	log.Info().Msg("New client connected")
// 	w.Header().Set("Content-Type", "video/mp2t")

// 	clientChan := make(chan []byte, 200)

// 	app.mu.Lock()
// 	app.subscribers = append(app.subscribers, clientChan)
// 	log.Info().Msgf("New subscriber added: %d", len(app.subscribers))

// 	app.mu.Unlock()

// 	defer func() {
// 		app.mu.Lock()
// 		for i, ch := range app.subscribers {
// 			log.Info().Msgf("Removing client %d", i)
// 			if ch == clientChan {
// 				app.subscribers = append(app.subscribers[:i], app.subscribers[i+1:]...)
// 				break
// 			}
// 		}
// 		app.mu.Unlock()
// 		close(clientChan)
// 		log.Info().Msg("Client disconnected")
// 	}()

// 	// Отслеживание отключения клиента
// 	ctx := r.Context()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case chunk, ok := <-clientChan:
// 			if !ok {
// 				return
// 			}
// 			_, err := w.Write(chunk)
// 			if err != nil {
// 				log.Error().Msgf("Write error: %s", err)
// 				return
// 			}

// 			flusher := w.(http.Flusher)
// 			flusher.Flush()
// 		}
// 	}
// }
