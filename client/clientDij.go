package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Nd struct {
	nom string
}

type Lien struct {
	dep   Nd
	fin   Nd
	poids int
}

type GraphSommet struct {
	Job     bool
	idGraph int
	Sommet  Nd
}

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

//faire une fonction qui lit le graph et le met dans un tableau (et attribue un poid random au graph)

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

		//Client envoie le graphe lu, ligne par ligne jusqua EOF au serveur

		/*for i := 0; i < len(tabGraph(graph)); i++ {
			if i < len((tabGraph(graph)))-1 {
				io.WriteString(conn, fmt.Sprintf(tabGraph(graph)[i]))
			} else {
				io.WriteString(conn, fmt.Sprintf("EOF"))
			}
		}*/

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
