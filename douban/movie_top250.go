package douban

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tealeg/xlsx"
)

const DOWNLOAD_URL = "https://movie.douban.com/top250/"

type Douban struct{}

func parseInfo(url string) (name, infoList []string, starCon, score []string, nextPage string) {
	fmt.Println(url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	ol := doc.Find("ol.grid_view")

	ol.Find("li").Each(func(i int, s *goquery.Selection) {
		detail := s.Find("div.hd")
		movieName := detail.Find("span.title")

		levelStar := s.Find("span.rating_num").Text()

		star := s.Find("div.star")
		starNum := strings.Split(strings.Fields(star.Text())[1], "人")[0]

		info := s.Find("span.inq").Text()
		if info != "" {
			infoList = append(infoList, info)
		} else {
			infoList = append(infoList, "无")
		}

		score = append(score, levelStar)
		name = append(name, movieName.Text())
		starCon = append(starCon, starNum)
	})

	href := doc.Find("span.next").Find("a")
	nextPage, exists := href.Attr("href")
	if exists {
		return name, infoList, starCon, score, DOWNLOAD_URL + nextPage
	}
	return name, infoList, starCon, score, ""
}

func (dou *Douban) Run() {
	var names []string
	var startCons []string
	var scores []string
	var infos []string

	for url := DOWNLOAD_URL; url != ""; {
		name, infoList, starCon, score, nextPage := parseInfo(url)

		names = append(names, name...)
		infos = append(infos, infoList...)
		startCons = append(startCons, starCon...)
		scores = append(scores, score...)
		url = nextPage
	}

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("my sheet")

	row := sheet.AddRow()
	row.AddCell().Value = "片名"
	row.AddCell().Value = "点赞数目"
	row.AddCell().Value = "得分"
	row.AddCell().Value = "短评"

	for i := 0; i < len(names); i += 1 {
		row := sheet.AddRow()

		row.AddCell().Value = names[i]
		row.AddCell().Value = startCons[i]
		row.AddCell().Value = scores[i]
		row.AddCell().Value = infos[i]
	}
	_ = file.Save("doubanTop250.xlsx")
}
