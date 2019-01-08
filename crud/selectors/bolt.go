package selectors

//func boltMultiArgs(args []interface{}, oper string, data map[string]string) (bool, error) {
//	for _, arg := range args {
//		//log.Println("arg:", arg)
//		res, err := CheckMap(arg, data)
//		if err != nil {
//			return false, err
//		}
//		if strings.ToLower(oper) == "or" && res {
//			return true, nil
//		}
//		if strings.ToLower(oper) == "and" && !res {
//			return false, nil
//		}
//	}
//
//	return true, nil
//}
//
//func BoltPrepare(originalSelector Selector) (keys [][]byte, filteredSelector Selector, err error) {
//	switch s := originalSelector.(type) {
//	case *multi:
//		if strings.ToLower(s.oper) == "and" {
//			var newValues []interface{}
//			for _, v := range s.values {
//				if s1, ok := v.(Selector); ok {
//					switch s2 := s1.(type) {
//					case *in:
//						if s2.field == "id" {
//							keys, err = getPrepareKeys(s2.values)
//							if err != nil {
//								return nil, originalSelector, errors.Wrapf(err, "can't prepare selector: %v", originalSelector)
//							}
//						} else {
//							newValues = append(newValues, v)
//						}
//					default:
//						newValues = append(newValues, v)
//					}
//				} else if v == nil {
//					continue
//				} else {
//					return nil, originalSelector, errors.Wrapf(errors.New("can't prepare"), "selector: %v, %v is not selector", originalSelector, v)
//				}
//			}
//			filteredSelector = And(newValues...)
//		}
//
//	case *in:
//		if s.field == "id" {
//			keys, err = getPrepareKeys(s.values)
//			if err != nil {
//				return nil, originalSelector, errors.Wrapf(err, "can't prepare selector: %v", originalSelector)
//			}
//			filteredSelector = nil
//		}
//
//	default:
//		return nil, originalSelector, nil
//	}
//
//	return keys, filteredSelector, nil
//}
//
//func getPrepareKeys(values []interface{}) (keys [][]byte, err error) {
//	for _, v := range values {
//		if s, ok := v.(string); ok {
//			b, err := strconv.ParseUint(s, 10, 64)
//			if err != nil {
//				return nil, errors.Wrapf(err, "can't parse string '%v' to uint64", s)
//			}
//			keys = append(keys, basis.Uint64ToByte(b))
//		} else if b, ok := v.(uint64); ok {
//			keys = append(keys, basis.Uint64ToByte(b))
//		} else {
//			err = errors.Wrapf(errors.New("is not string or uint64"), "%v, typeOf=%v", v, reflect.TypeOf(v))
//			return nil, err
//		}
//	}
//	return keys, nil
//}
