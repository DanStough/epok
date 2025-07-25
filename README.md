# epok 
_Pronounced "epic-but-with-a-k"_

A CLI for working with Unix Timestamps inspired by [epochconverter.com](https://www.epochconverter.com).
Also kind of like Unix `date`, but for smoother brains and folks who like nanosecond timestamps.

![Terminal prompt showing generating a timestamp and parsing it with epok](./docs/assets/epok.gif)

Main commands:
1. **`parse`** - read a unix timestamp and return the human readable form. Infers the precision.
2. **`now`** - generate a unix timestamp for the current instant. Multiple precisions supported.
2. **`timezone`** (_TBA_)- work with system timezones (view, list, search)
3. **`at`** (_TBA_)- convert human readable timestamps and expressions to unix timestamps
4. **`between`** (_TBA_)- find the human readable delta between two timestamps

Built with great open source libraries:
* [spf13/cobra](https://github.com/spf13/cobra)
* [spf13/viper](https://github.com/spf13/viper)
* [charmbracelet/fang](https://github.com/charmbracelet/fang)
* [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)
* [charmbracelet/vhs (for GIF recording)](https://github.com/charmbracelet/vhs)
* [goreleaser/goreleaser](https://github.com/goreleaser/goreleaser)

## Installation

## Homebrew

Install on MacOS and Linux using this homebrew tap:
```bash
brew install DanStough/tap/epok
```

### Install from Source

You can install the tool directly if you have a Go installation >= 1.21:
```bash
go install github.com/DanStough/epok@latest
```

## Documentation

Run `epok help`.

## Development

> [!IMPORTANT]  
> Tests require the use of the experimental `testing/synctest` package.
> You will need to configure your Go environment with `GOEXPERIMENT=synctest`.
> Setting `GOTRACEBACK=all` can also be helpful for debugging synctest errors ([Go #70911](https://github.com/golang/go/issues/70911)).

For now, you can use the vanilla Go commands for building and testing: 
```bash
# Run
go run . parse --help

# Test
go test -v ./...
```

## TODO

<details>

<summary>MVP</summary>

* [X] `parse` command
  * [X]  Outputs UTC and local system time from arg
  * [X]  Outputs UTC and local system time from stdin
* [X]  Lipgloss for styling
* [X] `now` command to return the instantaneous timestamp.
  * [X] `-p,--precision` to specify the precision
* [X] `go install` instructions
* [X] Version command (handled by `fang`)
* [X] Flair: GIF + Ascii Art
* [X] CI
  * [X] go releaser
* [X] homebrew tap

</details>

### Milestone 1
* [ ] Makefile, Taskfile or `Just` to build
* [ ] `parse` command
  * [ ]  timezone flag for parse command to specify additional output zone
* [ ] `timezone` command
  * [ ] list command (include source?)
  * [ ] show current system timezone
* [X] Output Mode: `simple` - something that can be copy/pasted
* [ ] CI
  * [X] go test
  * [ ] conventional commits check

### Future
* [ ] `between` command - compare two timestamps
* [X] Output Mode: `json`
* [ ] golintci + CI
* [ ]  `at` command for generating a unix timestamp from multiple formats.
* [ ] Add "preferred timezones" to the config file, which are used when outputting human-readable information.
* [ ] batch process multiple timestamps and return tabular delta 
* [X] ~built-in copy/paste functionality (yes, I know `pbcopy`/`pbpaste` is a thing)~ now I'm thinking this doesn't make much sense if you can read from stdin.
* [ ] "default" command - alias your favorite command in the tool