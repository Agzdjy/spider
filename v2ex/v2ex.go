package v2ex

import (
	"fmt"
	"spider/util"
	"strconv"
	"time"

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

func parseTablePage(tabUrl string, done chan<- string) {
	response := util.PorxyGet(tabUrl)
	doc, err := goquery.NewDocumentFromResponse(response)
	defer response.Body.Close()
	if err != nil {
		done <- "failed"
	}

	tags := doc.Find("div#Main").Find("div.box").Find("div.cell").Eq(0).Find("a")

	tags.Each(func(index int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		parseTagPage(domanURL + url)
		time.Sleep(3 * time.Second)
	})
	done <- "table" + tabUrl + "--over"
}

func parseTagPage(url string) {
	response := util.PorxyGet(url)
	doc, err := goquery.NewDocumentFromResponse(response)
	defer response.Body.Close()
	if err != nil {
		return
	}

	pages := doc.Find("div#Main").Find("div.box").Find("div.cell").Find("table").Find("td").Eq(0)
	as := pages.Find("a")
	maxPage, _ := strconv.Atoi(as.Eq(len(as.Nodes) - 1).Text())

	for page := 1; page <= maxPage; page += 1 {
		parseListPage(url + "?p=" + strconv.Itoa(page))
		time.Sleep(3 * time.Second)
	}
}

func parseListPage(url string) {
	response := util.PorxyGet(url)
	doc, err := goquery.NewDocumentFromResponse(response)
	defer response.Body.Close()
	if err != nil {
		return
	}

	topicsNodes := doc.Find("div#TopicsNode").Find("div.cell")
	topicsNodes.Each(func(index int, s *goquery.Selection) {
		title := s.Find("table").Find("tbody").Find("tr").Find("td").Eq(2).Find("a").Eq(0).Text()
		fmt.Println(url, "---->", title)
	})
	fmt.Println(url, "---------over")
}

func Run() {
	done := make(chan string, len(tabUrls))
	for _, tablUrl := range tabUrls {
		url := domanURL + tablUrl
		go parseTablePage(url, done)
		time.Sleep(3 * time.Second)
	}

	for i := 0; i < len(tabUrls); i += 1 {
		fmt.Println(<-done)
	}
}
