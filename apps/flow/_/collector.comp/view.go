package collector_comp

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/confidenter/auth"

	"github.com/pavlo67/punctum/collector/importer/html"
	"github.com/pavlo67/punctum/processor.old/news"
	"github.com/pavlo67/punctum/processor/sources"

	"github.com/pavlo67/punctum/basis/viewshtml"
	"github.com/pavlo67/punctum/collector/importer"
)

var title = "Новини"

var htmlLeft, htmlLeftNoUser string

const markImportedFlow = "noteID:"

func initHTML() {
	htmlLeft =
		`<div class="ut">` + viewshtml.My + ` Збирач даних` + "</div>\n" +
			`<div class="ul">` + viewshtml.My + ` <a href="` + endpoints["viewFounts"].ServerPath + `"> джерела даних</a>` + "</div>\n" +
			`<div class="ul">` + viewshtml.My + ` <a href="` + endpoints["viewFlowsFount"].Path("all") + `"> потоки новин</a>` + "</div>\n"

	htmlLeftNoUser =
		`<div class="ut gray">` + viewshtml.No + ` Збирач даних` + "</div>\n" +
			`<div class="ul gray">` + viewshtml.No + ` джерела даних</div>` + "</div>\n" +
			`<div class="ul gray">` + viewshtml.No + ` потоки новин` + "</div>\n"
}

func fountTemplator(r *http.Request, user *auth.User) map[string]string {
	if user == nil || user.ID == "" {
		return map[string]string{
			"left.comp": htmlLeftNoUser,
		}
	}
	return map[string]string{
		"left.comp": htmlLeft,
		"front":     htmlFront,
	}
}

func authorFilter(allAuthors map[string]int64, author string) string {
	corpus := ""
	if len(allAuthors) > 0 {
		keys := []string{}
		for key := range allAuthors {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		corpus += `
		<select id="` + listeners["authorsFromImport"].ID + `">
`
		for _, a := range keys {
			selected := ""
			if a == author {
				selected = " selected"
			}
			corpus += `
			<option value="` + html.EscapeString(a) + `"` + selected + `>` + html.EscapeString(a) + ` [` + strconv.FormatInt(allAuthors[a], 10) + `]</option>
`
		}
		corpus += `
		</select>
`
	}
	return corpus
}

func htmlFount(fount *sources.Item, id string, userIsOwner bool) string {
	context := `
	<div>
		<h5>` + fount.Title + `</h4>
		<p>` + fount.URL + `</p>
		<p> <a href="` + endpoints["viewFlowsFount"].Path(id) + `"> &lt;НОВИНИ&gt; </a>
`
	context += `
		<p>
` + showFountTags(fount.Tags) + `
		</p>`
	if userIsOwner {
		context += `
		<a href="` + endpoints["editFount"].Path(id) + `">Редагувати запис</a>
		<a href="#" id="` + listeners["deleteFount"].ID[0:len(listeners["deleteFount"].ID)-1] + id + `">Видалити запис</a>`
	}
	context += `
	</div>
`

	return context
}

func htmlFountFlows(flows []news.Item) string {
	var corpus string
	corpus += showFlows(flows, nil)
	return corpus
}

func htmlFountFree(tag string, flows []news.Item, flowsCnt, page uint64, allTags []sources.FountTag) string {
	var corpus string
	if allTags != nil {
		path := endpoints["freeFlows"].Path()
		corpus += showTags(tag, path, allTags, flowsCnt)
	}
	if flows != nil {
		corpus += showFlows(flows, nil)
	}
	return corpus
}

var reTimeSec = regexp.MustCompile(`:\d\d$`)
var reFlowInNote = regexp.MustCompile(`^` + markImportedFlow)

func showFlows(flows []news.Item, tagFounts *map[int64]string) string {
	corpus := ""
	if tagFounts != nil {
		corpus += `
	<br clear=all>
		<span style="float:right;padding:5px;min-width:190px;" class="border">
		Джерела:
		<p>
`
		for i, n := range *tagFounts {
			corpus += `
			<linker_server_http>` + n + ` &nbsp;<small><span style="float:right;">
				<a href="` + endpoints["editFount"].Path(strconv.FormatInt(i, 10)) + `">[ред.]</a>
				<a href="` + endpoints["viewFlowsFount"].Path(strconv.FormatInt(i, 10)) + `">[новини]</a></span></small>	
			</linker_server_http>
`
		}
		corpus += `<br><a href="` + endpoints["blankFount"].ServerPath + `">Додати нове</a>
		</span>
`
	}

	// corpus += "\n<div>\n"

	oldFoundAndDate := ""
	liTitle := 0
	for _, f := range flows {
		if f.Summary == "" {
			f.Summary = f.Content
		}
		if f.Title == "" {
			f.Title = f.Content
		}
		foundAndDate := `
			<small>Джерело: <i><a href="` + f.SourceURL + `" target=_blank>` + f.SourceURL + `</i></a>
			[` + reTimeSec.ReplaceAllString(f.CreatedAt.Format("02.01.2006 15:04:05"), "") + `]
			</small>
`
		if foundAndDate != oldFoundAndDate {
			if oldFoundAndDate != "" {
				corpus += "\n&nbsp;<br>\n"
			}

			corpus += foundAndDate
			oldFoundAndDate = foundAndDate
		}
		liTitle++

		corpus += "\n<linker_server_http>"

		if reFlowInNote.MatchString(f.ImportedTo) {
			f.ImportedTo = reFlowInNote.ReplaceAllString(f.ImportedTo, "")
			corpus += `[<a href="` + itemsEndpoints["view"].Path(f.ImportedTo) + `" target=_blank>...</a>] `
		} else {
			corpus += `<small><strong> [<a style="cursor:pointer;" id="` + listeners["importFlow"].ID[:len(listeners["importFlow"].ID)-1] + f.ID + `" target=_blank>&lt;&lt;&lt;</a>] </strong></small>`
		}
		corpus += `<a style="cursor:pointer;" target=_blank><img src="/images/fb.png" id="share_fb_` + strconv.Itoa(liTitle) + `" style="margin-top:2px;vertical-align: top;"></a> `

		corpus += `
		<a href="https://twitter.com/intent/tweet?text=` + url.QueryEscape(f.Title) + `&url=` + url.QueryEscape(f.URL) + `" target=_blank>
			<img src="/images/twitter.png" style="margin-top:2px;vertical-align: top;">
		</a>`

		hashTags := ""
		media := ""
		if f.Media != nil {
			if len(f.Media.HashTags) > 0 {
				hashTags += "<br>"
				for _, h := range f.Media.HashTags {
					if h == "" {
						continue
					}
					// todo!!! make universally
					//hashTags += `<a href="https://twitter.com/hashtag/` + h + `"  target=_blank>#` + h + "</a> "
					hashTags += `<a href="` + hrefForHashTag(f.SourceURL) + h + `"  target=_blank>#` + h + "</a> "
				}
			}
			if len(f.Media.Pictures) > 0 {
				media += "<br>"
				for _, p := range f.Media.Pictures {
					if p.ImageUrl == "" {
						continue
					}
					// todo!!! make universally
					media += `<a href="` + p.HREFUrl + `" target=_blank><img style="max-width:300px;" src="` + p.ImageUrl + `"></a> `
				}
			}

		}
		corpus += `
				<a href="` + f.URL + `" target=_blank  id="news_title_id_` + strconv.Itoa(liTitle) + `">` + f.Title + `</a>
			</linker_server_http>
			<div id="news_content_id_` + strconv.Itoa(liTitle) + `" style="display:none; position:absolute; background-color:#eeffee; padding: 10px 10px;">` +
			f.Summary + hashTags + media + `</div>
`
	}
	//corpus += "\n<p>" +
	//	admin.GetPagination(page, pagination, flowsCnt, path, "", []string{"created_at-"}) +
	//	"\n</div>"

	//corpus += "\n<small><b>" +
	//	crud.HTMLPagination(
	//		page,
	//		pagination,
	//		flowsCnt,
	//		path,
	//		"",
	//		[]string{"created_at-"}) + "</b>\n</small>\n"

	return corpus
}

func showTags(tag, path string, allTags []sources.FountTag, flowsCnt uint64) string {

	var previousTag, corpus string

	if tag == allMyFlows {
		corpus += "\n<tr><td><b>- " + allMyFlows + "</b></td></tr>\n"
		// " [" + strconv.FormatUint(flowsCnt, 10) + "]"

	} else {
		corpus += "\n<tr><td>" + `<a href="` + endpoints["viewFlowsFount"].Path(allMyFountsID) + `">- ` + allMyFlows + "</a></td></tr>\n"
	}

	for _, t := range allTags {
		if previousTag == t.Tag {
			continue
		}
		previousTag = t.Tag
		if t.Tag == tag {
			corpus += "\n<tr><td><b>- " + t.Tag + "</b></td></tr>\n"
			// ` [` + strconv.FormatUint(flowsCnt, 10) + "]"
		} else {
			corpus += "\n<tr><td>" + `<a href="` + path + t.Tag + `">- ` + t.Tag + "</a></td></tr>\n"
		}
	}

	return corpus
}

func showFountTags(tags string) string {

	context := ""
	for _, t := range strings.Split(tags, ";") {
		t = strings.Trim(t, " ")
		if t != "" {
			context += `
			<a href="` + endpoints["viewFlowsByTag"].Path(t) + `">&lt;` + t + `&gt;</a> 
`
		}
	}
	return context
}

func htmlFounts(founts []sources.Item, userInScannerGroup bool, userIS auth.ID) string {

	context := `
	+ <a href="` + endpoints["blankFount"].ServerPath + `"><i>Додати нове джерело</i></a><p>
	<div>
`
	for _, sc := range founts {
		context += `
			<p>
`
		if userInScannerGroup || sc.ImportType == importer.ImportType(htmlimporter.InterfaceKey) {
			context += `
					<button id="test_import_fount_` + strconv.FormatInt(sc.ID, 10) + `">закачати вже ...</button>
			`
		}
		var markNoActive string
		if sc.ToObject == false && sc.ToFlow == false {
			markNoActive = ` style="color:#cccccc;" `
		}
		context += `
			<a href="` + endpoints["viewFount"].Path(strconv.FormatInt(sc.ID, 10)) +
			`" ` + markNoActive + `><b>` + sc.Title + `</b></a> 
			<a href="` + endpoints["viewFlowsFount"].Path(strconv.FormatInt(sc.ID, 10)) +
			`"> &lt;НОВИНИ&gt; </a> 
			(` + sc.URL + `) 
`
		context += `
				<span style="float:right;">
`
		context += showFountTags(sc.Tags)
		if userIS == sc.ROwner {
			context += `
				<a href="` + endpoints["editFount"].Path(strconv.FormatInt(sc.ID, 10)) + `">[Редагувати запис] </a>
				<a href="#" id="` + listeners["deleteFount"].ID[0:len(listeners["deleteFount"].ID)-1] + strconv.FormatInt(sc.ID, 10) + `"> [Видалити запис]</a>`
		}
		context += `
				</span>
			</p>	
`
	}
	context += ` 
	</div>
`
	return context
}

const stopFount = "призупинити"

func setFields() []viewshtml.Field {
	return []viewshtml.Field{
		{"id", "", "hidden", "", nil, nil},
		{"title", "заголовок", "", "", nil, nil},
		{"url", "посилання", "", "", nil, nil},
		{"direct", "підключення", "select", "", nil, nil},
		{"tags", "теми", "", "", nil, nil},
		{},
		{"updated_at", "востаннє відредаґовано", "view", "datetime", viewshtml.NotEmpty, nil},
		{"importDetailsType", "", "hidden", "", nil, nil},
		{"importDetailsParams", "", "hidden", "", nil, nil},
	}

}

// [][2]string{{"щогодинний сканер", "flow"}, {stopFount, ""}}, ""

func setAdminFields() []viewshtml.Field {
	return []viewshtml.Field{
		{"separator", "НАЛАШТУВАННЯ", "view", "", nil, nil},
		{"importStartRegexp", "видалити все до", "", "", nil, nil},
		{"importFinishRegexp", "видалити все після", "", "", nil, nil},
		{"importSplitRegexp", "розбити на частини, використовуючи", "", "", nil, nil},
		{"importTagsList", "залишити тільки такі теги", "", "", nil, nil},
		{"titleRegexp", "шаблон для заголовка частини", "", "", nil, nil},
		{listeners["testRegexp"].ID, "Тест налаштувань", "button", "", nil, nil},
		{listeners["exportRegexp"].ID, "Експортувати налаштуваня", "button", "", nil, nil},
	}

}

func htmlBlankFount() string {
	var titleUrl, htmlCode string
	fieldsForCreate := setFields()
	var buttonSave = viewshtml.Field{
		listeners["createFount"].ID,
		"Зберегти запис",
		"button",
		"",
		nil,
		nil,
	}
	fieldsForCreate = append(fieldsForCreate, buttonSave)

	fieldsForCreate[5] = viewshtml.Field{"type", "", "hidden", "", nil, nil}
	data := map[string]string{}
	values := map[string]viewshtml.SelectString{
		"direct": {{"щогодинний сканер", "flow"}, {stopFount, ""}},
	}

	context := "\n<form>\n"
	for _, f := range fieldsForCreate {
		if f.Type != "view" {
			titleUrl, htmlCode = viewshtml.FieldEdit("fount_edit", f, data, values, nil)
			id := f.Key
			context += `<div>` + "\n"
			if titleUrl != "" {
				context += `<label for="` + id + `">` + titleUrl + ":</label>\n"
			}
			context += htmlCode + "\n</div>\n"
		}
	}
	context += "</form>\n"

	return context
}

func htmlEditFount(fount *sources.Item, id string, userIsAdmin bool) string {
	var titleUrl, htmlCode string
	fieldsForCreate := setFields()
	var buttonSave = viewshtml.Field{
		listeners["updateFount"].ID,
		"Зберегти запис",
		"button",
		"",
		nil,
		nil,
	}
	fieldsForCreate = append(fieldsForCreate, buttonSave)
	data := map[string]string{}

	var typeNames = viewshtml.SelectString{{"не вказано", ""}}
	for k := range importerOperators {
		typeNames = append(typeNames, [2]string{k, k})
	}
	values := map[string]viewshtml.SelectString{
		"type":   typeNames,
		"direct": {{"щогодинний сканер", "flow"}, {stopFount, ""}},
	}

	if userIsAdmin {

		fieldsForCreate[5] = viewshtml.Field{"type", "тип", "select", "", nil, nil}
		if fount.ImportType == importer.ImportType(htmlimporter.InterfaceKey) {
			fieldsForCreate = append(fieldsForCreate, setAdminFields()...)
			if fount.ImportDetailsParams != "" {
				var p htmlimporter.ImportParams
				err := json.Unmarshal([]byte(fount.ImportDetailsParams), &p)
				if err == nil {
					data["import_start_regexp"] = p.ImportStartRegexp
					data["import_finish_regexp"] = p.ImportFinishRegexp
					data["import_split_regexp"] = p.ImportSeparateRegexp
					data["import_tags_list"] = strings.Join(p.AcceptableTags, " ")
					data["title_regexp"] = p.TitleRegexp
				} else {
					log.Println("can't parse fount.ImportDetailsParams: ", fount.ImportDetailsParams, err)
				}
			}
		}
	} else {
		fieldsForCreate[5] = viewshtml.Field{"type", "", "hidden", "", nil, nil}
	}

	data["id"] = strconv.FormatInt(fount.ID, 10)
	data["type"] = string(fount.ImportType)
	data["import_details_type"] = fount.ImportDetailsType
	data["import_details_params"] = fount.ImportDetailsParams
	data["title"] = fount.Title
	data["url"] = fount.URL
	data[string(fount.ImportType)] = "select"
	if fount.ToFlow {
		data["flow"] = "select"
	} else {
		data[stopFount] = "select"
	}
	data["tags"] = fount.Tags
	data["updated_at"] = fount.CreatedAt.Format("02.01.2006 15:04:05")

	context := "\n<form>\n"
	for _, f := range fieldsForCreate {
		titleUrl, htmlCode = viewshtml.FieldEdit("fount_edit", f, data, values, nil)
		context += `<div>` + "\n"
		if titleUrl != "" {
			context += `<label for="` + f.Key + `">` + titleUrl + ":</label>\n"
		}
		context += htmlCode + "\n</div>\n"
	}
	context += "</form>\n"

	return context
}

var reTwitterHashTag = regexp.MustCompile(`(?i)twitter`)
var reInstagramHashTag = regexp.MustCompile(`(?i)instagram`)

func hrefForHashTag(fountURL string) string {
	if reTwitterHashTag.MatchString(fountURL) {
		return "https://twitter.com/hashtag/"
	}
	if reInstagramHashTag.MatchString(fountURL) {
		return "https://www.instagram.com/explore/tags.comp/"
	}
	return ""
}
