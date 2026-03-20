package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const dataFile = "expenses.json"

type Expense struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
}

type Store struct {
	Expenses []Expense `json:"expenses"`
	NextID   int       `json:"next_id"`
}

func loadStore() Store {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return Store{NextID: 1}
	}
	var store Store
	json.Unmarshal(data, &store)
	return store
}

func (s *Store) save() {
	data, _ := json.MarshalIndent(s, "", "  ")
	os.WriteFile(dataFile, data, 0644)
}

func (s *Store) add(desc, category string, amount float64) Expense {
	e := Expense{
		ID:          s.NextID,
		Description: desc,
		Amount:      amount,
		Category:    strings.ToLower(category),
		Date:        time.Now(),
	}
	s.Expenses = append(s.Expenses, e)
	s.NextID++
	s.save()
	return e
}

func (s *Store) delete(id int) bool {
	for i, e := range s.Expenses {
		if e.ID == id {
			s.Expenses = append(s.Expenses[:i], s.Expenses[i+1:]...)
			s.save()
			return true
		}
	}
	return false
}

func (s *Store) summary(month int) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tDate\tCategory\tAmount\tDescription")
	fmt.Fprintln(w, "--\t----\t--------\t------\t-----------")

	total := 0.0
	byCategory := map[string]float64{}
	count := 0

	for _, e := range s.Expenses {
		if month > 0 && int(e.Date.Month()) != month {
			continue
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t$%.2f\t%s\n",
			e.ID, e.Date.Format("2006-01-02"), e.Category, e.Amount, e.Description)
		total += e.Amount
		byCategory[e.Category] += e.Amount
		count++
	}
	w.Flush()

	if count == 0 {
		fmt.Println("No expenses found.")
		return
	}

	fmt.Printf("\n── Summary ─────────────────\n")
	fmt.Printf("  Total expenses: %d\n", count)
	fmt.Printf("  Total amount:   $%.2f\n", total)

	// Sort categories by spend
	type catSpend struct{ cat string; amount float64 }
	cats := make([]catSpend, 0, len(byCategory))
	for cat, amt := range byCategory {
		cats = append(cats, catSpend{cat, amt})
	}
	sort.Slice(cats, func(i, j int) bool { return cats[i].amount > cats[j].amount })

	fmt.Println("\n  By category:")
	for _, c := range cats {
		bar := strings.Repeat("█", int(c.amount/total*20))
		fmt.Printf("    %-12s $%7.2f  %s\n", c.cat, c.amount, bar)
	}
}

func printHelp() {
	fmt.Println(`
expense-tracker — manage your expenses from the terminal

Usage:
  add <description> <amount> [category]    Add an expense
  list                                     List all expenses
  summary [month]                          Show summary (optional month 1-12)
  delete <id>                              Delete an expense
  help                                     Show this help

Examples:
  expense-tracker add "Coffee" 4.50 food
  expense-tracker add "AWS bill" 42.00 cloud
  expense-tracker list
  expense-tracker summary
  expense-tracker summary 11
  expense-tracker delete 3
`)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		return
	}

	store := loadStore()

	switch args[0] {
	case "add":
		if len(args) < 3 {
			fmt.Println("Usage: add <description> <amount> [category]")
			os.Exit(1)
		}
		amount, err := strconv.ParseFloat(args[2], 64)
		if err != nil || amount <= 0 {
			fmt.Println("Error: amount must be a positive number")
			os.Exit(1)
		}
		category := "general"
		if len(args) >= 4 {
			category = args[3]
		}
		e := store.add(args[1], category, amount)
		fmt.Printf("✅ Added expense #%d: %s — $%.2f [%s]\n", e.ID, e.Description, e.Amount, e.Category)

	case "list":
		store.summary(0)

	case "summary":
		month := 0
		if len(args) >= 2 {
			m, err := strconv.Atoi(args[1])
			if err != nil || m < 1 || m > 12 {
				fmt.Println("Error: month must be 1-12")
				os.Exit(1)
			}
			month = m
		}
		store.summary(month)

	case "delete":
		if len(args) < 2 {
			fmt.Println("Usage: delete <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: id must be a number")
			os.Exit(1)
		}
		if store.delete(id) {
			fmt.Printf("🗑️  Deleted expense #%d\n", id)
		} else {
			fmt.Printf("Error: expense #%d not found\n", id)
			os.Exit(1)
		}

	default:
		printHelp()
	}
}
