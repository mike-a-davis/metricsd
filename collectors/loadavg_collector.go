package collectors

import "github.com/c9s/goprocinfo/linux"
import "github.com/mike-a-davis/metricsd/mappings"
import "github.com/mike-a-davis/metricsd/structs"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

type LoadAvgCollector struct {
	enabled bool
}

func (c *LoadAvgCollector) Enabled() bool {
	return c.enabled
}

func (c *LoadAvgCollector) State(state bool) {
	c.enabled = state
}

func (c *LoadAvgCollector) Setup(conf ini.File) {
	c.State(true)
}

func (c *LoadAvgCollector) Report() (structs.MetricSlice, error) {
	var report structs.MetricSlice
	values, _ := c.collect()

	if values != nil {
		for k, v := range values {
			metric := structs.BuildMetric("LoadAvgCollector", "loadavg", "gauge", k, v, structs.FieldsMap{
				"unit": "Load",
				"name": k,
				// "raw_value": v,
			})
			report = append(report, metric)
		}
	}

	return report, nil
}

func (c *LoadAvgCollector) collect() (mappings.MetricMap, error) {
	stat, err := linux.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		logrus.Fatal("stat read fail")
		return nil, err
	}

	// TODO: Add processes_running and processes_total,
	// unit:processes, type:(running|total)
	return mappings.MetricMap{
		"01": stat.Last1Min,
		"05": stat.Last5Min,
		"15": stat.Last15Min,
	}, nil
}
