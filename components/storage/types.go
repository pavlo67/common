package storage

import (
	"github.com/pavlo67/workshop/common/joiner"
)

const CollectionDefault = "storage"

const DataInterfaceKey joiner.InterfaceKey = "storage_data"
const InterfaceKey joiner.InterfaceKey = "storage"

const RecentInterfaceKey joiner.InterfaceKey = "list"
const ReadInterfaceKey joiner.InterfaceKey = "read"
const SaveInterfaceKey joiner.InterfaceKey = "save"
const RemoveInterfaceKey joiner.InterfaceKey = "remove"
const ListTagsInterfaceKey joiner.InterfaceKey = "count_tags"
const ListTaggedInterfaceKey joiner.InterfaceKey = "list_with_tag"
