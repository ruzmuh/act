package runner

import (
	"context"

	"github.com/nektos/act/pkg/model"
)

// This file is the actl soft-fork patch. It is intentionally self-contained so
// the diff against upstream act stays tiny: the only edits to existing files are
// one field on Config (see runner.go) and the call site in newJobExecutor (see
// job_executor.go). Everything else lives here and never conflicts on rebase.
//
// The goal (per actl's design) is a pause hook at step boundaries. act keeps the
// job container alive and execs each step into it, so blocking *between* execs
// yields a live workspace + env — exactly what a step-debugger needs.

// StepBarrier is actl's pause hook (wired via Config.StepBarrier). It fires at
// every step boundary — before a step's main executor and after it returns — and
// the job pipeline blocks until it returns. Returning nil resumes; returning a
// non-nil error aborts the job. A debugger front-end typically blocks inside the
// hook until the user issues a resume, and decides whether to actually halt based
// on its own breakpoint policy.
//
// "Before step N" and "after step N-1" denote the same live-container moment; the
// distinct value of BarrierAfter is the final boundary (after the last step, just
// before teardown) and break-on-error (BarrierAfter carries the step's error).
//
// Upstream act leaves Config.StepBarrier nil and behaves exactly as before.
type StepBarrier func(ctx context.Context, info StepBarrierInfo) error

// BarrierWhen marks which side of a step's main executor the barrier fired on.
type BarrierWhen int

const (
	// BarrierBefore fires immediately before the step's main executor runs.
	BarrierBefore BarrierWhen = iota
	// BarrierAfter fires immediately after the step's main executor returns,
	// before any post steps or container teardown.
	BarrierAfter
)

func (w BarrierWhen) String() string {
	if w == BarrierAfter {
		return "after"
	}
	return "before"
}

// StepBarrierInfo describes the boundary at which the barrier fired, plus the
// live state a debugger needs to inspect while paused.
type StepBarrierInfo struct {
	When  BarrierWhen // before or after the step's main executor
	Index int         // zero-based position of the step within the job
	Step  *model.Step // the step model at this boundary
	Err   error       // for When==BarrierAfter: the step's error, or nil on success

	// Env is the job's live, interpolated environment at this boundary (the same
	// map act mutates, safe to read while the pipeline is blocked here).
	Env map[string]string
	// ContainerName is the docker name of the running job container, so a
	// frontend can drop an interactive shell into it (e.g. `docker exec -it`).
	// Empty when no job container is in use.
	ContainerName string
}
