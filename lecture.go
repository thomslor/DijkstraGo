package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

//type nd struct {
//	nom string
//}

type Lien struct {
	dep   string
	fin   string
	poids int
	id    int
}

func getGraph(graph map[string][]Lien) {
	for i := range graph {
		fmt.Println(graph[i])
	}
}

func makeGraph() map[string][]Lien {
	graph := make(map[string][]Lien)
	idl := 0
	f, err := os.Open("graph.txt")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	rd := bufio.NewReader(f)

	for {

		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		tsep := strings.Split(line, ";")
		res1 := tsep[0]
		res2 := tsep[1]
		resq := strings.TrimSuffix(tsep[2], "\r\n")
		res3, err := strconv.Atoi(resq)
		if err == nil {
			//fmt.Println(res3)
		}

		lien := Lien{res1, res2, res3, idl}
		idl += 1
		graph[res1] = append(graph[res1], lien)

		//fmt.Println("Depart:", res1, "|", "Fin:", res2, "|", "ID: ", res3, "|", "Distance:")
	}
	return graph

}


func main() {
	getGraph(makeGraph())





}
