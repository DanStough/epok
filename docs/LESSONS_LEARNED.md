# Things We Learned Along the Way

## Reading from STDIN can block all calls to read it
Some more context in this [GH Issue](https://github.com/charmbracelet/fang/issues/60) filed against `fang`.
Not only is it nearly impossible to cancel an `io.Reader`, reading from a stream and not closing it might block other attempts to read from the same stream.
`fang` reads from STDIN on error to find out more about the terminal.

Here is a good article about why [readers aren't really cancelable](https://benjamincongdon.me/blog/2020/04/23/Cancelable-Reads-in-Go/).

## Writing tests with `synctest`
AS usual, the Go docs were pretty good, including this [blog post](https://go.dev/blog/synctest) about the feature.
I also found this [blog post](https://victoriametrics.com/blog/go-synctest/) helpful to understand some of the constraints.
