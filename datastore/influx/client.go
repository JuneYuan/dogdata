package influx

import (
	"context"
	"sync"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

var (
	bucket = "gospel"
	org    = "DogData"
	once   sync.Once
)

type WrapClient struct {
	client influxdb2.Client
	org    string
	bucket string
}

var wrapClient WrapClient

func NewWrapClient(url string, token string) *WrapClient {
	once.Do(func() {
		wrapClient = WrapClient{
			client: influxdb2.NewClient(url, token),
			org:    org,
			bucket: bucket,
		}
	})
	return &wrapClient
}

func (c *WrapClient) WritePoint(p *write.Point) error {
	writeAPI := c.client.WriteAPIBlocking(c.org, c.bucket)
	err := writeAPI.WritePoint(context.Background(), p)
	c.client.Close()
	return err
}
