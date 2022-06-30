package provider

import (
	"os"
)

func setEnv(envMap map[string]any) error {
	for k, v := range envMap {
		val := v.(string)
		err := os.Setenv(k, val)
		if err != nil {
			return err
		}
	}
	return nil
}
