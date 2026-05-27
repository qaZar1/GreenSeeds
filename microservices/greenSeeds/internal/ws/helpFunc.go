package ws

import (
	"errors"
	"strings"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/device"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func failIteration(iter *models.Iteration, err error) {
	if err != nil {
		iter.Err = append(iter.Err, err)
	}

	iter.Success = false
	iter.Finished = true
}

func isDeviceAlive(m *device.Manager) bool {
	return m.GetStatus() == device.ManagerStateConnected
}

func ErrorsToString(errs []error) string {
	builder := strings.Builder{}

	for _, err := range errs {
		builder.WriteString(err.Error() + ";\n")
	}

	return builder.String()
}

func checkStop(c *Client) bool {
	select {
	case <-c.Control:
		c.planting.Stop = true
		c.Send <- okResponse("STOP", "Stopped by user")
		return true

	default:
		return false
	}
}

func Err(iter *models.Iteration, err error) {
	failIteration(iter, err)
}

func emptyPhotoError() error {
	return errors.New("empty photo buffer")
}
