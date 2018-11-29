package telemetry

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/processors"
)

const metricName = "request_timing"

type KafkaMetrics struct {
	Metrics    map[string]interface{} `json:"Metrics"`
	Dimensions map[string]interface{} `json:"Dimensions"`
	Version    string                 `json:"V"`
}

type RequestTimingProcessor struct {
}

var sampleConfig = `
	## This processor should only accept request_timing metrics
	namepass = ["request_timing"]
`

func (p *RequestTimingProcessor) SampleConfig() string {
	return sampleConfig
}

func (p *RequestTimingProcessor) Description() string {
	return "This is a processor to parse request timing kafka messages"
}

func (p *RequestTimingProcessor) Apply(in ...telegraf.Metric) []telegraf.Metric {
	fmt.Println("=====> IN!")
	if len(in) == 0 {
		log.Printf("E! [processors.as_request_timing] did not receive an input")
		return make([]telegraf.Metric, 0)
	}

	kafkaInput := in[0]
	metrics := make([]telegraf.Metric, 0)

	kafkaMetric, err := p.unmarshalKafkaInput(kafkaInput)

	// Do not process input if structure is malformed
	if err != nil {
		log.Printf("E! [processors.as_request_timing] could not unmarshal input: %v", err)
		return make([]telegraf.Metric, 0)
	}

	total := p.getTotal(kafkaMetric)
	host := p.getHost(kafkaMetric)

	// Create the telegraf metrics from kafka_consumer input
	for requestTimingMetric, durationDiff := range kafkaMetric.Metrics {
		telegrafMetric, err := p.createTelegrafMetric(requestTimingMetric, durationDiff, total, host)

		if err != nil {
			log.Printf("E! [processors.as_request_timing] could not create metrics: %v", err)
			return make([]telegraf.Metric, 0)
		}
		metrics = append(metrics, telegrafMetric)
	}

	return metrics
}

func (p *RequestTimingProcessor) getTotal(kafkaMetric KafkaMetrics) interface{} {
	for key, value := range kafkaMetric.Metrics {
		if strings.Contains(key, "/total") {
			return value
		}
	}
	return nil
}

func (p *RequestTimingProcessor) getHost(kafkaMetric KafkaMetrics) string {
	for key, value := range kafkaMetric.Dimensions {
		if strings.Contains(key, "hostname") {
			return fmt.Sprint(value)
		}
	}
	return ""
}

// Transform kafka input into a kafkaMetric struct
func (p *RequestTimingProcessor) unmarshalKafkaInput(kafkaInput telegraf.Metric) (KafkaMetrics, error) {
	kafkaMetrics := &KafkaMetrics{}
	for _, jsonString := range kafkaInput.Fields() {
		jsonObj := fmt.Sprint(jsonString)
		err := json.Unmarshal([]byte(jsonObj), kafkaMetrics)
		if err != nil {
			return *kafkaMetrics, fmt.Errorf("could not unmarshal kafka input")
		}

		if kafkaMetrics.Version != "2.0" {
			return *kafkaMetrics, fmt.Errorf("invalid version number")
		}
	}
	return *kafkaMetrics, nil
}

// Create telegraf metrics from the kafka_consumer input
func (p *RequestTimingProcessor) createTelegrafMetric(requestTimingMetric string, diff interface{}, total interface{}, host string) (telegraf.Metric, error) {
	var t time.Time
	location, _ := time.LoadLocation("America/Montreal")
	t = time.Now().In(location)

	var durationDiff, metricTotal int
	var err error

	durationDiff, err = parseNumeric(diff)
	metricTotal, err = parseNumeric(total)

	if err == nil && total != nil && host != "" {
		telegrafMetric, _ := metric.New(metricName, map[string]string{"metric": requestTimingMetric, "host": host}, map[string]interface{}{"duration-diff": durationDiff, "total": metricTotal}, t)
		return telegrafMetric, nil
	} else {
		return nil, fmt.Errorf("could not create telegraf metric")
	}
}

func parseNumeric(a interface{}) (int, error) {
	var result uint
	switch v := a.(type) {
	case string:
		result64, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			// string is a "Float"
			resultF64, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return 0, err
			}
			resultF64 = math.Round(resultF64)
			result = uint(resultF64)
		} else {
			result = uint(result64)
		}
	case float32:
		v = float32(math.Round(float64(v)))
		result = uint(v)
	case float64:
		v = math.Round(v)
		result = uint(v)
	case int:
		result = uint(v)
	case int32:
		result = uint(v)
	case int64:
		result = uint(v)
	}
	return int(result), nil
}

func init() {
	processors.Add("as_request_timing", func() telegraf.Processor {
		return &RequestTimingProcessor{}
	})
}
