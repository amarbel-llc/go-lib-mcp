package purse

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteProject(t *testing.T) {
	dir := t.TempDir()

	mf := MappingFile{
		Server: "test-server",
		Mappings: []Mapping{
			{
				Replaces: BuiltinRead,
				Tools: []ToolSuggestion{
					{Name: "read_file", UseWhen: "reading files"},
				},
				Reason: "Use server's reader",
			},
		},
	}

	if err := WriteProject(dir, mf); err != nil {
		t.Fatalf("WriteProject: %v", err)
	}

	path := filepath.Join(dir, ".purse-first", "test-server.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	// Verify trailing newline
	if data[len(data)-1] != '\n' {
		t.Error("expected trailing newline")
	}

	var got MappingFile
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Server != "test-server" {
		t.Errorf("server = %q, want %q", got.Server, "test-server")
	}
	if len(got.Mappings) != 1 {
		t.Fatalf("mappings len = %d, want 1", len(got.Mappings))
	}
}

func TestWriteGlobal(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)

	mf := MappingFile{
		Server: "global-server",
		Mappings: []Mapping{
			{
				Replaces:   BuiltinGrep,
				Extensions: []string{".go"},
				Tools: []ToolSuggestion{
					{Name: "lsp_references", UseWhen: "finding usages"},
				},
				Reason: "Use LSP references",
			},
		},
	}

	if err := WriteGlobal(mf); err != nil {
		t.Fatalf("WriteGlobal: %v", err)
	}

	path := filepath.Join(dir, "purse-first", "global-server.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var got MappingFile
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Server != "global-server" {
		t.Errorf("server = %q, want %q", got.Server, "global-server")
	}
}
