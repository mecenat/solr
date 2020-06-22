package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mecenat/solr"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	slr := solr.New("http://localhost:8983", "films", http.DefaultClient)
	err := slr.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create a document entry for the first film of the array
	res, err := slr.Create(ctx, films[0], &solr.WriteOptions{Commit: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// Update the genre attribute
	uf := solr.NewUpdate("1")
	uf.Add("genre", "Slasher")
	uf.Add("directed_by", "some guy")
	res, err = slr.Update(ctx, uf, &solr.WriteOptions{Commit: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// Then delete that document
	res, err = slr.DeleteByID(ctx, "1", &solr.WriteOptions{Commit: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// Create a document entry for the first film of the array
	// but without commiting anything
	res, err = slr.Create(ctx, films[0], nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// RealTime Get the film that was just added
	res, err = slr.Get(ctx, films[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Doc)

	// rollback the latest changes
	res, err = slr.Rollback(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	// RealTime Get the film that was just added
	res, err = slr.Get(ctx, films[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	// should be null since we rollbacked
	fmt.Println(res.Doc)

	// insert all the documents at once
	res, err = slr.BatchCreate(ctx, films, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// RealTime Get the film that was just added
	res, err = slr.BatchGet(ctx, []string{films[0].ID, films[3].ID}, "")
	if err != nil {
		log.Fatal(err)
	}
	// should be 2
	fmt.Println(res.Data.NumFound)

	// commit the latest changes
	res, err = slr.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	// query for all items and filter only to horror films
	opts := &solr.ReadOptions{Rows: 10}
	q := solr.NewQuery(opts)
	q.SetQuery("*:*")
	q.AddFilter("genre", "horror")
	q.SetSort("year desc")
	fmt.Println(q.String())
	res, err = slr.Search(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Data.NumFound)
	fmt.Println(len(res.Data.Docs))

	// the original query had a limit to 10 rows, therefore we
	// can ask solr for the remainder rows by setting the start
	// param of the already existing query
	q.SetStart(10)
	fmt.Println(q.String())
	res, err = slr.Search(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Data.NumFound)
	fmt.Println(len(res.Data.Docs))

	// we can delete all horror films by specifying a query
	// for deletion
	res, err = slr.DeleteByQuery(ctx, "genre:horror", &solr.WriteOptions{Commit: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Header)

	// we query for all items
	q2 := solr.NewQuery(opts)
	q2.SetQuery("*:*")
	res, err = slr.Search(ctx, q2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Data.NumFound)
	fmt.Println(len(res.Data.Docs))

	remainingFilms := []*film{}
	resDocBytes, err := res.Data.Docs.ToBytes()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(resDocBytes, &remainingFilms)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(remainingFilms))
	fmt.Println(remainingFilms[0].Name)

	// the delete everything in films
	res, err = slr.Clear(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

type film struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Year     string   `json:"year"`
	Genre    []string `json:"genre"`
	Director []string `json:"directed_by"`
}

var films = []*film{
	{
		ID:       "1",
		Name:     "Lock, Stock and Two Smoking Barrels",
		Year:     "1998",
		Genre:    []string{"Crime", "Comedy", "Action"},
		Director: []string{"Guy Ritchie"},
	},
	{
		ID:       "2",
		Name:     "Reservoir Dogs",
		Year:     "1992",
		Genre:    []string{"Crime", "Drama", "Thriller"},
		Director: []string{"Quentin Tarantino"},
	},
	{
		ID:       "3",
		Name:     "Memento",
		Year:     "2000",
		Genre:    []string{"Mystery", "Thriller"},
		Director: []string{"Christopher Nolan"},
	},
	{
		ID:       "4",
		Name:     "Night of the Living Dead",
		Year:     "1968",
		Genre:    []string{"Horror"},
		Director: []string{"George A. Romero"},
	},
	{
		ID:       "5",
		Name:     "The Return of the Living Dead",
		Year:     "1985",
		Genre:    []string{"Horror", "Comedy", "Sci-Fi"},
		Director: []string{},
	},
	{
		ID:       "6",
		Name:     "The Evil Dead",
		Year:     "1981",
		Genre:    []string{"Horror"},
		Director: []string{"Sam Raimi"},
	},
	{
		ID:       "7",
		Name:     "Alien",
		Year:     "1979",
		Genre:    []string{"Horror", "Sci-Fi"},
		Director: []string{"Ridley Scott"},
	},
	{
		ID:       "8",
		Name:     "The Shining",
		Year:     "1980",
		Genre:    []string{"Drama", "Horror"},
		Director: []string{"Stanley Kubrick"},
	},
	{
		ID:       "9",
		Name:     "The Host",
		Year:     "2006",
		Genre:    []string{"Action", "Drama", "Horror"},
		Director: []string{"Bong Joon Ho"},
	},
	{
		ID:       "10",
		Name:     "The Grudge",
		Year:     "2004",
		Genre:    []string{"Mystery", "Thriller", "Horror"},
		Director: []string{"Takashi Shimizu"},
	},
	{
		ID:       "11",
		Name:     "The Thing",
		Year:     "1982",
		Genre:    []string{"Mystery", "Sci-Fi", "Horror"},
		Director: []string{"John Carpenter"},
	},
	{
		ID:       "12",
		Name:     "Låt den Rätte Komma in",
		Year:     "2008",
		Genre:    []string{"Drama", "Romance", "Horror"},
		Director: []string{"Tomas Alfredson"},
	},
	{
		ID:       "13",
		Name:     "REC",
		Year:     "2007",
		Genre:    []string{"Action", "Adventure", "Fantasy"},
		Director: []string{"Jaume Balagueró", "Paco Plaza"},
	},
	{
		ID:       "14",
		Name:     "Evil",
		Year:     "2005",
		Genre:    []string{"Action", "Comedy", "Horror"},
		Director: []string{"Giorgos Nousias"},
	},
	{
		ID:       "15",
		Name:     "Zombi 2",
		Year:     "1979",
		Genre:    []string{"Horror"},
		Director: []string{"Lucio Fulci"},
	},
	{
		ID:       "16",
		Name:     "Shaun of the Dead",
		Year:     "2004",
		Genre:    []string{"Comedy", "Horror"},
		Director: []string{"Edgar Wright"},
	},
	{
		ID:       "17",
		Name:     "Død Snø",
		Year:     "2009",
		Genre:    []string{"Comedy", "Horror"},
		Director: []string{"Tommy Wirkola"},
	},
	{
		ID:       "18",
		Name:     "Zombeavers",
		Year:     "2014",
		Genre:    []string{"Comedy", "Horror"},
		Director: []string{"Jordan Rubin"},
	},
	{
		ID:       "19",
		Name:     "Killdozer",
		Year:     "1974",
		Genre:    []string{"Sci-Fi", "Horror"},
		Director: []string{"Jerry London"},
	},
	{
		ID:       "20",
		Name:     "Busanhaeng",
		Year:     "2016",
		Genre:    []string{"Action", "Horror", "Thriller"},
		Director: []string{"Song-ho Yeon"},
	},
	{
		ID:       "21",
		Name:     "The House of 1000 Corpses",
		Year:     "2003",
		Genre:    []string{"Horror"},
		Director: []string{"Rob Zombie"},
	},
}
