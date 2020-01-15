package flow

import (
	"github.com/pavlo67/workshop/common/joiner"
)

const DataInterfaceKey joiner.InterfaceKey = "flow_data"

const InterfaceKey joiner.InterfaceKey = "flow"
const CleanerInterfaceKey joiner.InterfaceKey = "flow_cleaner"

const ImporterTaskInterfaceKey joiner.InterfaceKey = "flow_importer_task"
const CleanerTaskInterfaceKey joiner.InterfaceKey = "flow_cleaner_task"
const CopierTaskInterfaceKey joiner.InterfaceKey = "flow_copier_task"

const ListInterfaceKey joiner.InterfaceKey = "flow_list"
const ReadInterfaceKey joiner.InterfaceKey = "flow_read"
const ExportInterfaceKey joiner.InterfaceKey = "flow_export"

const CollectionDefault = "flow"
const FlowLimitDefault = 3000
