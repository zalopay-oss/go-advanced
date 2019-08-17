package main

import (
  "fmt"
  "regexp"
  "time"

  "github.com/gocolly/colly"
)

var visited = map[string]bool{}

func main() {
  // Instantiate default collector
  c := colly.NewCollector(
    colly.AllowedDomains("www.abcdefg.com"),
    colly.MaxDepth(1),
  )

  // We think the matching page is the details page of the site
  detailRegex, _ := regexp.Compile(`/go/go\?p=\d+$`)
  // Matching the following pattern is the list page of the site
  listRegex, _ := regexp.Compile(`/t/\d+#\w+`)

  // All a tags, set the callback function
  c.OnHTML("a[href]", func(e *colly.HTMLElement) {
    link := e.Attr("href")

    // Visited details page or list page, skipped
    if visited[link] && (detailRegex.Match([]byte(link)) || listRegex.Match([]byte(link))) {
      return
    }

    // is neither a list page nor a detail page
    // Then it's not what we care about, skip it
    if !detailRegex.Match([]byte(link)) && !listRegex.Match([]byte(link)) {
      println("not match", link)
      return
    }

    // Because most websites have anti-reptile strategies
    // So there should be sleep logic in the crawler logic to avoid being blocked
    time.Sleep(time.Second)
    println("match", link)

    visited[link] = true

    time.Sleep(time.Millisecond * 2)
    c.Visit(e.Request.AbsoluteURL(link))
  })

  err := c.Visit("https://www.abcdefg.com/go/go")
  if err != nil {fmt.Println(err)}
}