package router

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
