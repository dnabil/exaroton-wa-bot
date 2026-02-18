package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPagination(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		limit       int
		total       int
		expectPage  int
		expectTotal int
		expectStart int
		expectEnd   int
	}{
		{
			name:        "normal case",
			page:        2,
			limit:       10,
			total:       50,
			expectPage:  2,
			expectTotal: 5,
			expectStart: 10,
			expectEnd:   20,
		},
		{
			name:        "page exceeds total pages",
			page:        10,
			limit:       10,
			total:       25,
			expectPage:  3,
			expectTotal: 3,
			expectStart: 20,
			expectEnd:   25,
		},
		{
			name:        "limit <= 0 uses total",
			page:        1,
			limit:       0,
			total:       20,
			expectPage:  1,
			expectTotal: 1,
			expectStart: 0,
			expectEnd:   20,
		},
		{
			name:        "page <= 0 becomes page 1",
			page:        0,
			limit:       10,
			total:       20,
			expectPage:  1,
			expectTotal: 2,
			expectStart: 0,
			expectEnd:   10,
		},
		{
			name:        "total = 0",
			page:        1,
			limit:       10,
			total:       0,
			expectPage:  1,
			expectTotal: 0,
			expectStart: 0,
			expectEnd:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPagination(tt.page, tt.limit, tt.total)

			// assert
			assert.Equal(t, tt.expectPage, p.CurrentPage)
			assert.Equal(t, tt.expectTotal, p.TotalPage)
			assert.Equal(t, tt.expectStart, p.Start())
			assert.Equal(t, tt.expectEnd, p.End())
		})
	}
}
