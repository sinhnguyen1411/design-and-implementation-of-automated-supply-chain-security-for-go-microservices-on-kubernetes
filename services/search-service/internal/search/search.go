package search

import "strings"

type Symbol struct {
	Ticker   string `json:"ticker"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Sector   string `json:"sector"`
}

func Query(catalog []Symbol, q string) []Symbol {
	q = strings.ToLower(q)
	if q == "" { return catalog }
	var results []Symbol
	for _, s := range catalog {
		if strings.Contains(strings.ToLower(s.Ticker), q) || strings.Contains(strings.ToLower(s.Name), q) {
			results = append(results, s)
		}
	}
	return results
}
