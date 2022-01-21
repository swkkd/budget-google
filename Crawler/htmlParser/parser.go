package htmlParser

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

//func Parser(url string) (string, error) {
//	log.Printf("Parser received :: %s ...\n", url)
//	resp, err := http.Get(url)
//	// handle the error if there is one
//	if err != nil {
//		log.Printf("error get: %v", err)
//		return "", err
//	}
//	// do this now so it won't be forgotten
//	defer resp.Body.Close()
//	// reads html as a slice of bytes
//	htmll, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Printf("error readall: %v", err)
//		return "", err
//	}
//	// show the HTML code as a string %s
//	//fmt.Printf("%s\n ", html)
//	//HTMLToReadable(html) // convert html to readable
//	return string(htmll), nil
//}

func HtmlToReadable(s string) string {
	response, err := http.Get(s)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	textTags := []string{
		"a",
		"p", "span", "em", "string", "blockquote", "q", "cite",
		"h1", "h2", "h3", "h4", "h5", "h6",
	}

	tag := ""
	enter := false

	tokenizer := html.NewTokenizer(response.Body)

	var b bytes.Buffer

	for {
		tt := tokenizer.Next()
		token := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tt {
		case html.ErrorToken:
			log.Fatal(err)
		case html.StartTagToken, html.SelfClosingTagToken:
			enter = false

			tag = token.Data
			for _, ttt := range textTags {
				if tag == ttt {
					enter = true
					break
				}
			}
		case html.TextToken:
			if enter {
				data := strings.TrimSpace(token.Data)

				if len(data) > 0 {
					b.WriteString(data + " ")
					// fmt.Println(data)
				}
			}
		}
	}
	return b.String()
}
