/*
Name: Jabari Brown
Purpose: Find and display all the links within the fetched URL
*/

package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

func scrape(murl string) []string {
	temp := make([]string, 1)
	resp, err := http.Get(murl) //get the response and error
	if err != nil {             //if there is an error gettng the link
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}

	tokenizer := html.NewTokenizer(resp.Body) //Used to tokenize html

	for { //We iterate through the page
		ite := tokenizer.Next()

		switch {
		case ite == html.ErrorToken:
			// End of the document, we're done
			return temp
		case ite == html.StartTagToken:
			tok := tokenizer.Token()

			// Check if the token is an <a> tag
			isAnchor := tok.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, url := getHref(tok)
			if !ok {
				continue
			}

			// Handling relative links
			hasProto := strings.Index(url, "/") == 0 || strings.Index(url, "#") == 0
			if hasProto {
				//add to splice
				temp = append(temp, murl+url)
			} else {
				temp = append(temp, url)
			}
		}
	}
}

func main() {
	urls := make([]string, 1)
	for _, links := range os.Args[1:] { //gets all the arguments at runtime
		urls = scrape(links) //scraping
	}

	for _, items := range urls { //printing
		fmt.Println(items)
	}

}
