# DijkstraGo
Projet réalisé par Naël BOUNIA, Estelle MONIER et Thomas LORRAIN

## Structure du projet
### Client
* Appel du client : go run clientDij.go portnumber yourgraph.txt
* Récupère un graph sous le format .txt avec un lien par ligne sous la forme suivante : SommetDeDépart;SommetD'Arrivée;PoidsDuLien
* Envoie les lignes au serveur situé sur le port défini par l'utilisateur, le serveur sait combien de lignes il doit recevoir via l'envoi du nombre de ligne par le client
* Récupère les lignes envoyées par le serveur et les copie dans un fichier texte unique

### Serveur
* Appel du serveur : go run serverDij.go portnumber
* Récupère les lignes des clients
* Génère des structures graph à partir de ces lignes
* Applique Dijkstra sur les graphs via une worker pool, un worker travaille sur un sommet d'un graph à la fois
* Renvoie le résultat de Dijkstra au bon client sous forme d'un string composé de plusieurs lignes

### Lecture

* Fichier permettant de tester Dijkstra sans TCP

### GénérateurGraph
* Permet de créer des sous-graphs au bon format à partir d'un graph de 5000 sommets avec le nombre de Sommets que l'on veut
* Appel du générateur : go run generateurgraph.go nbDeSommetsVoulu

## Comment utiliser le code ?
### Avec un ou plusieurs clients lancés non-simultanément
* Tout d'abord, générer un sous-graph avec le "générateur"
* placer le fichier créé dans le répertoire client
* Lancer le serveur dans un terminal
* Lancer le client dans un terminal
* Récupérer le résultat dans le fichier situé dans le répertoire client dont le nom est affiché sur le terminal

### Clients en simultanée
* Même procédure mais lancer les clients via le script bash script.bash situé dans le répertoire client (ATTENTION, le script est écrit pour un serveur sur le port 4337 et avec 2 clients en simultanée avec les graphs graph.txt et graph100.txt)

## Bonus
Le code (Client, Serveur, Générateur) est aussi disponible en format executable
