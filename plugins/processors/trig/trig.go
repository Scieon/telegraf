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

type Printer struct {
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

func (p *Printer) SampleConfig() string {
	return sampleConfig
}

func (p *Printer) Description() string {
	return "this is a test trig processor"
}

func (p *Printer) Apply(in ...telegraf.Metric) []telegraf.Metric {

	metrics := in[0]
	var points []Point
	for _, kafkaInput := range in {
		points = p.makePoints(kafkaInput)
		//p.showMetrics(kafkaInput)
	}

	in = nil

	for _, point := range points {
		newMetric, _ := metric.New(point.MeasurementName, point.Tags, point.Fields, point.Timestamp)
		in = append(in, newMetric)
	}
	fmt.Println(metrics)

	fmt.Println("------------------")
	fmt.Println("Len of points: " + string(len(in)))
	for _, v := range in {
		fmt.Println(v)
	}

	fmt.Println("...........")

	return in
}

func (p *Printer) makePoints(kafkaInput telegraf.Metric) []Point {

	points := []Point{}
	var total interface{}

	host, _ := kafkaInput.GetTag("host")
	for k, v := range kafkaInput.Fields() {
		if strings.Contains(k, "Metrics") {
			if strings.Contains(k, "eb/total") {
				total = v
			} else {
				parsedKey := strings.TrimPrefix(k, "Metrics_")
				newPoint := NewPoint("request_timing", map[string]string{"metric": parsedKey, "host": host}, map[string]interface{}{"duration-diff": v, "total": total}, time.Time{}) //todo
				points = append(points, newPoint)
			}
		}
	}

	for _, v := range points {
		fmt.Println(v)
		fmt.Println("^^^^^^")
	}

	return points
}
func (p *Printer) showMetrics(metric telegraf.Metric) {
	for k, v := range metric.Fields() {
		fmt.Println("------FIELDS-----")
		fmt.Println("Key: " + k)
		fmt.Println(v)
		fmt.Println("--------------")
	}

	for k, v := range metric.Tags() {
		fmt.Println("-----TAGS------")
		fmt.Println("Key: " + k)
		fmt.Println(v)
		fmt.Println("--------------")
	}
}

func init() {
	processors.Add("trig", func() telegraf.Processor {
		return &Printer{DropOriginal: true}
	})
}
