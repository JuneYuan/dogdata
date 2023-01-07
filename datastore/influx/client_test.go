package influx

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

var (
	token = "TSmwhXA3hHFZy8i_xAevEDH0iNKfVm4YEw9Wu_mmrCjOcwGUIgNV4rX5VpXA5FLh8nN9zzUSKdyFpLbyOQOBuA=="
	// Store the URL of your InfluxDB instance
	url = "http://localhost:8086"
)

func TestWrapClient_WritePoint(t *testing.T) {
	tests := []struct {
		name string
		p    *write.Point
	}{
		{
			name: "trivial test",
			p: influxdb2.NewPoint("stat",
				map[string]string{"unit": "temperature"},
				map[string]interface{}{"avg": 24.5, "max": 45},
				time.Now()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewWrapClient(url, token).WritePoint(tt.p)
			assert.NoError(t, err)
		})
	}
}
