package main

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
)

func doList() {
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

func convert() {
	i := 10
	fmt.Printf("i convert string : %s", strconv.Itoa(i))

	s := "1000"
	// The bitSize argument specifies the integer type
	// that the result must fit into. Bit sizes 0, 8, 16, 32, and 64
	// correspond to int, int8, int16, int32, and int64.
	if v, err := strconv.ParseInt(s, 10, 0); err != nil {
		fmt.Errorf("\n strconv.ParseInt %s\n", err.Error())
	} else {
		fmt.Printf("string s convert int64 : %d", v)
	}

}
