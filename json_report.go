package main

import (
	"encoding/json"
	"os"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	keys := []string{}
	for k := range pages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedPages := []PageData{}
	for _, k := range keys {
		sortedPages = append(sortedPages, pages[k])
	}

	data, err := json.MarshalIndent(sortedPages, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, data, 0666)
	if err != nil {
		return err
	}
	return nil
}
