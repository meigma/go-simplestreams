package incus

import (
	"errors"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"

	simplestreams "github.com/meigma/go-simplestreams"
	"github.com/meigma/go-simplestreams/schema"
)

// ValidateRuntimeProductFile validates productFile against the embedded Incus product-file schema.
func ValidateRuntimeProductFile(productFile *simplestreams.ProductFile) error {
	if productFile == nil {
		return errors.New("incus: nil runtime product file")
	}

	data, err := simplestreams.MarshalJSONDocument(productFile)
	if err != nil {
		return fmt.Errorf("incus: marshal runtime product file: %w", err)
	}

	ctx := cuecontext.New()
	schemaValue, err := schema.IncusProductFileSchema(ctx)
	if err != nil {
		return fmt.Errorf("incus: load product-file schema: %w", err)
	}
	if _, err := schema.ValidateBytes(ctx, schemaValue, data, cue.Filename("runtime-product-file.json")); err != nil {
		return fmt.Errorf("incus: validate runtime product file: %w", err)
	}
	return nil
}
