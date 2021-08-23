package main

import (
	"mime"
	"net/http"
	"strings"
)

type Types []string

func NewTypes(r *http.Request, headerName string) Types {
	header := r.Header.Get(headerName)
	res := make([]string, 0)

	for _, part := range strings.Split(header, ",") {
		contentType, _, err := mime.ParseMediaType(part)
		if err == nil {
			res = append(res, contentType)
		}
	}

	return res
}

func (types Types) Has(s string) bool {
	for _, t := range types {
		if t == s {
			return true
		}
	}

	return false
}
