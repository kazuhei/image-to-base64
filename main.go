package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func downloadImage(url string) string {
	response, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(response.Body)

	imgBase64Str := base64.StdEncoding.EncodeToString(body)

	return imgBase64Str
}

func getExtension(path string) (string, error) {
	pos := strings.LastIndex(path, ".")
	switch path[pos:] {
	case ".jpeg", ".jpg":
		return "jpg", nil
	case ".png":
		return "png", nil
	case ".gif":
		return "gif", nil
	}

	return "", errors.New("未対応の拡張子です " + path[pos:])
}

func overwriteImageToBase64(fileName string) {
	fileInfos, _ := ioutil.ReadFile("src.html")
	stringReader := strings.NewReader(string(fileInfos))
	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		fmt.Println("scrapping failed")
	}

	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("src")
		extension, err := getExtension(url)
		if err != nil {
			fmt.Println(err)
		}
		base64String := downloadImage(url)
		s.SetAttr("src", "data:image/"+extension+";base64,"+base64String)
	})

	newHtml, err := doc.Selection.Html()
	if err != nil {
		fmt.Println("make new html failed")
	}

	ioutil.WriteFile("output.html", []byte(newHtml), os.ModePerm)
}

func main() {
	fileName := os.Args[1]
	overwriteImageToBase64(fileName)
}
