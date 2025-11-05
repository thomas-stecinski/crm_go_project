
# CRM GO PROJET

Un projet pour prendre un main le langage GO 


## Run Locally

Clone the project

```bash
  git clone https://github.com/thomas-stecinski/crm_go_project.git
```

Go to the project directory

```bash
  cd crm_go_project
```

Start the program   


```bash
go mod tidy
go run . interactive --data data/contacts.json
```
Non interative uses :

```bash
go run . contact add --name "Alice" --email "alice@mail.com"
go run . contact list
go run . contact update --id 1 --name "Alice Cooper"
go run . contact delete --id 1
```


## Authors

- [@Thomas](https://www.github.com/thomas-stecinski)
- [@Boris](https://www.github.com/Lowinne)
