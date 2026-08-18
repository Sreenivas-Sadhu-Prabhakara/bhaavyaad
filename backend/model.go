package backend

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Purchase records what you paid a supplier for an item on a given buy.
type Purchase struct {
	Supplier string  `json:"supplier"`
	Item     string  `json:"item"`
	Price    float64 `json:"price"` // price per unit paid
	Qty      float64 `json:"qty"`
}

// Validate reports whether the Purchase is well formed.
func (p Purchase) Validate() error {
	if strings.TrimSpace(p.Supplier) == "" || strings.TrimSpace(p.Item) == "" {
		return fmt.Errorf("supplier and item are required")
	}
	if p.Price < 0 || p.Qty < 0 {
		return fmt.Errorf("price and qty cannot be negative")
	}
	return nil
}

// LastPrice is the most recent price paid for one item at one supplier.
type LastPrice struct {
	Supplier string  `json:"supplier"`
	Item     string  `json:"item"`
	Price    float64 `json:"price"`
}

// Summarize returns the last price paid per (supplier, item). Records arrive
// newest-first, so the first occurrence of each key is the latest price.
func Summarize(records []Record) []LastPrice {
	seen := map[string]bool{}
	var out []LastPrice
	for _, r := range records {
		var p Purchase
		if json.Unmarshal(r.Input, &p) != nil {
			continue
		}
		key := p.Supplier + "\x00" + p.Item
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, LastPrice{Supplier: p.Supplier, Item: p.Item, Price: p.Price})
	}
	return out
}

// parseEntry decodes+validates a purchase, returning its headline (price) and
// label (item) for the generic store.
func parseEntry(raw []byte) (float64, string, error) {
	var p Purchase
	if err := json.Unmarshal(raw, &p); err != nil {
		return 0, "", fmt.Errorf("invalid json")
	}
	if err := p.Validate(); err != nil {
		return 0, "", err
	}
	return p.Price, p.Item, nil
}
