package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	resp, err := http.Get("https://image.baidu.com/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	body := doc.Find("body")
	background_url := get_background_url(body)

	download(background_url)
	for i := 0; i < 6; i++ {
		download(get_pictures_urls(body)[i])
	}
	return
}

func download(url string) bool {
	fmt.Print("下载" + url + "\n")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print("下载失败\n")
		panic(err)
	}
	defer resp.Body.Close()

	picture, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("读取图片失败\n")
		panic(err)
	}

	idx := strings.LastIndex(url, "/")
	if url[len(url)-4] != '.' {
		url += ".jpg"
	}
	err = os.WriteFile(url[idx+1:], picture, 0666)
	if err != nil {
		fmt.Print("保存失败\n")
		panic(err)
	} else {
		fmt.Print("下载完成\n")
		return true
	}
}

func get_background_url(body *goquery.Selection) string {
	bd_home_wrapper := body.Find("div.bd_home_wrapper")
	wrapper_main_box := bd_home_wrapper.Find("div.wrapper_main_box")
	wrapper_skin_box := wrapper_main_box.Find("div.wrapper_skin_box")
	style, _ := wrapper_skin_box.Attr("style")
	style = regexp.MustCompile("\\((.*?)\\)").FindStringSubmatch(style)[1]
	return style
}

func get_pictures_urls(body *goquery.Selection) []string {
	var urls [6]string

	bd_home_wrapper := body.Find("div.bd_home_wrapper")
	bd_home_content := bd_home_wrapper.Find("#bd-home-content")
	bd_home_content_wrapper := bd_home_content.Find("div.bd-home-content-wrapper")
	bd_home_content_item := bd_home_content_wrapper.Find("div.bd-home-content-item")
	bd_home_content_item_main := bd_home_content_item.Find("div.bd-home-content-item-main")
	bd_home_content_album := bd_home_content_item_main.Find("div.bd-home-content-album")
	bd_home_content_album.Find("a").Each(func(i int, s *goquery.Selection) {
		bd_home_content_album_item_pic := s.Find("div.bd-home-content-album-item-pic")
		style, _ := bd_home_content_album_item_pic.Attr("style")
		urls[i] = regexp.MustCompile("\\((.*?)\\)").FindStringSubmatch(style)[1]
	})
	return urls[:]
}
