package shippers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mike-a-davis/metricsd/structs"
	"github.com/vaughan0/go-ini"
)

// MlxShipper is an exported type that
// allows shipping metrics to Mlx
type MlxShipper struct {
	enabled bool
	debug   bool
	url     string
}

// Enabled allows checking whether the shipper is enabled or not
func (s *MlxShipper) Enabled() bool {
	return s.enabled
}

// State allows setting the enabled state of the shipper
func (s *MlxShipper) State(state bool) {
	s.enabled = state
}

// Setup configures the shipper
func (s *MlxShipper) Setup(conf ini.File) {
	s.State(true)

	useDebug, ok := conf.Get("MlxShipper", "debug")
	if ok && useDebug == "true" {
		s.debug = true
	} else {
		s.debug = false
	}

	s.url = "http://127.0.0.1:8888/udm"
	useMlxURL, ok := conf.Get("MlxShipper", "url")
	if ok {
		s.url = useMlxURL
	}
}

// Ship sends a list of MetricSlices to Mlx
func (s *MlxShipper) Ship(logs structs.MetricSlice) error {

	for _, item := range logs {
		mlxMetric := item.ToJSON()
		if s.debug {
			fmt.Printf("%s\n", string(mlxMetric))
		}

		jsonStr := []byte(mlxMetric)
		req, err := http.NewRequest("POST", s.url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		if s.debug {
			fmt.Println(req)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if s.debug {
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if s.debug {
			fmt.Println("response Body:", string(body))
		}
	}

	return nil
}
