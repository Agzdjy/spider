package v2ex

import (
	"fmt"
	"math/rand"
	"spider/util"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type V2ex struct {
	title string
	tag   string
	link  string
}

const domanURL = "https://www.v2ex.com"

var v2exs []V2ex = []V2ex{}
var nodeUrls = []string{
	"/?tab=tech",
	"/?tab=creative",
	"/?tab=play",
	"/?tab=apple",
	"/?tab=jobs",
	"/?tab=deals",
	"/?tab=city",
	"/?tab=qna",
	"/?tab=hot",
	"/?tab=all",
	"/?tab=r2",
	"/?tab=nodes",
	"/?tab=members"}

func (v *V2ex) Run() {
	fmt.Println("v2ex spider start")
	nodeChan := make(chan string, 2)
	childNodeChan := make(chan string, 10)
	listChan := make(chan string, 20)

	go func() {
		for _, nodeUrl := range nodeUrls {
			nodeChan <- domanURL + nodeUrl
		}
	}()

	for {
		timeout := time.After(2 * time.Minute)

		select {
		case nodeUrl := <-nodeChan:
			time.Sleep(time.Second * time.Duration(rand.Int63n(4)))
			go writeChildNodeUrl2Chan(nodeUrl, nodeChan, childNodeChan)

		case childUrl := <-childNodeChan:
			time.Sleep(time.Second * time.Duration(rand.Int63n(4)))
			go writeListUrl2Chan(childUrl, childNodeChan, listChan)

		case listUrl := <-listChan:
			time.Sleep(time.Second * time.Duration(rand.Int63n(4)))
			go parseList(listUrl, listChan)

		case <-timeout:
			fmt.Println("all over or timeout 2 min")
			return
		}
	}
}

func getDoc(requestUrl string) (*goquery.Document, error) {
	response := util.PorxyGet(requestUrl)
	doc, err := goquery.NewDocumentFromResponse(response)

	return doc, err
}

func writeChildNodeUrl2Chan(nodeUrl string, nodeChan chan string, childNodeChan chan string) {
	doc, err := getDoc(nodeUrl)
	if err != nil {
		nodeChan <- nodeUrl
		return
	}

	childNodes := doc.Find("div#Main").Find("div.box").Find("div.cell").Eq(0).Find("a")
	if len(childNodes.Nodes) == 0 {
		nodeChan <- nodeUrl
		return
	}

	childNodes.Each(func(index int, s *goquery.Selection) {
		childNodeUrl, _ := s.Attr("href")
		childNodeChan <- domanURL + childNodeUrl
	})

}

func writeListUrl2Chan(childUrl string, childNodeChan chan string, listChan chan string) {
	doc, err := getDoc(childUrl)
	if err != nil {
		childNodeChan <- childUrl
		return
	}

	pages := doc.Find("div#Main").Find("div.box").Find("div.cell").Find("table").Find("td").Eq(0)
	as := pages.Find("a")
	maxPage, _ := strconv.Atoi(as.Eq(len(as.Nodes) - 1).Text())
	for page := 1; page <= maxPage; page += 1 {
		listChan <- childUrl + "?p=" + strconv.Itoa(page)
	}
}

func parseList(listUrl string, listChan chan string) {
	doc, err := getDoc(listUrl)
	if err != nil {
		listChan <- listUrl
		return
	}

	topicsNodes := doc.Find("div#TopicsNode").Find("div.cell")
	if len(topicsNodes.Nodes) == 0 {
		listChan <- listUrl
		return
	}
	topicsNodes.Each(func(index int, s *goquery.Selection) {
		title := s.Find("table").Find("tbody").Find("tr").Find("td").Eq(2).Find("a").Eq(0).Text()
		// TODO save data
		fmt.Println(listUrl, "---->", title)
	})
}
