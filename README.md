代码源自[iostrovok/go-jsoncompare](https://github.com/iostrovok/go-jsoncompare) 增加根据特定的模板做过滤

## Comparing two json string as json structure. ##

### Installing ###
```bash
go get github.com/xdthesun/go-jsoncompare/jsoncompare

```
### How use: example ###

```go

package main

import (
	"fmt"
	"log"
	"mylocal/go-jsoncompare/jsoncompare"
)

var JSON0 []byte = []byte(`{"family":{"papa":"yes","mama":"too","children":{"count":2,"list":[{"type":"boy","name":"John"},{"type":"girl","name":"Linda"}]}}}`)
var JSON1 []byte = []byte(`{"family":{"papa":"yes","mama":"too","children":{"number":3,"list":[{"type":"boy","name":"John"},{"type":"girl","name":"Linda"},{"type":"boy","name":"Mike"}]}}}`)

func main() {
	ignoreRules := []string{"/family/children/list/.*?/1/.*?"}
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

```
