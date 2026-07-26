package helpers

import (
	"net/http"
	"strconv"
	"strings"
)

func ParseAcceptLanguage(header string, validCodes []string) string {
	if header == "" {
		return "en"
	}

	langs := strings.Split(header, ",")
	for _, entry := range langs {
		code := strings.TrimSpace(strings.Split(entry, ";")[0])
		if parts := strings.Split(code, "-"); len(parts) > 0 {
			code = parts[0]
		}
		code = strings.ToLower(code)
		for _, valid := range validCodes {
			if code == valid {
				return code
			}
		}
	}

	return "en"
}

func ParsePagination(r *http.Request, defaultPageSize int) (pageNumber int, pageSize int) {
	pageNumber = 1
	pageSize = defaultPageSize

	if pn, err := strconv.Atoi(r.URL.Query().Get("page[number]")); err == nil && pn > 0 {
		pageNumber = pn
	}
	if ps, err := strconv.Atoi(r.URL.Query().Get("page[size]")); err == nil && ps > 0 {
		pageSize = ps
	}

	return
}
