package collectors

import "github.com/c9s/goprocinfo/linux"
import "github.com/mike-a-davis/metricsd/mappings"
import "github.com/mike-a-davis/metricsd/structs"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

type VmstatCollector struct {
	enabled bool
}

func (c *VmstatCollector) Enabled() bool {
	return c.enabled
}

func (c *VmstatCollector) State(state bool) {
	c.enabled = state
}

func (c *VmstatCollector) Setup(conf ini.File) {
	c.State(true)
}

func (c *VmstatCollector) Report() (structs.MetricSlice, error) {
	var report structs.MetricSlice
	values, _ := c.collect()

	if values != nil {
		for k, v := range values {
			metric := structs.BuildMetric("VmstatCollector", "vmstat", "rate", k, v, structs.FieldsMap{
				"unit": "Page",
				"name": k,
				// "raw_value": v,
			})
			report = append(report, metric)
		}
	}

	return report, nil
}

func (c *VmstatCollector) collect() (mappings.MetricMap, error) {
	stat, err := linux.ReadVMStat("/proc/vmstat")
	if err != nil {
		logrus.Fatal("stat read fail")
		return nil, err
	}

	return mappings.MetricMap{
		"paging_in": stat.PagePagein,
		"pagingout": stat.PagePageout,
		"swap_in":   stat.PageSwapin,
		"swap_out":  stat.PageSwapout,
	}, nil
}
