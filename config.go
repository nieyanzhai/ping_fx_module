package ping_fx_module

import (
	"encoding/json"
	"os"
	"path"
)

var (
	cfgPath = path.Join("./config", "ping.json")
)

type config struct {
	IP       string `json:"ip"`
	Count    int
	Timeout  int
	MaxLoss  float64 `json:"max_loss"`
	Interval int
}

func loadConfig(path string) (config, error) {
	var t config
	bytes, err := os.ReadFile(path)
	if err != nil {
		return t, err
	}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return t, err
	}

	return t, nil
}
