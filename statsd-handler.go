package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"sideproject/dogdata/common"
	"sideproject/dogdata/datastore/influx"
)

func serveStatsd() {
	// TODO 真正工作的代码，不能只监听 localhost
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 7125})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Local: <%s> \n", listener.LocalAddr().String())
	data := make([]byte, 1024)
	// TODO 怎么处理多 client, 并发？
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		common.CheckErr(err, "listener.ReadFromUDP()")

		// fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
		// 打印内容改成了存储到 influx
		err = saveToInflux(string(data[:n]))
		if err != nil {
			// TODO 改为正经的日志
			fmt.Printf("saveToInflux(): %v\n", err)
		}

		// TODO Go network programming in action
		// 这部分还不确定有什么意义，是照着 go udp example 写的
		_, err = listener.WriteToUDP([]byte("world"), remoteAddr)
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
}

// TODO 应当依赖配置获取 url, token
var (
	token = "pShDwhGCHnWONE6m5jt07i3Sfj08iOXoEgQXZkuR4-kMlMny6udVi1--kXkFcDq7gRIyfBlXWGGcC88hX0DwnQ==" // 2023.6
	// token = "TSmwhXA3hHFZy8i_xAevEDH0iNKfVm4YEw9Wu_mmrCjOcwGUIgNV4rX5VpXA5FLh8nN9zzUSKdyFpLbyOQOBuA==" // 2023.1
	url = "http://localhost:8086"
)

func saveToInflux(data string) error {
	statsdMsg, err := parseStatsd(data)
	if err != nil {
		return err
	}
	p := convert(statsdMsg)
	return influx.NewWrapClient(url, token).WritePoint(p)
}

type StatsdMsgType string

const (
	COUNT        StatsdMsgType = "COUNT"
	GAUGE        StatsdMsgType = "GAUGE"
	TIMER        StatsdMsgType = "TIMER"
	HISTOGRAM    StatsdMsgType = "HISTOGRAM"
	SET          StatsdMsgType = "SET"
	DISTRIBUTION StatsdMsgType = "DISTRIBUTION"
	UNKNOWN      StatsdMsgType = "UNKNOWN"
)

var statsdTypes = map[string]StatsdMsgType{
	"c":  COUNT,
	"g":  GAUGE,
	"ms": TIMER,
	"h":  HISTOGRAM,
	"s":  SET,
	"d":  DISTRIBUTION,
}

func parseStatsd(data string) (ret StatsdMsg, err error) {
	// 按 `|` 分隔：
	// [0] - Required - <METRIC_NAME>:<VALUE>
	// [1] - Required - <TYPE>
	// [2] - Optioanl - @<SAMPLE_RATE>
	// [3] - Optioanl - #<TAG_KEY_1>:<TAG_VALUE_1>,<TAG_2>
	// 有内部结构的，再继续按分隔符解析。注意 `@` `#` 符号

	fmt.Printf("parseStatsd input: %q\n", data)

	elems := strings.Split(data, "|")
	for i, e := range elems {
		switch i {
		case 0:
			xs := strings.Split(e, ":")
			if len(xs) != 2 {
				return ret, fmt.Errorf("malformed data, cannot find `<METRIC_NAME>:<VALUE>`: %q", data)
			}
			ret.MetricName = xs[0]
			ret.Value, err = strconv.ParseFloat(xs[1], 64)
			if err != nil {
				return ret, fmt.Errorf("malformed data, cannot parse %q as metric value", xs[1])
			}
		case 1:
			ret.Type = UNKNOWN
			if typ, ok := statsdTypes[e]; ok {
				ret.Type = typ
			} else {
				return ret, fmt.Errorf("unknown type: %q", e)
			}
		case 2, 3:
			if strings.HasPrefix(e, "@") {
				ret.SampleRate, err = strconv.ParseFloat(e[1:], 64)
				if err != nil {
					return ret, fmt.Errorf("malformed data, cannot parse %q as sample rate", e)
				}
			}
			if strings.HasPrefix(e, "#") {
				// <TAG_KEY_1>:<TAG_VALUE_1>,<TAG_KEY_2>:<TAG_VALUE_2>
				for _, x := range strings.Split(e[1:], ",") {
					xxs := strings.Split(x, ":")
					if len(xxs) != 2 {
						return ret, fmt.Errorf("malformed data, cannot parse %q as <TAG_KEY>:<TAG_VALUE>", xxs)
					}
					ret.Tags = append(ret.Tags, Tag{xxs[0], xxs[1]})
				}
			}
		default:
			return ret, fmt.Errorf("malformed data, cannot parse %q", e)
		}
	}
	return ret, err
}

func convert(msg StatsdMsg) *write.Point {
	// influx 写接口是怎么调的， client lib 有定义结构体吗
	// -- 构造 Point, 然后调 WritePoint. 其中 Point 是 lib 定义好的结构
	ret := write.NewPointWithMeasurement(msg.MetricName)
	for _, tag := range msg.Tags {
		ret.AddTag(tag.Key, fmt.Sprintf("%v", tag.Value))
	}
	ret.AddField("value", msg.Value)
	return ret
}

// custom_metric:80|g|#port:8125
// TODO struct can be found in datadog-agent source code?
type StatsdMsg struct {
	MetricName string
	Value      float64 // TODO should be general
	Type       StatsdMsgType
	SampleRate float64
	Tags       []Tag
}
type Tag struct {
	Key   string
	Value interface{}
}
