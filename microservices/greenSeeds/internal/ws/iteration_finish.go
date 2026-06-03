package ws

import (
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func finishIteration(
	s *Server,
	c *Client,
	iter *models.Iteration,
	success bool,
	errStr string,
	solution string,
	mark string,
) {
	now := time.Now()

	report := buildReport(
		c,
		now,
		success,
		errStr,
		solution,
		mark,
		iter,
	)

	s.repo.RepRepo.UpdateReports(report)
}

func buildReport(
	c *Client,
	dt time.Time,
	success bool,
	err string,
	solution string,
	mark string,
	iter *models.Iteration,
) models.Reports {
	return models.Reports{
		Shift:    int64(iter.Shift),
		Number:   iter.Number,
		Recipe:   int64(iter.Recipe),
		Turn:     iter.Turn,
		Dt:       &dt,
		Success:  success,
		Error:    &err,
		Solution: &solution,
		Mark:     &mark,
	}
}
