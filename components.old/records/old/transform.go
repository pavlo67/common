package old

//type Transform func(object *items.Object) (*items.Object, error)
//
//func TransformAll(userIS basis.UserIS, obj Operator, options *content.ListOptions, method Transform, verbose bool) (numOk, numErr uint64, err error) {
//	// !!! to use for admin purposes only
//
//	options.ForAdmin = true
//	objs, _, err := obj.ReadList(userIS, options)
//	if err != nil {
//		return 0, 0, err
//	}
//
//	var o1 *items.Object
//
//	if verbose {
//		fmt.Println("total objects to transform: ", len(objs))
//	}
//	for i, o := range objs {
//		if i > 0 && verbose && i%10 == 0 {
//			fmt.Println(i)
//		}
//		o1, err = method(&o)
//		if err != nil {
//			numErr++
//			log.Printf("erron on transform (%#v): %s", o, err)
//			continue
//		} else if o1 == nil {
//			numErr++
//			log.Printf("erron on transform (%#v): %s", o, "no result")
//			continue
//		}
//
//		ownerIdentity := basis.UserIS(o1.ROwner).Identity()
//		res, err := obj.Update(&ownerIdentity, o1)
//		if err != nil {
//			numErr++
//			log.Printf("erron on update (%#v): %s", o1, err)
//			continue
//		} else if res.NumOk < 1 {
//			numErr++
//			log.Printf("erron on update (%#v): %s", o1, "no result")
//			continue
//		} else {
//			numOk++
//		}
//	}
//
//	return numOk, numErr, nil
//}
