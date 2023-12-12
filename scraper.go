package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func GetHttpHtmlContent(url string, selector string, sel interface{}) (string, error) {

	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))

	_ = chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 40*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector),
		chromedp.OuterHTML(sel, &htmlContent, chromedp.ByJSPath),
	)

	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

func GetSpecialData(htmlContent string, selector string) ([][]string, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var results [][]string

	// 使用選擇器尋找所有匹配的元素
	dom.Find(selector).Each(func(i int, selection *goquery.Selection) {
		// 在每個匹配的元素中查找"td"標籤並提取文本
		/*
			selection.Find("td").Each(func(j int, trSelection *goquery.Selection) {
				// 将每个"tr"标签的文本追加到结果数组
				results = append(results, trSelection.Text())
			})
		*/

		//tr下的第一個td
		/*
			firstTd := selection.Find("td:nth-child(1)")
			results = append(results, firstTd.Text())
			//tr下的第二個td
			secondTd := selection.Find("td:nth-child(2)")
			results = append(results, secondTd.Text())
			//tr下的最後一個
			lasttd := selection.Find("td:last-child")
			results = append(results, lasttd.Text())
		*/
		selection.Find("tr").Each(func(j int, trSelection *goquery.Selection) {
			var row []string
			trSelection.Find("td").Each(func(k int, tdSelection *goquery.Selection) {
				row = append(row, tdSelection.Text())
			})

			results = append(results, row)
		})
	})

	return results, nil
}
func main() {
	//selector := "body > div > div.banner > div.swiper-container-place > div > div.swiper-slide.swiper-slide-0.swiper-slide-visible.swiper-slide-active > a.item.item-big > div.item-bottom"
	//selector := "#MainContent > ul > li > article > div > div > div.tb-out > table > tbody > tr:nth-child(1) > th:nth-child(1)"
	selector := "#MainContent > ul > li > article > div > div > div.tb-out > table > tbody"

	param := `document.querySelector("body")`
	url := "https://www.cmoney.tw/finance/3293/f00027"
	html, _ := GetHttpHtmlContent(url, selector, param)
	//res, _ := GetSpecialData(html, ".tb-out")
	res, _ := GetSpecialData(html, selector)
	for i, row := range res {
		if len(row) == 0 {
			continue
		}
		fmt.Printf("Row %d: 年度:%v, 現金股利:%v,股票股利合計:%v, 股利合計:%v\n", i+1, row[0], row[1], row[len(row)-2], row[len(row)-1])
	}
	return
}
