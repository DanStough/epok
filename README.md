# epok
A CLI for working with Unix Timestamps inspired by [epochconverter.com](https://www.epochconverter.com).

Built with great open source libraries:
* [spf13/cobra](https://github.com/spf13/cobra)
* [spf13/viper](https://github.com/spf13/viper)
* [charmbracelet/fang](https://github.com/charmbracelet/fang)
* [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)

## Installation

### Install from Source

You can install the tool directly if you have a Go installation >= 1.24:
```bash
go install github.com/DanStough/epok@latest
```

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

### MVP
* [X] `parse` command
  * [X]  Outputs UTC and local system time from arg
  * [X]  Outputs UTC and local system time from stdin
* [ ] `now` command to return the instantaneous timestamp.
  * [ ] `-p,--precision` to specify the precision
* [ ]  Lipgloss for styling
* [X] `go install` instructions
* [X] Version command (handled by `fang`)
* [ ] Makefile or Just to build

### Part 1
* [ ] `parse` command
  * [ ]  timezone flag for parse command to specify additional output zone
* [ ] `timezone` command
  * [ ] list command (include source?)
  * [ ] show current system timezone
* [ ] Output Mode: `simple` - something that can be copy/pasted
* [ ] CI
  * [ ] go releaser
  * [ ] go test
  * [ ] conventional commits check
* [ ] homebrew tap

### Future
* [ ] Output Mode: `json`
* [ ]  `at` command for generating a unix timestamp from multiple formats.
* [ ] Add "preferred timezones" to the config file, which are used when outputing 
human readable information.
* [ ] batch process multiple timestamps and return tabular delta 
* [ ] built-in copy/paste functionality (yes, I know `pbcopy`/`pbpaste` is a thing)
* [ ] "default" command - alias your favorite command in the tool