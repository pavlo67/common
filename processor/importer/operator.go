package importer

import (
	"github.com/pavlo67/constructor/starter/joiner"

	"github.com/pavlo67/constructor/processor/news"
)

const InterfaceKey joiner.ComponentKey = "importer"

//var ErrNoFount = errors.New("no source is reachable")
//var ErrNoMoreItems = errors.New("no more items.comp")
//var ErrBadItemID = errors.New("bad item id")
//var ErrBadItem = errors.New("bad item")
//var ErrNilItem = errors.New("item is nil")

type Operator interface {
	// Run opens import session with selected data source
	Init(source string) error

	// Next gets the next data entity from the source
	Next() (*news.Item, error)

	Close() error
}

//func Run(impOp Operator, sources []string, newsOp news.Operator) error {
//	var errs basis.Errors
//	var cnt int
//
//	for _, src := range sources {
//		err := impOp.Run(src)
//		if err != nil {
//			errs = append(errs, errors.Wrapf(err, "on impOp.Run(%s)", src))
//			continue
//		}
//
//		for {
//			if cnt%100 == 0 {
//				fmt.Println(cnt)
//			}
//			cnt++
//
//			item, err := impOp.Next()
//			if err != nil {
//				errs = errs.Append(err)
//				log.Printf("on impOp.Next(): %s", err)
//			}
//			if item == nil {
//				continue
//			}
//
//			//if staged {
//			//	err = dsOpStaged.SaveStaged(*item)
//			//} else {
//			//	err = newsOp.Save(*item)
//			//}
//
//			err = newsOp.Save(item)
//
//			if err != nil {
//				errs = errs.Append(err)
//				log.Printf("on dsOp.Save(%#v): %s", *item, err)
//			}
//		}
//	}
//	//err := dsOp.Commit(nil)
//	//if err != nil {
//	//	errs = errs.Append(err)
//	//	log.Printf("on dsOp.Commit(nil): %s", err)
//	//}
//
//	return errs.Err()
//}
