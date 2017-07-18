package collectors

import "fmt"
import "regexp"
import "github.com/c9s/goprocinfo/linux"
import "github.com/mike-a-davis/metricsd/mappings"
import "github.com/mike-a-davis/metricsd/structs"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

type IostatCollector struct {
	enabled bool
}

func (c *IostatCollector) Enabled() bool {
	return c.enabled
}

func (c *IostatCollector) State(state bool) {
	c.enabled = state
}

func (c *IostatCollector) Setup(conf ini.File) {
	c.State(true)
}

func (c *IostatCollector) Report() (structs.MetricSlice, error) {
	var report structs.MetricSlice
	data, _ := c.collect()

	if data != nil {
		for device, values := range data {
			for k, v := range values {
				metric := structs.BuildMetric("IostatCollector", "iostat", "gauge", k, v, structs.FieldsMap{
					"device": device,
					"unit":   "IO",
					"name":   k,
					// "raw_value": v,
				})
				metric.Path = fmt.Sprintf("iostat.%s", device)
				report = append(report, metric)
			}
		}
	}

	return report, nil
}

func (c *IostatCollector) collect() (map[string]mappings.MetricMap, error) {
	stat, err := linux.ReadDiskStats("/proc/diskstats")
	if err != nil {
		logrus.Fatal("stat read fail")
		return nil, err
	}

	diskusageMapping := map[string]mappings.MetricMap{}
	pattern := "PhysicalDrive[0-9]+$|md[0-9]+$|sd[a-z]+[0-9]*$|x?vd[a-z]+[0-9]*$|disk[0-9]+$|dm\\-[0-9]+$"
	r, _ := regexp.Compile(pattern)

	for i := range stat {
		if !r.MatchString(stat[i].Name) {
			continue
		}

		diskusageMapping[stat[i].Name] = mappings.MetricMap{
			// "average_queue_length": TODO
			// "average_request_size_byte": TODO
			// "await": TODO
			// "concurrent_io": TODO
			"io": stat[i].ReadIOs + stat[i].WriteIOs,
			// "io_in_progress": TODO
			// "io_milliseconds": TODO
			// "io_milliseconds_weighted": TODO
			// "iops": TODO
			// "read_await": TODO
			// "read_byte_per_second": TODO
			// "read_requests_merged_per_second": TODO
			"reads": stat[i].ReadIOs,
			// "reads_byte": TODO
			"reads_merged": stat[i].ReadMerges,
			// "reads_milliseconds": TODO
			// "reads_per_second": TODO
			// "service_time": TODO
			// "util_percentage": TODO
			// "write_await": TODO
			// "write_byte_per_second": TODO
			// "write_requests_merged_per_second": TODO
			"writes": stat[i].WriteIOs,
			// "writes_byte": TODO
			"writes_merged": stat[i].WriteMerges,
			// "writes_milliseconds": TODO
			// "writes_per_second": TODO
		}
	}

	return diskusageMapping, nil
}
