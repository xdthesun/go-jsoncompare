package jsoncompare

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type PathDiff struct {
	PathLeft   string
	PathRight  string
	IsEqual    bool
	IsIgnored  bool
	ValueLeft  interface{}
	ValueRight interface{}
}

type path struct {
	path      string      // from top to element likes: "/user/email/0/login"
	valueType string      // "map", "slice", "value"
	value     interface{} // Value of element. There is len(value) if valueType is maps or slices.
}

func SplitBySide(list []*PathDiff) (leftOnly, rightOnly, noEqual, goodList, IgnoredList []*PathDiff) {
	for _, v := range list {
		if v.PathLeft == "" && !v.IsIgnored {
			rightOnly = append(rightOnly, v)
		} else if v.PathRight == "" && !v.IsIgnored {
			leftOnly = append(leftOnly, v)
		} else if !v.IsEqual && !v.IsIgnored{
			noEqual = append(noEqual, v)
		} else if v.IsIgnored{
			IgnoredList = append(IgnoredList, v)
		} else {
			goodList = append(goodList, v)
		}
	}
	return
}

// Compare returns list of diffs
func Compare(left, right []byte, ignoreRules []string) ([]*PathDiff, error) {

	var bodyLeft, bodyRight map[string]interface{}
	var err error

	if bodyLeft, err = getJson(left); err != nil {
		return nil, err
	}
	if bodyRight, err = getJson(right); err != nil {
		return nil, err
	}

	pathsLeft := allPaths(bodyLeft, "")
	pathsRight := allPaths(bodyRight, "")

	checkList, err := comparePaths(pathsLeft, pathsRight, ignoreRules)

	return checkList, err
}

func getJson(data []byte) (map[string]interface{}, error) {

	list := []byte(`{"top":`)
	list = append(list, data...)
	list = append(list, []byte(`}`)...)

	var out map[string]interface{}
	err := json.Unmarshal(list, &out)
	return out, err
}

func allPaths(body interface{}, way string) []path {

	var list []path

	switch body.(type) {
	case map[string]interface{}:
		list = mAppend(list, way, "map", len(body.(map[string]interface{})))
		for k, v := range body.(map[string]interface{}) {
			next := way + "/" + k
			p := allPaths(v, next)
			list = append(list, p...)
		}
	case []interface{}:
		list = mAppend(list, way, "slice", len(body.([]interface{})))
		for k, v := range body.([]interface{}) {
			next := way + "/" + strconv.Itoa(k)
			p := allPaths(v, next)
			list = append(list, p...)
		}
	default:
		list = mAppend(list, way, fmt.Sprintf("%T", body), body)
	}

	var out []path
	for _, v := range list {
		v.path = strings.TrimPrefix(v.path, "/top")
		if v.path != "" {
			out = append(out, v)
		}
	}

	return out
}

func mAppend(list []path, way, valueType string, body interface{}) []path {

	if way == "" {
		return list
	}

	p := path{
		path:      way,
		valueType: valueType,
		value:     body,
	}
	list = append(list, p)
	return list
}

// ignoreRule 为忽略的模板,一般为path或者path的正则
func comparePaths(left, right []path, ignoreRules []string) ([]*PathDiff, error) {

	var out []*PathDiff
	rightCheck := map[string]bool{}

	for _, vL := range left {
		diff, err0 := fundInPath(true, vL, right, ignoreRules)
		if err0 != nil{
			return out, err0
		}
		if diff.PathRight != "" {
			rightCheck[diff.PathRight] = true
		}
		out = append(out, diff)
	}

	for _, vR := range right {
		if !rightCheck[vR.path] {
			diff, err1 := fundInPath(false, vR, left, ignoreRules)
			if err1 != nil{
				return out, err1
			}
			out = append(out, diff)
		}
	}
	return out,nil
}

func fundInPath(leftOrRight bool, from path, p []path, ignoreRules []string) (*PathDiff, error) {

	diff := &PathDiff{
		PathLeft:  "",
		PathRight: "",
		IsEqual:   false,
		IsIgnored: false,
	}

	if leftOrRight {
		diff.PathLeft = from.path
		diff.ValueLeft = from.value
	} else {
		diff.ValueRight = from.value
		diff.PathRight = from.path
	}

	for _, to := range p {
		contain, err := isRegexpContain(ignoreRules, from.path)
		if err != nil{
			return diff, err
		}
		diff.IsIgnored = contain
		if from.path == to.path {

			if diff.PathLeft == "" {
				diff.PathLeft = to.path
				diff.ValueLeft = to.value
			} else {
				diff.PathRight = to.path
				diff.ValueRight = to.value
			}
			diff.IsEqual = getEqual(from, to)
			break
		}
	}
	return diff, nil
}

func getEqual(vL, vR path) bool {

	if vL.valueType != vR.valueType {
		return false
	}

	return vL.value == vR.value
}

// 判断target是否满足特定规则，如果满足的话返回true，否则返回false，支持rule支持正则匹配
func isRegexpContain(ruleList []string, target string) (bool, error){
	for _, it := range ruleList {
		matched, err := regexp.MatchString("^" + it + "$", target)
		if err != nil{
			return false, err
		}
		if  matched {
			return true, nil
		}
	}
	return false, nil
}