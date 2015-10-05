package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"html/template"
	"io/ioutil"
	"time"
	"encoding/json"
	"bufio"
	"regexp"
	"os"
)


type Page struct {
	Title string
	Time string
	Links template.HTML
}

type Feed struct {
	Name string
	Atom string
}

type FeedData struct {
	Updated time.Time
	Data []byte
}


var indexTemplate *template.Template
var indexPage *Page = &Page{Title: "Work Dashboard", Time: time.Now().Format(time.RFC3339)}
var feeds []Feed = []Feed{}
var feedData map[string]FeedData = make(map[string]FeedData)


func initWorkgate() {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	}
	indexTemplate = t

	f, err := os.Open("feeds.txt")
	if err != nil {
		log.Fatal(err)
	}

	rx, err := regexp.Compile("([a-zA-Z0-9-_]+) = (.*)")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		match := rx.FindStringSubmatch(line)
		if match != nil {
			feeds = append(feeds, Feed{Name: match[1], Atom: match[2]})
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
	f.Close()

	linkPattern := "<li><a href=\"%s\">%s</a></li>\n"
	links := ""

	f, err = os.Open("links.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		match := rx.FindStringSubmatch(line)
		if match != nil {
			links += fmt.Sprintf(linkPattern, match[2], match[1])
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
	f.Close()

	indexPage.Links = template.HTML(links)

}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	indexPage.Time = time.Now().Format(time.RFC3339)
	indexTemplate.Execute(w, indexPage)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	requestedFeed := r.URL.Path[len("/feed/"):]
	if requestedFeed == "list" {
		b, err := json.Marshal(feeds)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(b)
	} else {
		data, ok := feedData[requestedFeed]
		if ok {
			distance := time.Now().Sub(data.Updated)
			if distance.Minutes() > 15 {
				ok = false
			} else {
				fmt.Printf("Cached %s\n", requestedFeed)
				w.Write(data.Data)
			}
		}
		if !ok {
			for i := range feeds {
				if requestedFeed == feeds[i].Name {
					fmt.Printf("Fetching %s\n", feeds[i].Atom)
					resp, err := http.Get(feeds[i].Atom)
					if err != nil {
						log.Fatal(err)
					}
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Fatal(err)
					}
					w.Write(body)
					resp.Body.Close()
					feedData[requestedFeed] = FeedData{
						Updated: time.Now(),
						Data: body,
					}
					break
				}
			}
		}
	}
}

func main() {
	port := flag.String("p", "8080", "the port to bind on (ports below 1024 require root permissions)")
	flag.Parse()
	fmt.Printf("http://localhost:%s\n", *port)
	initWorkgate()
	mux := http.NewServeMux()
	mux.HandleFunc("/", viewHandler)
	mux.HandleFunc("/feed/", feedHandler)
	mux.Handle("/js/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":"+*port, mux))
}
