package test

//func TestCRUD(t *testing.T) {
//
//	//t.Skip()
//
//	conf, err := config.ReadList(filelib.CurrentPath() + "../../../cfg.json5")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if conf == nil {
//		log.Fatal(nil.New("no config data after setup.Run()"))
//	}
//	partKeys := config.PartKeys{
//		"mysql": "items",
//	}
//	mysqlConfig, errs := conf.MySQL("", nil)
//	err = errs
//	if err != nil {
//		log.Fatal(err)
//	}
//	//err := config.LoadContext("../../../cfg.json5")
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//// creating Operator
//	//mysqlConfig, ok := config.Mysql["items"]
//	//if !ok {
//	//	log.Fatal(nil.Errorf("no mysql[items.comp] section in config: %v", config.Mysql))
//	//}
//
//	domain := "aaa"
//	confidenter := auth.IDentity{domain, "user", "111", ""}
//	identityBAD := auth.IDentity{"aaa", "user", "999", ""}
//	//isBAD := identityBAD.String()
//	UserIS := confidenter.String()
//	//identityGroup := auth.IDentity{"aaa", "group", "222", ""}
//	//isGroup := identityGroup.String()
//	//testController, err := controller.NewCRUDController(
//	//	map[auth.IDentity]auth.IDentity{
//	//		UserIS:    isGroup,
//	//		isBAD: isGroup,
//	//	},
//	//)
//
//	mysqlGroup, _ := groupsmysql.NewGroupsMySQL(
//		confidenter,
//		"aaa",
//		mysqlConfig,
//		"group",
//		"group_member",
//		rights.Managers{},
//	)
//
//	fountmysql, err := fountsmysql.NewMySQLFount(
//		//testController,
//		mysqlGroup,
//		mysqlConfig,
//		"fount",
//		"fount_tags",
//		"fount_stat",
//		"scanner_stat",
//		rights.Managers{rights.Create: UserIS, rights.Change: UserIS, rights.View: UserIS, rights.DeleteList: UserIS},
//	)
//	if err != nil {
//		t.Fatalf("can't init MySQLFount for tests: %v", err)
//	}
//	//fountsmysqlCRUD := fount.FountCRUD{fountsmysql}
//
//	s1 := rand.NewSource(time.Now().UnixNano())
//	r1 := rand.New(s1)
//
//	descriptionFount := crud.Contentus{
//		Label: "",
//		Details: map[string]string{
//			"Label": "Label CRUD TEST",
//			"Url":   "https://rss.unian.net/site/news_ukr.rss#" + strconv.Itoa(r1.Intn(1000)),
//		},
//		Managers: rights.Managers{rights.Owner: UserIS, rights.View: UserIS},
//	}
//
//	testCases := []crud.OperatorNewTestCase{
//		{
//			//&fountsmysqlCRUD,
//			fountmysql,
//			confidenter,
//			auth.IDentity{},
//			auth.IDentity{},
//			identityBAD,
//			identityBAD,
//			descriptionFount,
//			"Label",
//			false,
//			nil,
//			nil,
//			false,
//			false,
//		},
//	}
//	for _, testCase := range testCases {
//		fmt.Println("\n =============== Operator TEST: fount =================")
//		crud.OperatorNewTest(t, testCase)
//	}
//}
