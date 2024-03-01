package httpie

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClockService(t *testing.T) {
	clockService := ClockService{}
	result := clockService.Now()
	assert.NotNil(t, result)
}

func TestClockServiceMock(t *testing.T) {
	clockServiceMock := ClockServiceMock{}
	now := time.Now().UTC()
	clockServiceMock.On("Now").Return(now)
	result := clockServiceMock.Now()
	assert.NotNil(t, result)
	assert.Equal(t, now, result)
	clockServiceMock.AssertExpectations(t)
}
