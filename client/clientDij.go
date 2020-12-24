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

		//Client lit le graphe texte et l'envoi en format string ligne par ligne
		f, err := os.Open("graph.txt")
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			io.WriteString(conn, fmt.Sprintf("%s\n", scanner.Text()))
		}

		//Apres l'envoi, le client attend une reponse du serveur avec les chemins les plus courts
		resultString, err := reader.ReadString('$')
		if err != nil {
			fmt.Printf("DEBUG MAIN could not read from server")
			os.Exit(1)
		}
		resultString = strings.TrimSuffix(resultString, "$")
		fmt.Printf("#DEBUG server replied :\n%s\n", resultString)
		time.Sleep(1000 * time.Millisecond)

		//Stockage des infos recues dans un fichier texte
		file, err := os.OpenFile("Dijsktra.txt", os.O_CREATE|os.O_WRONLY, 0600)
		//CREATE pour créer fichier s'il n'existe pas deja
		//WR ONLY pour  rendre le fichier (dans le programme) accessible en écriture seulement
		//0600 : permission -rw pour le fichier
		defer file.Close() // on ferme automatiquement à la fin de notre programme

		if err != nil {
			panic(err)
		}

		_, err = file.WriteString(resultString) // écrire dans le fichier
		if err != nil {
			panic(err)
		}
	}

}
