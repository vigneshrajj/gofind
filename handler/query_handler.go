package handler

import "errors"

func HandleQuery(query string) (string, error) {
	if query == "" {
		return "", errors.New("Missing 'query' parameter")
	}
	return "Query: " + query, nil
}
