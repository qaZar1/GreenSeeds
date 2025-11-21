package device

import (
	"fmt"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (m *SerialManager) buildGcode(msg models.WSMessage, sorrel string) string {
	const query = `
BEGIN %d/%d/%d
BUNKER %d
\x02%s\x03%s`

	data := fmt.Sprintf(query,
		msg.Params.Shift,
		msg.Params.Number,
		msg.Params.Turn,
		msg.Params.Bunker,
		msg.Params.Gcode,
		sorrel)
	return data
}

func (m *SerialManager) buildSorrel(msg models.WSMessage) string {
	const query = `
\x07Sorrel\x0A%d/%d\x0A%d/%d\x0A\x0D`

	data := fmt.Sprintf(query,
		msg.Params.Shift,
		msg.Params.Number,
		msg.Params.Turn,
		msg.Params.RequiredAmount)
	return data
}
