package middleware

import (
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func HtmlToReadable(s string) ([]string, error) {
	response, err := http.Get(s)
	if err != nil {
		//log.Fatal(err)
		return nil, err
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

	var b []string
	//var b bytes.Buffer
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
					//b.WriteString(" " + data)
					b = append(b, data)
					// fmt.Println(data)
				}
			}
		}
	}
	return b, nil
}
