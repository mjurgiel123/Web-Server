package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// struct for first JSON (names API)
type Name struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// struct for second JSON (jokes API)
type Joker struct {
	Type  string `json:"type"`
	Value struct {
		Categories []string `json:"categories"`
		ID         int      `json:"id"`
		Joke       string   `json:"joke"`
	} `json:"value"`
}

// writes joke onto HTTP Response
func JokeHandler(w http.ResponseWriter, r *http.Request) {
	person := retrieveName(w, r)                                          // retrieves a struct with first/last name
	joke := retrieveJoke(w, r)                                            // retrieves joke as a string
	finalOutput := strings.Replace(joke, "John", person.FirstName, 1)     // replaces "John" with person.FirstName
	finalOutput = strings.Replace(finalOutput, "Doe", person.LastName, 1) // replaces "Doe" with person.LastName
	w.Write([]byte(finalOutput))
}

// retrieves Name struct
func retrieveName(w http.ResponseWriter, r *http.Request) Name {
	resp, err := http.Get("https://names.mcquay.me/api/v0/") // retrieves data from Name API
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body) // converts resp (*http.Response) to bytes
	if err != nil {
		log.Fatal(err)
	}

	person := Name{} // convert JSON to Name struct
	err = json.Unmarshal(bytes, &person)
	if err != nil {
		log.Fatal(err)
	}

	return person
}

// retrieves Joke as a string
func retrieveJoke(w http.ResponseWriter, r *http.Request) string {
	resp, err := http.Get("http://joke.loc8u.com:8888/joke?limitTo=nerdy&firstName=John&lastName=Doe") //retrieves data from Joke API
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body) // converts resp (*http.Response) to JSON (bytes)
	if err != nil {
		log.Fatal(err)
	}

	components := Joker{} // converts JSON to Joker Struct
	err = json.Unmarshal(bytes, &components)
	if err != nil {
		log.Fatal(err)
	}

	joke := components.Value.Joke // retrieves Joke (String) from components

	return joke
}

func main() {
	// creates mux for resource path: /joke
	mux := http.NewServeMux()
	mux.HandleFunc("/joke", JokeHandler)

	// listens to port 5000
	addr := ":5000"

	// runs until program is aborted
	log.Fatal(http.ListenAndServe(addr, mux))
}
