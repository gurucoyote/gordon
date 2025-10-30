# Fade Effects Notes

This note captures the baseline requirements for adding fade-in/out support to the
multitrack player.

## Why
- Smoothly ramping tracks avoids clicks when entering/exiting playback.
- Per-track fades would let us create intro/outro automation without editing the raw files.

## Implementation Sketch
1. Use `github.com/gopxl/beep/v2/effects.Transition` to modulate track gain.  
   - Fade-in: start gain 0.0 → end gain 1.0.  
   - Fade-out: start gain 1.0 → end gain 0.0.  
   - Duration is expressed in samples via `sampleRate.N(time.Duration)`.
2. Wrap each track’s streamer in a small helper that applies optional `Transition`
   before we feed it into the `CompositeSeeker`.  
   - For fade-in/out offsets, keep per-track metadata (e.g. seconds or samples).
3. Expose CLI flags such as `--fade-in 2s` / `--fade-out 1.5s` on `load` (per track) or
   on `loop`/`save` if we want segment-level fades.
4. Ensure fades interact cleanly with markers and `save`: exporting a loop should
   preserve fades if they’re active.

## Open Questions
- Do we need per-track curves (linear vs. exponential)?  
- Should fades be defaulted for all tracks, or opt-in per file?  
- Interaction with speed changes: when playback speed varies, transitions are currently
  tied to samples, so doubling the speed halves real-time fade length—might need to
  recompute durations when `speed` changes.

Keep this note handy when we implement the feature so the docs and CLI UX stay aligned.
