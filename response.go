package solr

// Response represents the response from the solr server. It usually contains
// Header information, the response data or an error in case of erroneous
// response. Also it can contain Debug information when requested, a
// single document (in the case of realtimeGet) or just a status
// (in the case of the Ping request)
type Response struct {
	Header *ResponseHeader         `json:"responseHeader"`
	Data   *ResponseData           `json:"response"`
	Error  *ResponseError          `json:"error"`
	Debug  *map[string]interface{} `json:"debug"`
	Doc    *Doc                    `json:"doc"`
	Status *string                 `json:"status"`
}

// ResponseHeader is populated on every response from the solr server
// unless explicitly omitted. It contains the request status code
// the time it took as well as the params for the search query
// when applicable
type ResponseHeader struct {
	Status int64  `json:"status"`
	QTime  int64  `json:"QTime"`
	Params *Query `json:"params"`
}

// ResponseData is populated on a successful response from the solr
// server. It contains the number of documents found, the starting
// index (in case of a search) as well as the documents found
type ResponseData struct {
	NumFound int64 `json:"numFound"`
	Start    int64 `json:"start"`
	Docs     Docs  `json:"docs"`
}

// ResponseError is populated in the event the response from the solr
// server is erroneous. It contains the status code, a message
// and some metadata about the error's class
type ResponseError struct {
	Code    int64    `json:"code"`
	Message string   `json:"msg"`
	Meta    []string `json:"metadata"`
}

func (r *ResponseError) Error() string {
	return r.Message
}

// Docs represents an array of doc
type Docs []*Doc

// Doc is essentialy a map[string]interface{}
type Doc map[string]interface{}

// ToBytes returs a byte slice to simplify unmarshaling to JSON
func (d Docs) ToBytes() ([]byte, error) {
	return interfaceToBytes(d)
}

// ToBytes returs a byte slice to simplify unmarshaling to JSON
func (d *Doc) ToBytes() ([]byte, error) {
	return interfaceToBytes(d)
}
