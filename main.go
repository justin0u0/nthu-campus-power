package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	powerStation1 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "power_station_1",
	})
	powerStation2 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "power_station_2",
	})
	powerStation3 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "power_station_3",
	})
)

func main() {
	for _, m := range []struct {
		metric  prometheus.Gauge
		station int
	}{
		{powerStation1, 1},
		{powerStation2, 2},
		{powerStation3, 3},
	} {
		m := m
		go func() {
			for {
				power, err := getPower(m.station)
				if err != nil {
					log.Printf("Failed to get power station %d: %v", m.station, err)
					continue
				}

				m.metric.Set(float64(power))
				time.Sleep(5 * time.Second)
			}
		}()
	}

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

var re = regexp.MustCompile(`kW: ([\d,-]+)`)

func getPower(station int) (int, error) {
	url := fmt.Sprintf("http://140.114.188.57/nthu2020/fn1/kw%d.aspx", station)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	matches := re.FindSubmatch(b)
	power := bytes.ReplaceAll(matches[1], []byte(","), []byte(""))
	return strconv.Atoi(string(power))
}
