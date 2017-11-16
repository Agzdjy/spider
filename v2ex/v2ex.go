package v2ex

import (
	"fmt"
	"strconv"
	"time"

	"spider/util"

	"github.com/PuerkitoBio/goquery"
)

type v2ex struct {
	title string
	tag   string
	link  string
}

const domanURL = "https://www.v2ex.com"

var tabUrls = []string{
	"/?tab=tech"}

// "/?tab=creative",
// "/?tab=play",
// "/?tab=apple",
// "/?tab=jobs",
// "/?tab=deals",
// "/?tab=city",
// "/?tab=qna",
// "/?tab=hot",
// "/?tab=all",
// "/?tab=r2",
// "/?tab=nodes",
// "/?tab=members"}

var v2exs []v2ex = []v2ex{}

func parseTablePage(tabUrl string, done chan string) {
	response := util.PorxyGet(tabUrl)
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		done <- "failed"
	}

	tags := doc.Find("div#Main").Find("div.box").Find("div.cell").Eq(0).Find("a")
	numOfTags := len(tags.Nodes)
	listDone := make(chan string, numOfTags)

	tags.Each(func(index int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		go parseTagPage(domanURL+url, listDone)
	})

	for i := 0; i < numOfTags; i += 1 {
		fmt.Println(<-listDone)
	}
	close(listDone)
	done <- "table ok-->" + tabUrl
}

func parseTagPage(url string, done chan string) {
	fmt.Println("------->parseTagPage")
	time.Sleep(500 * time.Millisecond)

	response := util.PorxyGet(url)
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		done <- "failed"
	}
	pages := doc.Find("div#Main").Find("div.box").Find("div.cell").Find("table").Find("td").Eq(0)
	as := pages.Find("a")
	maxPage, _ := strconv.Atoi(as.Eq(len(as.Nodes) - 1).Text())

	pageDone := make(chan string, maxPage)

	for page := 1; page <= maxPage; page += 1 {
		go parseListPage(url+"?p="+strconv.Itoa(page), pageDone)
	}

	for i := 1; i <= maxPage; i += 1 {
		fmt.Println(<-pageDone)
	}
	close(pageDone)
	done <- "tag ok--->" + url
}

func parseListPage(url string, done chan string) {
	time.Sleep(2 * time.Second)

	response := util.PorxyGet(url)
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		done <- "failed"
	}

	topicsNodes := doc.Find("div#TopicsNode").Find("div.cell")
	topicsNodes.Each(func(index int, s *goquery.Selection) {
		title := s.Find("table").Find("tbody").Find("tr").Find("td").Eq(2).Find("a").Eq(0).Text()
		fmt.Println(url, "---->", title)
	})
	fmt.Println("------->, over")
	done <- "list ok-->" + url
}

func Run() {
	done := make(chan string, len(tabUrls))
	for _, tablUrl := range tabUrls {
		url := domanURL + tablUrl
		go parseTablePage(url, done)
	}

	for i := 0; i < len(tabUrls); i += 1 {
		fmt.Println(<-done)
	}
}
