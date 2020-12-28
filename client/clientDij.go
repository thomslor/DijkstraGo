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

func getArgs() []string {
	//on créé notre tableau de sortie
	var res []string
	//On verifie qu'on a bien 2 arguments
	if len(os.Args) != 3 {
		fmt.Printf("Usage: go run client.go <portnumber> <yourgraph.txt>\n")
		os.Exit(1)
	} else {

		//on verifie que le 1er argument est bien un int
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[1])
		_, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Usage: go run client.go <portnumber> <yourgraph.txt>\n")
			os.Exit(1)
		} else {
			//ajout du port au tableau de sortie
			res = append(res, os.Args[1])
		}
		//on vérifie si le 2e argument est un fichier du dossier courant
		fmt.Printf("#DEBUG ARGS Graph : %s\n", os.Args[2])
		//a coder

		//ajout du graph au tableau de sortie
		res = append(res, os.Args[2])
	}
	return res
}

func main() {
	//Get the port number & graph
	res := getArgs()
	port, _ := strconv.Atoi(res[0])
	graphName := res[1]
	//fmt.Println(port, graphName)
	fmt.Printf("#DEBUG DIALING TCP Server on port %d\n", port)
	//Create the target port string
	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))
	fmt.Printf("Connexion sur le port |%s|\n", portString)
	//Connect
	conn, err := net.Dial("tcp", portString)
	if err != nil {
		//Leave if connection does not work
		fmt.Printf("#DEBUG MAIN could not connect\n")
		os.Exit(1)
	} else {

		defer conn.Close()
		reader := bufio.NewReader(conn)
		envoi := ""
		fmt.Printf("Connexion réussie\n")

		//Client lit le graphe texte et l'envoi en format string ligne par ligne
		f, err := os.Open(graphName)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(f)
		nbLigne := 0
		for scanner.Scan() {
			//fmt.Println(scanner.Text())
			envoi += fmt.Sprintf("%s\n", scanner.Text())
			nbLigne++
			//io.WriteString(conn, fmt.Sprintf("%s\n", scanner.Text()))
		}
		io.WriteString(conn, fmt.Sprintf("%d\n", nbLigne))
		io.WriteString(conn, envoi)

		//Apres l'envoi, le client attend une reponse du serveur avec les chemins les plus courts
		resultString, err := reader.ReadString('$')
		if err != nil {
			fmt.Printf("Lecture impossible de la réponse du serveur")
			os.Exit(1)
		}
		resultString = strings.TrimSuffix(resultString, "$")
		// fmt.Printf("#DEBUG server replied :\n%s\n", resultString)
		ID := resultString[0]
		//fmt.Println(resultString[0])
		time.Sleep(1000 * time.Millisecond)

		//Stockage des infos recues dans un fichier texte
		file, err := os.OpenFile(fmt.Sprintf("Dijkstra%s.txt", ID), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		//CREATE pour créer fichier s'il n'existe pas deja
		//WR ONLY pour  rendre le fichier (dans le programme) accessible en écriture seulement
		//0600 : permission -rw pour le fichier
		defer file.Close() // on ferme automatiquement à la fin de notre programme
		if err != nil {
			panic(err)
		}

		_, err = file.WriteString(fmt.Sprintf("ID CLIENT : %s\n", resultString)) // écrire l'id du client + le graph
		if err != nil {
			panic(err)
		}

		fmt.Printf("Résultat disponible dans le fichier : Dijkstra%s.txt ", ID)

	}

}
