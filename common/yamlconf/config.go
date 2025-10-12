package yamlconf

import (
	"os"

	"gopkg.in/yaml.v2"
)

func Load(path string, cfg any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err = yaml.UnmarshalStrict(data, cfg); err != nil {
		return err
	}

	return err
}
