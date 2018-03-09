package eui

import (
	"encoding/csv"
	"net/http"
	"strings"
)

var Registries = []string{
	"https://standards.ieee.org/develop/regauth/oui36/oui36.csv",
	"https://standards.ieee.org/develop/regauth/oui28/mam.csv",
	"https://standards.ieee.org/develop/regauth/oui/oui.csv",
}

func Fetch(url string) ([]Registration, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	reader := csv.NewReader(res.Body)
	_, err = reader.Read() // Skip the header
	if err != nil {
		return nil, err
	}
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	registrations := make([]Registration, 0, len(records))
	for _, record := range records {
		if len(record) < 3 {
			continue
		}
		registrations = append(registrations, Registration{Prefix: strings.ToUpper(record[1]), Name: record[2]})
	}
	return registrations, nil
}
