package solr

type SolrClient interface {
	Ping() (int, error)
	Search(q *Query) (*SearchResponse, error)
	Get(ids []string, filter string) (*GetResponse, error)
}
