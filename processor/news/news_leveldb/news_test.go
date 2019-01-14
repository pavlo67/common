package news_leveldb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/punctum/processor/flow"
	"github.com/pavlo67/punctum/processor/news"
)

func TestCRUD(t *testing.T) {
	if err := os.Setenv("ENV", "test"); err != nil {
		t.Fatal("No test environment set!!!")
	}

	newsOp, err := New("_test")

	src := flow.Source{
		URL:      "abc",
		SourceID: "def",
	}
	item := news.Item{
		Source: src,
		Content: news.Content{
			Title:   "title",
			Summary: "summary",
			Text:    "text",
		},
	}

	err = newsOp.Save(&item)
	require.NoError(t, err)

	has, err := newsOp.Has(&item.Source)
	require.True(t, has)
	require.NoError(t, err)

	//operatorCRUD := founts.OperatorCRUD{newsOp}
	//testCases, err := operatorCRUD.TestCases(func() error { return newsOp.clean() })
	//require.NoError(t, err)
	//
	//crud.OperatorTest(t, testCases)

	err = newsOp.Close()
	require.NoError(t, err)
}
