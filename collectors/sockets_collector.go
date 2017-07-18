package collectors

import "strings"
import "github.com/c9s/goprocinfo/linux"
import "github.com/mike-a-davis/metricsd/mappings"
import "github.com/mike-a-davis/metricsd/structs"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

type SocketsCollector struct {
	enabled bool
}

func (c *SocketsCollector) Enabled() bool {
	return c.enabled
}

func (c *SocketsCollector) State(state bool) {
	c.enabled = state
}

func (c *SocketsCollector) Setup(conf ini.File) {
	c.State(true)
}

func (c *SocketsCollector) Report() (structs.MetricSlice, error) {
	var report structs.MetricSlice
	values, _ := c.collect()

	if values != nil {
		for k, v := range values {
			fieldsMap := structs.FieldsMap{
				"unit": "Sock",
				"name": k,
				// "raw_value": v,
			}

			splitted := strings.Split(k, "_")
			protocol, metricType := splitted[0], k
			if len(splitted) > 1 {
				fieldsMap["protocol"] = protocol
				metricType = splitted[1]
			}

			metric := structs.BuildMetric("SocketsCollector", "sockets", "gauge", metricType, v, fieldsMap)
			report = append(report, metric)
		}
	}

	return report, nil
}

func (c *SocketsCollector) collect() (mappings.MetricMap, error) {
	stat, err := linux.ReadSockStat("/proc/net/sockstat")
	if err != nil {
		logrus.Fatal("stat read fail")
		return nil, err
	}

	return mappings.MetricMap{
		"tcp_alloc":  stat.TCPAllocated,
		"tcp_inuse":  stat.TCPInUse,
		"tcp_mem":    stat.TCPMemory,
		"tcp_orphan": stat.TCPOrphan,
		"tcp_tw":     stat.TCPTimeWait,
		"udp_inuse":  stat.UDPInUse,
		"udp_mem":    stat.UDPMemory,
		"used":       stat.SocketsUsed,
	}, nil
}
