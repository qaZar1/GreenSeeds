package ws

import (
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/device"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func waitForDeviceReady(
	manager *device.Manager,
	c *Client,
	timeout time.Duration,
	iter *models.Iteration,
) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return ErrTimeoutWaitingReady

		case <-c.done:
			return ErrClientDisconnected

		case <-c.planting.StopChan:
			if iter != nil {
				iter.Finished = true
			}
			return ErrStopped

		case <-ticker.C:
			if manager.GetStatus() != device.ManagerStateConnected {
				return ErrDeviceDisconnected
			}

			if manager.GetState() == device.StateReady {
				return nil
			}
		}
	}
}

func isStopped(c *Client) bool {
	select {
	case <-c.planting.StopChan:
		return true
	default:
		return false
	}
}
