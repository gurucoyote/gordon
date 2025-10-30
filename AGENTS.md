# Repository Guidelines

## Project Structure & Module Organization
- `main.go` wires the CLI entrypoint and orchestrates command flow via `cmd/`.
- `cmd/` hosts Cobra commands and playback logic; keep feature helpers beside the command that consumes them.
- `examples/` and `streaming/` hold experiment binaries plus reference assets—prototype there before touching production code.
- Audio fixtures (`*.mp3`, `*.wav`, `midi/`) support manual verification; avoid committing large new media without prior agreement.

## Build, Test, and Development Commands
- `go run . play test.wav` launches the player against the bundled sample file for a quick smoke test.
- `go build -o gordon` produces a native binary for your platform; run it before every pull request.
- `./build.sh` cross-compiles the Windows binary expected for distributable releases.
- `go test ./...` executes all package tests; add `-run` filters or `-cover` when refining suites.

## Coding Style & Naming Conventions
- Format Go code with `gofmt`; prefer `goimports` to keep imports ordered.
- Use UpperCamelCase for exported identifiers and lowerCamelCase for internals, matching files such as `cmd/multitrack.go`.
- Place new commands in `cmd/<feature>.go`, keeping package-level state private unless shared by multiple commands.
- Reserve comments for non-obvious playback math or stream timing; avoid verbose narrative blocks.

## Testing Guidelines
- Create `_test.go` files beside the code they exercise; table-driven tests work well for command behaviors.
- Mock audio streams with short fixtures in `examples/` or generated buffers so tests stay fast.
- Cover parsing, state transitions, and multi-track timing helpers before proposing UI or asset changes.
- CI is manual—run `go test ./...` plus your local smoke run (`go run . play <file>`) before pushing.

## Commit & Pull Request Guidelines
- Follow the Conventional Commits style evident in history (`feat:`, `fix:`, `refactor:`); keep scopes lower-case and concise.
- Squash work-in-progress commits locally and write messages that summarize the observable change.
- Pull requests must describe behavior changes, include manual test notes (files played, commands used), and attach terminal captures for UI tweaks.
- Link related issues and call out any new sample media or configuration steps required for validation.
