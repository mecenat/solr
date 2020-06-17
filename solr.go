package solr

type SolrClient interface {
	Ping() (*Response, error)
	Search(q *Query) (*Response, error)
	Get(id string) (*Response, error)
	BatchGet(ids []string, filter string) (*Response, error)
	Create(item []byte) (*Response, error)
}
