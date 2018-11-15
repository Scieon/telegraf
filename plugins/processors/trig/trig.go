package trig

import (
	"encoding/json"
	"fmt"
	"strconv"
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
	Metrics map[string]string `json:"Metrics"`
	Dimensions map[string]interface{} `json:"Dimensions"`
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

//func (p *RequestTimingProcessor) getEBTotal(kafkaInput telegraf.Metric) interface{} {
//	for key, value := range kafkaInput.Fields() {
//		if strings.Contains(key, "eb/total") {
//			return value
//		}
//	}
//	return nil
//}

func (p *RequestTimingProcessor) getEBTotal(kafkaMetric KafkaMetrics) interface{} {
	for key, value := range kafkaMetric.Metrics{
		if strings.Contains(key, "eb/total") {
			return value
		}
	}
	return nil
}

func (p* RequestTimingProcessor) getHost(kafkaMetric KafkaMetrics) string {
	for key, value := range kafkaMetric.Dimensions{
		if strings.Contains(key, "hostname") {
			return fmt.Sprint(value)
		}
	}
	return ""
}

func (p *RequestTimingProcessor) MakeUnmarshal(kafkaInput telegraf.Metric) KafkaMetrics {
	kafkaMetrics := &KafkaMetrics{}

	for _, jsonObj := range kafkaInput.Fields() {
		jsonString := fmt.Sprint(jsonObj)
		err := json.Unmarshal([]byte(jsonString), kafkaMetrics)
		if err != nil {
			//return err
			fmt.Println(err)
			fmt.Println("===> Error")

		}
	}
	fmt.Println(kafkaMetrics)

	for k,v := range kafkaMetrics.Metrics {
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println("--------")
	}

	return *kafkaMetrics
}

func (p *RequestTimingProcessor) makePoints(kafkaInput telegraf.Metric) ([]Point, error) {
	points := []Point{}

	kafkaMetric := p.MakeUnmarshal(kafkaInput)
	total := p.getEBTotal(kafkaMetric)
	host := p.getHost(kafkaMetric)

	if total == nil {
		fmt.Println("eb total not found")
		return nil, fmt.Errorf("eb total not found")
	}

	for requestTimingMetric, durationDiff := range kafkaMetric.Metrics {
		if durationDiff, err := strconv.Atoi(durationDiff); err == nil {
			newPoint := NewPoint("request_timing", map[string]string{"metric": requestTimingMetric, "host": host}, map[string]interface{}{"duration-diff": durationDiff, "total": total})
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
