package main

import (
	"os"
	"log"
	"flag"
	"time"
	"io/ioutil"
)

func writeMetrics() {
	metrics := metricsGet()
	content := []byte(metrics)
	filename := config.Path + "/relmon.prom"
	err := ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		log.Printf("Error writing metrics file: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "/etc/prometheus/relmon.yml", "path to release monitor configuration file")
	flag.Parse()
	readConfigFile(configFile)

	for true {
		collectMetrics()
		writeMetrics()
		time.Sleep(time.Duration(config.Interval)*time.Second)
		metricsClear()
	}
}
