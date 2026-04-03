package circuitbreaker

import (
	"sync"
	"time"
)

type State string

const (
	StateClosed   State = "closed"    // funcionando normal
	StateOpen     State = "open"      // bloqueando requests
	StateHalfOpen State = "half_open" // probando si el servicio volvió
)

type CircuitBreaker struct {
	mu              sync.Mutex
	state           State
	failures        int
	maxFailures     int
	timeout         time.Duration
	lastFailureTime time.Time
	ServiceName     string
}

func New(serviceName string, maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:       StateClosed,
		maxFailures: maxFailures,
		timeout:     timeout,
		ServiceName: serviceName,
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

func (cb *CircuitBreaker) Success() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.state = StateClosed
}

func (cb *CircuitBreaker) Failure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	cb.lastFailureTime = time.Now()

	if cb.failures >= cb.maxFailures {
		cb.state = StateOpen
	}
}

func (cb *CircuitBreaker) GetState() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}
