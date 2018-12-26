package rest_flower_serverhttp_jsschmhr

//import (
//	"encoding/json"
//	"net/http"
//	"reflect"
//
//	"github.com/julienschmidt/httprouter"
//	"github.com/pkg/errors"
//
//	"github.com/pavlo67/punctum/basis"
//	"github.com/pavlo67/punctum/confidenter/auth"
//	"github.com/pavlo67/punctum/crud"
//	"github.com/pavlo67/punctum/fronthttp/serverhttp/serverhttp_jschmhr"
//)
//
//const pageLengthDefault = 200
//
//func (rcOp *rest_datastore_serverhttp_jschmhr) Save(user *auth.User, r *http.Request, params httprouter.Params) (serverhttp.RESTResponse, error) {
//	defer r.Body.Close()
//	crudType := params.ByName("type")
//
//	crudOp, err := rcOp.getOp(crudType)
//	if err != nil {
//		return serverhttp.RESTError(basis.ErrCantPerform), err
//	}
//
//	description, err := crudOp.Describe()
//	if err != nil {
//		return serverhttp.RESTError(basis.ErrCantPerform), err
//	}
//
//	data := reflect.New(reflect.ValueOf(description.Exemplar).Elem().Type()).Interface()
//	err = json.NewDecoder(r.Body).Decode(&data)
//	if err != nil {
//		return serverhttp.RESTError(basis.ErrJSONFormat), errors.Wrap(err, "on decode r.Body")
//	}
//
//	id, err := crudOp.IDFromNative(data)
//	if err != nil {
//		return serverhttp.RESTError(basis.ErrJSONFormat), errors.Wrapf(err, "on crudOp.IDFromNative(%#v)", data)
//	}
//
//	var info string
//	if id == "" {
//		id, err = crudOp.Create(user.Identity().String(), data)
//		info = `Створено запис (` + id + `)`
//
//	} else {
//		var res crud.Result
//		res, err = crudOp.Update(user.Identity().String(), data)
//		if res.NumOk < 1 {
//			info = `Запис (` + id + `) не зазнав змін.`
//		} else {
//			info = `Запис (` + id + `) оновлено.`
//		}
//	}
//	if err != nil {
//		return serverhttp.RESTError(basis.ErrCantPerform), err
//	}
//
//	return serverhttp.RESTResponse{
//		Data: crud.ResultData{
//			IDs:  []string{id},
//			Info: info,
//		},
//	}, nil
//
//}
//
//func (rcOp *rest_datastore_serverhttp_jschmhr) DeleteList(user *auth.User, r *http.Request, params httprouter.Params) (serverhttp.RESTResponse, error) {
//	defer r.Body.Close()
//	crudType := params.ByName("type")
//
//	crudOp, err := rcOp.getOp(crudType)
//	if err != nil {
//		return serverhttp.RESTError(basis.ErrCantPerform), err
//	}
//
//	id := params.ByName("id")
//
//	res, err := crudOp.DeleteList(user.Identity().String(), id)
//	if err != nil {
//		return serverhttp.RESTError(basis.ErrCantPerform), err
//	}
//
//	var info string
//	if res.NumOk < 1 {
//		info = `Запис (` + id + `) не зазнав змін.`
//	} else {
//		info = `Запис (` + id + `) вилучено.`
//	}
//
//	return serverhttp.RESTResponse{
//		Data: crud.ResultData{
//			IDs:  []string{id},
//			Info: info,
//		},
//	}, nil
//}
//
//func (rcOp *rest_datastore_serverhttp_jschmhr) ReadList(user *auth.User, r *http.Request, params httprouter.Params) (serverhttp.RESTResponse, error) {
//	defer r.Body.Close()
//	crudType := params.ByName("type")
//
//	crudOp, err := rcOp.getOp(crudType)
//	if err != nil {
//		return serverhttp.RESTError(basis.ErrCantPerform), err
//	}
//
//	description, err := crudOp.Describe()
//	if err != nil {
//		l.Errorf("can't crudOp.Describe(): %s", err)
//	}
//
//	options, _, err := crud.ReadOptionsFromRequest(r, pageLengthDefault, description.SortByDefault)
//	if err != nil || options == nil {
//		return serverhttp.RESTError(errors.New("не вдається прочитати параметри запиту")), err
//	}
//
//	items, allCnt, err := crudOp.ReadList(user.Identity().String(), options)
//	if err != nil {
//		return serverhttp.RESTError(basis.ErrCantPerform), err
//	}
//
//	return serverhttp.RESTResponse{
//		Data: crud.ReadListData{
//			Description: description,
//			Items:       items,
//			AllCnt:      allCnt,
//		},
//	}, nil
//}
