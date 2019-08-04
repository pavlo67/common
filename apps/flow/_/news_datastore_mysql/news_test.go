package news_datastore_mysql

//// INTEGRATION TESTS
//if joiner != nil {
//	flowOp, ok := joiner.Interface(flow.InterfaceKey).(flow.Operator)
//	if !ok {
//		log.Fatalf("can't joiner.Interface(%s).(flow.Operator)", flow.InterfaceKey)
//	}
//
//	item1 := flow.Census{
//		SourceURL:  "1",
//		OriginalID: "2",
//		Original:   "3",
//		Key:        "4",
//		Contentus: &flow.Contentus{
//			Banner:  "5",
//			Summary: "6",
//			Text:    "7",
//			Tags:    []string{"8", "9"},
//			Href:    "10",
//		},
//		Status:  "11",
//		History: "12",
//	}
//	err := flowOp.Save(item1)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	items, allCnt, err := flowOp.ReadList(nil)
//	log.Printf("items: %#v\n\nitems[0].Contentus: %#v\n\nallCnt: %d\n\nerr: %s", items, *items[0].Contentus, allCnt, err)
//}
