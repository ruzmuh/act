package container

import "context"

// This file is part of the actl soft-fork patch. It is self-contained so the
// diff against upstream act stays tiny and rebase-safe: it adds a way to ask a
// CopyDir to leave git submodules out of the copy, threaded through the context
// rather than the CopyDir signature (which has several callers and a test mock).
//
// actl uses it to make a default local `actions/checkout` faithful — checkout
// defaults to `submodules: false`, so the workspace copy should skip submodule
// paths instead of recursing into them (which is also far cheaper for repos with
// large vendored submodules). Callers that don't set the key keep act's original
// recurse-into-submodules behaviour.

type skipSubmodulesKey struct{}

// WithSkipSubmodules returns a context that asks the next CopyDir to leave git
// submodule paths out of the copy.
func WithSkipSubmodules(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipSubmodulesKey{}, true)
}

// skipSubmodules reports whether the context asked CopyDir to skip submodules.
func skipSubmodules(ctx context.Context) bool {
	v, _ := ctx.Value(skipSubmodulesKey{}).(bool)
	return v
}
