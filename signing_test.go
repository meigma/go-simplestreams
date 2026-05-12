package simplestreams_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	simplestreams "github.com/meigma/go-simplestreams"
)

func TestSignedJSONRequiresVerifierWhenPolicyRequiresSignature(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.sjson": signedIndexJSON(),
		},
		opens: map[string]int{},
	}
	mirror, err := simplestreams.NewMirror(
		source,
		simplestreams.WithIndexPath("streams/v1/index.sjson"),
		simplestreams.WithSignaturePolicy(simplestreams.SignatureRequired),
	)
	require.NoError(t, err)

	_, err = mirror.Index(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no verifier")
}

func TestSignatureRequiredDefaultsToSignedIndex(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.json":  `not selected`,
			"streams/v1/index.sjson": signedIndexJSON(),
		},
		opens: map[string]int{},
	}
	verifier := simplestreams.VerifierFunc(func(context.Context, []byte) ([]byte, error) {
		return []byte(`{"format":"index:1.0","index":{}}`), nil
	})
	mirror, err := simplestreams.NewMirror(
		source,
		simplestreams.WithSignaturePolicy(simplestreams.SignatureRequired),
		simplestreams.WithVerifier(verifier),
	)
	require.NoError(t, err)

	index, err := mirror.Index(context.Background())
	require.NoError(t, err)
	assert.Empty(t, index.Entries)
	assert.Equal(t, 1, source.opens["streams/v1/index.sjson"])
	assert.Zero(t, source.opens["streams/v1/index.json"])
}

func TestSignatureRequiredRejectsExplicitUnsignedJSON(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.json": `{"format":"index:1.0","index":{}}`,
		},
		opens: map[string]int{},
	}
	mirror, err := simplestreams.NewMirror(
		source,
		simplestreams.WithIndexPath("streams/v1/index.json"),
		simplestreams.WithSignaturePolicy(simplestreams.SignatureRequired),
	)
	require.NoError(t, err)

	_, err = mirror.Index(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "signature required")
	assert.Zero(t, source.opens["streams/v1/index.json"])
}

func TestSignaturePreferredDefaultsToSignedIndex(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.json":  `not selected`,
			"streams/v1/index.sjson": signedIndexJSON(),
		},
		opens: map[string]int{},
	}
	mirror, err := simplestreams.NewMirror(
		source,
		simplestreams.WithSignaturePolicy(simplestreams.SignaturePreferred),
	)
	require.NoError(t, err)

	index, err := mirror.Index(context.Background())
	require.NoError(t, err)
	assert.Empty(t, index.Entries)
	assert.Equal(t, 1, source.opens["streams/v1/index.sjson"])
	assert.Zero(t, source.opens["streams/v1/index.json"])
}

func TestSignaturePreferredFallsBackWhenSignedIndexIsMissing(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.json": `{"format":"index:1.0","index":{}}`,
		},
		opens: map[string]int{},
	}
	mirror, err := simplestreams.NewMirror(
		source,
		simplestreams.WithSignaturePolicy(simplestreams.SignaturePreferred),
	)
	require.NoError(t, err)

	index, err := mirror.Index(context.Background())
	require.NoError(t, err)
	assert.Empty(t, index.Entries)
	assert.Equal(t, 1, source.opens["streams/v1/index.sjson"])
	assert.Equal(t, 1, source.opens["streams/v1/index.json"])
}

func TestSignaturePreferredDoesNotFallbackOnSignedIndexParseError(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.json":  `{"format":"index:1.0","index":{}}`,
			"streams/v1/index.sjson": signedInvalidJSON(),
		},
		opens: map[string]int{},
	}
	mirror, err := simplestreams.NewMirror(
		source,
		simplestreams.WithSignaturePolicy(simplestreams.SignaturePreferred),
	)
	require.NoError(t, err)

	_, err = mirror.Index(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode index")
	assert.Equal(t, 1, source.opens["streams/v1/index.sjson"])
	assert.Zero(t, source.opens["streams/v1/index.json"])
}

func TestSignaturePreferredDoesNotFallbackOnSignedIndexOpenError(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.json": `{"format":"index:1.0","index":{}}`,
		},
		errs: map[string]error{
			"streams/v1/index.sjson": errors.New("temporary failure"),
		},
		opens: map[string]int{},
	}
	mirror, err := simplestreams.NewMirror(
		source,
		simplestreams.WithSignaturePolicy(simplestreams.SignaturePreferred),
	)
	require.NoError(t, err)

	_, err = mirror.Index(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "temporary failure")
	assert.Equal(t, 1, source.opens["streams/v1/index.sjson"])
	assert.Zero(t, source.opens["streams/v1/index.json"])
}

func TestSignedJSONUsesVerifier(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.sjson": signedIndexJSON(),
		},
		opens: map[string]int{},
	}
	verified := false
	verifier := simplestreams.VerifierFunc(func(_ context.Context, signed []byte) ([]byte, error) {
		verified = true
		assert.Contains(t, string(signed), "BEGIN PGP SIGNED MESSAGE")
		return []byte(`{"format":"index:1.0","index":{}}`), nil
	})
	mirror, err := simplestreams.NewMirror(
		source,
		simplestreams.WithIndexPath("streams/v1/index.sjson"),
		simplestreams.WithSignaturePolicy(simplestreams.SignatureRequired),
		simplestreams.WithVerifier(verifier),
	)
	require.NoError(t, err)

	index, err := mirror.Index(context.Background())
	require.NoError(t, err)
	assert.Empty(t, index.Entries)
	assert.True(t, verified)
}

func TestSignedJSONVerifierFailureIsFatalWhenPreferred(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.json":  `{"format":"index:1.0","index":{}}`,
			"streams/v1/index.sjson": signedIndexJSON(),
		},
		opens: map[string]int{},
	}
	verifier := simplestreams.VerifierFunc(func(context.Context, []byte) ([]byte, error) {
		return nil, errors.New("bad signature")
	})
	mirror, err := simplestreams.NewMirror(
		source,
		simplestreams.WithSignaturePolicy(simplestreams.SignaturePreferred),
		simplestreams.WithVerifier(verifier),
	)
	require.NoError(t, err)

	_, err = mirror.Index(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "verify signed metadata")
	assert.Equal(t, 1, source.opens["streams/v1/index.sjson"])
	assert.Zero(t, source.opens["streams/v1/index.json"])
}

func TestSignedJSONVerifierFailureIsFatalWhenRequired(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.sjson": signedIndexJSON(),
		},
		opens: map[string]int{},
	}
	verifier := simplestreams.VerifierFunc(func(context.Context, []byte) ([]byte, error) {
		return nil, errors.New("bad signature")
	})
	mirror, err := simplestreams.NewMirror(
		source,
		simplestreams.WithIndexPath("streams/v1/index.sjson"),
		simplestreams.WithSignaturePolicy(simplestreams.SignatureRequired),
		simplestreams.WithVerifier(verifier),
	)
	require.NoError(t, err)

	_, err = mirror.Index(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "verify signed metadata")
}

func TestSignedJSONCanBeUnwrappedWithoutVerificationWhenDisabled(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.sjson": signedIndexJSON(),
		},
		opens: map[string]int{},
	}
	mirror, err := simplestreams.NewMirror(source, simplestreams.WithIndexPath("streams/v1/index.sjson"))
	require.NoError(t, err)

	index, err := mirror.Index(context.Background())
	require.NoError(t, err)
	assert.Empty(t, index.Entries)
}

func TestVerifierFunc(t *testing.T) {
	expected := []byte(`{"format":"index:1.0","index":{}}`)
	verifier := simplestreams.VerifierFunc(func(_ context.Context, _ []byte) ([]byte, error) {
		return expected, nil
	})

	got, err := verifier.VerifyCleartext(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestSourceFunc(t *testing.T) {
	source := simplestreams.SourceFunc(func(_ context.Context, path simplestreams.RelativePath) (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(path.String())), nil
	})

	reader, err := source.Open(context.Background(), "streams/v1/index.json")
	require.NoError(t, err)
	defer reader.Close()
	body, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, "streams/v1/index.json", string(body))
}

func signedIndexJSON() string {
	return "-----BEGIN PGP SIGNED MESSAGE-----\nHash: SHA256\n\n{\"format\":\"index:1.0\",\"index\":{}}\n-----BEGIN PGP SIGNATURE-----\nplaceholder\n-----END PGP SIGNATURE-----\n"
}

func signedInvalidJSON() string {
	return "-----BEGIN PGP SIGNED MESSAGE-----\nHash: SHA256\n\nnot-json\n-----BEGIN PGP SIGNATURE-----\nplaceholder\n-----END PGP SIGNATURE-----\n"
}
