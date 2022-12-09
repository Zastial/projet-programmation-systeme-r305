# **Projet Programmation Systeme**


# Changements effectués dans le jeu  

## **Changement dans des fichiers existants :**

### **runner.go :**

- Ajout de plusieurs fonctions get et set (lignes 47 à 61)

### **main.go :**

- Rajout d'une commande afin de fournir l'IP pour la connexion serveur (ligne 27)

### **game.go :**

- Ajout de plusieurs variables dans le struct Game : 
    - id_runner (type : *int*)
    - nbRunner (*int*)
    - IP (*string*)
    - conn (*net.Conn*)
    - receiveChannel (*chan string*)
    - good (*bool*)
- Changement de la valeur de **maxFrameInterval** de 0 à **frameInterval** (ligne 65)
- Ajout d'un channel (ligne 77)
- Changement de la fonction **Update**. Désormais chaque case est géré par une fonction implémentée dans **reseau.go** (ligne 112)
- Affichage des joueurs connectés à l'écran d'accueil (ligne 40)

## **Ajout de nouveaux fichiers :**

### **reseau.go :** 

- se charge de l'interaction entre le client et le serveur

### **serveur.go :**

- se charge de la communication avec les joueurs


## **Signification des codes :**

- 100 : Le joueur est connecté.

- 200 : Tous les joueurs sont connectés, la partie peut débuter.

- 3 + id + num.Couleur : Code envoyé par le joueur au serveur pour notifier qu'il a choisi sa couleur.

- 400 : Tous les joueurs ont sélectionné une couleur.
    - 4 + id + num.Couleur : Renvoie aux clients les couleurs choisies par les autres.

- 50+id : Le joueur est arrivé, sa course est terminée.

- 51 + id + pos.Joueur : Le joueur envoie sa position au serveur.

- 9 + id + pos.Joueur : Le serveur envoie la position d'un joueur à tous les autres.

- 600 : Les joueurs sont tous arrivés, la course est terminée.

- 70 + id : le joueur dit qu'il veut rejouer.

- 800 : Tous les joueurs veulent rejouer, la course se relance.


