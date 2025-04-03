package pkg

import (
	"fmt"
	"log"

	"github.com/DataDog/datadog-go/v5/statsd"
)

var statsd_cli *statsd.Client

func InitStatsD(address, env, service string) {
	var err error
	statsd_cli, err = statsd.New(address,
		statsd.WithTags([]string{fmt.Sprintf("env:%s", env), fmt.Sprintf("service:%s", service)}),
		statsd.WithExtendedClientSideAggregation())
	if err != nil {
		log.Fatal(err)
	}
}

func SendLatency(rtt float64) {
	// gauge只会保留最新的一个数据
	// statsd_cli.Gauge("quant.ob.delay_gauge", float64(latency), nil, 1)
	// fmt.Println(rtt)
	statsd_cli.Histogram("blofin_lp.test.rtt", rtt, nil, 1)
	// statsd_cli.Histogram("blofin_lp.test.offset", offset, nil, 1)
}
