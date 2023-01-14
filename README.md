# DogData

## Run

+ Start datadog-agent
```
cd <datadog-agent directory>
./bin/agent/agent run -c bin/agent/dist/datadog.yaml
```

+ Start influxDB
`influxd`

+ Start dogdata

+ Send a statsD message through `datadog-agent` to `dogdata`.

`echo -n "custom_metric:82|g|#port:9125" | nc -4u -w0 127.0.0.1 9125`

+ Verify the data in influxDB.
