package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CheckandLoggingError(err error) {
	is_err := err != nil
	if is_err {
		log.Fatal(err)
	}
}

func DownloadImgFromURL(img_url string) {
	const base_url string = "http://animekabegami.com"
	const excusion_word string = "icon"

	// no thanks icon img
	if strings.Index(img_url, excusion_word) != -1 {
		return
	}

	// thumb img is mini size image
	// so i will download original image that has no "thumb" in url
	const thumb_str string = "thumb"
	if strings.LastIndex(img_url, thumb_str) == -1 {
		img_url = img_url[1:]
	} else {
		expanded := img_url[strings.LastIndex(img_url, "."):]
		img_url = img_url[:strings.LastIndex(img_url, thumb_str)-1] + expanded
	}
	res, err := http.Get(base_url + img_url)
	CheckandLoggingError(err)
	defer res.Body.Close()

	img_name := img_url[strings.LastIndex(img_url, "/")+1:]
	fmt.Println(img_name)
	file, err := os.Create(img_name)
	CheckandLoggingError(err)
	defer file.Close()

	// download image
	io.Copy(file, res.Body)
}

func GetImgURLfromWeb(url string) {
	res, err := http.Get(url)
	CheckandLoggingError(err)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal("status code error")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckandLoggingError(err)
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		img_url, _ := s.Attr("src")
		DownloadImgFromURL(img_url)
	})

}

func main() {
	const url string = "http://animekabegami.com/select?title=%E3%81%AF%E3%81%9F%E3%82%89%E3%81%8F%E7%B4%B0%E8%83%9E"
	GetImgURLfromWeb(url)
}
