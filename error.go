package solr

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ErrorDetail is an interface to interpret the details of an error. Solr
// tends to be inconsistent about the type of the detail, therefore
// an interface is needed to cover all possible scenarios.
type ErrorDetail interface {
	String() string
	Item() map[string]interface{}
}

// ErrorDetailObj provides detailed information on the errors
// that might arise when multiple commands are sent in a
// batch.
type ErrorDetailObj struct {
	Messages    []string               `json:"errorMessages"`
	Command     string                 `json:"command"`
	CommandItem map[string]interface{} `json:"item"`
}

func (d *ErrorDetailObj) String() string {
	return fmt.Sprintf("%s: %s", d.Command, d.Messages)
}

// Item returns the item causing the error
func (d *ErrorDetailObj) Item() map[string]interface{} {
	return d.CommandItem
}

// ErrorDetailString provides information about the details of the
// error.
type ErrorDetailString string

func (d *ErrorDetailString) String() string {
	return string(*d)
}

// Item returns an empty map here.
func (d *ErrorDetailString) Item() map[string]interface{} {
	return map[string]interface{}{}
}

// ResponseError is populated in the event the response from the solr
// server is erroneous. It contains the status code, a message
// and some metadata about the error's class
type ResponseError struct {
	Code    float64       `json:"code"`
	Message string        `json:"msg"`
	Meta    []string      `json:"metadata"`
	Details []ErrorDetail `json:"details"`
}

func (r *ResponseError) Error() string {
	if len(r.Details) > 0 {
		var msgs []string
		for _, detail := range r.Details {
			msgs = append(msgs, detail.String())
		}
		return fmt.Sprintf("%s: {%s}", r.Message, strings.Join(msgs, ", "))
	}
	return r.Message
}

// UnmarshalJSON implements the unmarshaler interface
func (r *ResponseError) UnmarshalJSON(b []byte) error {
	var temp map[string]interface{}
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	code, ok := temp["code"].(float64)
	if ok {
		r.Code = code
	}

	msg, ok := temp["msg"].(string)
	if ok {
		r.Message = msg
	}

	metadata, ok := temp["metadata"].([]interface{})
	if ok {
		for _, meta := range metadata {
			m, ok := meta.(string)
			if ok {
				r.Meta = append(r.Meta, m)
			}
		}
	}

	detailsArr, ok := temp["details"].([]interface{})
	if ok {
		for _, detailItem := range detailsArr {
			switch detail := detailItem.(type) {
			case string:
				dStr := ErrorDetailString(detail)
				r.Details = append(r.Details, &dStr)

			case map[string]interface{}:
				if len(detail) != 2 {
					return nil
				}
				var dObj ErrorDetailObj
				for key, val := range detail {
					if key == "errorMessages" {
						val, ok := val.([]interface{})
						if ok {
							for _, v := range val {
								msg, ok := v.(string)
								if ok {
									dObj.Messages = append(dObj.Messages, msg)
								}
							}
						}
						continue
					}
					dObj.Command = key
					item, ok := val.(map[string]interface{})
					if ok {
						dObj.CommandItem = item
					}
				}
				r.Details = append(r.Details, &dObj)
			}
		}
	}

	return nil
}
