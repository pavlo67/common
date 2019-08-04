package collector_comp

import (
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/partes/crud/selectors"
	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis/strlib"
	"github.com/pavlo67/punctum/filer"
	"github.com/pavlo67/punctum/notebook/notes"
	"github.com/pavlo67/punctum/processor.old/news"
	"github.com/pavlo67/punctum/processor/sources"
	"github.com/pavlo67/punctum/starter/joiner"

	"github.com/pavlo67/punctum/things_old/files"

	"github.com/pavlo67/punctum/_notebook/notebook.comp/note"
)

const maxBriefLength = 1000

func flowsByTag(userIS auth.ID, tag string, options *content.ListOptions) ([]news.Item, uint64, []sources.FountTag, *map[int64]string, error) {
	selector := selectors.FieldEqual("r_view", string(auth.String()))
	userFountsTags, err := fountOp.ReadTags(auth, selector)
	if err != nil {
		return nil, 0, nil, nil, err
	}
	if tag == "" {
		return nil, 0, userFountsTags, nil, nil
	}
	var val = []interface{}{}
	for _, v := range userFountsTags {
		if v.Tag == tag || tag == allMyFlows {
			val = append(val, v.FountID)
		}
	}
	var tagFounts = map[int64]string{}
	selector = selectors.FieldEqual("id", val...)
	founts, _, err := fountOp.ReadList(auth, nil, selector)
	if err != nil {
		return nil, 0, nil, nil, err
	}
	var fountsIS = []interface{}{}
	for _, f := range founts {
		fountsIS = append(fountsIS, joiner.SystemDomain()+"/fount/"+strconv.FormatInt(f.ID, 10))
		tagFounts[f.ID] = f.Title
	}
	selector = selectors.FieldEqual("fount_is", fountsIS...)
	flows, flowsCnt, err := flowOp.ReadList(auth, options, selector)
	if err != nil {
		return nil, 0, nil, nil, err
	}
	return flows, uint64(flowsCnt), userFountsTags, &tagFounts, nil
}

func flowToObject(userIS auth.ID, id string) (string, error) {

	item, err := flowOp.Read(userIS, auth.ID(id))
	if err != nil {
		return "", errors.Wrapf(err, "on flowOp.Read (id=%v):", id)
	}

	sel := selectors.FieldEqual("fount_id", item.FountIS.Identity().ID)
	fountTags, err := fountOp.ReadTags(userIS, sel)
	tags := ""
	for _, t := range fountTags {
		tags += t.Tag + "; "
	}

	globalIS := joiner.SystemDomain() + "/flow/" + item.ID
	o := notes.Item{
		GlobalIS:   globalIS,
		ROwner:     userIS.String(),
		RView:      userIS.String(),
		Visibility: things_old.Private,
		Genus:      note.GenusKey,
		Name:       item.Title,
		Content:    item.Summary + item.Content + "\n\n[" + item.URL + "]\n[" + item.SourceURL + "]",
		// Links: []items.Link{{Original: globalIS}},
	}

	runes := []rune(o.Content)
	if len(runes) > maxBriefLength {
		o.Brief = string(runes[0:maxBriefLength]) + "..."
	} else {
		o.Brief = o.Content
	}

	notes.AddTags(userIS, &o, tags)

	if item.Media != nil {
		for _, p := range item.Media.Pictures {
			if p.ImageUrl == "" {
				continue
			}
			fileName, err := getPictureFromFlow(userIS, p.ImageUrl)
			if err != nil {
				log.Println("can't get picture: ", p.ImageUrl, " from flow.id=", item.ID)
				continue
			}
			o.Links = append(o.Links, notes.Item{
				Type: files.LinkType,
				Name: fileName,
				To:   fileName,
			})
		}
	}

	res, err := objectsOp.Create(userIS, &o)
	if err != nil {
		return "", errors.Wrapf(err, "can't create import object: %v", o)
	}

	itemID, err := strconv.ParseInt(item.ID, 10, 64)
	if err != nil {
		log.Println("add flow to object; can't ParseInt:", item.ID, err)
	} else {
		err = flowOp.ImportTo(userIS, itemID, markImportedFlow+res)
		if err != nil {
			log.Println("add flow to object; can't mark as imported flow id:", item.ID, err)
		}
	}

	return res, nil
}

var imageMIME = regexp.MustCompile(`(?i)image/`)

func getPictureFromFlow(userIS auth.ID, url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//_, _, err = filer.comp.pathToFileOrDir(userIS, pathRepository, "")
	//if err != nil {
	//	return "", "", err
	//}

	contentType := resp.Header.Get("Content-Type")

	ext := ".jpg"
	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			return "", err
		}
		if imageMIME.MatchString(t) {
			//log.Println("UserIS MIMEType: ", t)
			ext = "." + imageMIME.ReplaceAllString(t, "")
		}
	}

	fileName := strlib.RandomString(12)
	//picFile := pathRepository + userIS.SystemDomain + "_" + userIS.LocalPath + "_" + userIS.Label + "/" +
	//	fileName + ext

	content, err := ioutil.ReadAll(resp.Body)
	info := files.Item{
		Data: filer.Data{
			Name:   fileName + ext,
			RView:  userIS.String(),
			ROwner: userIS.String(),
		},
		Content: content,
	}
	picFile, err := filesOp.Create(userIS, &info)
	if err != nil {
		return "", err
	}
	//out, err := os.Create(picFile)
	//if err != nil {
	//	return "", "", err
	//}
	//defer out.Close()
	//io.Copy(out, resp.Body)

	return picFile, nil
}
