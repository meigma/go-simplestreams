package simplestreams

import "fmt"

// SignaturePolicy controls how signed Simple Streams metadata is handled.
type SignaturePolicy int

const (
	// SignatureDisabled accepts unsigned JSON and unwraps .sjson without verification.
	SignatureDisabled SignaturePolicy = iota

	// SignaturePreferred verifies .sjson metadata when a verifier is configured.
	SignaturePreferred

	// SignatureRequired requires .sjson metadata to be verified by a configured verifier.
	SignatureRequired
)

// Option configures a Mirror.
type Option func(*mirrorOptions) error

type mirrorOptions struct {
	indexPath              RelativePath
	indexPathExplicit      bool
	allowContentIDMismatch bool
	signaturePolicy        SignaturePolicy
	verifier               Verifier
}

func defaultMirrorOptions() mirrorOptions {
	return mirrorOptions{
		indexPath:       DefaultIndexPath,
		signaturePolicy: SignatureDisabled,
	}
}

// WithIndexPath sets the metadata path used by Mirror.Index.
func WithIndexPath(path string) Option {
	return func(options *mirrorOptions) error {
		relativePath, err := ParseRelativePath(path)
		if err != nil {
			return err
		}
		if !relativePath.IsMetadataPath() {
			return fmt.Errorf("simplestreams: index path %q is not .json or .sjson", path)
		}
		options.indexPath = relativePath
		options.indexPathExplicit = true
		return nil
	}
}

// WithAllowContentIDMismatch controls whether product content_id mismatches are accepted.
func WithAllowContentIDMismatch(allow bool) Option {
	return func(options *mirrorOptions) error {
		options.allowContentIDMismatch = allow
		return nil
	}
}

// WithSignaturePolicy sets the metadata signature handling policy.
func WithSignaturePolicy(policy SignaturePolicy) Option {
	return func(options *mirrorOptions) error {
		switch policy {
		case SignatureDisabled, SignaturePreferred, SignatureRequired:
			options.signaturePolicy = policy
			return nil
		default:
			return fmt.Errorf("simplestreams: unknown signature policy %d", policy)
		}
	}
}

// WithVerifier sets the verifier used for cleartext signed metadata.
func WithVerifier(verifier Verifier) Option {
	return func(options *mirrorOptions) error {
		options.verifier = verifier
		return nil
	}
}
