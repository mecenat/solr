package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mecenat/solr"
)

// XXX: just for testing purposes, not part of actual repo
func main() {
	slr := solr.New("http://localhost:8983", "films", http.DefaultClient)
	err := slr.Ping()
	if err != nil {
		log.Fatal(err)
	}

	opts := &solr.QueryOptions{Rows: 1}
	q := solr.NewQuery(opts)
	q.SetQuery("*:*")
	// q.SetOperation("and")
	// q.AddFilter("published", "true")
	// q.AddFilter("deleted", "false")
	// q.SetSort("created_at desc")
	fmt.Println(q.String())
	_, err = slr.Search(q)
	if err != nil {
		log.Fatal(err)
	}

	q.SetStart(1)
	fmt.Println(q.String())
	_, err = slr.Search(q)
	if err != nil {
		log.Fatal(err)
	}

	// ffs := false
	// res3, err := slr.Create(films[1], &solr.WriteOptions{Commit: true})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(res3)

	// res5, err := slr.BatchCreate(films)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(res5)

	// res6, err := slr.DeleteByQuery("*:*")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(res6)

	// res7, err := slr.DeleteByID("c0589192-a996-11ea-b456-acde48001122")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(res7)

	// res4, err := slr.Commit()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(res4)

	// 	res2, err := slr.Get(films[0].ID)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	var f2 film
	// 	fby, err := res2.Doc.ToBytes()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	fmt.Println(string(fby))

	// 	err = json.Unmarshal(fby, &f2)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	fmt.Println(f2.Name)
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
}
