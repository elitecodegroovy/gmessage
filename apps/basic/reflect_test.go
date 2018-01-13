package main

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestDoReflect(t *testing.T) {
	doReflectStruct()
}

func TestDoReflectSlice(t *testing.T) {
	doReflectSlice()
}

func TestDoReplaceValue(t *testing.T) {
	doReplaceValue()
}

func TestDoReflectStructBasic(t *testing.T) {
	doReflectStructBasic()
}

func TestPrintValue(t *testing.T) {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Color           bool
		Actor           map[string]string
		Musics          []string
		Sequel          *string
	}
	//!-movie
	goodMovie := Movie{
		Title:    "芳华",
		Subtitle: "The Youth",
		Year:     2017,
		Color:    true,
		Actor: map[string]string{
			"刘峰":  "黄轩",
			"何小萍": "苗苗",
			"萧穗子": "钟楚曦",
			"林丁丁": "杨采钰",
			"郝淑雯": "李晓峰",
			`陈灿`:  "王天辰",
		},

		Musics: []string{
			"那些花儿",
			"想把你留在这里",
			"美好生活",
			"绒花",
		},
	}

	printValue("goodMovie", reflect.ValueOf(goodMovie))
	//fmt.Printf("\n %v", strangelove)
}

func TestChangeValue(t *testing.T) {
	changeValue()
}

func TestCallTypeFunc(t *testing.T) {
	fmt.Printf(" CallTypeFunc(func(x string) int { return len(x) }, []string{\"1\", \"10\", \"100\"}).([]int) return :%+v\n",
		CallTypeFunc(func(x string) int { return len(x) }, []string{"1", "10", "100"}))

	fmt.Printf(" CallTypeFunc(func(x int) int { return x*x }, []int{1, 10, 100}).([]int) return :%+v",
		CallTypeFunc(func(x int) int { return x * x }, []int{1, 10, 100}))
}

func TestBytesToString(t *testing.T) {
	/***************************************/
	byteArray1 := []byte{'J', 'O', 'H', 'N'}
	str1 := BytesToString(byteArray1)
	fmt.Println("String:", str1)

	/****************************************/

	str2 := string(byteArray1[:])
	fmt.Println("String:", str2)

	/****************************************/
	str3 := bytes.NewBuffer(byteArray1).String()
	fmt.Println("String:", str3)
}
