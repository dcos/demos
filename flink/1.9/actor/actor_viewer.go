package main

import (
    "fmt"
    "net/http"
    "os"

    "gopkg.in/redis.v4"
)

func getKeys(w http.ResponseWriter) {
    // Connect to redis.
    client := redis.NewClient(&redis.Options{
        Addr:     "redis.marathon.l4lb.thisdcos.directory:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    pong, _ := client.Ping().Result()
    fmt.Fprintf(w, "Pong " + pong +"\n")
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<html>")
    fmt.Fprintf(w, "<title>Welcome to DC/OS 101!</title><body>")
    fmt.Fprintf(w, "<h1>Welcome to DC/OS 101!</h1>")
    fmt.Fprintf(w, "<h1>Running on node '"+  os.Getenv("HOST") + "' and port '" + os.Getenv("PORT0") + "' </h1>")
    fmt.Fprintf(w, "<h1>Add a new key:value pair</h1>"+
        "<form action=\"/save\" method=\"POST\">"+
        "<textarea name=\"key\">Key</textarea><br>"+
	"<textarea name=\"value\">Value</textarea><br>"+
        "<input type=\"submit\" value=\"Save\">"+
        "</form>")
   fmt.Fprintf(w, "</body></html>")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    key := r.FormValue("key")
    value := r.FormValue("value")
    fmt.Println("Key: "  + key)
    fmt.Println("Value :" + value)
    fmt.Fprintln(w, "Key: "  + key)
    fmt.Fprintln(w, "Value :" + value)

    client := redis.NewClient(&redis.Options{
        Addr:     "redis.marathon.l4lb.thisdcos.directory:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    err := client.Set(key, value, 0).Err()
       if err != nil {
          fmt.Println("Value :" + value)
       }

    // http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
    fmt.Println("Starting DC/OS-101 App ")
    http.HandleFunc("/", handler)
    http.HandleFunc("/save", saveHandler)
    http.ListenAndServe(":" + os.Getenv("PORT0") , nil)
}

