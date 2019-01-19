package server

type RouteParam struct {
	Name  string
	Value string
}

type RouteParams []RouteParam

func (p RouteParams) ByName(name string) string {
	for i := range p {
		if p[i].Name == name {
			return p[i].Value
		}
	}
	return ""
}
