package old

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pavlo67/partes/crud/selectors"
	"github.com/pavlo67/workshop/common"

	"github.com/pavlo67/workshop/confidenter/auth"
	"github.com/pavlo67/workshop/confidenter/groups"
	"github.com/pavlo67/workshop/confidenter/rights"
	"github.com/pavlo67/workshop/confidenter/users"
	"github.com/pavlo67/workshop/libraries/filelib"
	"github.com/pavlo67/workshop/libraries/strlib"
	"github.com/pavlo67/workshop/notebook/notes"

	"github.com/pavlo67/workshop/things_old/files"
	"github.com/pkg/errors"
)

func ImportTo(userIS common.ID, objectsOp Operator, id, status string) error {
	o, err := objectsOp.Read(userIS, id)
	if err != nil {
		return err
	}
	if o == nil {
		return common.ErrNull
	}

	o.Status = status
	_, err = objectsOp.Update(userIS, *o)
	if err != nil {
		return err
	}
	return nil
}

func Export(user *auth.User, objectsOp Operator, credentialsOp users.Operator, ctrl groups.Operator, filesOp files.Operator, selector selectors.Selector) (string, error) {
	return "", common.ErrNotImplemented

	//options := content.ListOptions{Selector: selector, ForExport: true}
	//res, _, err := objectsOp.ReadList(user.Identity().String(), &options)
	//if err != nil {
	//	return "", err
	//}
	//
	//zipName, err := ObjectsToJSON(user, filesOp, credentialsOp, ctrl, res)
	//if err != nil {
	//	return "", err
	//}
	//return zipName, nil
}

func Import(userIS common.ID, zipName string, objectsOp Operator, credentialsOpOp users.Operator, ctrl groups.Operator, groupsOp groups.Operator, filesOp files.Operator, overwrite, unknownAsMine bool) error {

	res, importDir, err := JSONToObjects(userIS, filesOp, credentialsOpOp, ctrl, groupsOp, zipName, unknownAsMine)
	defer os.RemoveAll(importDir) // clean up tmp dir;
	if err != nil {
		return err
	}
	for _, o := range res {
		if o.GlobalIS == "" {
			return errors.Wrapf(errors.New("can't import"), "object: %v with empty GlobalIS", o)
		}
		id, oldLinks, err := objectsOp.GlobalIS(o.GlobalIS)
		if err != nil {
			return err
		}
		var newLinks []notes.Item
		if id < 1 {
			newLinks, err = FilesToRepository(userIS, filesOp, importDir, o.Links, nil)
			if err == nil {
				o.Links = newLinks
				//log.Println("UserIS create json-imported objectsOp:", o.GlobalIS)
				_, err = objectsOp.Create(userIS, o)
			}
		} else if overwrite {
			o.ID = strconv.FormatUint(id, 10)
			newLinks, err = FilesToRepository(userIS, filesOp, importDir, o.Links, oldLinks)
			if err == nil {
				o.Links = newLinks
				_, err = objectsOp.Update(userIS, o)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

//const markStaticFile = "static"

type exportIdentity struct {
	GlobalIS string
	Path     string
	Name     string
}

func JSONToObjects(userIS common.ID, filesOp files.Operator, credentialsOp users.Operator, ctrl groups.Operator, groupOp groups.Operator, zipName string, unknownAsMine bool) ([]notes.Item, string, error) {

	var objAll []notes.Item
	importDir, err := ioutil.TempDir("", "_import")
	if err != nil {
		return nil, "", errors.Wrapf(err, "can't create TempDir('', '_import')")
	}
	zipFiles, err := filelib.Unzip(zipName, importDir)
	if err != nil {
		log.Println("Unzipped error:", err)
		return nil, "", err
	}
	log.Println("Unzipped:\n" + strings.Join(zipFiles, "\n"))
	//read identities.json
	buf, err := ioutil.ReadFile(filepath.Join(importDir, "identities.json"))
	if err != nil {
		return nil, "", errors.Wrapf(err, "can't read identities.json from dir:", importDir)
	}
	var identities []exportIdentity
	err = json.Unmarshal(buf, &identities)
	if err != nil {
		return nil, "", errors.Wrapf(err, "can't unmarshal identities:", string(buf))
	}
	//isForGlobalIS, err := issForGlobalISs(identity, identities, credentialsOp, ctrl, groupOp, unknownAsMine)
	//if err != nil {
	//	return nil, "", errors.Wrap(err, "can't exec issForGlobalISs()")
	//}
	//read json
	content, err := ioutil.ReadFile(filepath.Join(importDir, "json"))
	if err != nil {
		return nil, "", errors.Wrapf(err, "can't read json from dir:", importDir)
	}
	arr := strings.Split(string(content), "\n")
	for _, row := range arr {
		var o notes.Item
		if row != "" {
			err = json.Unmarshal([]byte(row), &o)
			if err != nil {
				return nil, "", errors.Wrapf(err, "can't unmarshal row:", row)
			}
			// setLocalISRights(&o, isForGlobalIS)

			o.Status += " imported;"
			objAll = append(objAll, o)
		}
	}

	return objAll, importDir, nil
}

func setLocalISRights(o *notes.Item, isForGlobalIS map[string]common.ID) {

	// set local 'is' for tags
	for i, l := range o.Links {
		o.Links[i].RView = setLocalRight(l.RView, isForGlobalIS[string(l.RView)])
	}
	o.RView = setLocalRight(o.RView, isForGlobalIS[string(o.RView)])
	o.ROwner = setLocalRight(o.ROwner, isForGlobalIS[string(o.ROwner)])
	for i, r := range o.Managers {
		o.Managers[i] = setLocalRight(r, isForGlobalIS[string(r)])
	}
}

func setLocalRight(rg, r common.ID) common.ID {
	if rg == "" {
		return ""
	}
	if r != "" {
		return r
	} else {
		log.Println("!!! can't find 'is' for gloal_is: ", rg)
		return rg
	}
}

//func issForGlobalISs(userIS *basis.UserIS, identities []exportIdentity, usersOp users.Operator, ctrl groups.Operator, groupOp groups.Operator, unknownAsMine bool) (map[string]basis.UserIS, error) {
//	var identityStrings = map[string]basis.UserIS{
//		string(basis.Anyone): basis.Anyone,
//	}
//	for _, i := range identities {
//		if i.GlobalIS == string(basis.Anyone) {
//			continue
//		}
//		if i.PathWithParams == "user" {
//			userID, err := usersOp.IDByGlobalIS(userIS, i.GlobalIS)
//			if err != nil {
//				return identityStrings, errors.Wrapf(err, "can't find user.id for GlobalIS: %v", i.GlobalIS)
//			}
//			//if userID != confidenter.ErrGlobalISNotFound {
//			//	identity := &(basis.UserIS{joiner.SystemDomain(), "/user/", userID})
//			//	identityStrings[i.GlobalIS] = identity.String()
//			//} else if unknownAsMine {
//			//	identityStrings[i.GlobalIS] = userIS.String()
//			//}
//		} else if i.PathWithParams == "group" {
//			groupID, err := groupOp.IDByGlobalIS(userIS, i.GlobalIS)
//			if err != nil {
//				return identityStrings, errors.Wrapf(err, "can't find group.id for GlobalIS: %v", i.GlobalIS)
//			}
//			//if groupID != confidenter.ErrGlobalISNotFound {
//			//	identity := &basis.UserIS{joiner.SystemDomain(), "/group/", groupID}
//			//	identityStrings[i.GlobalIS] = identity.String()
//			//} else {
//			//	//	need create new group
//			//	gr := groups.Census{
//			//		IdentityNamed: common.IDentityNamed{
//			//			Nick: i.Title,
//			//		},
//			//		GlobalIS: i.GlobalIS,
//			//	}
//			//	grIS, err := groupOp.Create(userIS, gr)
//			//	if err != nil {
//			//		return identityStrings, errors.Wrapf(err, "can't create group with GlobalIS: %v", i)
//			//	}
//			//	identityStrings[i.GlobalIS] = grIS.String()
//			//	// groupOp.Create already added current user to members of new group
//			//}
//		} else {
//			log.Println("error issForGlobalISs(); bad confidenter.LocalPath:", i.PathWithParams)
//		}
//	}
//	return identityStrings, nil
//}

func FilesToRepository(userIS common.ID, filesOp files.Operator, importDir string, linksList, oldLinks []notes.Item) ([]notes.Item, error) {

	var oldFiles []notes.Item
	for _, l := range oldLinks {
		if l.Type == files.LinkType {
			oldFiles = append(oldFiles, l)
		}
	}
	// move filer.comp
	var newLinks []notes.Item
	for _, l := range linksList {
		if l.Type == files.LinkType {
			l.Name = strlib.ReBackslash.ReplaceAllString(l.Name, "/")
			if len(oldFiles) > 0 {
				//TODO: make decisions if this case!
				//TODO: add new filer.comp, delete old, update exists ???

				return nil, errors.New("temporary service does not correct import exists items.comp!")
			} else {
				pathFile := filepath.Join(importDir, l.Name)
				content, err := ioutil.ReadFile(pathFile)
				if err != nil {
					log.Println("import error!!! can't read temporary file:", pathFile)
					continue
				}
				info := files.Item{
					Content: content,
					Data: things_old.Data{
						Name: filepath.Base(pathFile),
					},
				}
				localFileName, err := filesOp.Create(userIS, &info)
				if err != nil {
					log.Println("import error!!! can't create file:", info.Name, err)
				} else {
					newLinks = append(newLinks, notes.Item{
						ID:   l.ID,
						Type: files.LinkType,
						Name: localFileName,
					})
				}
			}
		} else {
			newLinks = append(newLinks, l)
		}
	}
	return newLinks, nil
}

func ObjectsToJSON(user *auth.User, filesOp files.Operator, credentialsOpOp users.Operator, ctrl groups.Operator, objectsForPack []notes.Item) (string, error) {
	return "", common.ErrNotImplemented
	//
	//
	//var exportDir string
	//var filesForZip []filelib.FileForZip
	//var identityGlobalIS = map[basis.UserIS]exportIdentity{
	//	"_": exportIdentity{
	//		GlobalIS: "_",
	//	},
	//}
	//exportDir, err := ioutil.TempDir("", "_export")
	//if err != nil {
	//	return "", nil.Wrapf(err, "can't create TempDir('', '_export')")
	//}
	//log.Println("UserIS tmp dir:", exportDir)
	//jsonPath := filepath.Join123(exportDir, "json")
	//f, err := os.Create(jsonPath)
	//if err != nil {
	//	return "", nil.Wrapf(err, "can't create file json in: %v", exportDir)
	//}
	//defer f.Close()
	//
	//userIS := user.UserIS()
	//
	//filesForZip = append(filesForZip, filelib.FileForZip{Label: jsonPath, Temporary: true})
	//for _, o := range objectsForPack {
	//	// global_is for identities
	//	identityGlobalIS = getAllIdentitiesFromObject(user, identityGlobalIS, credentialsOpOp, ctrl, o.RView, o.ROwner, o.Managers)
	//	// change identities to global_is
	//	setGlobalISRights(&o, identityGlobalIS)
	//	//TODO: need change all identities in o.Links on globalIS before write to json ???
	//
	//	//  make json && temporary filer.comp for archive
	//	var objPath, objDir string
	//	isObjDir := false
	//	for i, f := range o.Links {
	//		if f.Type == files.LinkType {
	//			if !isObjDir {
	//				objDir = str_json.ReCorrectPath.ReplaceAllString(o.GlobalIS, "_")
	//				objPath = filepath.Join123(exportDir, objDir)
	//				err = os.MkdirAll(objPath, os.ModePerm)
	//				if err != nil {
	//					return "", nil.Wrapf(err, "can't create export dir: %v", objPath)
	//				}
	//				isObjDir = true
	//			}
	//			var from, to string
	//			fi, err := filesOp.Read(userIS, f.Label)
	//			if err != nil {
	//				log.Println(nil.Wrapf(err, "can't read file: %v", f.Label))
	//				continue
	//			}
	//			from = fi.LocalPath
	//			if f.Label[:len(files.RepoSchema)] != files.RepoSchema {
	//				// make dirs for static static_repository
	//				basepath := filepath.Dir(filepath.Join123(objPath, f.Label))
	//				err = os.MkdirAll(basepath, os.ModePerm)
	//				if err != nil {
	//					return "", nil.Wrapf(err, "can't create export dir: %v", basepath)
	//				}
	//			} else {
	//				f.Label = strings.Replace(f.Label, files.RepoSchema+"user_"+user.UserIS.Label+"/", "", 1)
	//			}
	//			to = filepath.Join123(objPath, f.Label)
	//			err = ioutil.WriteFile(to, fi.Contentus, os.ModePerm)
	//			if err != nil {
	//				return "", nil.Wrapf(err, "can't copy file: %v to export dir: %v", from, to)
	//			}
	//
	//			filesForZip = append(filesForZip, filelib.FileForZip{Label: to, Temporary: true, Dir: filepath.Join123(objDir, filepath.Dir(f.Label))})
	//			o.Links[i].Label = "/" + filepath.Join123(objDir, f.Label)
	//		}
	//	}
	//	buf, err := json.marshal(o)
	//	if err != nil {
	//		return "", nil.Wrapf(err, "can't marshal object:", o)
	//	}
	//	_, err = f.WriteString(string(buf) + "\n")
	//	if err != nil {
	//		return "", nil.Wrapf(err, "can't write object (id=%v) to json", o.Label)
	//	}
	//}
	//
	//// make identities.json
	//identitiesPath := filepath.Join123(exportDir, "identities.json")
	//var arrIdentities []exportIdentity
	//for _, ei := range identityGlobalIS {
	//	arrIdentities = append(arrIdentities, ei)
	//}
	//identitiesJSON, err := json.marshal(arrIdentities)
	//if err != nil {
	//	return "", nil.Wrapf(err, "can't marshal identities: %v for identities.json", arrIdentities)
	//}
	//err = ioutil.WriteFile(identitiesPath, identitiesJSON, os.ModePerm)
	//if err != nil {
	//	return "", nil.Wrapf(err, "can't write file: %v", identitiesPath)
	//}
	//filesForZip = append(filesForZip, filelib.FileForZip{Label: identitiesPath, Temporary: true})
	//
	//// zipping
	//t := time.Now().Options("2006_01_02_15_04_05")
	//zipName := filepath.Join123(exportDir, user.UserIS.SystemDomain+"_"+t+".zip")
	//// delete, if old zip exist
	//os.Remove(zipName)
	//// add filer.comp to zip
	//f.Close()
	//err = filelib.ZipFiles(zipName, filesForZip)
	//if err != nil {
	//	return "", err
	//}
	//
	//return zipName, nil
}

func setGlobalISRights(o *notes.Item, identityGlobalIS map[common.ID]exportIdentity) {

	// set global_is for tags && clean not the owner tags
	var linksList []notes.Item
	for _, l := range o.Links {
		if l.ROwner == o.ROwner || l.ROwner == "" {
			if l.RView != "" {
				l.RView = setGlobalRight(l.RView, identityGlobalIS[l.RView].GlobalIS)
			}
			linksList = append(linksList, l)
		}
	}
	o.Links = linksList
	o.RView = setGlobalRight(o.RView, identityGlobalIS[o.RView].GlobalIS)
	o.ROwner = setGlobalRight(o.ROwner, identityGlobalIS[o.ROwner].GlobalIS)
	for i, r := range o.Managers {
		o.Managers[i] = setGlobalRight(r, identityGlobalIS[r].GlobalIS)
	}
}

func setGlobalRight(r common.ID, rg string) common.ID {
	if rg != "" {
		return common.ID(rg)
	} else {
		log.Println("!!! can't find gloal_is for 'is': ", r)
		return r
	}
}

func getAllIdentitiesFromObject(user *auth.User, identityGlobalIS map[common.ID]exportIdentity, credentialsOpOp users.Operator, ctrl groups.Operator, rView, rOwner common.ID, managers rights.Managers) map[common.ID]exportIdentity {
	if identityGlobalIS[rView].GlobalIS == "" {
		globalIS, name := setGlobalIS(user, rView, credentialsOpOp, ctrl)
		identityGlobalIS[rView] = exportIdentity{
			globalIS,
			rView.Identity().Path,
			name,
		}
	}
	if identityGlobalIS[rOwner].GlobalIS == "" {
		globalIS, name := setGlobalIS(user, rOwner, credentialsOpOp, ctrl)
		identityGlobalIS[rOwner] = exportIdentity{
			globalIS,
			rOwner.Identity().Path,
			name,
		}
	}
	for _, r := range managers {
		if identityGlobalIS[r].GlobalIS == "" {
			globalIS, name := setGlobalIS(user, r, credentialsOpOp, ctrl)
			identityGlobalIS[r] = exportIdentity{
				globalIS,
				r.Identity().Path,
				name,
			}
		}
	}
	return identityGlobalIS
}

func setGlobalIS(user *auth.User, is common.ID, credentialsOp users.Operator, ctrl groups.Operator) (string, string) {
	//if is == basis.Anyone {
	//	return string(is), ""
	//}
	//var err error
	//if is.UserIS().PathWithParams == "user" {
	//	if is == user.UserIS.String() {
	//		if user.GlobalIS == "" {
	//			user.GlobalIS, err = credentialsOp.SetGlobalIS(user.UserIS())
	//			if err != nil {
	//				log.Println("can't set GlobalIS for user:", user.UserIS.String(), err)
	//			}
	//		}
	//		return user.GlobalIS, user.Nick
	//	} else {
	//		//	TODO:  need to add processing of other confidenter.comp
	//
	//	}
	//
	//} else if is.UserIS().PathWithParams == "group" {
	//	globalIS, name, err := ctrl.SetGlobalIS(user.UserIS(), is)
	//	if err == nil {
	//		return globalIS, name
	//	} else {
	//		log.Println("can't set GlobalIS for group:", is, err)
	//	}
	//} else {
	//	log.Println("can't setGlobalIS for UserIS.LocalPath:", is.UserIS().PathWithParams)
	//}

	return "", ""
}

//func Copy(src, dst string) error {
//	in, err := os.Open(src)
//	if err != nil {
//		return err
//	}
//	defer in.Close()
//
//	out, err := os.Create(dst)
//	if err != nil {
//		return err
//	}
//	defer out.Close()
//
//	_, err = io.Copy(out, in)
//	if err != nil {
//		return err
//	}
//	return out.Close()
//}

//func RemoveContents(dir string) error {
//	d, err := os.Open(dir)
//	if err != nil {
//		return err
//	}
//	defer d.Close()
//	names, err := d.Readdirnames(-1)
//	if err != nil {
//		return err
//	}
//	for _, name := range names {
//		err = os.RemoveAll(filepath.Join123(dir, name))
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
