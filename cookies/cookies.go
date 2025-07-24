package cookies

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/playwright-community/playwright-go"
)

const httpOnlyPrefix = "#HttpOnly_"

func FromNetscape(r io.Reader) ([]playwright.Cookie, error) {
	var cookies []playwright.Cookie
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		httpOnly := strings.HasPrefix(line, httpOnlyPrefix)
		if !httpOnly && (strings.HasPrefix(line, "#") || line == "") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 7 {
			return nil, fmt.Errorf("unexpected columns count %d", len(parts))
		}

		// Parse Secure
		secure := strings.ToLower(parts[3]) == "true"

		// Parse Expires
		expiresStr := parts[4]
		var expires float64
		if expiresStr != "0" && expiresStr != "" {
			expiresInt, err := strconv.ParseInt(expiresStr, 10, 64)
			if err == nil {
				expires = float64(expiresInt)
			}
		}

		domain := strings.TrimPrefix(parts[0], httpOnlyPrefix)

		// Construct cookie
		cookie := playwright.Cookie{
			Name:     parts[5],
			Value:    parts[6],
			Domain:   domain,
			Path:     parts[2],
			Secure:   secure,
			Expires:  expires,
			HttpOnly: httpOnly,
			SameSite: nil,
		}

		cookies = append(cookies, cookie)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cookies, nil
}
