package runner

import (
	"context"

	"github.com/nektos/act/pkg/model"
)

// This file is the actl soft-fork patch. It is intentionally self-contained so
// the diff against upstream act stays tiny: the only edits to existing files are
// one field on Config (see runner.go) and one call site in newJobExecutor (see
// job_executor.go). Everything else lives here and never conflicts on rebase.
//
// The goal (per actl's design) is a pause hook between steps. act keeps the job
// container alive and execs each step into it, so blocking *between* execs yields
// a live workspace + env — exactly what a step-debugger needs.

// StepBarrier is actl's pause hook (wired via Config.StepBarrier). It is invoked
// once per step, immediately before that step's main executor runs, and the job
// pipeline blocks until it returns. Returning nil resumes the step; returning a
// non-nil error aborts the job. A debugger front-end typically blocks inside the
// hook until the user issues a resume.
//
// Upstream act leaves Config.StepBarrier nil and behaves exactly as before.
type StepBarrier func(ctx context.Context, info StepBarrierInfo) error

// StepBarrierInfo describes the step that is about to run when the barrier fires.
type StepBarrierInfo struct {
	Index int         // zero-based position of the step within the job
	Step  *model.Step // the step model about to execute
}
