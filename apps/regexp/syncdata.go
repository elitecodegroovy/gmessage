package main

import (
	"os"
	"log"
	"bufio"
	"fmt"
	"regexp"
)

func loadData(){
	file, err := os.Open("./data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	//data := []string{}
	for scanner.Scan() {
		//get(scanner.Text())
		i++
		r, _ :=regexp.Compile(`\{[\S\s]+}`)
		matches := r.FindString(scanner.Text())

		fmt.Println("", i, ":", matches)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("total number: ", i)
}

func sync(){
	loadData()
}
func main(){
	sync()
}
