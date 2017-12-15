package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func SwapCase(str string) string {
	buf := &bytes.Buffer{}
	for _, r := range str {
		if unicode.IsUpper(r) {
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(unicode.ToUpper(r))
		}
	}
	return buf.String()
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

//Defer is not the same as deferred (or futures/promises) in other languages.
func StartDemo() {
	aValue := new(int)

	defer fmt.Println(*aValue)
	for i := 0; i < 100; i++ {
		*aValue++
	}
	FormatTime()
}

func ReadFile(filePath string) {
	file, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		// skip all line starting without line 'http'
		//		if equal := strings.Index(line, "http"); equal < 0 {
		//			fmt.Print(line)
		//		}

		//alternatively, only print line starting with 'http'
		if equal := strings.Index(line, "http"); equal >= 0 {
			fmt.Print(line)
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}

	}

}

func Substring(s string, limit int) string {
	reader := strings.NewReader(s)

	//create buffer with specified limit of characters
	buff := make([]byte, limit)

	n, err := io.ReadAtLeast(reader, buff, limit)
	if err != nil {
		fmt.Println("error in ReadAtLeast", err)
	}
	if n != 0 {
		return string(buff)
	} else {
		//nothing happens, return original string
		return s
	}
}
