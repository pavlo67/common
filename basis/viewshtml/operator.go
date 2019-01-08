package viewshtml

type Operator interface {
	// HTMLToEdit returns an html code to view the element with selected modelID
	HTMLToView(field Field, data string, values SelectString, params map[string]string) string

	// HTMLToEdit returns an html code to edit the element with selected modelID
	HTMLToEdit(field Field, data string, values SelectString, params map[string]string) string

	// HTMLToLoad returns an html code to load all required js/css partes
	HTMLToLoad() string
}

//func New(htmlToJoin HTMLToJoin, path, componentKey string, endpoints map[string]config.Endpoint, listeners map[string]config.Listener) Operator {
//	return &frontComponent{
//		htmlFront:  middleware.HandleJS(path, componentKey, endpoints, listeners),
//		htmlToJoin: htmlToJoin,
//	}
//}
//

func New(htmlToView, htmlToEdit HTMLToShow, htmlFront string) Operator {
	return &frontComponent{
		htmlToView: htmlToView,
		htmlToEdit: htmlToEdit,
		htmlFront:  htmlFront,
	}
}

// -----------------------------------------------------------------------------------------------------

var _ Operator = &frontComponent{}

type HTMLToShow func(field Field, data string, values SelectString, params map[string]string) string

type frontComponent struct {
	htmlToView HTMLToShow
	htmlToEdit HTMLToShow
	htmlFront  string
}

// HTMLToEdit returns an html code to join this front component to the element with selected modelID
func (fc frontComponent) HTMLToEdit(field Field, data string, values SelectString, params map[string]string) string {
	return fc.htmlToEdit(field, data, values, params)
}

// HTMLToView returns an html code to join this front component to the element with selected modelID
func (fc frontComponent) HTMLToView(field Field, data string, values SelectString, params map[string]string) string {
	return fc.htmlToView(field, data, values, params)
}

// HTMLToLoad returns an html code to load all required js/css partes
func (fc frontComponent) HTMLToLoad() string {
	return fc.htmlFront
}
