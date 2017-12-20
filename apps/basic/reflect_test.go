package main

import (
	"testing"
	"reflect"
)

func TestDoReflect(t *testing.T){
	doReflectStruct()
}

func TestDoReflectSlice(t *testing.T){
	doReflectSlice()
}

func TestDoReplaceValue(t *testing.T) {
	doReplaceValue()
}

func TestDoReflectStructBasic(t *testing.T){
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
			"刘峰"					:"黄轩",
			"何小萍"				:"苗苗",
			"萧穗子"				:"钟楚曦",
			"林丁丁"				:"杨采钰",
			"郝淑雯"				:"李晓峰",
			`陈灿`					:"王天辰",
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

func TestChangeValue(t *testing.T){
	changeValue()
}