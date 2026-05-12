package simplestreams

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
)

const (
	cleartextHeader    = "-----BEGIN PGP SIGNED MESSAGE-----"
	cleartextSignature = "-----BEGIN PGP SIGNATURE-----"
	cleartextFooter    = "-----END PGP SIGNATURE-----"
)

// Verifier verifies cleartext signed metadata and returns its JSON payload.
type Verifier interface {
	VerifyCleartext(ctx context.Context, signed []byte) ([]byte, error)
}

// VerifierFunc adapts a function into a Verifier.
type VerifierFunc func(context.Context, []byte) ([]byte, error)

// VerifyCleartext calls f(ctx, signed).
func (f VerifierFunc) VerifyCleartext(ctx context.Context, signed []byte) ([]byte, error) {
	return f(ctx, signed)
}

func (mirror *Mirror) decodeSigned(ctx context.Context, signed []byte) ([]byte, error) {
	if mirror.signaturePolicy == SignatureRequired && mirror.verifier == nil {
		return nil, errors.New("simplestreams: signature verification required but no verifier configured")
	}
	if mirror.verifier != nil && mirror.signaturePolicy != SignatureDisabled {
		payload, err := mirror.verifier.VerifyCleartext(ctx, signed)
		if err != nil {
			return nil, fmt.Errorf("simplestreams: verify signed metadata: %w", err)
		}
		return payload, nil
	}
	return unwrapCleartextSigned(signed)
}

func unwrapCleartextSigned(signed []byte) ([]byte, error) {
	lines := strings.Split(string(signed), "\n")
	if len(lines) == 0 || strings.TrimRight(lines[0], "\r") != cleartextHeader {
		return nil, errors.New("simplestreams: signed metadata missing cleartext header")
	}

	mode := "headers"
	var body bytes.Buffer
	foundSignature := false
	for _, rawLine := range lines[1:] {
		line := strings.TrimRight(rawLine, "\r")
		switch mode {
		case "headers":
			if line == "" {
				mode = "body"
			}
			continue
		case "body":
			if line == cleartextSignature {
				foundSignature = true
				mode = "signature"
				continue
			}
			line = strings.TrimPrefix(line, "- ")
			body.WriteString(line)
			body.WriteByte('\n')
		case "signature":
			if line == cleartextFooter {
				return body.Bytes(), nil
			}
		}
	}
	if !foundSignature {
		return nil, errors.New("simplestreams: signed metadata missing signature block")
	}
	return nil, errors.New("simplestreams: signed metadata missing signature footer")
}
