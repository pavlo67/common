package workspace_v1

import (
	"net/http"

	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"

	"github.com/pavlo67/workshop/components/auth"

	"github.com/pavlo67/workshop/apps/workspace/workspace_routes"
)

var _ = server_http.InitEndpoint(&endpoints, "GET", filelib.RelativePath(filelib.CurrentFile(true), workspace_routes.PathBase, workspace_routes.Prefix), nil, workerList, "", l)

var _ server_http.WorkerHTTP = workerList

func workerList(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	items, err := workspaceOp.List(nil, nil)

	l.Infof("%#v", items)

	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, err)
	}

	//items = []data.Item{
	//	{
	//		ID:         "1",
	//		Title:      "2",
	//		Summary:    "3",
	//		URL:        "4",
	//		Embedded:   nil,
	//		Tags:       nil,
	//		Details:    nil,
	//		DetailsRaw: nil,
	//		Status:     crud.Status{},
	//		Origin:     flow.Origin{},
	//	},
	//}

	//items = append(items,
	//	data.Brief{
	//		Brief: crud.Brief{
	//			ID: "111",
	//			// Type:      "",
	//			Title:     "bbbbbbb1111111 bbb!!!",
	//			Summary:   "dfs/;m sasffg dsf-09-0dfg--- -009-",
	//			OriginURL: "http://abc.ru",
	//		},
	//		Embedded: []crud.Brief{
	//			{
	//				Type:      "href",
	//				Title:     "мама мила раму",
	//				Summary:   "а то!",
	//				OriginURL: "http://abc.ru/1",
	//			},
	//			{
	//				Type:      "img",
	//				Title:     "мама мила раму",
	//				Summary:   "а то!",
	//				OriginURL: "http://abc.ru/1.png",
	//			},
	//		},
	//		SavedAt: savedAt,
	//	},
	//	data.Brief{
	//		Brief: crud.Brief{
	//			ID: "222",
	//			// Type:      "",
	//			Title:     "2222222 2222222222",
	//			Summary:   "dfs/sfgncdfjh wtwaert fdthr-",
	//			OriginURL: "http://abc.ru",
	//		},
	//		SavedAt: savedAt,
	//	},
	//	data.Brief{
	//		Brief: crud.Brief{
	//			ID: "333",
	//			// Type:      "",
	//			Title:     "333333 333333 33",
	//			Summary:   "dfs/sfgncdfjh wtwaert fdthr-",
	//			OriginURL: "http://stolica.ru",
	//		},
	//		SavedAt: time.Now(),
	//	},
	//)

	return server.ResponseRESTOk(items)
}
