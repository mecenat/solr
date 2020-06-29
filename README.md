# solr
A Solr client written in Go

Designed for Solr 8.5 (Should support earlier versions as well)

Currently supports only JSON and basic CRUDL actions.


## Installation
```
go get -u github.com/mecenat/solr
```

## Usage

To create a new solr Client you need the host location (e.g. http://localhost:8983), the core name, and a http client (e.g. http.DefaultClient). Sending your own client could be useful when you need to wrap the client with another service, for example if you want to use AWS's X-Ray service to trace your API's calls.

When using a single server:
```
package main

import "github.com/mecenat/solr"

func main() {
	ctx := context.Background()
	slr, err := solr.NewSingleClient(ctx, "host", "core", http.DefaultClient)
	if err != nil {
	      ...
	}
```
When using a the Primary-Replica paradigm:
```
package main

import "github.com/mecenat/solr"

func main() {
	ctx := context.Background()
	slr, err := solr.NewPrimaryReplicaClient(ctx, "primaryHost", "primaryCore", "replicaHost", "replicaCore", primaryClient, replicatClient)
	if err != nil {
	      ...
	}
```

## License
MIT