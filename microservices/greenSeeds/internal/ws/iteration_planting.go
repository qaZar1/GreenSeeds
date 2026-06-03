package ws

import (
	"fmt"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func Begin(s *Server, c *Client, iter *models.Iteration) error {
	bunkers, err := s.repo.SeedRepo.GetBestBunker(iter.Seed)
	if err != nil {
		err = fmt.Errorf("bunkers with seed %s are empty", iter.Seed)

		return err
	}

	iter.Bunker = bunkers.Bunker

	command := buildGcode(c, iter)

	if err := s.dClient.Begin(c.SessionId, command, iter.Turn); err != nil {
		return err
	}

	if err := s.repo.PlcRepo.DecrementSeed(iter.Bunker); err != nil {
		return fmt.Errorf("bunker is empty: %w", err)
	}

	return nil
}

func buildGcode(c *Client, iter *models.Iteration) string {
	const query = `
BEGIN %d/%d/%d
BUNKER %d
\x02%s\x03`

	data := fmt.Sprintf(query,
		iter.Shift,
		iter.Number,
		iter.Turn,
		iter.Bunker,
		iter.Gcode)

	if iter.ExtraMode {
		sorrel := buildSorrel(iter)
		data = fmt.Sprintf("%s%s", data, sorrel)
	}

	return data
}

func buildSorrel(iter *models.Iteration) string {
	const query = `
\x07Sorrel\x0A%d/%d\x0A%d/%d\x0A\x0D`

	data := fmt.Sprintf(query,
		iter.Shift,
		iter.Number,
		iter.Turn,
		iter.Required)
	return data
}
