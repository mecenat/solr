package solr

import "testing"

func TestIsJSONInvalid(t *testing.T) {
	input := "key"
	err := isJSON([]byte(input))
	if err == nil {
		t.Fatal("did not get error while input is invalid")
	}
}

func TestIsJSONValid(t *testing.T) {
	input := `{"key": "value"}`
	err := isJSON([]byte(input))
	if err != nil {
		t.Fatal("got error while input is valid")
	}
}

func TestIsArrayOfJSONInvalid(t *testing.T) {
	input := "key"
	err := isArrayOfJSON([]byte(input))
	if err == nil {
		t.Fatal("did not get error while input is invalid")
	}
}

func TestIsArrayOfJSONInvalid2(t *testing.T) {
	input := `{"key": "value"}`
	err := isArrayOfJSON([]byte(input))
	if err == nil {
		t.Fatal("did not get error while input is invalid")
	}
}

func TestIsArrayOfJSONValid(t *testing.T) {
	input := `[{"key": "value"}, {"key": "value"}]`
	err := isArrayOfJSON([]byte(input))
	if err != nil {
		t.Fatal("got error while input is valid")
	}
}
