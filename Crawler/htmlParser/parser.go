package htmlParser

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//todo обработать ошибку с невалидным url?

func Parser(url []byte) ([]byte, error) {
	fmt.Printf("HTML code of %s ...\n", url)
	resp, err := http.Get(string(url))
	// handle the error if there is one
	if err != nil {
		log.Printf("error get: %v", err)
		return nil, err
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error readall: %v", err)
		return nil, err
	}
	// show the HTML code as a string %s
	//fmt.Printf("%s\n ", html)
	HTMLToReadable(html) // convert html to readable
	return html, nil
}

func HTMLToReadable(s []byte) {
	domDocTest := html.NewTokenizer(strings.NewReader(string(s)))
	previousStartTokenTest := domDocTest.Token()
loopDomTest:
	for {
		tt := domDocTest.Next()
		switch {
		case tt == html.ErrorToken:
			break loopDomTest // End of the document,  done
		case tt == html.StartTagToken:
			previousStartTokenTest = domDocTest.Token()
		case tt == html.TextToken:
			if previousStartTokenTest.Data == "script" {
				continue
			}
			TxtContent := strings.TrimSpace(html.UnescapeString(string(domDocTest.Text())))
			if len(TxtContent) > 0 {
				fmt.Printf("%s\n", TxtContent)
			}
		}
	}
}
