package ws

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func Photo(s *Server, c *Client, iter *models.Iteration) error {
	// buf, err := s.camera.GetBytesFromPhoto("./Proj_img/Amarant/photo_2024-06-18_15-22-18.jpg") // todo change to cam
	// if err != nil {
	// 	return err
	// }

	buf, err := s.camera.TakePhoto()
	if err != nil {
		return err
	}

	if buf == nil || buf.Len() == 0 {
		return fmt.Errorf("photo is empty")
	}

	iter.LastBuf = buf

	return nil
}

func Control(s *Server, c *Client, iter *models.Iteration) (bool, error) {
	response, err := s.api.CheckAI(iter.Seed, *iter.LastBuf)
	if err != nil {
		return false, err
	}

	if response.PercentOfMatch < s.config.MinMatchPercent {
		return false, fmt.Errorf(
			"low match percentage %d", int(s.config.MinMatchPercent*100),
		)
	}

	if response.Seed != iter.Seed {
		return false, fmt.Errorf(
			"wrong seed type. Expected %s, but got %s", iter.Seed, response.Seed,
		)
	}

	count := s.opencv.Counter(iter.LastBuf.Bytes())
	if count == 0 {
		return false, fmt.Errorf("counting failed")
	}

	seed, err := s.repo.SeedRepo.GetSeedsBySeed(iter.Seed)
	if err != nil {
		return false, fmt.Errorf("failed to get seed: %w", err)
	}

	iter.Count = count
	iter.MinDensity = seed.MinDensity
	iter.MaxDensity = seed.MaxDensity

	if count > iter.MaxDensity {
		return false, fmt.Errorf("too many seeds: %d", count)
	}

	if count < iter.MinDensity {
		return false, nil
	}

	return true, nil
}

func EditGcode(iter *models.Iteration) {
	re := regexp.MustCompile(`OPEN_TIME=([0-9]+(?:\.[0-9]+)?)`)
	matches := re.FindStringSubmatch(iter.Gcode)
	if len(matches) < 2 {
		return
	}

	openTime, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return
	}

	newOpenTime := CalculateNewOpenTime(iter, openTime)

	newGCode := re.ReplaceAllString(
		iter.Gcode,
		fmt.Sprintf("OPEN_TIME=%.2f", newOpenTime),
	)
	iter.Gcode = newGCode
}

func CalculateNewOpenTime(iter *models.Iteration, openTime float64) float64 {
	avgDensity := (iter.MinDensity + iter.MaxDensity) / 2
	needSeeds := avgDensity - iter.Count
	return float64(needSeeds) * openTime / float64(avgDensity)
}
