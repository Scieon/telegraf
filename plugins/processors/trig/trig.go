package trig

import (
	"fmt"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/processors"
)

// Point represents an InfluxDB point/row, and defines a "schema"
// Tags and Fields hold the column names of all tags and fields in the measurement/table
type Point struct {
	MeasurementName string
	Tags            map[string]string
	Fields          map[string]interface{}
	Timestamp       time.Time
}

type KafkaMetrics struct {
	Total   string            `json:"total"`
	Metrics map[string]string `json:"metrics"`
	Time    time.Time
	Host    string
}

func NewPoint(name string, tags map[string]string, fields map[string]interface{}) Point {
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
	DropOriginal bool `toml:"drop_original"`
}

var sampleConfig = `
  ## If true, incoming metrics are not emitted.
	drop_original = true
`

func (p *RequestTimingProcessor) SampleConfig() string {
	return sampleConfig
}

func (p *RequestTimingProcessor) Description() string {
	return "This is a processor to parse request timing kafka messages"
}

func (p *RequestTimingProcessor) Apply(in ...telegraf.Metric) []telegraf.Metric {
	var points []Point
	var err error
	
	// Create the points from kafka_consumer input
	for _, kafkaInput := range in {
		//p.showMetrics(kafkaInput)
		points, err = p.makePoints(kafkaInput)
		if err != nil {
			fmt.Println("something went wrong")
			return nil
		}
	}

	in = nil

	// Create an array of Telegraf metrics from all the points
	for _, point := range points {
		newMetric, _ := metric.New(point.MeasurementName, point.Tags, point.Fields, point.Timestamp)
		in = append(in, newMetric)
	}

	fmt.Println("DONE!")
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

func (p *RequestTimingProcessor) makePoints(kafkaInput telegraf.Metric) ([]Point, error) {

	points := []Point{}
	total := p.getEBTotal(kafkaInput)
	host, _ := kafkaInput.GetTag("Dimensions_hostname")

	if total == nil {
		fmt.Println("eb total not found")
		return nil, fmt.Errorf("eb total not found")
	}

	for key, value := range kafkaInput.Fields() {
		if strings.Contains(key, "Metrics") {
			parsedKey := strings.TrimPrefix(key, "Metrics_")
			newPoint := NewPoint("request_timing", map[string]string{"metric": parsedKey, "host": host}, map[string]interface{}{"duration-diff": value, "total": total})
			points = append(points, newPoint)
		}
	}
	return points, nil
}

func (p *RequestTimingProcessor) showMetrics(metric telegraf.Metric) {
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
		return &RequestTimingProcessor{DropOriginal: true}
	})
}
