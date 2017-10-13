package ziroom

import (
	"fmt"
	"regexp"
	"spider/util"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tealeg/xlsx"
)

const DOWNLOAD_URL = "http://www.ziroom.com/z/nl/z2.html?qwd="

type ziroom struct {
	hoseName       string
	subWay         string
	subWayDistance string
	area           string
	rooms          string
	price          string
	rentType       string
	tags           string
	link           string
}

func parseInfo(url string) (zirooms []ziroom, nextPage string) {
	fmt.Println(url)
	doc, err := goquery.NewDocument(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	hostList := doc.Find("ul#houseList")

	hostList.Find("li.clearfix").Each(func(i int, s *goquery.Selection) {
		ziroom := ziroom{}
		numReg, _ := regexp.Compile("[0-9]+")

		txt := s.Find("div.txt")
		detail := txt.Find("div.detail").Find("p")

		ziroom.link, _ = txt.Find("h3").Find("a").Attr("href")
		ziroom.link = strings.Replace(ziroom.link, "//", "", 1)

		ziroom.rooms = numReg.FindString(detail.Eq(0).Find("span").Eq(2).Text())
		ziroom.area = numReg.FindString(detail.Eq(0).Find("span").Eq(0).Text())
		ziroom.rentType = detail.Eq(0).Find("span.icons").Text()
		ziroom.hoseName = txt.Find("h3").Find("a").Text()
		ziroom.price = numReg.FindString(s.Find("div.priceDetail").Find("p.price").Text())

		subWay := strings.Split(detail.Eq(1).Find("span").Text(), "站")
		if len(subWay) > 1 {
			ziroom.subWay = strings.Split(subWay[0], "距")[1]
			ziroom.subWayDistance = numReg.FindString(strings.Join(subWay[1:], ""))
		}

		var tags []string
		txt.Find("p.room_tags").Find("span").Each(func(index int, s *goquery.Selection) {
			tags = append(tags, s.Text())
		})

		ziroom.tags = strings.Join(tags, "| ")
		zirooms = append(zirooms, ziroom)
	})
	nextPage, exists := doc.Find("div.pages").Find("a.next").Attr("href")

	if !exists {
		nextPage = ""
	} else {
		nextPage = "http:" + nextPage
	}
	return zirooms, nextPage
}

func Run() (err error) {
	var allZirooms []ziroom
	query := util.GetEnv("query")

	url := DOWNLOAD_URL
	if query != "" {
		url = DOWNLOAD_URL + query
	}

	for url != "" {
		zirooms, nextPage := parseInfo(url)
		url = nextPage
		allZirooms = append(allZirooms, zirooms...)

	}

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("my sheet")

	row := sheet.AddRow()
	row.AddCell().Value = "房屋名"
	row.AddCell().Value = "地铁"
	row.AddCell().Value = "距离地铁(m)"
	row.AddCell().Value = "面积(㎡)"
	row.AddCell().Value = "居室"
	row.AddCell().Value = "价格(元)"
	row.AddCell().Value = "类型"
	row.AddCell().Value = "标签"
	row.AddCell().Value = "链接"

	for i := 0; i < len(allZirooms); i += 1 {
		ziroom := allZirooms[i]

		row := sheet.AddRow()
		row.AddCell().Value = ziroom.hoseName
		row.AddCell().Value = ziroom.subWay
		row.AddCell().Value = ziroom.subWayDistance
		row.AddCell().Value = ziroom.area
		row.AddCell().Value = ziroom.rooms
		row.AddCell().Value = ziroom.price
		row.AddCell().Value = ziroom.rentType
		row.AddCell().Value = ziroom.tags
		row.AddCell().Value = ziroom.link
	}
	fileName := "ziroom.xlsx"
	if query != "" {
		fileName = query + ".xlsx"
	}
	err = file.Save(fileName)
	return err
}
