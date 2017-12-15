package util

import (
	"fmt"
	"strconv"
)

var a string = "1"

func OuterFunc(strToInt func(s string) int, b int) string {
	c := strToInt(a) + b
	//int is converted to string type
	return strconv.Itoa(c)
}

func StartFormatInt() {
	strToInt := func(s string) int {
		//Atoi is shorthand for ParseInt(s, 10, 0).
		num, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("error :", s, ", info :", err)
		}
		return num
	}
	a = "12"
	fmt.Println("output 14 is ", OuterFunc(strToInt, 2) == "14")
}
