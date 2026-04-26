package device

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/rs/zerolog"
)

type DeviceClient struct {
	Manager *Manager
	ctx     context.Context
	RespCh  chan models.WSResponse

	state DeviceState

	pollDur   time.Duration
	pollTimer *time.Timer
	pollMu    sync.Mutex
}

func NewClient(ctx context.Context, port string, baud int, log zerolog.Logger) *DeviceClient {
	manager := NewManager(ctx, port, baud, log)
	dur := 1 * time.Minute

	d := &DeviceClient{
		Manager: manager,
		ctx:     ctx,
		RespCh:  make(chan models.WSResponse, 100),

		state: StateStandby,

		pollDur:   dur,
		pollTimer: time.NewTimer(dur),
	}

	d.pollingDevice()
	d.ListenStatusCh()

	return d
}

func (c *DeviceClient) ListenStatusCh() {
	go func() {
		for {
			select {
			case st := <-c.Manager.StatusChangeCh:
				resp := models.WSResponse{
					Type: "DEVICE",
				}

				stStr := fmt.Sprintf("%v", st)
				if st == ManagerStateDisconnected {
					resp.Status = "ERROR"
					resp.Message = stStr
				} else {
					resp.Status = "OK"
					resp.Message = stStr
				}

				select {
				case c.RespCh <- resp:
				default:
				}

			case <-c.ctx.Done():
				return
			}
		}
	}()
}

func (c *DeviceClient) RefreshPolling() {
	c.pollMu.Lock()
	defer c.pollMu.Unlock()

	if !c.pollTimer.Stop() {
		select {
		case <-c.pollTimer.C:
		default:
		}
	}

	c.pollTimer.Reset(c.pollDur)
}

func (c *DeviceClient) PausePolling() {
	c.pollMu.Lock()
	defer c.pollMu.Unlock()

	if !c.pollTimer.Stop() {
		select {
		case <-c.pollTimer.C:
		default:
		}
	}
}

func (c *DeviceClient) pollingDevice() {
	go func() {
		for {
			select {
			case <-c.pollTimer.C:
				if c.Manager.dispatcher.HasActive() {
					c.RefreshPolling()
					continue
				}

				if c.Manager.GetStatus() != ManagerStateConnected {
					c.RefreshPolling()
					continue
				}

				if c.Manager.GetState() == StateBusy {
					c.RefreshPolling()
					continue
				}

				if err := c.Ping(InternalSessionID); err != nil {
					fmt.Println("ERROR!!") // todo
				}
				c.RefreshPolling()

			case <-c.Manager.ctx.Done():
				c.pollMu.Lock()
				defer c.pollMu.Unlock()

				c.pollTimer.Stop()
				return
			}
		}
	}()
}
