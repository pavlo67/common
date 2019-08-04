package collector_comp

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/pkg/errors"

	"github.com/julienschmidt/httprouter"
	"github.com/pavlo67/partes/crud"
	"github.com/pavlo67/partes/crud/selectors"
	"github.com/pavlo67/partes/serverhttp"
	"github.com/pavlo67/partes/serverhttp/serverhttp_jschmhr"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/collector/importer"
	"github.com/pavlo67/punctum/collector/importer/importer_factory"
	"github.com/pavlo67/punctum/confidenter/auth"
	"github.com/pavlo67/punctum/confidenter/controller"
	"github.com/pavlo67/punctum/notebook/notes"
	"github.com/pavlo67/punctum/processor/sources"
	"github.com/pavlo67/punctum/starter/joiner"
)

var htmlHandlers = map[string]serverhttp_jschmhr.HTMLHandler{
	"viewFounts":     viewFounts,
	"blankFount":     blankFount,
	"editFount":      editFount,
	"viewFount":      viewFount,
	"viewFlowsByTag": viewFlowsByTag,
	"viewFlows":      viewFlowsFount,
	"freeFlows":      freeFlows,
	"viewFlowsFount": viewFlowsFount,
}

var restHandlers = map[string]serverhttp_jschmhr.RESTHandler{
	"createFount":      createFount,
	"updateFount":      updateFount,
	"deleteFount":      deleteFount,
	"importTest":       importTest,
	"fountSettings":    fountSettings,
	"fountTestSetting": fountTestSetting,
	"importFlow":       importFlow,
}

const allMyFlows = "всі мої новини"
const allMyFountsID = "all"

var sortFlowsByDefault = "created_at-"

type dataFount struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Url    string `json:"url"`
	Direct string `json:"direct"`
	Tags   string `json:"tags"`
	Type   string `json:"type"`

	ImportDetailsType   string `json:"import_details_type"`
	ImportDetailsParams string `json:"import_details_params"`
}

type fountSetting struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Type string `json:"type"`

	ImportStartRegexp  string `json:"import_start_regexp"`
	ImportFinishRegexp string `json:"import_finish_regexp"`
	ImportSplitRegexp  string `json:"import_split_regexp"`
	ImportTagsList     string `json:"import_tags_list"`
	TitleRegexp        string `json:"title_regexp"`
}

func blankFount(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.HTMLResponse, error) {
	var responseData serverhttp.HTMLResponse
	if user == nil {
		responseData = serverhttp.HTMLResponse{
			Status: http.StatusUnauthorized,
			Data: map[string]string{
				"corpus": `<h4 style="color:red;">` + basis.ErrAuthenticated.Error() + "</h4>",
			},
		}
		return responseData, basis.ErrAuthenticated
	}

	responseData = serverhttp.HTMLResponse{
		Status: http.StatusOK,
		Data: map[string]string{
			"caput":  "Нове джерело",
			"corpus": htmlBlankFount(),
		},
	}
	return responseData, nil
}

func editFount(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.HTMLResponse, error) {

	responseData := serverhttp.HTMLResponse{
		Data: map[string]string{},
	}
	if user == nil {
		responseData.Status = http.StatusUnauthorized
		responseData.Data["corpus"] = `<h4 style="color:red;">` + basis.ErrAuthenticated.Error() + "</h4>"
		return responseData, basis.ErrAuthenticated
	}
	id := params.ByName("id")
	var fount *sources.Item
	fount, err := fountOp.Read(user.Identity(), id)
	if err != nil {
		responseData.Status = http.StatusNonAuthoritativeInfo
		responseData.Data["corpus"] = `<h4 style="color:red;">` + errors.Cause(err).Error() + "</h4>"
		return responseData, err
	}
	IS := user.Identity.String()
	if IS != fount.ROwner {
		err := errors.New("no rights")
		responseData.Status = http.StatusNonAuthoritativeInfo
		responseData.Data["corpus"] = `<h4 style="color:red;">` + err.Error() + "</h4>"
		return responseData, err
	}
	responseData.Status = http.StatusOK
	responseData.Data["corpus"] = htmlEditFount(fount, id, controller.IsAdmin(user))
	responseData.Data["caput"] = "Редаґування джерела"
	return responseData, nil
}

func importFlow(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.RESTResponse, error) {
	//var err error
	if user == nil {
		return serverhttp.RESTResponse{Status: http.StatusUnauthorized, Data: serverhttp.RESTErrorOld{basis.ErrAuthenticated.Error()}}, basis.ErrAuthenticated
	}
	res, err := flowToObject(user.Identity(), params.ByName("id"))
	if err != nil {
		return serverhttp.RESTResponse{Status: http.StatusNonAuthoritativeInfo, Data: serverhttp.RESTErrorOld{errors.Cause(err).Error()}}, err
	}
	return serverhttp.RESTResponse{Status: http.StatusOK, Data: serverhttp.RESTDataMessage{"Імпортовано!", itemsEndpoints["edit"].Path(res)}}, nil
}

func deleteFount(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.RESTResponse, error) {

	var err error
	if user == nil {
		return serverhttp.RESTResponse{Status: http.StatusUnauthorized, Data: serverhttp.RESTErrorOld{basis.ErrAuthenticated.Error()}}, basis.ErrAuthenticated
	}

	id := params.ByName("id")
	_, err = fountOp.Delete(user.Identity(), id)
	if err != nil {
		return serverhttp.RESTResponse{Status: http.StatusInternalServerError, Data: serverhttp.RESTErrorOld{errors.Cause(err).Error()}}, err
	}
	return serverhttp.RESTResponse{Status: http.StatusOK, Data: serverhttp.RESTDataMessage{"Джерело видалено!", endpoints["viewFounts"].ServerPath}}, nil
}

func fountTestSetting(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.RESTResponse, error) {
	if user == nil || !controller.IsAdmin(user) {
		return serverhttp.RESTResponse{Status: http.StatusUnauthorized, Data: serverhttp.RESTErrorOld{basis.ErrAuthenticated.Error()}}, basis.ErrAuthenticated
	}
	testContent := ""
	var data fountSetting
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	defer r.Body.Close()
	if err != nil {
		return serverhttp.RESTResponse{
			Status: http.StatusInternalServerError,
			Data:   serverhttp.RESTErrorOld{err.Error()},
		}, basis.ErrJSONFormat
	}
	p, err := json.Marshal(htmlimporter.ImportParams{
		ImportStartRegexp:    data.ImportStartRegexp,
		ImportFinishRegexp:   data.ImportFinishRegexp,
		ImportSeparateRegexp: data.ImportSplitRegexp,
		TitleRegexp:          data.TitleRegexp,
		AcceptableTags:       reAcceptableTags.Split(data.ImportTagsList, -1),
	})
	if err == nil {
		err = importerOperators["htmlimporter"].Init(data.URL, string(p), false)
		if err == nil {
			var e importer.Entity
			var o *notes.Item
			for {
				e, err = importerOperators["htmlimporter"].Next()
				if err == importer.ErrNoMoreItems {
					err = nil
					break
				}
				if err != nil {
					log.Printf("error reading imported item: %s", err)
					continue
				}
				o, _ = e.Object()
				testContent += "<br>\n" + "Label: <h3>" + o.Name + "</h3>Contentus:<br>" + o.Content + "<hr>"
			}
		}
	}
	if err != nil {
		return serverhttp.RESTResponse{Status: http.StatusNonAuthoritativeInfo, Data: serverhttp.RESTErrorOld{errors.Cause(err).Error()}}, err
	}

	return serverhttp.RESTResponse{Status: http.StatusOK, Data: serverhttp.RESTDataMessage{testContent, ""}}, nil
}

var reAcceptableTags = regexp.MustCompile(`\W+`)

func fountSettings(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.RESTResponse, error) {
	if user == nil || !controller.IsAdmin(user) {
		return serverhttp.RESTResponse{Status: http.StatusUnauthorized, Data: serverhttp.RESTErrorOld{basis.ErrAuthenticated.Error()}}, basis.ErrAuthenticated
	}
	var data fountSetting
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	defer r.Body.Close()
	if err != nil {
		return serverhttp.RESTResponse{
			Status: http.StatusInternalServerError,
			Data:   serverhttp.RESTErrorOld{err.Error()},
		}, basis.ErrJSONFormat
	}
	ip := htmlimporter.ImportParams{
		ImportStartRegexp:    data.ImportStartRegexp,
		ImportFinishRegexp:   data.ImportFinishRegexp,
		ImportSeparateRegexp: data.ImportSplitRegexp,
		TitleRegexp:          data.TitleRegexp,
	}
	if data.ImportTagsList != "" {
		ip.AcceptableTags = reAcceptableTags.Split(data.ImportTagsList, -1)
	}
	p, err := json.Marshal(ip)
	var affected int64
	if err == nil {
		affected, err = fountOp.ExportSettings(data.URL, data.Type, string(p))
	}
	if err != nil {
		return serverhttp.RESTResponse{Status: http.StatusNonAuthoritativeInfo, Data: serverhttp.RESTErrorOld{errors.Cause(err).Error()}}, err
	}
	return serverhttp.RESTResponse{Status: http.StatusOK, Data: serverhttp.RESTDataMessage{"Налаштування збережено для " + strconv.FormatInt(affected, 10) + " запису(ів)!", endpoints["viewFount"].Path(data.ID)}}, nil
}

func updateFount(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.RESTResponse, error) {
	if user == nil {
		return serverhttp.RESTResponse{Status: http.StatusUnauthorized, Data: serverhttp.RESTErrorOld{basis.ErrAuthenticated.Error()}}, basis.ErrAuthenticated
	}

	var data dataFount
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	defer r.Body.Close()
	if err != nil {
		return serverhttp.RESTResponse{
			Status: http.StatusInternalServerError,
			Data:   serverhttp.RESTErrorOld{err.Error()},
		}, basis.ErrJSONFormat
	}
	toFlow := false
	toObject := false
	if data.Direct == "flow" {
		toFlow = true
	}
	if data.Type == "" {
		//	check url type
		data.Type = importer_factory.CheckURLType(data.Url)
	}
	if data.Type == string(htmlimporter.InterfaceKey) {
		toObject = false
		toFlow = true

	} else if data.Type == importer_factory.NoneImporter {
		toObject = false
		toFlow = false
	}
	if data.Direct == "" {
		toObject = false
		toFlow = false
	}

	var id int64
	if err == nil {
		id, err = strconv.ParseInt(data.ID, 10, 64)
	}
	if err == nil {
		fountUpd := sources.Item{
			ID:                  id,
			URL:                 data.Url,
			Title:               data.Title,
			ToFlow:              toFlow,
			ToObject:            toObject,
			Tags:                data.Tags,
			ImportType:          importer.ImportType(data.Type),
			ImportDetailsType:   data.ImportDetailsType,
			ImportDetailsParams: data.ImportDetailsParams,
		}
		_, err = fountOp.Update(user.Identity(), fountUpd)
	}
	if err != nil {
		return serverhttp.RESTResponse{Status: http.StatusNonAuthoritativeInfo, Data: serverhttp.RESTErrorOld{errors.Cause(err).Error()}}, err
	}
	return serverhttp.RESTResponse{Status: http.StatusOK, Data: serverhttp.RESTDataMessage{"Дані джерела змінено!", endpoints["viewFount"].Path(data.ID)}}, nil
}

const defaultFlowSortBy = "created_at-"

func freeFlows(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.HTMLResponse, error) {

	responseData := serverhttp.HTMLResponse{
		Data: map[string]string{},
	}
	tag := params.ByName("tag")
	options := crud.GetHTTPOptions(r, 0, defaultFlowSortBy)

	ident := auth.IDentity{joiner.SystemDomain(), "group", "1"}
	flows, flowsCnt, allTags, _, err := flowsByTag(&ident, tag, &options.ReadOptions)
	if err != nil {
		responseData.Status = http.StatusNonAuthoritativeInfo
		responseData.Data["corpus"] = `<h4 style="color:red;">` + errors.Cause(err).Error() + "</h4>"
		return responseData, err
	}
	responseData.Status = http.StatusOK
	responseData.Data["corpus"] = htmlFountFree(tag, flows, flowsCnt, options.PageNum, allTags)
	responseData.Data["caput"] = title
	return responseData, nil
}

func viewFlowsFount(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.HTMLResponse, error) {

	responseData := serverhttp.HTMLResponse{
		Data: map[string]string{},
	}
	if user == nil {
		responseData.Status = http.StatusUnauthorized
		responseData.Data["corpus"] = `<h4 style="color:red;">` + basis.ErrAuthenticated.Error() + "</h4>"
		return responseData, basis.ErrAuthenticated
	}
	var err error
	id := params.ByName("id")
	if id == "" {
		id = allMyFountsID
	}

	var selector selectors.Selector
	var title1 string

	if id != allMyFountsID {
		f, err := fountOp.Read(user.Identity(), id)
		if err != nil {
			responseData.Status = http.StatusNonAuthoritativeInfo
			responseData.Data["corpus"] = `<h4 style="color:red;">` + errors.Cause(err).Error() + "</h4>"
			return responseData, err
		}
		fountIS := joiner.SystemDomain() + "/fount/" + id
		selector = selectors.Or(
			selectors.FieldEqual("fount_is", fountIS),
			selectors.FieldEqual("fount_url", f.URL),
		)
		title1 = " » Записи від: " + f.Title
	} else {
		title1 = " » всі мої новини"
	}

	options := crud.GetHTTPOptions(r, 0, defaultFlowSortBy)
	flows, flowsCnt, err := flowOp.ReadList(user.Identity(), &options.ReadOptions, selector)

	if err != nil {
		responseData.Status = http.StatusNonAuthoritativeInfo
		responseData.Data["corpus"] = `<h4 style="color:red;">` + errors.Cause(err).Error() + "</h4>"
		return responseData, err
	}
	responseData.Status = http.StatusOK

	if id == allMyFountsID {
		userFountsTags, err := fountOp.ReadTags(user.Identity(), selector)
		if err != nil {
			// log!!!
		}
		responseData.Data["index"] = showTags(allMyFlows, endpoints["viewFlows"].ServerPath+"/tag/", userFountsTags, uint64(flowsCnt))
	} else {
		responseData.Data["index"] = `<tr><td><a href="` + endpoints["viewFlowsFount"].Path(allMyFountsID) + `">- всі мої новини</a></td></tr>`
	}

	responseData.Data["corpus"] = htmlFountFlows(flows)
	responseData.Data["caput"] = title
	responseData.Data["caput"] += title1

	return responseData, nil
}

func viewFlowsByTag(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.HTMLResponse, error) {

	responseData := serverhttp.HTMLResponse{
		Data: map[string]string{},
	}
	if user == nil {
		responseData.Status = http.StatusUnauthorized
		responseData.Data["corpus"] = `<h4 style="color:red;">` + basis.ErrAuthenticated.Error() + "</h4>"
		return responseData, basis.ErrAuthenticated
	}
	tag := params.ByName("tag")
	//if tag == "" {
	//
	//}

	options := crud.GetHTTPOptions(r, 0, defaultFlowSortBy)

	flows, flowsCnt, allTags, tagFounts, err := flowsByTag(user.Identity(), tag, &options.ReadOptions)
	if err != nil {
		responseData.Status = http.StatusNonAuthoritativeInfo
		responseData.Data["corpus"] = `<h4 style="color:red;">` + errors.Cause(err).Error() + "</h4>"
		return responseData, err
	}

	responseData.Status = http.StatusOK
	responseData.Data["index"] = showTags(tag, endpoints["viewFlows"].ServerPath+"/tag/", allTags, flowsCnt)
	if flows != nil {
		responseData.Data["corpus"] = showFlows(flows, tagFounts)
	}

	responseData.Data["caput"] = title
	if tag != "" {
		responseData.Data["caput"] += " » Записи з міткою: '" + tag + "'"
	}

	return responseData, nil
}

func viewFounts(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.HTMLResponse, error) {

	responseData := serverhttp.HTMLResponse{
		Data: map[string]string{},
	}

	var founts []sources.Item
	if user == nil {
		responseData.Status = http.StatusUnauthorized
		responseData.Data["corpus"] = `<h4 style="color:red;">` + basis.ErrAuthenticated.Error() + "</h4>"
		return responseData, basis.ErrAuthenticated
	}
	IS := user.Identity.String()
	selector := selectors.FieldEqual("r_owner", string(IS))
	founts, _, err := fountOp.ReadList(user.Identity(), nil, selector)
	if err != nil {
		responseData.Status = http.StatusNonAuthoritativeInfo
		responseData.Data["corpus"] = `<h4 style="color:red;">` + errors.Cause(err).Error() + "</h4>"
		return responseData, err
	}

	responseData.Status = http.StatusOK
	responseData.Data["corpus"] = htmlFounts(founts, user.IsGroupMember(idFountScannerGroup), user.String())
	responseData.Data["caput"] = title
	return responseData, nil
}

func viewFount(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.HTMLResponse, error) {

	responseData := serverhttp.HTMLResponse{
		Data: map[string]string{},
	}
	if user == nil {
		responseData.Status = http.StatusUnauthorized
		responseData.Data["corpus"] = `<h4 style="color:red;">` + basis.ErrAuthenticated.Error() + "</h4>"
		return responseData, basis.ErrAuthenticated
	}
	id := params.ByName("id")
	var fount *sources.Item
	fount, err := fountOp.Read(user.Identity(), id)
	if err != nil {
		responseData.Status = http.StatusNonAuthoritativeInfo
		responseData.Data["corpus"] = `<h4 style="color:red;">` + errors.Cause(err).Error() + "</h4>"
		return responseData, err
	}
	userIsOwner := false
	IS := user.Identity.String()
	if IS == fount.ROwner {
		userIsOwner = true
	}
	responseData.Status = http.StatusOK
	responseData.Data["corpus"] = htmlFount(fount, id, userIsOwner)
	responseData.Data["caput"] = title
	return responseData, nil
}

func createFount(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.RESTResponse, error) {

	if user == nil {
		return serverhttp.RESTResponse{Status: http.StatusUnauthorized, Data: serverhttp.RESTErrorOld{basis.ErrAuthenticated.Error()}}, basis.ErrAuthenticated
	}
	var data dataFount
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	defer r.Body.Close()
	if err != nil {
		return serverhttp.RESTResponse{
			Status: http.StatusInternalServerError,
			Data:   serverhttp.RESTErrorOld{err.Error()},
		}, basis.ErrJSONFormat
	}
	//log.Println("fount data: ", data)
	toFlow := false
	toObject := false
	if data.Direct == "flow" {
		toFlow = true
	}
	if data.Type == "" {
		//	check url type
		data.Type = importer_factory.CheckURLType(data.Url)
		if data.Type == string(htmlimporter.InterfaceKey) {
			toObject = false
			toFlow = true
		} else if data.Type == importer_factory.NoneImporter {
			toObject = false
			toFlow = false
		}
	}
	if data.Direct == "" {
		toObject = false
		toFlow = false
	}
	if data.Tags == "" && data.Type != importer_factory.NoneImporter {
		data.Tags = data.Type
	}
	fountNew := sources.Item{URL: data.Url, Title: data.Title, ToFlow: toFlow, ToObject: toObject, Tags: data.Tags, ImportType: importer.ImportType(data.Type)}
	res, err := fountOp.Create(user.Identity(), fountNew)
	if err != nil {
		return serverhttp.RESTResponse{Status: http.StatusNonAuthoritativeInfo, Data: serverhttp.RESTErrorOld{errors.Cause(err).Error()}}, err
	}
	return serverhttp.RESTResponse{Status: http.StatusOK, Data: serverhttp.RESTDataMessage{"Зареєстровано нове джерело. id: " + res.ID, endpoints["viewFounts"].ServerPath}}, nil
}

var reMYSQL = regexp.MustCompile(`^mysql://(\w+)\?table=(\w+)`)

func importTest(r *http.Request, user *auth.User, params httprouter.Params) (serverhttp.RESTResponse, error) {
	if user == nil {
		return serverhttp.RESTResponse{Status: http.StatusUnauthorized, Data: serverhttp.RESTErrorOld{basis.ErrAuthenticated.Error()}}, basis.ErrAuthenticated
	}
	type dataFount struct {
		ID string `json:"id"`
	}
	var data dataFount
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	defer r.Body.Close()
	if err != nil {
		return serverhttp.RESTResponse{
			Status: http.StatusInternalServerError,
			Data:   serverhttp.RESTErrorOld{err.Error()},
		}, basis.ErrJSONFormat
	}

	f, err := fountOp.Read(user.Identity(), data.ID)
	if err != nil {
		return serverhttp.RESTResponse{Status: http.StatusNonAuthoritativeInfo, Data: serverhttp.RESTErrorOld{errors.Cause(err).Error()}}, err
	}
	importTo := "flow"
	source := f.URL
	dbKey := ""
	if f.ImportType == importer_factory.NoneImporter {
		importTo = "object_import"
	} else if f.ImportType == importer.ImportType(htmlimporter.InterfaceKey) {
		dbKey = f.ImportDetailsParams
	} else if !user.IsGroupMember(idFountScannerGroup) {
		err = errors.New("no rights for this action")
		return serverhttp.RESTResponse{Status: http.StatusNonAuthoritativeInfo, Data: serverhttp.RESTErrorOld{errors.Cause(err).Error()}}, err
	} else if f.ToObject {
		importTo = "object_import"
		mySQLData := reMYSQL.FindStringSubmatch(f.URL)
		if len(mySQLData) > 0 {
			source = mySQLData[2]
			dbKey = mySQLData[1]
		} else {
			err = errors.New("can't parse url: " + f.URL)
			return serverhttp.RESTResponse{Status: http.StatusNonAuthoritativeInfo, Data: serverhttp.RESTErrorOld{errors.Cause(err).Error()}}, err
		}

	}

	importerOp, ok := importerOperators[string(f.ImportType)]
	if !ok {
		return serverhttp.RESTError(errors.New("не можу обробити джерело типу " + string(f.ImportType))), err

	}

	cnt, err := importer.Task(user.Identity(), importerOp, objectsOp, flowOp, string(f.ImportType), source, importTo, dbKey, data.ID, user.Nickname, f.Tags, f.RView, true)
	if err != nil {
		return serverhttp.RESTResponse{Status: http.StatusNonAuthoritativeInfo, Data: serverhttp.RESTErrorOld{errors.Cause(err).Error()}}, err
	}
	redirectTo := endpoints["viewFlowsFount"].Path(data.ID)
	if f.ImportType == importer_factory.NoneImporter {
		redirectTo = endpoints["importFlows"].ServerPath
	}
	return serverhttp.RESTResponse{Status: http.StatusOK, Data: serverhttp.RESTDataMessage{"import data finished; add " + strconv.Itoa(cnt) + " new items.comp", redirectTo}}, nil
}
