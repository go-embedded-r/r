package oracle

import (
	"encoding/json"
	"os"
	"strings"
)

type goldenRecord struct {
	Src string `json:"src"`
	Hex string `json:"hex"`
}

// readGolden reads the fields of rds.jsonl this package's tests use. (The gen
// package cannot be imported here: it imports oracle.)
func readGolden(path string) ([]goldenRecord, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []goldenRecord
	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		var r goldenRecord
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}
