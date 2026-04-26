package device

import (
	"context"
	"sync"
)

type MatchResult struct {
	Matched bool // сообщение подходит запросу
	Done    bool // запрос завершён этим сообщением
}

type Request struct {
	Match func([]byte) MatchResult

	Ch   chan []byte // канал с входящими сообщениями
	Done chan error
}

type Dispatcher struct {
	mu       sync.Mutex
	requests []*Request
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		requests: make([]*Request, 0),
	}
}

func (d *Dispatcher) Add(r *Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.requests = append(d.requests, r)
}

func (d *Dispatcher) FailAll(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, r := range d.requests {
		r.Done <- err
		closeSafe(r.Done)
	}

	d.requests = nil
}

func (d *Dispatcher) Handle(data []byte) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var remaining []*Request

	for _, r := range d.requests {
		res := r.Match(data)

		if res.Matched {
			select {
			case r.Ch <- data:
			default:
				// дроп, если канал переполнен
			}
		}

		if res.Done {
			close(r.Ch)
			// сигналим завершение
			r.Done <- nil
			closeSafe(r.Done)
			continue
		}

		remaining = append(remaining, r)
	}

	d.requests = remaining
}

// Do — helper для выполнения запроса с ожиданием завершения или таймаута
func (d *Dispatcher) Do(ctx context.Context, r *Request) error {
	d.Add(r)

	select {
	case err := <-r.Done:
		return err

	case <-ctx.Done():
		d.remove(r)
		return ctx.Err()
	}
}

func (d *Dispatcher) remove(target *Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var filtered []*Request
	for _, r := range d.requests {
		if r != target {
			filtered = append(filtered, r)
		}
	}

	d.requests = filtered

	closeSafe(target.Done)
}

func (d *Dispatcher) HasActive() bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	return len(d.requests) > 0
}

func closeSafe(ch chan error) {
	select {
	case <-ch:
	default:
		close(ch)
	}
}
