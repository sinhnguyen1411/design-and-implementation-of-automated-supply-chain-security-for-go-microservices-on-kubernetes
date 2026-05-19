package watchlist

import "testing"

func TestAddRemove(t *testing.T) {
	w := New("u1")
	if !w.Add("AAPL") { t.Fatal("add failed") }
	if w.Add("AAPL") { t.Fatal("duplicate add should return false") }
	if !w.Remove("AAPL") { t.Fatal("remove failed") }
	if len(w.Symbols) != 0 { t.Fatal("should be empty") }
}
