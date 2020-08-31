package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mecenat/solr"
)

func main() {
	ctx := context.Background()
	// Initialize a new solr Core Admin API
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

	reqID := fmt.Sprintf("actors%d", time.Now().Unix())

	res, err = ca.Create(ctx, "actors", &solr.CoreCreateOpts{
		InstanceDir: "/var/solr/data/actors",
		DataDir:     "data",
		Config:      "conf/solrconfig.xml",
		AsyncID:     reqID,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ca.Create(ctx, "films2", &solr.CoreCreateOpts{
		InstanceDir: "/var/solr/data/films2",
		DataDir:     "data",
		Config:      "conf/solrconfig.xml",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ca.Swap(ctx, "actors", "films2", "")
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

	res, err = ca.Rename(ctx, "films2", "filmer", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ca.Merge(ctx, "actors", &solr.CoreMergeOpts{
		SrcCore: []string{"films", "filmer"},
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

	res, err = ca.Unload(ctx, "filmer", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ca.Unload(ctx, "films", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ca.Unload(ctx, "actors", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	res, err = ca.RequestStatus(ctx, reqID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.ReqStatus)
}
