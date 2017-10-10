package douban

import (
	//"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tealeg/xlsx"
	"strings"
)

const DOWNLOAD_URL = "https://movie.douban.com/top250/"

func parseInfo(uri string) (name, infoList []string, starCon, score []string, nextPage string) {
	doc, err := goquery.NewDocument(uri)

	if err != nil {
		// handle error
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
	//
	href := doc.Find("span.next").Find("a")

	nextPage, exists := href.Attr("href")

	if exists {
		return name, infoList, starCon, score, nextPage
	}

	return name, infoList, starCon, score, ""
}

func Run() {
	url := DOWNLOAD_URL
	var names []string
	var startCons []string
	var scores []string
	var infos []string

	for url != "" {
		name, infoList, starCon, score, nextPage := parseInfo(url)
		names = append(names, name...)
		infos = append(infos, infoList...)
		startCons = append(startCons, starCon...)
		scores = append(scores, score...)

		if nextPage != "" {
			url = DOWNLOAD_URL + nextPage
		} else {
			url = ""
		}
	}

	var file *xlsx.File
	var sheet *xlsx.Sheet

	file = xlsx.NewFile()
	sheet, _ = file.AddSheet("my sheet")

	row := sheet.AddRow()
	row.AddCell().Value = "片名"
	row.AddCell().Value = "点赞数目"
	row.AddCell().Value = "得分"
	row.AddCell().Value = "简介"

	for i := 0; i < len(names); i += 1 {
		row := sheet.AddRow()
		c1 := row.AddCell()
		c1.Value = names[i]

		c2 := row.AddCell()
		c2.Value = startCons[i]

		c3 := row.AddCell()
		c3.Value = scores[i]

		c4 := row.AddCell()
		c4.Value = infos[i]
	}

	file.Save("test.xlsx")

	//fmt.Println(names)
}
