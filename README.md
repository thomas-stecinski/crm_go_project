# 🧩 CRM GO PROJECT

Projet d’apprentissage du langage **Go** à travers la création d’un **CRM en ligne de commande** complet.

Ce projet met en pratique :
- 🧱 les **interfaces** et l’architecture modulaire  
- ⚙️ la **concurrence** et les **goroutines**  
- 💬 un **CLI interactif** avec **Cobra**  
- 🗄️ une **persistance JSON** ou **SQLite (via GORM)**  
- 🧩 une **configuration centralisée avec Viper**  
- ✨ une **interface console stylisée** avec couleurs, confirmation, et barre de progression  

---

## 🚀 Démarrage rapide

### 1️⃣ Installation

```bash
git clone https://github.com/thomas-stecinski/crm_go_project.git
cd crm_go_project
go mod tidy
```

---

## 🧭 Utilisation

### ▶️ Mode interactif

Lance le CRM avec un backend choisi au démarrage :

#### JSON (par défaut)
```bash
go run . interactive --backend json --data data/contacts.json
```

#### SQLite (GORM)
```bash
go run . interactive --backend gorm --dsn data/contacts.db
```

💡 Le mode interactif affiche :
- Le backend actif (JSON ou SQLite)
- Un menu coloré et ergonomique
- Des confirmations de suppression
- Des notifications automatiques (email / SMS)
- Une barre de progression pour les exports

---

### ⚙️ Mode non-interactif (CLI)

Tu peux aussi exécuter directement des commandes sans le menu :

```bash
# Ajouter un contact
go run . contact add --name "Alice" --email "alice@mail.com"

# Lister les contacts
go run . contact list

# Mettre à jour un contact
go run . contact update --id 1 --name "Alice Cooper"

# Supprimer un contact
go run . contact delete --id 1
```

Sélection du backend :
```bash
--backend json --data data/contacts.json
--backend gorm --dsn data/contacts.db
```

Exemple :
```bash
go run . --backend gorm contact list
```

---

## 🧩 Fonctionnalités interactives

| Fonction | Description |
|-----------|--------------|
| ➕ **Ajouter un contact** | Ajoute un nouveau contact (nom + email) |
| 📋 **Lister les contacts** | Affiche tous les contacts enregistrés |
| 🔍 **Rechercher un contact** | Filtre par nom ou email |
| ✏️ **Mettre à jour** | Modifie un contact existant |
| 🗑️ **Supprimer** | Supprime un contact après confirmation |
| 📤 **Exporter** | Exporte les contacts (CSV ou JSON) avec barre de progression |
| 🔔 **Notifications** | Envoi automatique via `internal/notifier` (email/SMS simulés) |
| 📂 **Changer de backend** | Bascule à chaud entre JSON ↔ GORM |
| ⚙️ **Quitter** | Ferme le programme proprement |

---

## ⚙️ Configuration (Viper)

Le comportement du CLI est défini via un fichier `config.yaml`, les variables d’environnement, ou les flags.

### Exemple `config.yaml`
```yaml
backend: json
data: data/contacts.json
dsn: data/contacts.db
```

### Variables d’environnement
```bash
export CRM_BACKEND=gorm
export CRM_DSN=data/prod.db
```

### Flags CLI (prioritaires)
```bash
go run . --backend json --data data/test.json contact list
```

🔁 **Ordre de priorité :**
1. Flags CLI  
2. Variables d’environnement (`CRM_…`)  
3. Fichier `config.yaml`  
4. Valeurs par défaut (`json` + `data/contacts.json`)

---

## 🗄️ Backends de stockage

### JSON Store
- Persistance locale dans `data/contacts.json`
- Sauvegarde atomique (fichier `.tmp` + rename)
- Simple, sans dépendance

### GORM Store (SQLite)
- ORM via **GORM**
- Auto-migration du schéma
- Driver multiplateforme **glebarez/sqlite**

Exemples :
```bash
go run . --backend gorm --dsn data/contacts.db contact add --name "Boris" --email "boris@mail.com"
go run . --backend gorm --dsn data/contacts.db contact list
```

---

## 🧰 Structure du projet

```
crm_go_project/
├── cmd/                   # CLI avec Cobra
│   ├── root.go            # Commande principale
│   ├── interactive.go     # Lancement du menu interactif
│   ├── contact_*.go       # Sous-commandes CRUD
│
├── internal/
│   ├── cli/menu.go        # Mode interactif stylisé
│   ├── contacts/          # Modèle et Stores (JSON + GORM)
│   ├── db/db.go           # Connexion SQLite
│   └── notifier/          # Email / SMS simulés
│
├── config.yaml            # Configuration par défaut
├── data/                  # Données persistées
└── main.go                # Point d’entrée du programme
```

---

## 📤 Export des contacts

### Formats supportés :
- **CSV**
- **JSON (tableau complet)**

Exemple :
```json
[
  {
    "id": 1,
    "nom": "Boris Prince",
    "email": "boris@gmail.com"
  },
  {
    "id": 2,
    "nom": "Stecinski",
    "email": "thomas@gmail.com"
  }
]
```

Fichiers générés dans `export/` :
```
export/contacts_1730900000.json
export/contacts_1730900000.csv
```

Barre de progression :
```
Export :  ██████████████████████ 100%
✅ Export terminé → export/contacts_1730900000.json
```

---

## 💡 Conseils d’utilisation

| Action | Commande |
|--------|-----------|
| Démarrer sur SQLite | `go run . interactive --backend gorm --dsn data/contacts.db` |
| Démarrer sur JSON | `go run . interactive --backend json --data data/contacts.json` |
| Changer de backend à chaud | Option 7️⃣ dans le menu interactif |
| Exporter les données | Option 6️⃣ |
| Rechercher un contact | Option 3️⃣ |

---

## 👥 Auteurs

- [@Thomas Stecinski](https://github.com/thomas-stecinski)
- [@Boris Prince](https://github.com/Lowinne)

---

## 📄 Licence

Projet pédagogique — EFREI 2024-2025  
Usage libre à des fins d’apprentissage.
