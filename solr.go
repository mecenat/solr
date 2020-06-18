package solr

type SolrClient interface {
	Ping() error
	Search(q *Query) (*Response, error)
	Get(id string) (*Response, error)
	BatchGet(ids []string, filter string) (*Response, error)
	Create(item interface{}, opts *WriteOptions) (*Response, error)
	BatchCreate(items interface{}) (*Response, error)
	DeleteByID(id string) (*Response, error)
	DeleteByQuery(query string) (*Response, error)
	Commit() (*Response, error)
}
