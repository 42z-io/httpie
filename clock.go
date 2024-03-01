package httpie

import "time"

type IClockService interface {
	Now() time.Time
}

type ClockService struct{}

func (c *ClockService) Now() time.Time {
	return time.Now().UTC()
}
