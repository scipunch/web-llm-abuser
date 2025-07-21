package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

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
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto("https://news.ycombinator.com"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	entries, err := page.Locator(".athing").All()
	if err != nil {
		log.Fatalf("could not get entries: %v", err)
	}
	for i, entry := range entries {
		title, err := entry.Locator("td.title > span > a").TextContent()
		if err != nil {
			log.Fatalf("could not get text content: %v", err)
		}
		fmt.Printf("%d: %s\n", i+1, title)
	}
	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
