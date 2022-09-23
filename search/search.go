package search

import (
	"bufio"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
)

var (
	tree map[byte] *node = make(map[byte] *node)
)

type node struct{
	end bool
	folw map[byte] *node
}

func reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
func fillIn(tld string) (has bool) {
	bytes := []byte(tld)
	reverse(bytes)
	var preNode *node = nil
	var ok bool = false
	has = true
	for _, b := range bytes {
		if preNode == nil {
			preNode, ok = tree[b]
			if !ok {
				preNode = &node{
					end:  false,
					folw: make(map[byte] *node),
				}
				tree[b] = preNode
				has = false
			}
		} else {
			preNode2, ok := preNode.folw[b]
			if !ok {
				preNode2 = &node{
					end:  false,
					folw: make(map[byte] *node),
				}
				preNode.folw[b]= preNode2
				has = false
			}
			preNode = preNode2
		}

	}
	if preNode != nil {
		if !preNode.end {
			preNode.end = true
			has = false
		}
	}
	return
}

func Init(tldFile string)  {
	f, err := os.OpenFile(tldFile, os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, "\r")
		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}
		if strings.Index(line, "//") == 0 {
			continue
		}
		fillIn(line)
	}
}

func Search(url string) string  {
	urlLen := len(url)
	if urlLen == 0 {
		return ""
	}
	from := 0
	doubleSlash := strings.Index(url, "//")
	if doubleSlash != -1 {
		from += doubleSlash + 2
	}
	poundIdx := strings.Index(url[from:], "#")
	questIdx := strings.Index(url[from:], "?")
	slashIdx := strings.Index(url[from:], "/")
	minIdx := urlLen - from
	if poundIdx != -1 && minIdx > poundIdx {
		minIdx = poundIdx
	}
	if questIdx != -1 && minIdx > questIdx {
		minIdx = questIdx
	}
	if slashIdx != -1 && minIdx > slashIdx {
		minIdx = slashIdx
	}
	subDomain := url[from:from + minIdx]
	bytes := []byte(subDomain)
	reverse(bytes)

	var preNode *node = nil
	var ok bool = false
	var mostRightSuffixIdx = 0
	var preDotIdx = -1
	for _, b := range bytes {
		mostRightSuffixIdx++
		if preNode == nil {
			preNode, ok = tree[b]
			if !ok {
				return ""
			}
		} else {
			preNode2, ok := preNode.folw[b]
			isDot := b == '.'
			if !ok {
				if !isDot {
					if preDotIdx != -1 {
						mostRightSuffixIdx = preDotIdx
						break
					} else {
						return ""
					}
				}
				if !preNode.end {
					return ""
				}
				break
			} else {
				if isDot {
					if preNode.end {
						preDotIdx = mostRightSuffixIdx
					}
				}
				preNode = preNode2
			}
		}
	}
	subDomainLen := len(subDomain)
	if mostRightSuffixIdx == subDomainLen {
		return ""
	}
	mostLeftSuffixIdx := subDomainLen - mostRightSuffixIdx
	for dotIdx:=mostLeftSuffixIdx -1; dotIdx >= 0; dotIdx -- {
		if subDomain[dotIdx] == '.' {
			return subDomain[dotIdx + 1:]
		}
	}
	return subDomain

}
