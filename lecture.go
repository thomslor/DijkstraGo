package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Lien struct {
	dep   string
	fin   string
	poids int
}

func main() {
	f, err := os.Open("graph.txt")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	rd := bufio.NewReader(f)

	for {

		line, err := rd.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}
		tsep := strings.Split(line, ";")
		res1 := tsep[0]
		res2 := tsep[1]
		res3 := tsep[2]
		//res4 := tsep[3]
		//res4, err := strconv.ParseInt("55", 10, 64)
		//if err == nil {
		//	fmt.Println(res4)
		//}
		fmt.Println("Depart: ", res1, "|", "Fin: ", res2, "|", "ID: ", res3, "|", "Distance: ", res4)
	}
}
