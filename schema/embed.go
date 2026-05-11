package schema

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
)

const (
	// ModulePath is the embedded CUE module path.
	ModulePath = "github.com/meigma/go-simplestreams/schema@v0"

	embeddedModuleDir = "go-simplestreams-schema-embedded"
)

//go:embed cue.mod/module.cue *.cue incus/*.cue linuxcontainers/*.cue lxd/*.cue
var embeddedModule embed.FS

// ModuleFS returns the embedded CUE module files.
func ModuleFS() fs.FS {
	return embeddedModule
}

// LoadSchema loads the embedded root Simple Streams schema package.
func LoadSchema(ctx *cue.Context) (cue.Value, error) {
	return LoadPackage(ctx, ".")
}

// LoadPackage loads one package from the embedded CUE module.
//
// packagePath is a CUE loader path relative to the embedded module root, such as
// ".", "./linuxcontainers", "./incus", or "./lxd".
func LoadPackage(ctx *cue.Context, packagePath string) (cue.Value, error) {
	if ctx == nil {
		return cue.Value{}, errors.New("schema: nil CUE context")
	}
	if packagePath == "" {
		return cue.Value{}, errors.New("schema: empty CUE package path")
	}

	overlay, err := embeddedOverlay()
	if err != nil {
		return cue.Value{}, err
	}

	moduleRoot := embeddedModuleRoot()
	insts := load.Instances([]string{packagePath}, &load.Config{
		Dir:        moduleRoot,
		ModuleRoot: moduleRoot,
		Module:     ModulePath,
		Overlay:    overlay,
		Env:        []string{},
	})
	if len(insts) != 1 {
		return cue.Value{}, fmt.Errorf("schema: expected one CUE instance for %q, got %d", packagePath, len(insts))
	}
	if err := insts[0].Err; err != nil {
		return cue.Value{}, fmt.Errorf("schema: load embedded CUE package %q: %w", packagePath, err)
	}

	value := ctx.BuildInstance(insts[0])
	if err := value.Err(); err != nil {
		return cue.Value{}, fmt.Errorf("schema: build embedded CUE package %q: %w", packagePath, err)
	}
	return value, nil
}

// IndexFileSchema returns the embedded root #IndexFile definition.
func IndexFileSchema(ctx *cue.Context) (cue.Value, error) {
	return schemaDefinition(ctx, ".", "#IndexFile")
}

// ProductFileSchema returns the embedded root #ProductFile definition.
func ProductFileSchema(ctx *cue.Context) (cue.Value, error) {
	return schemaDefinition(ctx, ".", "#ProductFile")
}

// LinuxContainersProductFileSchema returns the embedded linuxcontainers.#ProductFile definition.
func LinuxContainersProductFileSchema(ctx *cue.Context) (cue.Value, error) {
	return schemaDefinition(ctx, "./linuxcontainers", "#ProductFile")
}

// IncusProductFileSchema returns the embedded incus.#ProductFile definition.
func IncusProductFileSchema(ctx *cue.Context) (cue.Value, error) {
	return schemaDefinition(ctx, "./incus", "#ProductFile")
}

// LXDProductFileSchema returns the embedded lxd.#ProductFile definition.
func LXDProductFileSchema(ctx *cue.Context) (cue.Value, error) {
	return schemaDefinition(ctx, "./lxd", "#ProductFile")
}

// ValidateBytes validates CUE or JSON bytes against schemaValue.
//
// ValidateBytes returns the unified value so callers can decode defaulted or
// normalized data from the same CUE context.
func ValidateBytes(
	ctx *cue.Context,
	schemaValue cue.Value,
	data []byte,
	options ...cue.BuildOption,
) (cue.Value, error) {
	if ctx == nil {
		return cue.Value{}, errors.New("schema: nil CUE context")
	}
	if err := schemaValue.Err(); err != nil {
		return cue.Value{}, fmt.Errorf("schema: invalid schema value: %w", err)
	}

	input := ctx.CompileBytes(data, options...)
	if err := input.Err(); err != nil {
		return cue.Value{}, fmt.Errorf("schema: parse input: %w", err)
	}

	value := schemaValue.Unify(input)
	if err := value.Validate(cue.Concrete(true)); err != nil {
		return cue.Value{}, fmt.Errorf("schema: validate input: %w", err)
	}
	return value, nil
}

func schemaDefinition(ctx *cue.Context, packagePath string, definition string) (cue.Value, error) {
	value, err := LoadPackage(ctx, packagePath)
	if err != nil {
		return cue.Value{}, err
	}

	schemaValue := value.LookupPath(cue.ParsePath(definition))
	if err := schemaValue.Err(); err != nil {
		return cue.Value{}, fmt.Errorf("schema: lookup %s in %q: %w", definition, packagePath, err)
	}
	return schemaValue, nil
}

func embeddedOverlay() (map[string]load.Source, error) {
	moduleRoot := embeddedModuleRoot()
	overlay := map[string]load.Source{}

	err := fs.WalkDir(embeddedModule, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".cue") {
			return nil
		}

		data, err := fs.ReadFile(embeddedModule, path)
		if err != nil {
			return err
		}
		overlay[filepath.Join(moduleRoot, filepath.FromSlash(path))] = load.FromBytes(data)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("schema: read embedded CUE module: %w", err)
	}
	return overlay, nil
}

func embeddedModuleRoot() string {
	return filepath.Join(os.TempDir(), embeddedModuleDir)
}
