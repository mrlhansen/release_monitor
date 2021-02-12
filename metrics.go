package main

import (
	"fmt"
	"time"
)

var metrics string

func metricsClear() {
	metrics = ""
}

func metricsAppend(name string, version string, info string, value int64) {
	str := "release_monitor"
	ts := time.Unix(value ,0)

	str = str + `{name="` + name + `",version="` + version
	if len(info) > 0 {
		str = str + `",url="` + info
	}
	str = str + `",date="` + ts.Format("January 2, 2006") + `"} ` + fmt.Sprintf("%d\n", value)

	metrics += str
}

func metricsGet() string {
	return metrics
}
