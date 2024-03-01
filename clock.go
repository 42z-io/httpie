package httpie

import "time"

// Clock service interface
type IClockService interface {
	Now() time.Time
}

// ClockService is a service for getting the current time
type ClockService struct{}

// Now returns the current time
func (c *ClockService) Now() time.Time {
	return time.Now().UTC()
}
