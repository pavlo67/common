package storage

import (
	"github.com/pavlo67/workshop/common/joiner"
)

const CollectionDefault = "storage"

const DataInterfaceKey joiner.InterfaceKey = "storage_data"
const InterfaceKey joiner.InterfaceKey = "storage"

const RecentInterfaceKey joiner.InterfaceKey = "storage_list"
const ReadInterfaceKey joiner.InterfaceKey = "storage_read"
const SaveInterfaceKey joiner.InterfaceKey = "storage_save"
const RemoveInterfaceKey joiner.InterfaceKey = "storage_remove"
const ListTagsInterfaceKey joiner.InterfaceKey = "storage_count_tags"
const ListTaggedInterfaceKey joiner.InterfaceKey = "storage_list_with_tag"

const ExportInterfaceKey joiner.InterfaceKey = "storage_export"
const ImportInterfaceKey joiner.InterfaceKey = "storage_import"
