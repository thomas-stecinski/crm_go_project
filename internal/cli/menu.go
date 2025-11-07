package cli

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"

	"github.com/thomas-stecinski/crm_go_project/internal/contacts"
	"github.com/thomas-stecinski/crm_go_project/internal/db"
	"github.com/thomas-stecinski/crm_go_project/internal/notifier"
)

// App interactive CRM
type App struct {
	Store    contacts.ContactStore
	Notifier []notifier.Notifier
	Backend  string
}

// === Couleurs ===
var (
	title      = color.New(color.FgHiCyan, color.Bold).SprintFunc()
	label      = color.New(color.FgHiWhite, color.Bold).SprintFunc()
	success    = color.New(color.FgHiGreen).SprintFunc()
	warning    = color.New(color.FgHiYellow).SprintFunc()
	errorText  = color.New(color.FgHiRed, color.Bold).SprintFunc()
	info       = color.New(color.FgHiMagenta).SprintFunc()
	sectionBar = color.New(color.FgHiBlack).SprintFunc()
)

// === Entrée principale ===
func (a *App) RunInteractive() {
	reader := bufio.NewReader(os.Stdin)
	printBanner(a)

	for {
		printMenu(a)
		fmt.Print(label("Votre choix: "))

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			a.add(reader)
		case "2":
			a.list()
		case "3":
			a.search(reader)
		case "4":
			a.update(reader)
		case "5":
			a.delete(reader)
		case "6":
			a.exportMenu(reader)
		case "7":
			a.switchBackend(reader)
		case "8":
			fmt.Println(success("\n👋 Merci d’avoir utilisé Mini CRM !"))
			return
		default:
			fmt.Println(errorText("❌ Choix invalide."))
		}

		fmt.Println()
		waitEnter(reader)
	}
}

// === Bannière d’accueil ===
func printBanner(a *App) {
	fmt.Println()
	fmt.Println(title("📇 MINI CRM — Mode Interactif"))
	fmt.Println(sectionBar("───────────────────────────────────────────────"))
	fmt.Printf("%s %s\n", info("Backend actuel :"), detectBackend(a.Store))
	fmt.Println(sectionBar("───────────────────────────────────────────────"))
}

// Détection du backend via reflection
func detectBackend(store contacts.ContactStore) string {
	if store == nil {
		return warning("Inconnu")
	}
	t := reflect.TypeOf(store).Elem().Name()
	switch t {
	case "Store":
		return success("JSON Store (contacts.json)")
	case "GormStore":
		return success("GORM Store (SQLite)")
	default:
		return warning("Store personnalisé : " + t)
	}
}

// === Menu principal ===
func printMenu(a *App) {
	total := len(a.Store.All())
	fmt.Println()
	fmt.Println(title("=== MENU PRINCIPAL ==="))
	fmt.Println(sectionBar("───────────────────────────────────────────────"))
	fmt.Println("1️⃣  Ajouter un contact")
	fmt.Println("2️⃣  Lister les contacts")
	fmt.Println("3️⃣  Rechercher un contact")
	fmt.Println("4️⃣  Mettre à jour un contact")
	fmt.Println("5️⃣  Supprimer un contact")
	fmt.Println("6️⃣  Exporter les contacts (CSV / JSON)")
	fmt.Println("7️⃣  Changer de backend (JSON ↔ GORM)")
	fmt.Println("8️⃣  Quitter")
	fmt.Println(sectionBar("───────────────────────────────────────────────"))
	fmt.Printf("%s %d contact(s) enregistré(s)\n\n", info("📊 Total :"), total)
}

// === Ajouter un contact ===
func (a *App) add(r *bufio.Reader) {
	fmt.Println(title("\n➕ AJOUT D’UN CONTACT"))
	fmt.Print(label("Nom: "))
	nom, _ := r.ReadString('\n')
	fmt.Print(label("Email: "))
	email, _ := r.ReadString('\n')

	nom = strings.TrimSpace(nom)
	email = strings.TrimSpace(email)

	c := contacts.Contact{Nom: nom, Email: email}
	if !c.IsValid() {
		fmt.Println(errorText("❌ Nom et Email obligatoires."))
		return
	}
	if _, err := a.Store.Add(c); err != nil {
		fmt.Println(errorText("❌ Erreur:"), err)
		return
	}
	fmt.Println(success("✅ Contact ajouté avec succès."))

	// Notification automatique
	for _, n := range a.Notifier {
		go n.Send(fmt.Sprintf("Nouveau contact ajouté : %s (%s)", c.Nom, c.Email))
	}
}

// === Lister ===
func (a *App) list() {
	fmt.Println(title("\n📋 LISTE DES CONTACTS"))
	all := a.Store.All()
	if len(all) == 0 {
		fmt.Println(warning("💡 Aucun contact enregistré."))
		return
	}
	for _, c := range all {
		fmt.Printf("%s %-3d | %s %-15s | %s %s\n",
			color.HiBlackString("ID:"), c.ID,
			color.HiBlackString("Nom:"), c.Nom,
			color.HiBlackString("Email:"), c.Email)
	}
	fmt.Printf("\n%s %d contact(s)\n", success("✅ Total :"), len(all))
}

// === Recherche ===
func (a *App) search(r *bufio.Reader) {
	fmt.Println(title("\n🔍 RECHERCHE DE CONTACT"))
	fmt.Print(label("Terme à rechercher (nom ou email) : "))
	query, _ := r.ReadString('\n')
	query = strings.ToLower(strings.TrimSpace(query))

	results := []contacts.Contact{}
	for _, c := range a.Store.All() {
		if strings.Contains(strings.ToLower(c.Nom), query) || strings.Contains(strings.ToLower(c.Email), query) {
			results = append(results, c)
		}
	}

	if len(results) == 0 {
		fmt.Println(warning("💡 Aucun résultat trouvé."))
		return
	}

	fmt.Printf("%s %d résultat(s) trouvé(s) :\n", success("✅"), len(results))
	for _, c := range results {
		fmt.Printf(" - [%d] %s <%s>\n", c.ID, c.Nom, c.Email)
	}
}

// === Supprimer ===
func (a *App) delete(r *bufio.Reader) {
	fmt.Println(title("\n🗑️ SUPPRESSION D’UN CONTACT"))
	fmt.Print(label("ID à supprimer: "))
	raw, _ := r.ReadString('\n')
	raw = strings.TrimSpace(raw)
	id, err := strconv.Atoi(raw)
	if err != nil {
		fmt.Println(errorText("❌ ID invalide."))
		return
	}

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Confirmer la suppression du contact %d (y/n)", id),
		IsConfirm: true,
	}
	_, confirmErr := prompt.Run()
	if confirmErr != nil {
		fmt.Println(warning("❎ Suppression annulée."))
		return
	}

	if err := a.Store.Delete(id); err != nil {
		fmt.Println(errorText("❌ Erreur:"), err)
		return
	}
	fmt.Println(success("🗑️  Contact supprimé avec succès."))
}

// === Mettre à jour ===
func (a *App) update(r *bufio.Reader) {
	fmt.Println(title("\n✏️ MISE À JOUR D’UN CONTACT"))
	fmt.Print(label("ID à mettre à jour: "))
	raw, _ := r.ReadString('\n')
	raw = strings.TrimSpace(raw)
	id, err := strconv.Atoi(raw)
	if err != nil {
		fmt.Println(errorText("❌ ID invalide."))
		return
	}

	c, err := a.Store.Get(id)
	if err != nil {
		fmt.Println(errorText("❌ Erreur:"), err)
		return
	}

	fmt.Printf("%s (%s): ", label("Nouveau nom"), c.Nom)
	newNom, _ := r.ReadString('\n')
	fmt.Printf("%s (%s): ", label("Nouvel email"), c.Email)
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
		fmt.Println(errorText("❌ Erreur:"), err)
		return
	}
	fmt.Println(success("✏️  Contact mis à jour avec succès."))
}

// === Export (CSV / JSON) avec barre de progression ===
func (a *App) exportMenu(r *bufio.Reader) {
	fmt.Println(title("\n📤 EXPORT DES CONTACTS"))
	fmt.Print(label("Format (csv/json): "))
	format, _ := r.ReadString('\n')
	format = strings.ToLower(strings.TrimSpace(format))
	all := a.Store.All()
	if len(all) == 0 {
		fmt.Println(warning("💡 Aucun contact à exporter."))
		return
	}

	p := mpb.New()
	bar := p.New(int64(len(all)),
		mpb.BarStyle().Rbound("▏").Filler("█").Tip("▎"),
		mpb.PrependDecorators(
			decor.Name("Export : "),
			decor.Percentage(),
		),
	)
	outDir := "export"
	_ = os.MkdirAll(outDir, 0o755)

	var fileName string
	switch format {
	case "csv":
		fileName = filepath.Join(outDir, fmt.Sprintf("contacts_%d.csv", time.Now().Unix()))
		f, _ := os.Create(fileName)
		defer f.Close()
		w := csv.NewWriter(f)
		defer w.Flush()
		w.Write([]string{"ID", "Nom", "Email"})
		for _, c := range all {
			_ = w.Write([]string{strconv.Itoa(c.ID), c.Nom, c.Email})
			bar.Increment()
			time.Sleep(20 * time.Millisecond)
		}
	case "json":
		fileName = filepath.Join(outDir, fmt.Sprintf("contacts_%d.json", time.Now().Unix()))
		f, _ := os.Create(fileName)
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")

		_ = enc.Encode(all)
		for range all {
			bar.Increment()
			time.Sleep(15 * time.Millisecond)
		}
	default:
		fmt.Println(errorText("❌ Format invalide."))
		return
	}
	p.Wait()
	fmt.Println(success("✅ Export terminé → " + fileName))
}

// === Changer de backend (JSON <-> GORM) ===
func (a *App) switchBackend(r *bufio.Reader) {
	fmt.Println(title("\n📂 CHANGER DE BACKEND"))
	fmt.Print(label("Nouveau backend (json/gorm): "))
	newBackend, _ := r.ReadString('\n')
	newBackend = strings.TrimSpace(strings.ToLower(newBackend))

	switch newBackend {
	case "json":
		a.Store = contacts.NewStore("data/contacts.json")
		a.Backend = "json"
		fmt.Println(success("✅ Passage au JSON Store (data/contacts.json)"))
	case "gorm":
		dbConn, err := db.OpenSQLite("data/contacts.db")
		if err != nil {
			fmt.Println(errorText("❌ Erreur SQLite:"), err)
			return
		}
		newStore, err := contacts.NewGormStore(dbConn)
		if err != nil {
			fmt.Println(errorText("❌ Erreur GORM:"), err)
			return
		}
		a.Store = newStore
		a.Backend = "gorm"
		fmt.Println(success("✅ Passage au GORM Store (data/contacts.db)"))
	default:
		fmt.Println(errorText("❌ Backend inconnu."))
		return
	}

	fmt.Println(sectionBar("───────────────────────────────────────────────"))
	fmt.Printf("%s %s\n", info("Backend actif :"), detectBackend(a.Store))
	fmt.Println(sectionBar("───────────────────────────────────────────────"))
}

// === Attente avant retour menu ===
func waitEnter(r *bufio.Reader) {
	fmt.Print(sectionBar("\n↩️  Appuyez sur Entrée pour revenir au menu..."))
	r.ReadString('\n')
}
