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

func CheckandLoggingError(err error) {
	is_err := err != nil
	if is_err {
		log.Fatal(err)
	}
}

type Anime struct {
	title string
	url   string
}

func DownloadImgFromURL(img_url string, anime Anime) {
	const base_url string = "http://animekabegami.com"

	// no thanks icon, header img
	if strings.Index(img_url, "icon") != -1 {
		return
	}
	if strings.Index(img_url, "header") != -1 {
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
	//	fmt.Println(img_name)
	file, err := os.Create("../../AnimeKabegami/" + anime.title + "/" + img_name)
	CheckandLoggingError(err)
	defer file.Close()

	// download image
	io.Copy(file, res.Body)
}

func GetImgfromWeb(url string, anime Anime) {
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
		fmt.Println(img_url)
		time.Sleep(5 * time.Second)
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
	res, err := http.Get(url)
	CheckandLoggingError(err)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal("status code error")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckandLoggingError(err)

	doc.Find(".side-menu-body ul li a").Each(func(_ int, s *goquery.Selection) {
		anime := Anime{}
		anime.url, _ = s.Attr("href")
		anime.title = s.Text()
		fmt.Println(anime.title, anime.url)
		os.Mkdir("../../AnimeKabegami/"+anime.title, 0777)
		GetImgfromWeb(url+anime.url[2:], anime)
	})
}

func main() {
	os.Mkdir("../../AnimeKabegami", 0777)
	url := "http://animekabegami.com/"
	//url = "http://animekabegami.com/select?title=C"
	GetAnimeList(url)
	//anime := Anime{"C (26)", "http://animekabegami.com/select?title=C"}
	//GetImgfromWeb(url, anime)
}
