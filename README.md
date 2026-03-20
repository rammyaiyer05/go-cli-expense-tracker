# go-cli-expense-tracker

A fast, minimal **command-line expense tracker** written in Go. Add expenses, categorize them, and view monthly summaries with a visual category breakdown — all from your terminal.

## Features

- ✅ Add, list, and delete expenses
- ✅ Category tagging
- ✅ Monthly filter on summary
- ✅ Visual category bar chart in the terminal
- ✅ Persistent JSON storage (no database)
- ✅ Zero external dependencies

## Quick Start

```bash
git clone https://github.com/yourusername/go-cli-expense-tracker
cd go-cli-expense-tracker
go build -o expense-tracker
./expense-tracker help
```

## Usage

```bash
# Add expenses
./expense-tracker add "Coffee" 4.50 food
./expense-tracker add "AWS EC2" 58.20 cloud
./expense-tracker add "Gym membership" 35.00 health
./expense-tracker add "Netflix" 15.99 entertainment

# List all
./expense-tracker list

# Monthly summary
./expense-tracker summary
./expense-tracker summary 11   # November only

# Delete
./expense-tracker delete 2
```

## Example Output

```
ID   Date         Category        Amount   Description
--   ----         --------        ------   -----------
1    2024-11-01   food            $4.50    Coffee
2    2024-11-01   cloud           $58.20   AWS EC2
3    2024-11-01   health          $35.00   Gym membership
4    2024-11-01   entertainment   $15.99   Netflix

── Summary ─────────────────
  Total expenses: 4
  Total amount:   $113.69

  By category:
    cloud        $ 58.20  ██████████
    health       $ 35.00  ██████
    entertainment$ 15.99  ██
    food         $  4.50  █
```

## Data Storage

Expenses are stored in `expenses.json` in the current directory. You can back it up, commit it to Git, or move it freely.

## Tech Stack

- **Go** 1.21+
- Zero external dependencies (stdlib only)
- `encoding/json` · `text/tabwriter` · `os`

## License

MIT
