package httpie

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type ClockServiceMock struct {
	mock.Mock
}

func (c *ClockServiceMock) Now() time.Time {
	ret := c.Called()
	r0 := ret.Get(0).(time.Time)
	return r0
}
