package patch

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

const (
	patchPage   = "https://learn.microsoft.com/en-us/troubleshoot/sql/releases/download-and-install-latest-updates"
	SQL2017     = "SQL Server 2017"
	SQL2019     = "SQL Server 2019"
	SQL2022     = "SQL Server 2022"
	PatchCU     = "CU"
	PatchGDR    = "GDR"
	PatchCU_GDR = "CU+GDR"
)

var (
	cuLinkHeader       = "https://learn.microsoft.com/en-us/troubleshoot/sql/releases/"
	patchList          []Patch
	BrowserPageOptions = playwright.BrowserNewPageOptions{Locale: playwright.String("en-US")}
)

type Patch struct {
	// CU ,GDR,CU+GDR 名称
	PatchType string `json:"patch_type"`
	// 微软文档给的GDR，CU等的名称
	Name    string `json:"name"`
	Version string `json:"version"`
	// sql server latest patch page的链接
	// FirstLink   string    `json:"first_link"`
	FileName    string `json:"file_name"`
	Link        string `json:"link"`
	SizeBytes   int    `json:"size_bytes"`
	ProdVersion string `json:"prod_version"`
	ReleaseDate string `json:"release_date"`
}

// type Tabler interface {
// 	TableName() string
// }

// TableName 会将 User 的表名重写为 `profiles`
func (Patch) TableName() string {
	return "mssql_patch"
}

func GetSqlPatch() []Patch {
	browser, pw := newBrowser()
	defer closeBrowser(browser, pw)
	page, err := browser.NewPage(BrowserPageOptions)
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	patchLinkPage, err := browser.NewPage(BrowserPageOptions)
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto(patchPage); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	locator := page.GetByLabel("Latest updates available for currently supported versions of SQL Server")
	versions, _ := locator.Locator("tbody > tr").All()
	var patch Patch
	for i, version := range versions {
		if i == 0 {
			links, _ := version.Locator("td > a").All()
			for _, link := range links {
				patchType, _ := link.TextContent()
				l, _ := link.GetAttribute("href")
				if patchType == PatchGDR || strings.Contains(patchType, " for ") || strings.Contains(patchType, "+ GDR") {
					if patchType == PatchGDR {
						patch = Patch{Version: SQL2022, PatchType: PatchGDR, Name: patchType, Link: l}
					} else if strings.Contains(patchType, " for ") {
						l = cuLinkHeader + l
						patch = Patch{Version: SQL2022, PatchType: PatchCU, Name: patchType, Link: l}
					} else if strings.Contains(patchType, "+ GDR") {
						patch = Patch{Version: SQL2022, PatchType: PatchCU_GDR, Name: patchType, Link: l}
					}
					patch.GetDlLink(patchLinkPage)
					patchList = append(patchList, patch)
				}
			}

		}
		if i == 1 {
			links, _ := version.Locator("td > a").All()
			for _, link := range links {
				patchType, _ := link.TextContent()
				l, _ := link.GetAttribute("href")
				if patchType == PatchGDR || strings.Contains(patchType, " for ") || strings.Contains(patchType, "+ GDR") {
					if patchType == PatchGDR {
						patch = Patch{Version: SQL2019, PatchType: PatchGDR, Name: patchType, Link: l}
					} else if strings.Contains(patchType, " for ") {
						l = cuLinkHeader + l
						patch = Patch{Version: SQL2019, PatchType: PatchCU, Name: patchType, Link: l}

					} else if strings.Contains(patchType, "+ GDR") {
						patch = Patch{Version: SQL2019, PatchType: PatchCU_GDR, Name: patchType, Link: l}
					}
					patch.GetDlLink(patchLinkPage)
					patchList = append(patchList, patch)
				}

			}

		}
		if i == 2 {
			links, _ := version.Locator("td > a").All()
			for _, link := range links {
				patchType, _ := link.TextContent()
				l, _ := link.GetAttribute("href")
				if patchType == PatchGDR || strings.Contains(patchType, " for ") || strings.Contains(patchType, "+ GDR") {
					if patchType == PatchGDR {
						patch = Patch{Version: SQL2017, PatchType: PatchGDR, Name: patchType, Link: l}
					} else if strings.Contains(patchType, " for ") {
						l = cuLinkHeader + l
						patch = Patch{Version: SQL2017, PatchType: PatchCU, Name: patchType, Link: l}
					} else if strings.Contains(patchType, "+ GDR") {
						patch = Patch{Version: SQL2017, PatchType: PatchCU_GDR, Name: patchType, Link: l}
					}
					patch.GetDlLink(patchLinkPage)
					patchList = append(patchList, patch)
				}
			}

		}
		if i >= 3 {
			break
		}
	}
	return patchList

}

func (p *Patch) GetDlLink(page playwright.Page) {
	var err error
	// 处理CU下载链接
	if p.PatchType == PatchCU {
		if _, err = page.Goto(p.Link); err != nil {
			log.Fatalf("could not goto: %v", err)
		}
		p.Link, _ = page.GetByText("Download the latest cumulative update package for SQL Server").GetAttribute("href")
		p.GetDlInfo(page)

	} else {
		//  Download the package now Page
		if _, err := page.Goto(p.Link); err != nil {
			log.Fatalf("could not goto: %v", err)
		}
		// redirect link
		p.Link, _ = page.GetByText("Download the package now").GetAttribute("href")
		p.GetDlInfo(page)

	}
}

// 实际下载页面,获取下载链接,版本号和发布日期
func (p *Patch) GetDlInfo(page playwright.Page) {
	if _, err := page.Goto(p.Link); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	s, _ := page.Locator("script").All()
	for _, ss := range s {
		jsonText, _ := ss.TextContent()
		if strings.Contains(jsonText, "window.__DLCDetails__=") {
			jsonText = strings.ReplaceAll(jsonText, "window.__DLCDetails__=", "")
			var data map[string]interface{}
			err := json.Unmarshal([]byte(jsonText), &data)
			if err != nil {
				log.Fatal(err)
			}

			// 逐层访问
			dlcDetailsView := data["dlcDetailsView"].(map[string]interface{})
			downloadFile := dlcDetailsView["downloadFile"].([]interface{})

			var np = p
			for i := range downloadFile {
				file := downloadFile[i].(map[string]interface{})
				p.FileName = file["name"].(string)
				p.Link = file["url"].(string)
				p.SizeBytes, _ = strconv.Atoi(file["size"].(string))
				p.ProdVersion = file["version"].(string)
				rd, _ := time.Parse("1/2/2006 3:04:05 PM", file["datePublished"].(string))
				p.ReleaseDate = rd.UTC().Format("2006-01-02 3:04:05 PM")
				if i+1 < len(downloadFile) {
					np = p
					patchList = append(patchList, *np)
				}
			}

		}

	}

}
