package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/thomas-stecinski/crm_go_project/internal/contacts"
)

type App struct {
	Store contacts.ContactStore
}

func (a *App) RunInteractive() {
	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		fmt.Print("Votre choix: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			a.add(reader)
		case "2":
			a.list()
		case "3":
			a.delete(reader)
		case "4":
			a.update(reader)
		case "5":
			fmt.Println("Au revoir!")
			return
		default:
			fmt.Println("Choix invalide.")
		}
		fmt.Println()
	}
}

func printMenu() {
	fmt.Println("=== Mini-CRM ===")
	fmt.Println("1) Ajouter un contact (ID, Nom, Email)")
	fmt.Println("2) Lister tous les contacts")
	fmt.Println("3) Supprimer un contact par son ID")
	fmt.Println("4) Mettre à jour un contact")
	fmt.Println("5) Quitter")
}

func (a *App) add(r *bufio.Reader) {
	fmt.Print("Nom: ")
	nom, _ := r.ReadString('\n')
	fmt.Print("Email: ")
	email, _ := r.ReadString('\n')
	fmt.Print("ID (laisser vide pour auto): ")
	rawID, _ := r.ReadString('\n')

	nom = strings.TrimSpace(nom)
	email = strings.TrimSpace(email)
	rawID = strings.TrimSpace(rawID)

	var id int
	if rawID != "" {
		i, err := strconv.Atoi(rawID)
		if err != nil {
			fmt.Println("ID invalide.")
			return
		}
		id = i
	}

	c := contacts.Contact{ID: id, Nom: nom, Email: email}
	if !c.IsValid() {
		fmt.Println("Nom et Email sont obligatoires.")
		return
	}
	if _, err := a.Store.Add(c); err != nil {
		fmt.Println("Erreur:", err)
		return
	}
	fmt.Println("Contact ajouté.")
}

func (a *App) list() {
	all := a.Store.All()
	if len(all) == 0 {
		fmt.Println("Aucun contact.")
		return
	}
	for _, c := range all {
		fmt.Printf("ID: %d | Nom: %s | Email: %s\n", c.ID, c.Nom, c.Email)
	}
}

func (a *App) delete(r *bufio.Reader) {
	fmt.Print("ID à supprimer: ")
	raw, _ := r.ReadString('\n')
	raw = strings.TrimSpace(raw)
	id, err := strconv.Atoi(raw)
	if err != nil {
		fmt.Println("ID invalide.")
		return
	}
	if err := a.Store.Delete(id); err != nil {
		fmt.Println("Erreur:", err)
		return
	}
	fmt.Println("Contact supprimé.")
}

func (a *App) update(r *bufio.Reader) {
	fmt.Print("ID à mettre à jour: ")
	raw, _ := r.ReadString('\n')
	raw = strings.TrimSpace(raw)
	id, err := strconv.Atoi(raw)
	if err != nil {
		fmt.Println("ID invalide.")
		return
	}

	c, err := a.Store.Get(id)
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	fmt.Printf("Nouveau nom (actuel: %s, Enter = garder): ", c.Nom)
	newNom, _ := r.ReadString('\n')
	fmt.Printf("Nouvel email (actuel: %s, Enter = garder): ", c.Email)
	newEmail, _ := r.ReadString('\n')

	newNom = strings.TrimSpace(newNom)
	newEmail = strings.TrimSpace(newEmail)

	if newNom != "" {
		c.Nom = newNom
	}
	if newEmail != "" {
		c.Email = newEmail
	}

	if err := a.Store.Update(id, c); err != nil {
		fmt.Println("Erreur:", err)
		return
	}
	fmt.Println("Contact mis à jour.")
}
