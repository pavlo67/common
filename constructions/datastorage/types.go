package datastorage

import (
	"github.com/pavlo67/workshop/common/joiner"
)

const CollectionDefault = "storage"

const DataInterfaceKey joiner.InterfaceKey = "storage_data"
const InterfaceKey joiner.InterfaceKey = "storage"

const ListInterfaceKey joiner.InterfaceKey = "list"
const ReadInterfaceKey joiner.InterfaceKey = "read"
const SaveInterfaceKey joiner.InterfaceKey = "save"
const RemoveInterfaceKey joiner.InterfaceKey = "remove"
const CountTagsInterfaceKey joiner.InterfaceKey = "count_tags"
const ListWithTagInterfaceKey joiner.InterfaceKey = "list_with_tag"
