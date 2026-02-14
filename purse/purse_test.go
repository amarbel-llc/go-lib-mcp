package purse

import (
	"encoding/json"
	"testing"
)

func TestMappingFileJSON(t *testing.T) {
	mf := MappingFile{
		Server: "test-server",
		Mappings: []Mapping{
			{
				Replaces:   BuiltinRead,
				Extensions: []string{".go", ".py"},
				Tools: []ToolSuggestion{
					{Name: "lsp_hover", UseWhen: "getting type info"},
				},
				Reason: "Use LSP tools for reading",
			},
		},
	}

	data, err := json.MarshalIndent(mf, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got MappingFile
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Server != mf.Server {
		t.Errorf("server = %q, want %q", got.Server, mf.Server)
	}
	if len(got.Mappings) != 1 {
		t.Fatalf("mappings len = %d, want 1", len(got.Mappings))
	}
	if got.Mappings[0].Replaces != BuiltinRead {
		t.Errorf("replaces = %q, want %q", got.Mappings[0].Replaces, BuiltinRead)
	}
	if len(got.Mappings[0].Extensions) != 2 {
		t.Errorf("extensions len = %d, want 2", len(got.Mappings[0].Extensions))
	}
}

func TestExtensionsOmitEmpty(t *testing.T) {
	mf := MappingFile{
		Server: "test-server",
		Mappings: []Mapping{
			{
				Replaces: BuiltinBash,
				Tools: []ToolSuggestion{
					{Name: "run_command", UseWhen: "running shell commands"},
				},
				Reason: "Use server's run_command",
			},
		},
	}

	data, err := json.Marshal(mf)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	json.Unmarshal(data, &raw)

	var mappings []map[string]json.RawMessage
	json.Unmarshal(raw["mappings"], &mappings)

	if _, ok := mappings[0]["extensions"]; ok {
		t.Error("extensions should be omitted when nil")
	}
}

func TestBuilderSortsByReplaces(t *testing.T) {
	builder := NewMappingBuilder("test-server")
	builder.Replaces(BuiltinWrite).
		WithTool("write_file", "writing files").
		Because("Use server's writer")
	builder.Replaces(BuiltinRead).
		WithTool("read_file", "reading files").
		Because("Use server's reader")
	builder.Replaces(BuiltinGrep).
		WithTool("search", "searching code").
		Because("Use server's search")

	mf := builder.Build()

	if len(mf.Mappings) != 3 {
		t.Fatalf("mappings len = %d, want 3", len(mf.Mappings))
	}

	want := []string{BuiltinGrep, BuiltinRead, BuiltinWrite}
	for i, w := range want {
		if mf.Mappings[i].Replaces != w {
			t.Errorf("mappings[%d].replaces = %q, want %q", i, mf.Mappings[i].Replaces, w)
		}
	}
}

func TestBuilderChaining(t *testing.T) {
	builder := NewMappingBuilder("lux")
	builder.Replaces(BuiltinRead).
		ForExtensions(".go", ".py").
		WithTool("lsp_hover", "getting type info").
		WithTool("lsp_definition", "finding definitions").
		Because("Use LSP tools instead of reading raw files").
		Replaces(BuiltinGrep).
		ForExtensions(".go").
		WithTool("lsp_references", "finding usages").
		Because("Use LSP references instead of grep")

	mf := builder.Build()

	if mf.Server != "lux" {
		t.Errorf("server = %q, want %q", mf.Server, "lux")
	}
	if len(mf.Mappings) != 2 {
		t.Fatalf("mappings len = %d, want 2", len(mf.Mappings))
	}

	// Sorted: Grep before Read
	if mf.Mappings[0].Replaces != BuiltinGrep {
		t.Errorf("mappings[0].replaces = %q, want %q", mf.Mappings[0].Replaces, BuiltinGrep)
	}
	if len(mf.Mappings[0].Tools) != 1 {
		t.Errorf("mappings[0].tools len = %d, want 1", len(mf.Mappings[0].Tools))
	}
	if mf.Mappings[1].Replaces != BuiltinRead {
		t.Errorf("mappings[1].replaces = %q, want %q", mf.Mappings[1].Replaces, BuiltinRead)
	}
	if len(mf.Mappings[1].Tools) != 2 {
		t.Errorf("mappings[1].tools len = %d, want 2", len(mf.Mappings[1].Tools))
	}
	if len(mf.Mappings[1].Extensions) != 2 {
		t.Errorf("mappings[1].extensions len = %d, want 2", len(mf.Mappings[1].Extensions))
	}
}

func TestBuilderWireFormat(t *testing.T) {
	builder := NewMappingBuilder("my-server")
	builder.Replaces(BuiltinRead).
		ForExtensions(".go").
		WithTool("lsp_hover", "getting type info").
		Because("Use my-server's LSP tools")

	mf := builder.Build()
	data, err := json.MarshalIndent(mf, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// Re-parse as generic JSON to verify wire format
	var wire map[string]interface{}
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatalf("unmarshal wire: %v", err)
	}

	if wire["server"] != "my-server" {
		t.Errorf("wire server = %v", wire["server"])
	}

	mappings := wire["mappings"].([]interface{})
	m := mappings[0].(map[string]interface{})
	if m["replaces"] != "Read" {
		t.Errorf("wire replaces = %v", m["replaces"])
	}
	if m["reason"] != "Use my-server's LSP tools" {
		t.Errorf("wire reason = %v", m["reason"])
	}

	tools := m["tools"].([]interface{})
	tool := tools[0].(map[string]interface{})
	if tool["name"] != "lsp_hover" {
		t.Errorf("wire tool name = %v", tool["name"])
	}
	if tool["use_when"] != "getting type info" {
		t.Errorf("wire tool use_when = %v", tool["use_when"])
	}
}
