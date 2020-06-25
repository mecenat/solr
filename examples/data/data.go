package data

import "encoding/json"

// No example here, just the data to be used in the other examples
// For the examples to properly work you used create the correct
// fields in your solr core's schema. Solr's Dynamically generated
// fields tend to think everything is a multiValue and that doesn't
// sit well with the JSON Unmarshaler...

// Film ...
type Film struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Year        string   `json:"year"`
	Genre       []string `json:"genre"`
	Director    []string `json:"directed_by"`
	SeenCounter int      `json:"seen_counter"`
}

// ToMap ...
func (f *Film) ToMap() (doc map[string]interface{}, err error) {
	filmBytes, err := json.Marshal(f)
	if err != nil {
		return
	}
	err = json.Unmarshal(filmBytes, &doc)
	return
}

// Films ...
var Films = []*Film{
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
