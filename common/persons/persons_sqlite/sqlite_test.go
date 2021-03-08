package persons_sqlite

//func TestCRUD(t *testing.T) {
//	_, cfgService, l := apps.PrepareTests(t, "test_service", "../../../apps/", "test", "persons_sqlite")
//	require.NotNil(t, cfgService)
//
//	var cfg config.Access
//	err := cfgService.Value("files_fs", &cfg)
//	require.NoErrorf(t, err, "%#v", cfgService)
//
//	components := []starter.Starter{
//		{connect_sqlite.Starter(), nil},
//		{Starter(), nil},
//	}
//
//	joinerOp, err := starter.Run(components, cfgService, "CLI BUILD FOR TEST", l)
//	require.NoError(t, err)
//	require.NotNil(t, joinerOp)
//	defer joinerOp.CloseAll()
//
//	persons.OperatorTestScenarioNoRBAC(t, joinerOp, l)
//}
