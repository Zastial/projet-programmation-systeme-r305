# Changements effectués dans le jeu  

## Changement dans des fichiers existants :

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

## Ajout de nouveaux fichiers :

### **reseau.go :** 

- se charge de l'interaction entre le client et le serveur

### **serveur.go :**

- se charge de la communication avec les joueurs