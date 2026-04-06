package mock

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadResponse(subdir, filename string) ([]byte, error) {
	path := filepath.Join(TestdataDir(), "responses", subdir, filename)
	return os.ReadFile(path)
}

func LoadResponseJSON(subdir, filename string) (map[string]interface{}, error) {
	data, err := LoadResponse(subdir, filename)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s/%s: %w", subdir, filename, err)
	}
	return result, nil
}

func ResponsesDir(mode string) string {
	return filepath.Join(TestdataDir(), "responses", mode)
}
