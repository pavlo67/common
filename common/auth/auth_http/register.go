package auth_http

//func workerRegister(_ *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
//
//	var testModePath string
//
//	if confidence_routes.Cfg != nil && strlib.In(confidence_routes.Cfg.ServerHTTP.Testers, req.RemoteAddr) {
//		testModePath = req.Header.Get("Test-Mode-Path")
//	}
//
//	if testModePath != "" {
//
//	}
//
//	credsJSON, err := ioutil.ReadAll(req.Body)
//	if err != nil {
//		return server.ResponseRESTError(http.StatusBadRequest, errors.Wrap(err, "can't read body"))
//	}
//
//	// log.Printf("%s", credsJSON)
//
//	var toSet auth.Creds
//	err = json.Unmarshal(credsJSON, &toSet)
//	if err != nil {
//		return server.ResponseRESTError(http.StatusBadRequest, errors.Wrapf(err, "can't unmarshal body: %s", credsJSON))
//	}
//
//	//
//	//user, errs := auth.GetUser(toSet, confidence_routes.authOps, nil)
//	//if len(errs) > 0 {
//	//	return server.ResponseRESTError(http.StatusForbidden, errs.Err())
//	//}
//	//if user == nil {
//	//	return server.ResponseRESTError(http.StatusForbidden, errors.New("no user authorized"))
//	//}
//	//
//	//toAddModified, err := confidence_routes.authOpToSetToken.SetCreds(*user, auth.Creds{}) // TODO!!! add custom toAddModified
//	//if err != nil {
//	//	return server.ResponseRESTError(http.StatusInternalServerError, errors.Wrap(err, "can't create JWT"))
//	//}
//	//
//	//if toAddModified != nil {
//	//	for t, c := range toAddModified.Values {
//	//		user.creds[t] = c
//	//	}
//	//}
//	//
//	return server.ResponseRESTOk(map[string]interface{}{"user": nil})
//}
//
//func workerModify(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
//	if user == nil {
//		return server.ResponseRESTError(http.StatusForbidden, errors.New("no user authorized"))
//	}
//
//	credsJSON, err := ioutil.ReadAll(req.Body)
//	if err != nil {
//		return server.ResponseRESTError(http.StatusBadRequest, errors.Wrap(err, "can't read body"))
//	}
//
//	var toReplace auth.Creds
//	err = json.Unmarshal(credsJSON, &toReplace)
//	if err != nil {
//		return server.ResponseRESTError(http.StatusBadRequest, errors.Wrapf(err, "can't unmarshal body: %s", credsJSON))
//	}
//
//	// toReplace = append(user.Creds, toReplace...)
//
//	// !!! previous user.Creds are ignored here
//	toReplaceModified, err := confidence_routes.AuthOpToSetToken.SetCreds(*user, toReplace) // TODO!!! add custom toReplace
//	if err != nil {
//		return server.ResponseRESTError(http.StatusInternalServerError, errors.Wrap(err, "can't create JWT"))
//	}
//
//	if toReplaceModified != nil {
//		for t, c := range toReplaceModified.Values {
//			user.creds[t] = c
//		}
//	}
//	return server.ResponseRESTOk(map[string]interface{}{"user": user})
//}
