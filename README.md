# nothp

`nothp` is a tiny Go module that disables Transparent Huge Pages (THP) for
the current process.

## Why this exists

Linux's THP feature transparently backs anonymous memory with 2MB pages
instead of regular 4KB pages, which can help large, memory-bound
workloads by reducing TLB pressure. Most distributions ship it in
`madvise` mode: huge pages are only used for memory regions that are
explicitly marked with `madvise(MADV_HUGEPAGE)`, so it stays opt-in
system-wide.

The catch is that many language runtimes, including Go's, mark their
heap this way unconditionally, regardless of how much memory the program
actually needs. For small, non-heap-heavy processes this rounds up
allocations to 2MB boundaries, inflating RSS by ~15-20MB per process for
essentially no benefit; on a host running many small services, that adds
up. You may want the system-wide `madvise` policy to stay in place (other,
larger processes might benefit from it) while opting *your* process out of
it individually, without needing root or a kernel-wide config change.

That per-process opt-out is what `prctl(2)`'s `PR_SET_THP_DISABLE` gives
you, but it only affects memory mapped *after* the call. By the time your
own `main()` runs, the runtime has already mapped and madvise'd its
initial heap arenas, so calling it there is too late. `PR_SET_THP_DISABLE`
is documented as preserved across `execve(2)` though, so `nothp` disables
it and then re-execs the running binary in place (once), so the setting
is active from the very first byte of the new process image, without
running two copies of your program's runtime/heap at once.

If disabling THP or re-exec'ing fails for any reason (unsupported kernel,
permissions, sandboxed environment), it logs a warning and continues, so
it never blocks your program from starting. On non-Linux platforms,
`Relaunch` is a no-op.

## Usage

```go
import "nothp"

func main() {
    nothp.Relaunch() // must be the first line of main
    // ... rest of your program
}
```

`Relaunch` re-execs the binary at most once (guarded by an environment
variable sentinel), so it's safe to call unconditionally on every run.

## Testing

```bash
go test ./...
```
