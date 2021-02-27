package extractor

type SchemaElement struct {
	// One of "schema.ValueType" or "SchemaElements" or "SchemaInfo"
	Type string `json:",omitempty"`
	// Set for simple types (from ValueType)
	Value string `json:",omitempty"`
	// Set if Type == "SchemaElements"
	ElementsType string `json:",omitempty"`
	// Set if Type == "SchemaInfo"
	Info SchemaInfo `json:",omitempty"`
}

type SchemaDefinition struct {
	Type               string `json:",omitempty"`
	Optional           bool   `json:",omitempty"`
	Required           bool   `json:",omitempty"`
	Description        string `json:",omitempty"`
	InputDefault       string `json:",omitempty"`
	Computed           bool   `json:",omitempty"`
	MaxItems           int    `json:",omitempty"`
	MinItems           int    `json:",omitempty"`
	PromoteSingle      bool   `json:",omitempty"`
	IsBlock            bool   `json:",omitempty"`
	ConfigImplicitMode string `json:",omitempty"`

	ComputedWhen  []string `json:",omitempty"`
	ConflictsWith []string `json:",omitempty"`

	Deprecated string `json:",omitempty"`
	Removed    string `json:",omitempty"`

	Default *SchemaElement `json:",omitempty"`
	Elem    *SchemaElement `json:",omitempty"`
}

type SchemaInfo map[string]SchemaDefinition
type SchemaInfoWithTimeouts map[string]interface{}

//{
//	SchemaInfo `json:""`
//	Timeouts []string `json:"__timeouts__,omitempty"`
//}

// ResourceProviderSchema
type ResourceProviderSchema struct {
	Name          string                            `json:"name"`
	Type          string                            `json:"type"`
	Version       string                            `json:"version"`
	SDKType       string                            `json:".sdk_type"`
	SchemaVersion string                            `json:".schema_version"`
	Provider      SchemaInfo                        `json:"provider"`
	Resources     map[string]SchemaInfoWithTimeouts `json:"resources"`
	DataSources   map[string]SchemaInfoWithTimeouts `json:"data-sources"`
}

type ProviderInfo struct {
	Name     string
	Revision string
}

const (
	TimeoutCreate  = "create"
	TimeoutRead    = "read"
	TimeoutUpdate  = "update"
	TimeoutDelete  = "delete"
	TimeoutDefault = "default"
)

func timeoutKeys() []string {
	return []string{
		TimeoutCreate,
		TimeoutRead,
		TimeoutUpdate,
		TimeoutDelete,
		TimeoutDefault,
	}
}

func shortenType(value string) string {
	if len(value) > 4 && value[0:4] == "Type" {
		return value[4:]
	}
	return value
}
