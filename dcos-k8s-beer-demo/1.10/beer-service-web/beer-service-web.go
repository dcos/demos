package main

import (
    "encoding/json"
    "fmt"
    "html/template"
    "net/http"
    "os"
)

var beersPage = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Beer demo</title>
    <style>
	body {background-color: powderblue;}
	h1   {color: blue;}
	h2   {color: green;}
	p    {color: black;}
    </style>
</head>
<body>
    <h1>Beer: {{.BeerName}}</h1>
    <h2>Style: {{.BeerStyle}}</h2>
    <h3>Description:</p/n>
    <h4>{{.BeerDescription}}</p>
</body>
</html>`

type Beer struct {
    BeerName        string
    BeerStyle       string
    BeerDescription string
}

func handler(w http.ResponseWriter, r *http.Request) {
  res, err := http.Get(os.Getenv("BEER_URL"))

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer res.Body.Close()
	b := &Beer{}
	// decoding request body
	if err := json.NewDecoder(res.Body).Decode(b); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl, err := template.New("beers").Parse(beersPage)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = tmpl.Execute(w, b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Listening on: http://localhost:8090")
	http.ListenAndServe(":8090", nil)
}
