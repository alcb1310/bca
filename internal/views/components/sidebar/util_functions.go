package sidebar

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func hasActiveMenu(menu string, doc *goquery.Document) bool {
	testId := fmt.Sprintf(`[data-testid="menu-%s"]`, menu)
	return doc.Find(testId).HasClass("text-blue-gray-100")
}
