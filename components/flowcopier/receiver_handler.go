package flowcopier

// TODO!!! add parameters info into responces

// FlowList --------------------------------------------------------------------------------------

//const AfterIDParam = "after_id"
//
//var ExportFlowEndpoint = server_http.Endpoint{Method: "GET", PathParams: nil, QueryParams: []string{AfterIDParam}, WorkerHTTP: ExportFlow}
//
//func ExportFlow(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
//	afterIDStr := req.URL.Query().Get(AfterIDParam)
//
//	items, err := dataTaggedOp.Export(afterIDStr, nil)
//	if err != nil {
//		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET storage/...ExportFlow: ", err))
//	}
//
//	return server.ResponseRESTOk(packs.Pack{
//		// SourceURL: "",
//		// SentAt:  time.Now(),
//		TypeKey: data.TypesKeyDataItems,
//		Content: items,
//	})
//}
