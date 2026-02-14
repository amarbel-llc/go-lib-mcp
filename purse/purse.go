package purse

import "sort"

// Built-in tool names that purse-first can intercept.
const (
	BuiltinRead  = "Read"
	BuiltinEdit  = "Edit"
	BuiltinWrite = "Write"
	BuiltinGrep  = "Grep"
	BuiltinGlob  = "Glob"
	BuiltinBash  = "Bash"
)

// ToolSuggestion is an MCP tool that can replace a built-in tool.
type ToolSuggestion struct {
	Name    string `json:"name"`
	UseWhen string `json:"use_when"`
}

// Mapping is a single replacement rule declaring that an MCP server's tools
// should be used instead of a built-in tool.
type Mapping struct {
	Replaces   string           `json:"replaces"`
	Extensions []string         `json:"extensions,omitempty"`
	Tools      []ToolSuggestion `json:"tools"`
	Reason     string           `json:"reason"`
}

// MappingFile is the top-level structure written to disk and read by purse-first.
type MappingFile struct {
	Server   string    `json:"server"`
	Mappings []Mapping `json:"mappings"`
}

// MappingBuilder provides an ergonomic API for constructing a MappingFile.
type MappingBuilder struct {
	server   string
	mappings []*mappingEntry
}

type mappingEntry struct {
	replaces   string
	extensions []string
	tools      []ToolSuggestion
	reason     string
}

// NewMappingBuilder creates a builder for the given MCP server name.
func NewMappingBuilder(server string) *MappingBuilder {
	return &MappingBuilder{server: server}
}

// MappingEntryBuilder builds a single mapping within a MappingBuilder.
type MappingEntryBuilder struct {
	parent *MappingBuilder
	entry  *mappingEntry
}

// Replaces begins a new mapping that replaces the given built-in tool.
func (b *MappingBuilder) Replaces(builtinTool string) *MappingEntryBuilder {
	e := &mappingEntry{replaces: builtinTool}
	b.mappings = append(b.mappings, e)
	return &MappingEntryBuilder{parent: b, entry: e}
}

// ForExtensions limits the mapping to files with the given extensions.
func (eb *MappingEntryBuilder) ForExtensions(exts ...string) *MappingEntryBuilder {
	eb.entry.extensions = append(eb.entry.extensions, exts...)
	return eb
}

// WithTool adds an MCP tool suggestion to the current mapping.
func (eb *MappingEntryBuilder) WithTool(name, useWhen string) *MappingEntryBuilder {
	eb.entry.tools = append(eb.entry.tools, ToolSuggestion{
		Name:    name,
		UseWhen: useWhen,
	})
	return eb
}

// Because sets the human-readable reason for this mapping.
func (eb *MappingEntryBuilder) Because(reason string) *MappingEntryBuilder {
	eb.entry.reason = reason
	return eb
}

// Replaces begins a new mapping on the parent builder, allowing chained declarations.
func (eb *MappingEntryBuilder) Replaces(builtinTool string) *MappingEntryBuilder {
	return eb.parent.Replaces(builtinTool)
}

// Build produces the final MappingFile with mappings sorted by Replaces name.
func (b *MappingBuilder) Build() MappingFile {
	mappings := make([]Mapping, len(b.mappings))
	for i, e := range b.mappings {
		mappings[i] = Mapping{
			Replaces:   e.replaces,
			Extensions: e.extensions,
			Tools:      e.tools,
			Reason:     e.reason,
		}
	}

	sort.Slice(mappings, func(i, j int) bool {
		return mappings[i].Replaces < mappings[j].Replaces
	})

	return MappingFile{
		Server:   b.server,
		Mappings: mappings,
	}
}
