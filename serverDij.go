package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func getArgs() int {
	//On verifie qu'on a bien 2 arguments
	if len(os.Args) != 2 {
		fmt.Printf("Usage: go run server.go <portnumber>\n")
		os.Exit(1)
	} else {
		//on verifie que l'argument est bien un int
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[1])
		portNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Usage: go run server.go <portnumber>\n")
			os.Exit(1)
		} else {
			return portNumber
		}

	}
	//PFR should never be reached
	return -1
}

func main() {
	port := getArgs()
	fmt.Printf("#DEBUG MAIN Creating TCP Server on port %d\n", port)
	//Create a port string that lets us accept connection on all interfaces of the host
	portString := fmt.Sprintf(":%s", strconv.Itoa(port))
	fmt.Printf("#DEBUG MAIN PORT STRING |%s|\n", portString)

	ln, err := net.Listen("tcp", portString)
	if err != nil {
		fmt.Printf("#DEBUG MAIN Could not create listener\n")
		panic(err)
	}

	//If we're here, we did not panic and ln is a valid listener
	//connum = id de connexion
	connum := 1

	for {
		fmt.Printf("#DEBUG MAIN Accepting next connection\n")
		conn, errconn := ln.Accept()

		if errconn != nil {
			fmt.Printf("DEBUG MAIN Error when accepting next connection\n")
			panic(errconn)

		}

		//If we're here, we did not panic and conn is a valid handler to the new connection

		go handleConnection(conn, connum)
		connum += 1

	}
}

func handleConnection(connection net.Conn, connum int) {
	//PFR !!!
	defer connection.Close()
	connReader := bufio.NewReader(connection)
	//    if err !=nil{
	//        fmt.Printf("#DEBUG %d handleConnection could not create reader\n", connum)
	//        return
	//    }

	for {
		//on lit la ligne recue du client
		inputLine, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Printf("#DEBUG %d RCV ERROR no panic, just a client\n", connum)
			fmt.Printf("#DEBUG Error :|%s|\n", err.Error())
			break
		}

		//print la ligne recue
		inputLine = strings.TrimSuffix(inputLine, "\n")
		fmt.Printf("#DEBUG %d RCV |%s|\n", connum, inputLine)

		//Stocke la ligne recue

		//Construit un graphe grace a la ligne recue

		//Applique Dijkstra au graphe

		//ENVOI DU GRAPH AU CLIENT
		//print la ligne a envoyer au client

		//envoie les plus courts chemins au client

		//truc du prof pour exemple
		splitLine := strings.Split(inputLine, " ")
		returnedString := splitLine[len(splitLine)-1]
		fmt.Printf("#DEBUG %d SND |%s|\n", connum, returnedString)

		io.WriteString(connection, fmt.Sprintf("%s\n", returnedString))
	}

}
