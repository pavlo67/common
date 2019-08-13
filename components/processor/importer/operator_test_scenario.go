package importer

import (
	"log"
	"testing"
)

type ImporterTestCase struct {
	Operator Operator
	Source   string
}

func TestImporterWithCases(t *testing.T, testCases []ImporterTestCase) {
	var err error

	for _, tc := range testCases {
		err = tc.Operator.Init(tc.Source)
		if err != nil {
			t.Fatalf("can't init importer.Operator (%+v): %v", tc.Operator, err)
		}

		for {
			entity, err := tc.Operator.Next()
			if err != nil {
				t.Fatalf("can't get next item: %v", err)
			}
			if entity == nil {
				break
			}

			log.Println("/nID:", entity.Origin.Key, "\nItem:", entity.Content)

			//if object != nil {
			//	f, err := os.Create(`obj.html`)
			//	if err != nil {
			//		t.Fatalf("can't create file for write: %s", err)
			//	}
			//	_, err = f.Write([]byte(object.Contentus))
			//	if err != nil {
			//		t.Fatalf("error write to file: %s", err)
			//	}
			//}

		}
	}
}
