# solr
A Solr client written in Go

Designed for Solr 8.5 (Should support earlier versions as well)

Provides clients for Solr's Request API, Schema API & Core Admin API 

Currently supports only JSON and basic CRUDL actions.


## Installation
```
go get -u github.com/mecenat/solr
```

## Usage

To create a new solr Client you need first to create a Connection. To create a Connection you need the host location (e.g. http://localhost:8983), the core name, and a http client (e.g. http.DefaultClient). Sending your own client could be useful when you need to wrap the client with another service, for example if you want to use AWS's X-Ray service to trace your API's calls.

When using a single server:
```
package main

import "github.com/mecenat/solr"

func main() {
	conn, err := solr.NewConnection("host", "core", http.DefaultClient)
	if err != nil {
				...
	}
	slr, err := solr.NewSingleClient(conn)
	if err != nil {
	      ...
	}
```
When using a the Primary-Replica paradigm:
```
package main

import "github.com/mecenat/solr"

func main() {
	primaryConn, err := solr.NewConnection("primaryHost", "core", http.DefaultClient)
	if err != nil {
				...
	}
	replicaConn, err := solr.NewConnection("replicaHost", "core", http.DefaultClient)
	if err != nil {
				...
	}
	slr, err := solr.NewPrimaryReplicaClient(primaryConn, replicaConn)
	if err != nil {
	      ...
	}
```

Aside from the normal Connection you can you a RetryableConnection which implements [Hashicorp's retryable HttpClient](https://github.com/hashicorp/go-retryablehttp) specifying the max timeout and provide that connection to the clients.
```
package main

import "github.com/mecenat/solr"

func main() {
	retConn, err := solr.NewRetryableConnection("host", "core", http.DefaultClient, 500*time.Millisecond)
	if err != nil {
				...
	}
	slr, err := solr.NewSingleClient(retConn)
	if err != nil {
	      ...
	}
```

To access Solr's Core Admin API you need to create a separate client as follows:
```
package main

import "github.com/mecenat/solr"

func main() {
	ctx := context.Background()
	ca, err := solr.NewCoreAdmin(ctx, "host", http.DefaultClient)
	if err != nil {
				...
	}
```

To access Solr's Schema API you also need a separate client as follows:
```
package main

import "github.com/mecenat/solr"

func main() {
	ctx := context.Background()
	sa, err := solr.NewSchemaAPI(ctx, "host", "core", http.DefaultClient)
	if err != nil {
				...
	}
```

## License
MIT