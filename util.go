package solr

import (
	"encoding/json"
	"fmt"
	"strings"
)

func formatBasePath(host, core string) string {
	if strings.HasSuffix(host, "/solr") {
		return fmt.Sprintf("%s/%s", host, core)
	}
	return fmt.Sprintf("%s/solr/%s", host, core)
}

func isJSON(input []byte) error {
	var js map[string]interface{}
	return json.Unmarshal(input, &js)
}

func isArrayOfJSON(input []byte) error {
	var js []*map[string]interface{}
	return json.Unmarshal(input, &js)
}
