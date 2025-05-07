package patch

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

func newBrowser() (playwright.Browser, *playwright.Playwright) {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false), // 设置为 true 以无头模式运行
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	return browser, pw
}

func closeBrowser(browser playwright.Browser, pw *playwright.Playwright) {

	if err := browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err := pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
