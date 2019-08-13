package main

import (
  "net/url"
	"time"
	"github.com/gocolly/colly"
	nats "github.com/nats-io/nats.go"
	"os"
	"regexp"
	"fmt"
)

var domain2Collector = map[string]*colly.Collector{}
var nc *nats.Conn
var maxDepth = 10
var natsURL = "nats://localhost:4222"

func factory(urlStr string) *colly.Collector {
  u, _ := url.Parse(urlStr)
  return domain2Collector[u.Host]
}

func initABCDECollector() *colly.Collector {
  c := colly.NewCollector(
    colly.AllowedDomains("www.abcdefg.com"),
    colly.MaxDepth(maxDepth),
  )

  c.OnResponse(func(resp *colly.Response) {
    // Do some aftercare work after climbing
    // For example, the confirmation that the page has been crawled is stored in MySQL.
  })

	detailRegex, _ := regexp.Compile(`/go/go\?p=\d+$`)
	listRegex, _ := regexp.Compile(`/t/\d+#\w+`)

  c.OnHTML("a[href]", func(e *colly.HTMLElement) {
    // Basic anti-reptile strategy
    link := e.Attr("href")
    time.Sleep(time.Second * 2)

    // regular match list page, then visit
    if listRegex.Match([]byte(link)) {
      c.Visit(e.Request.AbsoluteURL(link))
    }
    // Regular match landing page, send message queue
    if detailRegex.Match([]byte(link)) {
			err := nc.Publish("tasks", []byte(link))
			fmt.Println(err)
      nc.Flush()
    }
  })
  return c
}

func initHIJKLCollector() *colly.Collector {
  c := colly.NewCollector(
    colly.AllowedDomains("www.hijklmn.com"),
    colly.MaxDepth(maxDepth),
  )

  c.OnHTML("a[href]", func(e *colly.HTMLElement) {
  })

  return c
}

func init() {
  domain2Collector["www.abcdefg.com"] = initABCDECollector()
  domain2Collector["www.hijklmn.com"] = initHIJKLCollector()

  var err error
  nc, err = nats.Connect(natsURL)
  if err != nil {os.Exit(1)}
}

func main() {
  urls := []string{"https://www.abcdefg.com", "https://www.hijklmn.com"}
  for _, url := range urls {
    instance := factory(url)
    instance.Visit(url)
  }
}