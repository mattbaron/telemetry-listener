package listener

import (
	"fmt"
	"strings"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/parsers/influx/influx_upstream"
	"github.com/influxdata/telegraf/plugins/serializers/influx"
	"github.com/mattbaron/telemetry-listener/client"
)

var serializer = influx.Serializer{}

type Batch struct {
	InputData      string
	Metrics        []telegraf.Metric
	ProcessedLines []string
	DroppedLines   []string
	Client         client.Client
}

func NewBatch(c client.Client) *Batch {
	return &Batch{
		Metrics:      make([]telegraf.Metric, 0),
		DroppedLines: make([]string, 0),
		Client:       c,
	}
}

func (b *Batch) processMetrics(metrics []telegraf.Metric) {
	for _, metric := range metrics {
		metric.AddTag("influxdb_node_group", b.Client.NodeGroup)
		metric.AddTag("influxdb_database", b.Client.Database)
	}
	b.Metrics = append(b.Metrics, metrics...)
}

func (b *Batch) ProcessEachLine(data []byte) error {
	parser := influx_upstream.Parser{}
	parser.Init()

	metrics := make([]telegraf.Metric, 0)
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		metric, err := parser.ParseLine(line)
		if err == nil {
			metrics = append(metrics, metric)
		} else {
			b.DroppedLines = append(b.DroppedLines, line)
		}
	}

	b.processMetrics(metrics)

	return nil
}

func (b *Batch) ProcessAll(data []byte) error {
	parser := influx_upstream.Parser{}
	parser.Init()

	metrics, err := parser.Parse(data)

	if err != nil {
		fmt.Printf("PARSE ERR: %v\n", err)
		return err
	}

	b.processMetrics(metrics)

	return nil
}

func (b *Batch) Serialize() {
	bytes, _ := serializer.SerializeBatch(b.Metrics)
	fmt.Println(string(bytes))
}
