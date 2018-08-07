package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Anime struct {
	title string
	url   string
}

func HttpGet(url string) *http.Response {
	res, err := http.Get(url)
	CheckandLoggingError(err)
	defer res.Body.Close()

	return res
}

func CheckandLoggingError(err error) {
	is_err := err != nil
	if is_err {
		log.Fatal(err)
	}
}

func URLHasExclusionWord(img_url string) bool {
	// no thanks icon, header img
	if strings.Index(img_url, "icon") != -1 {
		return true
	}
	if strings.Index(img_url, "header") != -1 {
		return true
	}

	return false
}

func RemoveThumb(img_url string) string {
	// thumb img is mini size image
	// so i will download original image that has no "thumb" in url

	const thumb_str string = "thumb"
	if strings.LastIndex(img_url, thumb_str) == -1 {
		img_url = img_url[1:]
	} else {
		expanded := img_url[strings.LastIndex(img_url, "."):]
		img_url = img_url[:strings.LastIndex(img_url, thumb_str)-1] + expanded
	}
	fmt.Println(img_url)

	return img_url
}

func FileIsExist(file_path string) bool {
	_, err := os.Stat(file_path)
	return !os.IsNotExist(err)
}

func CheckStatusCode(status_code int) {
	if status_code != 200 {
		log.Fatal("status code error")
	}
}

func DownloadImgFromURL(img_url string, anime Anime) {
	if URLHasExclusionWord(img_url) {
		return
	}

	img_url = RemoveThumb(img_url)

	img_path := "../../AnimeKabegami/" + anime.title + "/" + img_url[strings.LastIndex(img_url, "/")+1:]
	if FileIsExist(img_path) {
		return
	}

	// GET img
	const base_url string = "http://animekabegami.com"
	res := HttpGet(base_url + img_url)

	// download image
	file, err := os.Create(img_path)
	CheckandLoggingError(err)
	defer file.Close()
	io.Copy(file, res.Body)

	time.Sleep(5 * time.Second)
}

func GetImgfromWeb(url string, anime Anime) {
	res := HttpGet(url)
	CheckStatusCode(res.StatusCode)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckandLoggingError(err)
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		img_url, _ := s.Attr("src")
		DownloadImgFromURL(img_url, anime)
	})

	doc.Find(".paging .blk2 a").Each(func(_ int, s *goquery.Selection) {
		if s.Text() == "次へ" {
			next_url, _ := s.Attr("href")
			GetImgfromWeb("http://animekabegami.com/"+next_url[1:], anime)
		}
	})
}

func GetAnimeList(url string) {
	res := HttpGet(url)
	CheckStatusCode(res.StatusCode)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckandLoggingError(err)

	doc.Find(".side-menu-body ul li a").Each(func(i int, s *goquery.Selection) {
		// すでにダウンロードしてるの飛ばすため１０からスタート
		if i > 9 {
			anime := Anime{}
			anime.url, _ = s.Attr("href")
			anime.title = s.Text()
			fmt.Println(anime.title, anime.url)
			anime.title = strings.Replace(anime.title, "/", "_", -1)
			anime.title = strings.Replace(anime.title, "?", "_", -1)
			anime.title = strings.Replace(anime.title, ":", "_", -1)
			os.Mkdir("../../AnimeKabegami/"+anime.title, 0777)
			GetImgfromWeb(url+anime.url[2:], anime)
		}
	})
}

func main() {
	os.Mkdir("../../AnimeKabegami", 0777)
	url := "http://animekabegami.com/"
	GetAnimeList(url)
}
