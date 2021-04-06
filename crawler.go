package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type crawler struct {
	html string
	url  string
}

var ImgFormats = []string{"jpg", "png", "webp", "jpeg"}

func makeGetRequest(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error : ", err)
	}
	return resp.Body
}

func (c crawler) toString(reader io.Reader) string {
	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	c.html = string(bodyBytes)
	return string(bodyBytes)
}

func (c crawler) allImages() []string {
	formats := strings.Join(ImgFormats, "|")
	re := regexp.MustCompile(`(http(s?):)([/|.|\w|\s|-])*\.(?:` + formats + `)`)
	return removeDuplicates(re.FindAllString(c.toString(makeGetRequest(c.url)), -1))
}

func saveFile(fileName string, bs []byte) error {
	return os.WriteFile(fileName, bs, 0777)
}

func (c crawler) saveAllImages(dir string) int {
	images := c.allImages()
	os.Mkdir(dir, 0777)
	channel := make(chan string)
	for _, image := range images {
		go saveImageLink(image, dir, channel)
	}

	for i := 0; i < len(images); i++ {
		fmt.Println(<-channel)
	}
	return len(images)
}

func saveImageLink(link string, dir string, channel chan string) {
	req := makeGetRequest(link)
	data, err := io.ReadAll(req)
	req.Close()
	if err != nil {
		log.Fatal(err)
	}
	slashIndex := strings.LastIndex(link, `/`)
	fileName := link[slashIndex+1:]
	saveFile(dir+`/`+fileName, data)
	channel <- fileName + " saved!"
}

func removeDuplicates(slice []string) []string {
	check := make(map[string]int)
	res := make([]string, 0)
	for _, val := range slice {
		check[val] = 1
	}
	for letter, _ := range check {
		res = append(res, letter)
	}
	return res
}
