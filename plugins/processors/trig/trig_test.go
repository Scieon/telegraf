package trig
//
//import (
//	"github.com/influxdata/telegraf"
//	"github.com/influxdata/telegraf/metric"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"testing"
//	"time"
//)
//
//const measurement_name = "request_timing"
////const lestring = "{\"V\":\"2.0\",\"SamplingRate\":1000,\"Metrics\":{\"eb/core-time\":95,\"eb/pre-core\":0,\"eb/processRequest\":93,\"eb/processRequest/getProfileData\":1,\"eb/processRequest/proxyRTBBid\":91,\"eb/processRequest/proxyRTBBid/ajaxAuction\":90,\"eb/processRequest/proxyRTBBid/ajaxAuction/getWinningResponse\":74,\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction\":13,\"eb/processRequest/error15\":1,\"eb/total\":131},\"Dimensions\":{\"hostname\":\"myhost\",\"code-version\":\"code-version\",\"rules-serial-number\":\"rules-serial-number\",\"session-id\":\"session-id\",\"features-toggle1\":\"on\",\"features-toggle2\":\"off\"}}"
//const lestring = "{\"V\":\"2.0\",\"SamplingRate\":1000,\"Metrics\":{\"eb/core-time\":\"131\",\"eb/post-core\":\"1\",\"eb/pre-core\":\"1\",\"eb/processRequest\":\"121\",\"eb/processRequest/getProfileData\":\"1\",\"eb/processRequest/proxyRTBBid\":\"121\",\"eb/processRequest/proxyRTBBid/ajaxAuction\":\"121\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/handleShowdowBids\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/handleTimedoutDSP\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/init\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/GetWinningResponse/collectResponses/readDSPResponse\":\"101\",\"eb/processRequest/proxyRTBBid/ajaxAuction/getWinningResponse\":\"101\",\"eb/processRequest/proxyRTBBid/ajaxAuction/initialize\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/prepareIBReply\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/processSmartQPS\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction\":\"11\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/RTBDSPLoop\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/buildGlobalRequestData\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/preRTBDSPLoop\":\"1\",\"eb/processRequest/proxyRTBBid/initializeRTB\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/getWinningResponse/collectResponse\":\"101\",\"eb/processRequest/proxyRTBbid/ajaxAuction/getWinningResponse/performAuction\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/buildSharedRequest\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/dspEligibility\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/oneTimePublisher\":\"1\",\"eb/processRequest/proxyRTBbid/ajaxAuction/setupAuction/sendRequest\":\"1\",\"eb/processRequest/proxyRTBBid/ajaxAuction/setupAuction/setupBidCache\":\"1\",\"eb/processRequest/updateIBReplayWithUserInfo\":\"1\",\"eb/total\":\"131\"},\"Dimensions\":{\"hostname\":\"111\",\"fakeDimension\":221}}"
//
////compares metrics without comparing time
//func compareMetrics(t *testing.T, expected, actual []telegraf.Metric) {
//	assert.Equal(t, len(expected), len(actual))
//	for i, metric := range actual {
//		require.Equal(t, expected[i].Name(), metric.Name())
//		require.Equal(t, expected[i].Fields(), metric.Fields())
//		require.Equal(t, expected[i].Tags(), metric.Tags())
//	}
//}
//
//func Metric(v telegraf.Metric, err error) telegraf.Metric {
//	if err != nil {
//		panic(err)
//	}
//	return v
//}
//
//func createMetric() {
//
//}
//func TestApply(t *testing.T) {
//
//}
//
//func TestFoo(t *testing.T) {
//	var timeNow time.Time
//	location, _ := time.LoadLocation("America/Montreal")
//	timeNow = time.Now().In(location)
//
//
//	tag := map[string]string {
//
//	}
//	field := map[string]interface{} {
//	"blob": lestring,
//	}
//
//	metric,_ := metric.New(measurement_name, tag, field, timeNow)
//	f := RequestTimingProcessor{}
//	//kafka := f.MakeUnmarshal(metric)
//
//	//fmt.Println(metric)
//	//assert.Equal(kafka.Dimensions)
//
//	//assert.Equal(t,2, 3)
//}
//
//func TestGetHost(t *testing.T) {
//
//}
//
//func TestBadApply(t *testing.T) {
//}
//
//// Benchmarks
//
//func getMetricFields(metric telegraf.Metric) interface{} {
//	key := "field3"
//	if value, ok := metric.Fields()[key]; ok {
//		return value
//	}
//	return nil
//}
//
//func getMetricFieldList(metric telegraf.Metric) interface{} {
//	key := "field3"
//	fields := metric.FieldList()
//	for _, field := range fields {
//		if field.Key == key {
//			return field.Value
//		}
//	}
//	return nil
//}
