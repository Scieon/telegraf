package telemetry

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const measurementName = "request_timing"
const kafkaConsumerInput = "{\"V\":\"2.0\",\"SamplingRate\":1000,\"Metrics\":{\"eb/core-time\":131,\"eb/post-core\":\"1\",\"eb/pre-core\":\"1\",\"eb/processRequest\":\"121\",\"eb/processRequest/getProfileData\":\"1\",\"eb/processRequest/proxyRTBBid\":\"121\",\"eb/processRequest/proxyRTBBid/ajaxAuction\":\"121\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/handleShowdowBids\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/handleTimedoutDSP\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/init\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/readDSPResponse\":\"101\",\"eb/processRequest/proxyRTBBid/ajaxAuction/getWinningResponse\":\"101\",\"eb/processRequest/proxyRTBBid/ajaxAuction/initialize\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/prepareIBReply\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/processSmartQPS\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction\":\"11\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/RTBDSPLoop\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/buildGlobalRequestData\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/preRTBDSPLoop\":\"1\",\"eb/processRequest/proxyRTBBid/initializeRTB\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/getWinningResponse/collectResponse\":\"101\",\"eb/processRequest/proxyRTBbid/ajaxAuction/getWinningResponse/performAuction\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/buildSharedRequest\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/dspEligibility\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/oneTimePublisher\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/sendRequest\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/setupBidCache\":\"1\",\"eb/processRequest/updateIBReplayWithUserInfo\":\"1\",\"eb/total\":\"131\"},\"Dimensions\":{\"hostname\":\"111\",\"fakeDimension\":221}}"
const KafkaConsumerInputFloat = "{\"V\":\"2.0\",\"SamplingRate\":1000,\"Metrics\":{\"eb/core-time\":0.5,\"eb/total\":\"140.5\"},\"Dimensions\":{\"hostname\":\"111\",\"fakeDimension\":221}}"
const invalidKafkaConsumerInputNoVersion ="{\"SamplingRate\":1000,\"Metrics\":{\"eb/core-time\":\"131\",\"eb/post-core\":\"1\",\"eb/pre-core\":\"1\",\"eb/processRequest\":\"121\",\"eb/processRequest/getProfileData\":\"1\",\"eb/processRequest/proxyRTBBid\":\"121\",\"eb/processRequest/proxyRTBBid/ajaxAuction\":\"121\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/handleShowdowBids\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/handleTimedoutDSP\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/init\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/readDSPResponse\":\"101\",\"eb/processRequest/proxyRTBBid/ajaxAuction/getWinningResponse\":\"101\",\"eb/processRequest/proxyRTBBid/ajaxAuction/initialize\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/prepareIBReply\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/processSmartQPS\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction\":\"11\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/RTBDSPLoop\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/buildGlobalRequestData\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/preRTBDSPLoop\":\"1\",\"eb/processRequest/proxyRTBBid/initializeRTB\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/getWinningResponse/collectResponse\":\"101\",\"eb/processRequest/proxyRTBbid/ajaxAuction/getWinningResponse/performAuction\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/buildSharedRequest\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/dspEligibility\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/oneTimePublisher\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/sendRequest\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/setupBidCache\":\"1\",\"eb/processRequest/updateIBReplayWithUserInfo\":\"1\",\"eb/total\":\"131\"},\"Dimensions\":{\"hostname\":\"111\",\"fakeDimension\":221}}"
const invalidKafkaConsumerInputVersion = "{\"V\":\"1.0\",\"SamplingRate\":1000,\"Metrics\":{\"eb/core-time\":\"131\",\"eb/post-core\":\"1\",\"eb/pre-core\":\"1\",\"eb/processRequest\":\"121\",\"eb/processRequest/getProfileData\":\"1\",\"eb/processRequest/proxyRTBBid\":\"121\",\"eb/processRequest/proxyRTBBid/ajaxAuction\":\"121\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/handleShowdowBids\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/handleTimedoutDSP\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/init\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/readDSPResponse\":\"101\",\"eb/processRequest/proxyRTBBid/ajaxAuction/getWinningResponse\":\"101\",\"eb/processRequest/proxyRTBBid/ajaxAuction/initialize\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/prepareIBReply\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/processSmartQPS\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction\":\"11\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/RTBDSPLoop\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/buildGlobalRequestData\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/preRTBDSPLoop\":\"1\",\"eb/processRequest/proxyRTBBid/initializeRTB\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/getWinningResponse/collectResponse\":\"101\",\"eb/processRequest/proxyRTBbid/ajaxAuction/getWinningResponse/performAuction\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/buildSharedRequest\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/dspEligibility\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/oneTimePublisher\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/sendRequest\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/setupBidCache\":\"1\",\"eb/processRequest/updateIBReplayWithUserInfo\":\"1\",\"eb/total\":\"131\"},\"Dimensions\":{\"hostname\":\"111\",\"fakeDimension\":221}}"
const invalidKafkaConsumerMessageEmptyDimension = "{\"V\":\"2.0\",\"Metrics\":{\"eb/core-time\":131,\"eb/total\":132},\"FakeDimension\":{}}"


func createMetric(kafkaInput string) telegraf.Metric {
	var timeNow time.Time
	location, _ := time.LoadLocation("America/Montreal")
	timeNow = time.Now().In(location)
	field := map[string]interface{}{
		"value": kafkaInput,
	}

	telegrafMetric, _ := metric.New(measurementName, nil, field, timeNow)

	return telegrafMetric
}

func TestUnmarshalValidKafkaInput(t *testing.T) {
	processor := RequestTimingProcessor{}
	telegrafMetric := createMetric(kafkaConsumerInput)

	kafkaMetric, _ := processor.unmarshalKafkaInput(telegrafMetric)
	assert.Equal(t, kafkaMetric.Dimensions["hostname"], "111")
	assert.Equal(t, len(kafkaMetric.Metrics), 29)
}

func TestUnmarshalInvalidMessageVersion(t *testing.T) {
	processor := RequestTimingProcessor{}
	telegrafMetric := createMetric(invalidKafkaConsumerInputVersion)

	_, err := processor.unmarshalKafkaInput(telegrafMetric)
	assert.Error(t, err)

	telegrafMetric = createMetric(invalidKafkaConsumerInputNoVersion)
	_, err = processor.unmarshalKafkaInput(telegrafMetric)
	assert.Error(t, err)
}

func TestGetTotal(t *testing.T) {
	processor := RequestTimingProcessor{}
	telegrafMetric := createMetric(kafkaConsumerInput)

	kafkaMetric, _ := processor.unmarshalKafkaInput(telegrafMetric)
	total := processor.getTotal(kafkaMetric)
	assert.Equal(t, total, "131")
}

func TestApply(t *testing.T) {
	processor := RequestTimingProcessor{}
	telegrafMetric := createMetric(kafkaConsumerInput)

	metrics := processor.Apply(telegrafMetric)
	assert.Equal(t, 29, len(metrics))

	// Test for tags
	tag, _ := metrics[0].GetTag("host")
	assert.Equal(t, tag, "111")

	tag, _ = metrics[0].GetTag("metric")
	assert.Contains(t, tag, "eb/")

	// Test for metric field
	field, _ := metrics[0].GetField("total")
	assert.Equal(t, field, int64(131))

	_, present := metrics[0].GetField("duration-diff")
	assert.Equal(t, present, true)
}

func TestApplyFloatInput(t *testing.T) {
	processor := RequestTimingProcessor{}
	telegrafMetric := createMetric(KafkaConsumerInputFloat)

	metrics := processor.Apply(telegrafMetric)
	assert.Equal(t, len(metrics), 2)

	// Test for metric field
	field, _ := metrics[0].GetField("total")
	assert.Equal(t, field, int64(141))
}


func TestBadApply(t *testing.T) {
	processor := RequestTimingProcessor{}
	telegrafMetric := createMetric(invalidKafkaConsumerMessageEmptyDimension)
	metrics := processor.Apply(telegrafMetric)
	expectedResult := make([]telegraf.Metric, 0)

	// If apply was not successful the empty telegraf metrics will be returned
	assert.Equal(t, len(metrics), 0)
	assert.Equal(t, metrics, expectedResult)
}
