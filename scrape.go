package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gocolly/colly"
)

type article struct {
	Title  string `json:"title"`
	Url    string `json:"Url"`
	Score  string `json:"Score"`
	Poster string `json:"Poster"`
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	var articles []article
	var pageCount int
	i := 0

	selectors := "tbody > tr:nth-child(3) > td > table > tbody"
	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML(selectors, func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, h *colly.HTMLElement) {
			var art article
			title := h.ChildText("td.title > a")
			score := h.ChildText("td.subtext > span.score")
			if title == "More" {
				c.Visit("news.ycombinator.com/" + h.ChildAttr("td.title > a", "href"))
			} else if score != "" {
				articles[i].Score = score
				articles[i].Poster = h.ChildText("td.subtext > a.hnuser")
				i++
			} else if title != "" {
				art.Title = title
				art.Url = h.ChildAttr("td.title > a", "href")
				articles = append(articles, art)
			}
		})
	})

	c.OnResponse(func(r *colly.Response) {
		pageCount++
		urlVisited := r.Request.URL
		fmt.Println(fmt.Sprintf("%d  Finished Visiting : %s", pageCount, urlVisited))
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://news.ycombinator.com/")
	fmt.Println(articles[0])

	// Marshal instances of articles and conert to JSON
	articleJSON, _ := json.Marshal(articles)
	// fmt.Println(string(articleJSON))
	articleJSONString := string(articleJSON)
	FileWriter(articleJSONString)
}

func FileWriter(data string) error {
	file, err := os.Create("output.json")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}
