package solr

import "encoding/json"

type Response struct {
	Header *ResponseHeader         `json:"responseHeader"`
	Data   *ResponseData           `json:"response"`
	Error  *ResponseError          `json:"error"`
	Debug  *map[string]interface{} `json:"debug"`
	Doc    *Doc                    `json:"doc"`
}

type ResponseHeader struct {
	Status int64  `json:"status"`
	QTime  int64  `json:"QTime"`
	Params *Query `json:"params"`
}

type ResponseData struct {
	NumFound int64 `json:"numFound"`
	Start    int64 `json:"start"`
	Docs     Docs  `json:"docs"`
}

type ResponseError struct {
	Code    int64    `json:"code"`
	Message string   `json:"msg"`
	Meta    []string `json:"metadata"`
}

func (r *ResponseError) Error() string {
	return r.Message
}

type Docs []*Doc
type Doc map[string]interface{}

func (d Docs) ToBytes() ([]byte, error) {
	return interfaceToBytes(d)
}

func (d *Doc) ToBytes() ([]byte, error) {
	return interfaceToBytes(d)
}

func interfaceToBytes(a interface{}) ([]byte, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return b, err
}
