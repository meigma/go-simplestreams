package schema

import (
	"bytes"
	"io/fs"
	"os"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

func TestModuleFS(t *testing.T) {
	moduleFile, err := fs.ReadFile(ModuleFS(), "cue.mod/module.cue")
	if err != nil {
		t.Fatalf("read module file: %v", err)
	}
	if !bytes.Contains(moduleFile, []byte(ModulePath)) {
		t.Fatalf("module file does not contain module path %q", ModulePath)
	}

	for _, path := range []string{
		"schema.cue",
		"linuxcontainers/schema.cue",
		"incus/schema.cue",
		"lxd/schema.cue",
	} {
		if _, err := fs.Stat(ModuleFS(), path); err != nil {
			t.Fatalf("embedded module missing %q: %v", path, err)
		}
	}
}

func TestSchemaHelpers(t *testing.T) {
	tests := []struct {
		name   string
		schema func(*cue.Context) (cue.Value, error)
	}{
		{name: "index", schema: IndexFileSchema},
		{name: "product", schema: ProductFileSchema},
		{name: "linuxcontainers", schema: LinuxContainersProductFileSchema},
		{name: "incus", schema: IncusProductFileSchema},
		{name: "lxd", schema: LXDProductFileSchema},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := cuecontext.New()
			schemaValue, err := tt.schema(ctx)
			if err != nil {
				t.Fatalf("load schema: %v", err)
			}
			if got := schemaValue.IncompleteKind(); got != cue.StructKind {
				t.Fatalf("schema kind = %v, want %v", got, cue.StructKind)
			}
		})
	}
}

func TestValidateBytes(t *testing.T) {
	tests := []struct {
		name        string
		schema      func(*cue.Context) (cue.Value, error)
		fixturePath string
	}{
		{
			name:        "index",
			schema:      IndexFileSchema,
			fixturePath: "testdata/minimal/index.json",
		},
		{
			name:        "product",
			schema:      ProductFileSchema,
			fixturePath: "testdata/minimal/products.json",
		},
		{
			name:        "linuxcontainers",
			schema:      LinuxContainersProductFileSchema,
			fixturePath: "linuxcontainers/testdata/products.json",
		},
		{
			name:        "incus",
			schema:      IncusProductFileSchema,
			fixturePath: "incus/testdata/products.json",
		},
		{
			name:        "lxd",
			schema:      LXDProductFileSchema,
			fixturePath: "lxd/testdata/products.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := cuecontext.New()
			schemaValue, err := tt.schema(ctx)
			if err != nil {
				t.Fatalf("load schema: %v", err)
			}

			data, err := os.ReadFile(tt.fixturePath)
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}

			value, err := ValidateBytes(ctx, schemaValue, data, cue.Filename(tt.fixturePath))
			if err != nil {
				t.Fatalf("validate fixture: %v", err)
			}
			if err := value.Validate(cue.Concrete(false)); err != nil {
				t.Fatalf("validated value became invalid: %v", err)
			}
		})
	}
}

func TestValidateBytesRejectsInvalidInput(t *testing.T) {
	ctx := cuecontext.New()
	schemaValue, err := ProductFileSchema(ctx)
	if err != nil {
		t.Fatalf("load product schema: %v", err)
	}

	_, err = ValidateBytes(ctx, schemaValue, []byte(`{"format":"products:1.0"}`), cue.Filename("invalid.json"))
	if err == nil {
		t.Fatal("ValidateBytes accepted product file without content_id or products")
	}
}
