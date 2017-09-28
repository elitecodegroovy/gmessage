package main

import (
	"testing"
	"regexp"
	"fmt"
)

//Compile/MustCompile , MatchString
func TestBasicRegExp(t *testing.T){
	r := regexp.MustCompile(`Hello`)
	// Will print 'Match'
	if r.MatchString("Hello Regular Expression.") == true {
		fmt.Printf("Match ")
	} else {
		fmt.Printf("No match ")
	}

	s := "ABCDEEEEE"
	rr := regexp.MustCompile(`ABCDE{2}|ABCDE{4}`)
	rp := regexp.MustCompilePOSIX(`ABCDE{2}|ABCDE{4}`)
	fmt.Println(rr.FindAllString(s, 2))
	fmt.Println(rp.FindAllString(s, 2))

	//Character class '\w' represents any character from the class [A-Za-z0-9_], mnemonic: 'word'
	r, _ = regexp.Compile(`H\wllo`)
	// Will print 'true'.
	fmt.Printf("%v", r.MatchString("Hello Regular Expression."))

	//Character class '\d' represents any numeric digit.

	//Character class '\s' represents any of the following whitespaces:
	// TAB, SPACE, CR, LF. Or more precisely [\t\n\f\r ].
	r, _ = regexp.Compile(`\s`)
	// Will print 'true':
	fmt.Printf("%v", r.MatchString("/home/bill/My Documents"))

	//Character classes can be negated by using the uppercase '\D', '\S', '\W'. Thus, '\D' is any character that is not a '\d'.
	r, _ = regexp.Compile(`\S`) // Not a whitespace
	// Will print 'true', obviously there are non-whitespaces here:
	fmt.Printf("%v", r.MatchString("/home/bill/My Documents"))

	test := "[nsqplus-consumer] 2017/07/14 17:32:12.768419 eq body: nsqplus.Event{Platform:0, PointCode:\"0150101\", Compaign:\"\", Url:\"http%3A%2F%2Fwww.mama100.com%2Fwmall%2Fcenter%2Flottery%2Fhtml5%2Flottery.html\", " +
		"CreatedTime:1500024723891, CustomerId:\"35801214\", UserMark:\"osFJ4jtf87ClFMFNG9SVcmUSbgqo\", UserSourceId:\"\", Mobile:\"\", MobileInfo:\"vivoX5ProD\", " +
		"AppVersion:\"6.0.0\", SourceFrom:\"\", Sku:\"\", Spu:\"\", CourseId:\"\", AccountId:\"\", TerminalCode:\"\", CouponDefid:\"\", ExpertId:\"\", QuestionId:\"\", " +
		"Ip:\"157.122.164.17\", VarX:\"\", VarY:\"\", VarZ:\"\", Sign:\"bd2d6c979fa42a81e6e1696fb2f3f653\"} ,total count : 6679475, isSuccess: false"
	r, _ = regexp.Compile(`\{[\S\s]+}`)
	fmt.Println("\n match :", r.FindString(test))
}

func TestAdvanceRegExp(t *testing.T){
	//[[cat] [sat] [mat]]
	re, _ := regexp.Compile(`.at`)
	res := re.FindAllStringSubmatch("The cat sat on the mat.", -1)
	fmt.Printf("%v", res)

	//Parentheses allow to capture that piece of the string that you are actually interested in, instead of the entire regex.
	//[[cat c] [sat s] [mat m]]
	re, _ = regexp.Compile(`(.)at`) // want to know what is in front of 'at'
	res = re.FindAllStringSubmatch("The cat sat on the mat.", -1)
	fmt.Printf("%v", res)
}