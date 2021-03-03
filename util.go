package solr

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidConfig is returned when the hostname or corename are empty
var ErrInvalidConfig = errors.New("invalid configuration: no host or core provided")

func formatBasePath(host, core string) string {
	if strings.HasSuffix(host, "/solr") {
		return fmt.Sprintf("%s/%s", host, core)
	}
	return fmt.Sprintf("%s/solr/%s", host, core)
}

func formatDocEntry(doc Doc) map[string]interface{} {
	return map[string]interface{}{"doc": doc}
}

func formatDeleteByID(id string) Doc {
	return Doc{"id": id}
}

func formatDeleteByQuery(query string) Doc {
	return Doc{"query": query}
}

func isJSON(input []byte) error {
	var js map[string]interface{}
	return json.Unmarshal(input, &js)
}

func isArrayOfJSON(input []byte) error {
	var js []*map[string]interface{}
	return json.Unmarshal(input, &js)
}

func interfaceToBytes(a interface{}) ([]byte, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return b, err
}

// BoostField is a helper function to properly format field boosting
func BoostField(field string, boost float64) string {
	return fmt.Sprintf("%s^%f", field, boost)
}
