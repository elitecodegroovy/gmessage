package main

import (
	"container/list"
	"strings"
	"fmt"
)

func doList(){
	a := "ABCDEFGH"
	items := list.New()
	for _, x := range strings.Split(a, "") {
		items.PushFront(x)
	}
	e1 := items.PushBack(len(a))
	e2 := items.PushBack(11)
	items.InsertAfter(9, e1)
	items.InsertBefore(10, e2)
	for element := items.Front(); element != nil; element = element.Next() {
		switch value := element.Value.(type) {
		case string:
			fmt.Printf("%s ", value)
		case int:
			fmt.Printf("%d ", value)
		}
	}
	fmt.Println() // prints: H G F E D B A 9
}
