package server_http

type Params map[string]string

//type Param struct {
//	Name  string
//	Left string
//}
//
//type Params []Param
//
//func (p Params) ByName(name string) string {
//	for i := range p {
//		if p[i].Name == name {
//			return p[i].Left
//		}
//	}
//	return ""
//}
//
//func (p Params) ByNum(num uint) string {
//	if int(num) >= len(p) {
//		return ""
//	}
//
//	return p[num].Left
//}

//func (p Info) AllExcept(names ...string) []string {
//	var values []string
//
//PARAM:
//	for _, param := range p {
//		for _, name := range names {
//			if param.Title == name {
//				continue PARAM
//			}
//			values = append(values, param.Left)
//		}
//	}
//	return values
//}
