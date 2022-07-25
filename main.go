package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type UrlDatabase struct {
	ShortURL      string    `json:"short_url"`
	LongURL       string    `json:"long_url"`
	CreateTime    time.Time `json:"create_time"`
	RedirectCount int32     `json:"redirect_count"`
}

var Urls map[string]*UrlDatabase

const Alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

// to rand the string
func makeRandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = Alphabet[rand.Intn(len(Alphabet))]
	}
	return string(b)
}

func shorten(w http.ResponseWriter, r *http.Request) {
	// read request url
	url, ok := r.URL.Query()["url"]
	// validate not empty url
	if !ok || len(url) == 0 {
		fmt.Println("Empty URL")
		return
	}

	// generate unique url
	if !isValidURL(url[0]) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "URL is not valid")
		return
	}

	var linkString, randString string

	createTime := time.Now()
	if dataURL, ok := Urls[url[0]]; !ok {
		randString = makeRandString(6)
		for ok {
			if _, ok := Urls[randString]; !ok {
				randString = makeRandString(6)
			}
		}
		Urls[randString] = &UrlDatabase{
			LongURL:       url[0],
			CreateTime:    createTime,
			RedirectCount: 0,
		}
		Urls[url[0]] = &UrlDatabase{
			LongURL:       url[0],
			ShortURL:      randString,
			CreateTime:    createTime,
			RedirectCount: 0,
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusAccepted)
		linkString = fmt.Sprintf("<a href=\"http://localhost:8090/url/%s\">http://localhost:8090/url/%s</a>", randString, randString)
		fmt.Fprintf(w, "Successfully shorten the link\n")
	} else {
		randString = dataURL.ShortURL
		linkString = fmt.Sprintf("<a href=\"http://localhost:8090/url/%s\">http://localhost:8090/url/%s</a>", randString, randString)

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "Same link has been added before\n")
	}

	// return unique url, create time, redirect count
	fmt.Fprintln(w, "Short URL:", linkString)
	fmt.Fprintln(w, "Long URL:", Urls[randString].LongURL)
	fmt.Fprintln(w, "Created Time:", Urls[randString].CreateTime)
	fmt.Fprintln(w, "Redirect Count:", Urls[randString].RedirectCount)
	return
}

// getShortenedURL will get the corresponding url in map
func getShortenedURL(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	pathArgs := strings.Split(path, "/")
	if len(pathArgs) < 2 {
		fmt.Fprintf(w, "No URL input")
		return
	}
	if dataURL, ok := Urls[pathArgs[2]]; ok {
		if dataLongURL, ok := Urls[dataURL.LongURL]; ok {
			dataURL.RedirectCount = dataURL.RedirectCount + 1
			dataLongURL.RedirectCount = dataLongURL.RedirectCount + 1
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "URL has never been added before\n")
			return
		}
		http.Redirect(w, r, dataURL.LongURL, http.StatusPermanentRedirect)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "URL has never been added before\n")
	return
}

// isValidURL will check if URL is valid
func isValidURL(url string) bool {
	r, err := regexp.Compile("^(http|https)://")
	if err != nil {
		return false
	}
	url = strings.TrimSpace(url)
	// Check if string matches the regex
	if r.MatchString(url) {
		return true
	}
	return false
}

// just for testing
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Endpoint Hit: Index...")
}

func handleRequests() {
	http.HandleFunc("/", index)
	http.HandleFunc("/shorten", shorten)
	http.HandleFunc("/url/", getShortenedURL)

	log.Fatal(http.ListenAndServe(":8090", nil))
}

func main() {
	Urls = map[string]*UrlDatabase{}
	handleRequests()
}
