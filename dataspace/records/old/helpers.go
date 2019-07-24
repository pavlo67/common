package old

import (
	"fmt"

	"github.com/pkg/errors"

	"strconv"

	"github.com/pavlo67/partes/crud/selectors"
	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/notebook/links"
	"github.com/pavlo67/punctum/things_old/files"
)

//func RViewDefault(dataManagers rights.Managers) auth.ID {
//	return dataManagers[rights.ViewDefault]
//}

func FillFilesIDs(linksList []links.Item) ([]links.Item, bool) {
	var emptyFilesID []int
	var maxFilesID = 0
	var needUpdate = false
	for i, l := range linksList {
		if l.Type == files.LinkType {
			if l.ID == "" {
				emptyFilesID = append(emptyFilesID, i)
			} else {
				id, err := strconv.Atoi(l.ID)
				if err == nil {
					if id > maxFilesID {
						maxFilesID = id
					}
				}
			}
		}
	}
	for _, i := range emptyFilesID {
		needUpdate = true
		maxFilesID++
		linksList[i].ID = strconv.Itoa(maxFilesID)
	}
	return linksList, needUpdate
}

func GetUnique(userIS auth.ID, objectsOp Operator, selector selectors.Selector) (*Item, error) {
	options := &content.ListOptions{Selector: selector}
	objs, _, err := objectsOp.ReadList(userIS, options)
	if err != nil {
		return nil, fmt.Errorf("error on objectsOp.ReadList(%#v, %#v): %s", userIS, options, err)
	}

	if len(objs) < 1 {
		return nil, nil
	} else if len(objs) > 1 {
		return nil, fmt.Errorf("too many objects for objects.GetUnique(%#v): %#v", selector, objs)
	}

	return &objs[0], nil
}

func PutUnique(userIS auth.ID, o Item, objectsOp Operator, selector selectors.Selector) (string, error) {
	o0, err := GetUnique(userIS, objectsOp, selector)
	if err != nil {
		return "", err
	}
	if o0 == nil {
		return objectsOp.Create(userIS, o)
	}

	o.ID = o0.ID
	_, err = objectsOp.Update(userIS, o)

	return o.ID, err
}

const onContentByID = "on objects.ContentByID(%s, %#v)"

func ContentByID(userIS auth.ID, objectsOp Operator, genusKey string, content Content) error {
	id, err := content.ID()
	if err != nil {
		return errors.Wrapf(err, onContentByID, genusKey, content)
	}

	o, err := objectsOp.Read(userIS, id)
	if err != nil {
		return errors.Wrapf(err, onContentByID, genusKey, content)
	}

	if o.Genus != genusKey {
		// TODO: wrap!!!
		return errors.Wrapf(basis.ErrBadGenus, onContentByID+": %s", genusKey, content, o.Genus)
	}

	return content.FromObject(o)
}

//func DeleteWithFiles(user *auth.User, objectsOp Operator, filesOp files.Operator, id string) (*Item, error) {
//	o, err := objectsOp.Read(user.Identity().String(), id)
//	if err != nil {
//		return o, err
//	}
//	_, err = objectsOp.DeleteList(user.Identity().String(), id)
//	if err != nil {
//		return o, err
//	}
//
//	for _, f := range o.Links {
//		if f.Type != files.LinkType {
//			continue
//		}
//		if f.Name[:len(files.RepoSchema)] != files.RepoSchema {
//			continue
//		}
//		_, err = filesOp.DeleteList(user.Identity().String(), f.Name)
//		if err != nil {
//			log.Println("Err delete file: " + f.Name + " after delete item: " + id)
//		}
//	}
//	return o, nil
//
//}
