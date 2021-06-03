package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSessionExpire(t *testing.T) {
	now := time.Date(2021, time.May, 29, 10, 0, 0, 0, time.UTC)
	tBefore := time.Date(2021, time.May, 29, 9, 0, 0, 0, time.UTC)
	tAfter := time.Date(2021, time.May, 29, 11, 0, 0, 0, time.UTC)

	tests := []struct {
		name               string
		now, exIdle, exAbs time.Time
		isExpired          bool
	}{
		{"not expired", now, tAfter, tAfter, false},
		{"idle expire", now, tBefore, tAfter, true},
		{"absolute expire", now, tAfter, tBefore, true},
		{"both expire", now, tBefore, tBefore, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Session{ExpireIdle: test.exIdle, ExpireAbs: test.exAbs}
			assert.Equal(t, test.isExpired, s.Expired(test.now))
		})
	}
}
