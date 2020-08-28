package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/mecenat/solr"
)

// TODO (KK): sanely fix this shit
func main() {
	ctx := context.Background()
	// Initialize a new solr schema API
	ca, err := solr.NewCoreAdmin(ctx, "http://localhost:8983", http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}

	res, err := ca.Create(ctx, "films", &solr.CoreCreateOpts{
		InstanceDir: "/var/solr/data/films",
		DataDir:     "data",
		Config:      "conf/solrconfig.xml",
	})
	fmt.Println(err)
	fmt.Println(res.Header)

	res, err = ca.Status(ctx, "", false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Status["films"].Name, res.Status["films"].Uptime, res.Status["films"].Index.NumDocs)

	res, err = ca.Create(ctx, "actors", &solr.CoreCreateOpts{
		InstanceDir: "/var/solr/data/actors",
		DataDir:     "data",
		Config:      "conf/solrconfig.xml",
		// AsyncID:     "actorCreate",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ca.Swap(ctx, "actors", "films", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ca.Rename(ctx, "films", "filmer", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ca.Split(ctx, "actors", &solr.CoreSplitOpts{
		TargetCore: []string{"films", "films2"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ca.Merge(ctx, "actors", &solr.CoreMergeOpts{
		SrcCore: []string{"films", "films2"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ca.Reload(ctx, "actors")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ca.Unload(ctx, "actors", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ca.RequestStatus(ctx, "actorCreate")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.ReqStatus)
}
