package apputils

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig(target interface{}, path, mode string) {
	confFile, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var configs map[string]interface{}
	if err := yaml.Unmarshal(confFile, &configs); err != nil {
		log.Fatal(err)
	}

	config, ok := configs[mode]
	if !ok {
		modes := []string{}
		for mode := range configs {
			modes = append(modes, mode)
		}
		log.Fatalf("invalid mode %v. expected modes: %v", mode, modes)
	}

	// 🙏 https://github.com/go-yaml/yaml/issues/13#issuecomment-428952604
	b, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(b, target); err != nil {
		log.Fatal(err)
	}
}
