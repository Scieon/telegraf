package trig

// printer.go

import (
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
)

type Printer struct {
	DropOriginal bool     `toml:"drop_original"`
	Tag          []string `toml:"tag"`
	Field        []string `toml:"field"`
}

type ParsedObj struct {
	Name string
	duration string
	metric string
	// total
}

type ProcessorMetric struct {
	Tags []telegraf.Tag
	Fields []telegraf.Field
}

func (m *ProcessorMetric) addTag(key, value string) {
	m.Tags = append(m.Tags, telegraf.Tag{Key: key, Value: value})
}

func (m *ProcessorMetric) addField(key string, value interface{}) {
	m.Fields = append(m.Fields, telegraf.Field{Key: key, Value: value})
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

	//var metrics [] telegraf.Metric
	metrics := in[0]
	for _, kafkaInput := range in {

		p.makeMetrics(kafkaInput)
		//p.showMetrics(kafkaInput)
		//for k, _ := range kafkaInput.Fields() {
		//	kafkaInput.RemoveField(k)
		//	kafkaInput.AddField("duration-diff", "37")
		//}
		//fmt.Println(kafkaInput)
	}

	fmt.Println(metrics)
	return in
}

func (p *Printer) makeMetrics(kafkaInput telegraf.Metric) {

	tempHolder := ProcessorMetric{}

	for k, v := range kafkaInput.Fields() {



		tag := telegraf.Tag{Key: "metric", Value: k}
		field := telegraf.Field{Key: "duration-diff", Value: v}

		//kafkaInput.AddTag("metric", k)
		//kafkaInput.AddField("duration-diff", v)

		tempHolder.Tags = append(tempHolder.Tags, tag)
		tempHolder.Fields = append(tempHolder.Fields, field)

		// append metric obj to array


		//fmt.Println("Removing...: " + k)
		//kafkaInput.AddTag("metric", k)
		//kafkaInput.AddField("duration-diff", v)

		//kafkaInput.RemoveField(k)
		//fmt.Println("--------------")
	}

	fmt.Println("+++++++++++++++++")
	for k,v := range tempHolder.Tags {
		fmt.Println("///////")
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println("///////")
	}
	for k,v := range tempHolder.Fields {
		fmt.Println("///---////")
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println("///---////")
	}

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
