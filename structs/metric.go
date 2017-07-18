package structs

import "encoding/json"
import "fmt"
import "os"
import "time"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

type Metric struct {
	Collector string
	Path      string
	From      string
	Name      string
	Value     interface{}
	// RawValue   interface{}
	Timestamp  time.Time
	Precision  int
	Host       string
	MetricType string
	TTL        int
	Data       FieldsMap
	Tags       FieldsMap
}

var Hostname string

func init() {
	var err error
	Hostname, err = os.Hostname()

	if err != nil {
		logrus.Warning("error retrieving hostname, using `localhost` for now")
		Hostname = "localhost"
	}
}

type FieldsMap map[string]interface{}

type MetricSlice []*Metric

func BuildMetric(collector string, from string, metricType string, name string, value interface{}, data FieldsMap) *Metric {
	return &Metric{
		Timestamp:  time.Now(),
		Value:      value,
		MetricType: metricType,
		Host:       Hostname,
		Collector:  collector,
		From:       from,
		Name:       name,
		Precision:  0,
		TTL:        0,
		Data:       data,
	}
}

func (m *Metric) Process(conf ini.File) {
	if hostname, ok := conf.Get("metricsd", "hostname"); ok {
		m.Host = hostname
	}

	if hostname, ok := conf.Get(m.Collector, "hostname"); ok {
		m.Host = hostname
	}
}

func (m *Metric) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	tags := make(map[string]interface{})

	mlxName := m.Data["name"]
	mlxCollector := m.From
	data["timestamp"] = int32(m.Timestamp.Unix())
	data["unit"] = m.Data["unit"]
	data["name"] = fmt.Sprintf("%s.%s", mlxCollector, mlxName)
	data["target_type"] = m.MetricType

	data["result"] = m.Value
	if _, ok := data["host"]; !ok {
		data["host"] = m.Host
	}

	tags["version"] = "1"
	// tags["collector"] = m.From
	// tags["type"] = m.Name

	for k, v := range m.Data {
		_, exists := data[k]
		if exists {
			// data[fmt.Sprintf("%s", k)] = v
		} else {
			tags[k] = v
		}
	}
	// tags["raw_value"] = "true"

	data["tags"] = tags
	return data
}

func (m *Metric) ToJSON() []byte {
	data := m.ToMap()

	mlxMetric, err := json.Marshal(data)
	if err != nil {
		fmt.Errorf("Failed to marshal fields to JSON, %v", err)
		return nil
	}
	return mlxMetric
}

func (m *Metric) ToGraphite(prefix string) (response string) {
	path := m.From
	if m.Path != "" {
		path = m.Path
	}
	key := fmt.Sprintf("%s.%s.%s", m.Host, path, m.Name)
	if prefix != "" {
		key = fmt.Sprintf("%s%s", prefix, key)
	}
	return fmt.Sprintf("%s %v %d", key, m.Value, int32(m.Timestamp.Unix()))
}
