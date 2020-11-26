package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)
/*
type Nd struct {
	nom string
}

type Lien struct {
	dep   Nd
	fin   Nd
	poids int
}

type GraphSommet struct {
	Job bool
	idGraph int
	Sommet Nd
}

var graph map[Nd][]Lien = makeGraph()
*/
func getArgs() int {
	//Make sure we have an argument
	if len(os.Args) != 2 {
		fmt.Printf("Usage: go run client.go <portnumber>\n")
		os.Exit(1)
	} else {
		//Make sure the argument is a valid integer, return it
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[1])
		portNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Usage: go run client.go <portnumber>\n")
			os.Exit(1)
		} else {
			return portNumber
		}
	}
	//Should never be reached
	return -1
}

//fonction qui permet de créer un graphe avec un fichier .txt
func makeGraph() map[Nd][]Lien {
	graph := make(map[Nd][]Lien)
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
		res1 := Nd{nom: tsep[0]}
		res2 := Nd{nom: tsep[1]}
		resq := strings.TrimSuffix(tsep[2], "\r\n")
		res3, err := strconv.Atoi(resq)
		if err == nil {
			//fmt.Println(res3)
		}

		lien := Lien{res1, res2, res3}

		graph[res1] = append(graph[res1], lien)

		//fmt.Println("Depart:", res1, "|", "Fin:", res2, "|", "ID: ", res3, "|", "Distance:")
	}
	return graph

}

//fonction qui nous donne la liste des noeuds du graphe en entrée
func ListeNd(graph map[Nd][]Lien) []Nd {
	keys := make([]Nd, 0, len(graph))
	for k := range graph {
		keys = append(keys, k)
	}
	return keys
}

func main() {
	//Get the port number
	port := getArgs()
	fmt.Printf("#DEBUG DIALING TCP Server on port %d\n", port)
	//Create the target port string
	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))
	fmt.Printf("#DEBUG MAIN PORT STRING |%s|\n", portString)
	//Connect
	conn, err := net.Dial("tcp", portString)
	if err != nil {
		//Leave if connection does not work
		fmt.Printf("#DEBUG MAIN could not connect\n")
		os.Exit(1)
	} else {

		defer conn.Close()
		reader := bufio.NewReader(conn)
		fmt.Printf("#DEBUG MAIN connected\n")
		//Client lit le graphe texte et créé un graphe
		makeGraph()

		//Client envoie le graphe lu, ligne par ligne jusqua EOF au serveur
		for i := 0; i < len(ListeNd(graph)); i++ { //checker si ListeNd nous donne bien le nb de lignes du graphe
			if i < len((ListeNd(graph)))-1 {
				io.WriteString(conn, fmt.Sprintf(graph[i]))
			} else {
				io.WriteString(conn, fmt.Sprintf("EOF"))
			}
		}
		//Apres l'envoi, le client attend une reponse du serveur avec les chemins les plus courts
		resultString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("DEBUG MAIN could not read from server")
			os.Exit(1)
		}
		resultString = strings.TrimSuffix(resultString, "\n")
		fmt.Printf("#DEBUG server replied : |%s|\n", resultString)
		time.Sleep(1000 * time.Millisecond)
			
		//Stockage des infos recues dans un fichier texte

		
		}

	}

}
