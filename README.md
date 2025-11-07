# CRM GO PROJET

Mini projet pour apprendre Go à travers un CRM en ligne de commande.  
Le but : manipuler les concepts clés du langage (interfaces, CLI, persistance, configuration).

## Structure du projet (pour Axelle)

```bash
cmd/               # CLI avec Cobra
internal/contacts/ # Logique métier + stores JSON/GORM
internal/db/       # Connexion SQLite
internal/cli/      # Menu interactif
config.yaml        # Configuration par défaut
main.go            # Entrée du programme
```

## Run Locally

Clone the project

```bash
  git clone https://github.com/thomas-stecinski/crm_go_project.git
```

Go to the project directory

```bash
  cd crm_go_project
```

### Start the program in interactive mode

```bash
go mod tidy
go run . interactive --backend json --data data/contacts.json
```

### Use the program in non interative mode JSON

```bash
go run . contact add --name "Alice" --email "alice@mail.com"
go run . contact list
go run . contact update --id 1 --name "Alice Cooper"
go run . contact delete --id 1
```

## Backend Storage

### JSON (par défaut)

Local persistence via data/contacts.json

Non Interactive :
```bash
go run . contact add --name "Bob" --email "bob@mail.com"
```

Interactive :
```bash
go run . interactive --backend json --data data/contacts.json
```

### SQLite (via GORM)

Backend via GORM.
Interactive :

```bash
go run . interactive --backend gorm --dsn data/contacts.db
```
Non interactive :
```bash
go run . --backend gorm --dsn data/contacts.db contact add --name "Charlie" --email "charlie@mail.com"
go run . --backend gorm --dsn data/contacts.db contact list
```

## Configuration (Viper)

automaticaly read file config.yaml :

```bash
backend: gorm
dsn: data/contacts.db
data: data/contacts.json
```

Overloading possible

- Flags CLI

```bash
export CRM_BACKEND=gorm
export CRM_DSN=data/prod.db
```

- Env variables

```bash
go run . --backend json --data data/test.json contact list
```

## Authors

- [@Thomas](https://www.github.com/thomas-stecinski)
- [@Boris](https://www.github.com/Lowinne)
