# MIDI Playback Integration

This document captures how the Gordon CLI loads and plays MIDI assets after the
beep v2 upgrade.

## Dependencies & Assets

- Runtime decoding relies on `github.com/gopxl/beep/v2/midi` and a General MIDI
  compatible SoundFont (`.sf2` or `.sf3`).  
- Provide the SoundFont location at launch via the `--soundfont` persistent
  flag, e.g.:

```bash
go run . --soundfont ./assets/Florestan.sf2 play 0 intro.wav 8 intro.mid
```

The application keeps the SoundFont in memory using a `sync.Once` guarded
initializer. Restart the process to switch fonts.

## Loader Flow

1. `cmd/root.go` exposes the `--soundfont` option so every subcommand can reuse
   it.  
2. `cmd/play.go` augments the existing `load` multitrack logic:
- Files ending with `.mid`/`.midi` are decoded via
     `midi.Decode(file, soundFont, defaultSampleRate)`.
   - The decoded stream is immediately materialized into a `beep.Buffer` so it
     becomes a stable `StreamSeeker`; this avoids crashes inside the underlying
     synthesizer when users seek repeatedly.
   - The resulting buffer-backed streamer is wrapped in the same offset-aware
     `CompositeSeeker` that WAV/MP3 inputs use, so markers, looping, and saving
     continue to work.
   - The resampler still normalizes everything to the speaker’s sample rate, so
     mixed audio + MIDI tracks stay tempo aligned.

## Behaviour Notes

- MIDI playback fails fast when no soundfont is supplied; the loader prints a
  descriptive error so the user can re-run with `--soundfont`.  
- SoundFont handles are cached for subsequent tracks, but the original file is
  closed right after decoding to avoid leaking file descriptors.  
- `speed`, `loop`, `drop`, and `save` commands operate identically on MIDI and
  audio tracks because they work with the common `beep.StreamSeeker` interface.

## Follow-up Ideas

- Persist a default soundfont path in config so repeated runs do not require the
  flag.  
- Wrap `StreamSeekCloser` instances so removing a track closes the underlying
  reader, keeping long sessions from exhausting file handles.  
- Add smoke tests that load a tiny MIDI clip plus a WAV file into the
  multitrack player to guard against regressions.
