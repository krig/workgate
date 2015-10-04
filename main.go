package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"html/template"
	"io/ioutil"
	"time"
)

// relevant feeds:
// http://git.haproxy.org/?p=haproxy-1.5.git;a=atom
// https://github.com/ClusterLabs/crmsh/commits/master.atom
// https://github.com/ClusterLabs/resource-agents/commits/master.atom
// https://github.com/ClusterLabs/hawk/commits/master.atom
// https://github.com/ClusterLabs/pacemaker/commits/master.atom
// https://github.com/ClusterLabs/fence-agents/commits/master.atom

type Page struct {
	Title string
	Time string
}

type Feed struct {
	Name string
	Atom string
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		// handle error
	}
	p := &Page{
		Title: "Workdash",
		Time: time.Now().Format(time.RFC3339),
	}
	t.Execute(w, p)
}

var feeds []Feed = []Feed{
	Feed{
		Name: "haproxy",
		Atom: "http://git.haproxy.org/?p=haproxy-1.5.git;a=atom",
	},
	Feed{
		Name: "crmsh",
		Atom: "https://github.com/ClusterLabs/crmsh/commits/master.atom",
	},
	Feed{
		Name: "hawk",
		Atom: "https://github.com/ClusterLabs/hawk/commits/master.atom",
	},
	Feed{
		Name: "resource-agents",
		Atom: "https://github.com/ClusterLabs/resource-agents/commits/master.atom",
	},
	Feed{
		Name: "pacemaker",
		Atom: "https://github.com/ClusterLabs/pacemaker/commits/master.atom",
	},
	Feed{
		Name: "fence-agents",
		Atom: "https://github.com/ClusterLabs/fence-agents/commits/master.atom",
	},
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	requestedFeed := r.URL.Path[len("/feed/"):]
	for i := range feeds {
		if requestedFeed == feeds[i].Name {
			resp, err := http.Get(feeds[i].Atom)
			if err != nil {
				// handle error
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// handle error
			}
			w.Write(body)
			resp.Body.Close()
			break
		}
	}
}

// compile scss files here..
func main() {
	port := flag.String("p", "8080", "the port to bind on (ports below 1024 require root permissions)")
	flag.Parse()
	fmt.Printf("http://localhost:%s\n", *port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", viewHandler)
	mux.HandleFunc("/feed/", feedHandler)
	mux.Handle("/js/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":"+*port, mux))
}
