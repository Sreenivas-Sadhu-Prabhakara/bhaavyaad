package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type memStore struct{ items []Record }

func (m *memStore) Save(r Record) (Record, error) {
	r.ID = int64(len(m.items) + 1)
	m.items = append([]Record{r}, m.items...) // newest-first
	return r, nil
}
func (m *memStore) List(limit int) ([]Record, error) { return m.items, nil }

func mk(sup, item string, price float64) Record {
	in, _ := json.Marshal(Purchase{Supplier: sup, Item: item, Price: price, Qty: 1})
	return Record{Input: in, Headline: price, Label: item}
}

func TestSummarize_LastPricePerKey(t *testing.T) {
	// newest-first: latest ACME/soap is 22, earlier 20 must be ignored.
	recs := []Record{mk("ACME", "soap", 22), mk("BEST", "soap", 25), mk("ACME", "soap", 20)}
	out := Summarize(recs)
	var acme float64
	n := 0
	for _, lp := range out {
		if lp.Supplier == "ACME" && lp.Item == "soap" {
			acme = lp.Price
			n++
		}
	}
	if n != 1 || acme != 22 {
		t.Fatalf("expected single ACME/soap=22, got n=%d price=%v", n, acme)
	}
}

func TestPurchaseValidate(t *testing.T) {
	if err := (Purchase{Supplier: "A", Item: "x", Price: 1}).Validate(); err != nil {
		t.Fatalf("valid rejected: %v", err)
	}
	for i, bad := range []Purchase{{Item: "x"}, {Supplier: "A"}, {Supplier: "A", Item: "x", Price: -1}} {
		if err := bad.Validate(); err == nil {
			t.Fatalf("bad %d accepted", i)
		}
	}
}

func TestLogAndSummaryEndpoint(t *testing.T) {
	srv := NewServer(&memStore{})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/log",
		strings.NewReader(`{"supplier":"ACME","item":"soap","price":22,"qty":10}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("log %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/summary", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"price":22`) {
		t.Fatalf("summary=%s", rec.Body.String())
	}
}
