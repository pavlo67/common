package importer_rss

import (
	"testing"

	"github.com/pavlo67/constructor/processor/importer"
)

var testCases = []importer.ImporterTestCase{
	{
		Operator: &RSS{},
		Source:   "https://rss.unian.net/site/news_ukr.rss",
	},
}

func TestRSS(t *testing.T) {
	importer.TestImporterWithCases(t, testCases)
}
