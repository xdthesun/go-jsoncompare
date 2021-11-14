package main

import (
	//"github.com/iostrovok/go-jsoncompare/jsoncompare"

	"fmt"
	"io/ioutil"
	"log"
	"mylocal/go-jsoncompare/jsoncompare"
	"net/http"
)

var url1 string = "http://ergast.com/api/f1/2011.json"
var url2 string = "http://ergast.com/api/f1/2012.json"

var JSON0 []byte = []byte(`{"family":{"papa":"yes","mama":"too","children":{"count":2,"list":[{"type":"boy","name":"John"},{"type":"girl","name":"Linda"}]}}}`)
var JSON1 []byte = []byte(`{"family":{"papa":"yes","mama":"too","children":{"number":3,"list":[{"type":"boy","name":"John"},{"type":"girl","name":"Linda"},{"type":"boy","name":"Mike"}]}}}`)

func main() {
	ignoreRules := []string{"/family/children/list/.*?/1/.*?"}
	//b1, e1 := loadUrl(url1)
	//if e1 != nil {
	//	log.Fatalln(e1)
	//}
	//
	//b2, e2 := loadUrl(url2)
	//if e2 != nil {
	//	log.Fatalln(e2)
	//}

	//list, err := jsoncompare.Compare(b1, b2, ignoreRules)
	list, err := jsoncompare.Compare(JSON0, JSON1, ignoreRules)
	if err != nil {
		log.Fatalln(err)
	}

	leftOnly, rightOnly, noEqual, goodList,IgnoredList  := jsoncompare.SplitBySide(list)

	printList("Matched: ", goodList)
	printList("Left Only: ", leftOnly)
	printList("Right Only: ", rightOnly)
	printList("Ignored: ", IgnoredList)
	printList("No Equal: ", noEqual)

}

func printList(prefix string, list []*jsoncompare.PathDiff) {
	for i, v := range list {

		viewPath := v.PathRight
		if v.PathLeft != "" {
			viewPath = v.PathLeft
		}

		fmt.Printf("%d. %s path: %s \n", i, prefix, viewPath)

		if !v.IsEqual {
			if v.ValueLeft != nil && v.ValueRight != nil {
				fmt.Printf("   value: %+v != %+v isignored: %v\n,", v.ValueLeft, v.ValueRight, v.IsIgnored)
			} else if v.ValueLeft != nil {
				fmt.Printf("   value: %+v isignored: %v\n", v.ValueLeft, v.IsIgnored)
			} else if v.ValueRight != nil {
				fmt.Printf("   value: %+v isignored: %v\n", v.ValueRight, v.IsIgnored)
			}
		}
	}
}

func loadUrl(url string) ([]byte, error) {
	resp, err_get := http.Get(url)
	if err_get != nil {
		return nil, err_get
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
