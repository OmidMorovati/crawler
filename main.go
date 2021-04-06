package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("please enter url: ")
	url, _ := reader.ReadString('\n')
	fmt.Print("please enter directory to save images: ")
	dir, _ := reader.ReadString('\n')
	crawler := crawler{}
	crawler.url = strings.TrimSuffix(url, "\n")
	crawler.saveAllImages(strings.TrimSuffix(dir, "\n"))
}
