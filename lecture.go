package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type nd struct {
	nom string
}

type Lien struct {
	dep   string
	fin   string
	poids int
	id    int
}

func getGraph() {

}

var maMap map[string]Lien

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
		resq := strings.TrimSuffix(tsep[3], "\r\n")
		res4, err := strconv.Atoi(resq)
		if err == nil {
			//fmt.Println(res4)
		}
		fmt.Println("Depart:", res1, "|", "Fin:", res2, "|", "ID: ", res3, "|", "Distance:", res4)
	}
	if err != nil {
		fmt.Printf("DEBUG ERROR TYPE %d", err)
	}

}
