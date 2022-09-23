package main

import (
	"github.com/1121170088/find-domain/search"
	"log"
)

func main()  {
	search.Init("C:\\Users\\tanmingxin\\Desktop\\qwqw\\click-href2\\public_suffix_list.dat")
	log.Printf(search.Search("http://www.hao123.com/"))
	log.Printf(search.Search("http://www.hao123.com"))
}
