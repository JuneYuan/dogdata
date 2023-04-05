package influx

import (
	"context"
	"fmt"
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

func (c *WrapClient) Query() ([]float64, error) {
	queryAPI := c.client.QueryAPI(c.org)
	influxRet, err := queryAPI.Query(context.Background(), `from(bucket:"noaa")|> range(start: -480h) |> filter(fn: (r) => r._measurement == "average_temperature")`)
	if err != nil {
		return nil, fmt.Errorf("queryAPI.Query(): %v", err)
	}

	var ret []float64
	for influxRet.Next() {
		// TODO 这段逻辑，留待后续处理
		// if result.TableChanged() {

		// }
		ret = append(ret, influxRet.Record().Value().(float64))
	}
	if influxRet.Err() != nil {
		return nil, fmt.Errorf("influxRet.Err(): %v", influxRet.Err().Error())
	}
	return ret, nil
}
