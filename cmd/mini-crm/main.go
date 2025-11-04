package main

import (
	"flag"
	"fmt"

	"github.com/thomas-stecinski/crm_go_project/internal/cli"
	"github.com/thomas-stecinski/crm_go_project/internal/contacts"
	"github.com/thomas-stecinski/crm_go_project/internal/notifier"
)

func main() {

	// slice de notifier.Notifier
	notifiers := []notifier.Notifier{
		notifier.EmailNotifier{
			From:     "no-reply@shop.com",
			To:       "client@example.com",
			Subject:  "Confirmation de commande",
			SMTPHost: "smtp.shop.com",
			SMTPPort: 587,
		},
		notifier.SmsNotifier{
			From:     "ShopBot",
			To:       "+33612345678",
			Provider: "Twilio",
		},
	}

	message := "Votre commande a été expédiée !"

	for _, n := range notifiers {
		n.Send(message)
	}

	add := flag.Bool("add", false, "Ajouter un contact en mode non interactif")
	id := flag.Int("id", 0, "ID du contact (optionnel, auto si 0)")
	name := flag.String("name", "", "Nom du contact")
	email := flag.String("email", "", "Email du contact")
	data := flag.String("data", "data/contacts.json", "Chemin du fichier JSON de stockage")
	flag.Parse()

	store := contacts.NewStore(*data)

	if *add {
		c := contacts.Contact{ID: *id, Nom: *name, Email: *email}
		if !c.IsValid() {
			fmt.Println("Erreur: les flags -name et -email sont obligatoires avec -add")
			return
		}
		if _, err := store.Add(c); err != nil {
			fmt.Println("Erreur:", err)
			return
		}
		for _, c := range store.All() {
			fmt.Printf("ID: %d | Nom: %s | Email: %s\n", c.ID, c.Nom, c.Email)
		}
		return
	}

	app := &cli.App{Store: store}
	app.RunInteractive()
}
