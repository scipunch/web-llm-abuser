package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"

	"github.com/scipunch/wla/cookies"
)

func main() {
	var cookiesPath string
	flag.StringVar(&cookiesPath, "cookies", "", "path to cookies.txt in netscape format")
	flag.Parse()

	var cookieJar []playwright.Cookie
	if cookiesPath != "" {
		slog.Info("reading cookies", "from", cookiesPath)
		f, err := os.Open(cookiesPath)
		if err != nil {
			log.Fatalf("failed to read cookies file with %v", err)
		}
		cookieJar, err = cookies.FromNetscape(f)
		f.Close()
		if err != nil {
			log.Fatalf("failed to parse cookies with %v", err)
		}
		slog.Info("cookie jar parsed", "amount", len(cookieJar), "value", cookieJar)
	}
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	browserCtx, err := browser.NewContext()
	if err != nil {
		log.Fatalf("failed to create browser context with %v", err)
	}
	defer browserCtx.Close()
	optCookies := make([]playwright.OptionalCookie, len(cookieJar))
	for i, cookie := range cookieJar {
		optCookies[i] = cookie.ToOptionalCookie()
	}
	browserCtx.AddCookies(optCookies)
	page, err := browserCtx.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto("https://chatgpt.com"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	time.Sleep(120 * time.Second)
	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
