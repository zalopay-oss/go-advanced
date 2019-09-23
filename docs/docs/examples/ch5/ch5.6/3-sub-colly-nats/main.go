package main

import (
  "net/url"
	"time"
	"github.com/gocolly/colly"
	nats "github.com/nats-io/nats.go"
	"os"
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
  return c
}
  
func initHIJKLCollector() *colly.Collector {
  c := colly.NewCollector(
    colly.AllowedDomains("www.hijklmn.com"),
    colly.MaxDepth(maxDepth),
  )
  return c
}

func init() {
  domain2Collector["www.abcdefg.com"] = initABCDECollector()
  domain2Collector["www.hijklmn.com"] = initHIJKLCollector()

  var err error
  nc, err = nats.Connect(natsURL)
  if err != nil {os.Exit(1)}
}

func startConsumer() {
  nc, err := nats.Connect(nats.DefaultURL)
  if err != nil {return}

  sub, err := nc.QueueSubscribeSync("tasks", "workers")
  if err != nil {return}

  var msg *nats.Msg
  for {
    msg, err = sub.NextMsg(time.Hour * 10000)
    if err != nil {break}

    urlStr := string(msg.Data)
    ins := factory(urlStr)
    // Because the most downstream one must be the landing page of the corresponding website.
    // So don’t have to make extra judgments, just climb the content directly.
    ins.Visit(urlStr)
    // prevent being blocked
    time.Sleep(time.Second)
  }
}

func main() {
  startConsumer()
}