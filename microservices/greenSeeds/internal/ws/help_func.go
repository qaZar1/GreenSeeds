package ws

import (
	"fmt"
	"strings"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/device"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/rs/zerolog"
)

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

func AddError(log *zerolog.Logger, iter *models.Iteration, err error, stage string) {
	if err != nil {
		iter.Err = err
		iter.ErrStage = stage
	}
	log.Error().Err(err).Msg(fmt.Sprintf("[%s] %s", stage, err.Error()))

	iter.Success = false
	iter.Finished = true
}
