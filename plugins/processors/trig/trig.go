package trig

// printer.go

import (
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/processors"
	"strings"
	"time"
)

// Point represents an InfluxDB point/row, and defines a "schema"
// Tags and Fields hold the column names of all tags and fields in the measurement/table
type Point struct {
	MeasurementName string
	Tags            map[string]string
	Fields          map[string]interface{}
	Timestamp       time.Time
}

func NewPoint(name string, tags map[string]string, fields map[string]interface{}, timestamp ...time.Time) Point {
	var t time.Time
	location, _ := time.LoadLocation("America/Montreal")
	t = time.Now().In(location)
	p := Point{
		MeasurementName: name,
		Tags:            tags,
		Fields:          fields,
		Timestamp:       t,
	}
	return p
}

type RequestTimingProcessor struct {
	DropOriginal bool     `toml:"drop_original"`
	Tag          []string `toml:"tag"`
	Field        []string `toml:"field"`
}

var sampleConfig = `
  ## If true, incoming metrics are not emitted.
	drop_original = true
	tag = []
	field = []
`

func (p *RequestTimingProcessor) SampleConfig() string {
	return sampleConfig
}

func (p *RequestTimingProcessor) Description() string {
	return "This is a processor to parse request timing kafka messages"
}

func (p *RequestTimingProcessor) Apply(in ...telegraf.Metric) []telegraf.Metric {

	var points []Point
	for _, kafkaInput := range in {
		points = p.makePoints(kafkaInput)
	}

	in = nil

	for _, point := range points {
		newMetric, _ := metric.New(point.MeasurementName, point.Tags, point.Fields, point.Timestamp)
		in = append(in, newMetric)
	}

	return in
}

func (p *RequestTimingProcessor) getEBTotal(kafkaInput telegraf.Metric) interface{} {
	for key, value := range kafkaInput.Fields() {
		if strings.Contains(key, "eb/total") {
			return value
		}
	}
	return nil
}
func (p *RequestTimingProcessor) makePoints(kafkaInput telegraf.Metric) []Point {

	points := []Point{}
	total := p.getEBTotal(kafkaInput)
	host, _ := kafkaInput.GetTag("host")

	if total == nil {
		fmt.Println("ERROR NO EB TOTAL")
		return nil
	}

	for k, v := range kafkaInput.Fields() {
		if strings.Contains(k, "Metrics") {
			parsedKey := strings.TrimPrefix(k, "Metrics_")
			newPoint := NewPoint("request_timing", map[string]string{"metric": parsedKey, "host": host}, map[string]interface{}{"duration-diff": v, "total": total}, time.Time{})
			points = append(points, newPoint)
		}
	}
	return points
}

func init() {
	processors.Add("trig", func() telegraf.Processor {
		return &RequestTimingProcessor{DropOriginal: true}
	})
}
