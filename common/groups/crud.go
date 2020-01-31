package groups

//// Create ...
//func (gr *GroupsMySQL) Create(confidenter auth.IDentity, description crud.Contentus) (*auth.IDentity, error) {
//	gr.DescriptionToData(description)
//	res, err := gr.Create(confidenter, *gr.crudBuffer)
//	return &res, err
//}
//
//// Read returns object's Contentus data (accordingly to requester's rights).
//func (gr *GroupsMySQL) Read(confidenter auth.IDentity, groupIS basis.UserIS) (*crud.Contentus, error) {
//	//data, err := gr.Read(confidenter, groupIS)
//	var err error
//	gr.crudBuffer, err = gr.Read(confidenter, groupIS)
//	if err != nil {
//		return nil, err
//	}
//	return gr.DataToDescription()
//}
//
//// ReadListCRUD returns array of object's Contentus data (accordingly to requester's rights).
//func (gr *GroupsMySQL) ReadListCRUD(confidenter auth.IDentity, selector selectors.Selector, options *content.ListOptions) ([]crud.Contentus, int64, error) {
//
//	data, allCount, err := gr.ReadList(confidenter, options, selector)
//	if err != nil {
//		return nil, 0, err
//	}
//	var description []crud.Contentus
//	for _, g := range data {
//		gr.crudBuffer = &g
//		desc, _ := gr.DataToDescription()
//		description = append(description, *desc)
//	}
//	return description, allCount, nil
//}
//
//// Update changes object's Contentus data (accordingly to requester's rights).
//func (gr *GroupsMySQL) Update(confidenter auth.IDentity, groupIS basis.UserIS, description crud.Contentus) (crud.Result, error) {
//	gr.DescriptionToData(description)
//	return gr.Update(confidenter, groupIS, *gr.crudBuffer)
//}
//
//// DeleteList ...
//func (gr *GroupsMySQL) DeleteList(confidenter auth.IDentity, groupIS basis.UserIS) (crud.Result, error) {
//	return gr.DeleteList(confidenter, groupIS)
//}
//
//// Count ...
//func (gr *GroupsMySQL) CountCRUD(selector selectors.Selector, joinTo crud.JoinTo, groupBy, sortBy []string) ([]crud.Count, error) {
//	if Tables[joinTo.ToTable] != "" {
//		joinTo.ToTable = Tables[joinTo.ToTable]
//	} else {
//		return nil, nil.New("can't find table code: " + joinTo.ToTable)
//	}
//	return mysqllib.Count(gr.dbh, selector, joinTo, groupBy, sortBy)
//}
//
//// Describe ... read crud.json5
//func (gr *GroupsMySQL) DescribeCRUD() (*crud.Description, error) {
//	return crud.Describe(filelib.CurrentPath() + "../")
//}
//
//func (gr *GroupsMySQL) DataToDescription() (*crud.Contentus, error) {
//	g := gr.crudBuffer
//	od := crud.Contentus{
//		Label: string(g.UserIS),
//		Details: map[string]string{
//			"Nick":    g.Nick,
//			"Code":    g.Code,
//			"Details": g.Details,
//		},
//		Managers: g.Managers,
//	}
//
//	if g.Updated != nil {
//		od.Details["Updated"] = *g.Updated
//	}
//	return &od, nil
//}
//
//func (gr *GroupsMySQL) DescriptionToData(o crud.Contentus) error {
//	gr.crudBuffer =
//		&groups.Census{
//			Type: groups.Type{
//				Code: o.Details["Code"],
//			},
//			Nick:     o.Details["Nick"],
//			Details:  o.Details["Details"],
//			Managers: o.Managers,
//		}
//	return nil
//}
