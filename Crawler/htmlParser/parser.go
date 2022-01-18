package htmlParser

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	fmt.Printf("%s\n ", html)
	return html, nil
}
