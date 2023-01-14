package main

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func Test_parseStatsd(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantRet StatsdMsg
		wantErr bool
	}{
		{
			name: "trivial test",
			data: `custom_metric:80|g|#port:8125`,
			wantRet: StatsdMsg{
				MetricName: "custom_metric",
				Value:      80,
				Type:       GAUGE,
				Tags: []Tag{
					{"port", "8125"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := parseStatsd(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseStatsd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("parseStatsd() = \n%#v, \nwant \n%#v", gotRet, tt.wantRet)
			}
		})
	}
}

func Test_convert(t *testing.T) {
	tests := []struct {
		name string
		msg  StatsdMsg
		want *write.Point
	}{
		{
			name: "trivial test",
			msg: StatsdMsg{
				MetricName: "custom_metric",
				Value:      80,
				Type:       GAUGE,
				Tags: []Tag{
					{"port", "8125"},
				},
			},
			want: write.NewPoint(
				"custom_metric",
				map[string]string{"port": "8125"},
				map[string]interface{}{"value": float64(80)},
				time.Now(),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert(tt.msg)
			assert.Equal(t, tt.want.Name(), got.Name())
			assert.Equal(t, tt.want.TagList(), got.TagList())
			assert.Equal(t, len(tt.want.FieldList()), len(got.FieldList()))
			if len(tt.want.FieldList()) > 0 {
				assert.Equal(t, tt.want.FieldList(), got.FieldList())
			}
		})
	}
}
