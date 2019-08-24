package common

type Param struct {
	Name  string
	Value string
}

type Params []Param

func (p Params) ByName(name string) string {
	for i := range p {
		if p[i].Name == name {
			return p[i].Value
		}
	}
	return ""
}

func (p Params) ByNum(num uint) string {
	if int(num) >= len(p) {
		return ""
	}

	return p[num].Value
}

//func (p Info) AllExcept(names ...string) []string {
//	var values []string
//
//PARAM:
//	for _, param := range p {
//		for _, name := range names {
//			if param.Title == name {
//				continue PARAM
//			}
//			values = append(values, param.Value)
//		}
//	}
//	return values
//}
