[![Lint | Test](https://github.com/jordan-rash/slog-handler/actions/workflows/ltb.yml/badge.svg)](https://github.com/jordan-rash/slog-handler/actions/workflows/ltb.yml)
[![Dependabot Updates](https://github.com/jordan-rash/slog-handler/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/jordan-rash/slog-handler/actions/workflows/dependabot/dependabot-updates)
[![CodeQL](https://github.com/jordan-rash/slog-handler/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/jordan-rash/slog-handler/actions/workflows/github-code-scanning/codeql)
[![Go Report Card](https://goreportcard.com/badge/disorder.dev/shandler)](https://goreportcard.com/report/disorder.dev/shandler)
[![Go Reference](https://pkg.go.dev/badge/disorder.dev/shandler.svg)](https://pkg.go.dev/disorder.dev/shandler)

# SLOG Handler with more knobs

As great as `log/slog` is, the provided handlers don't have enough customization knobs. This tries to
provide more flexibility to the user.

## Installation

```shell
go get disorder.dev/shandler
```

## Features

#### WithJSON

Enables JSON output for the log message. This is useful for structured logging.

#### WithLogLevel

Controls the log level for the message. This is useful for filtering messages.

#### WithLineInfo(short)

Adds the file and line number to a `slog_info` attribute within the log message  
if `short == true`; the file basename and line number are included.  
if `short == false`; the function name, full file path, and line number are included.

`short` is `true` by default. This is helpful in smaller programs where the full path is not needed.

#### WithTimeFormat

Controls the time format for the messages.

#### WithTextOutputFormat

This is a format string that gets used in text based logs. It takes 3 strings: time, level, and message (in that order). Include a newline at the end of your string.

#### WithStdOut

Controls which `io.Writer` is used for non-error log messages.

#### WithStdErr

Controls which `io.Writer` is used for error messages.

#### WithColor

Adds color to the log levels in text mode

#### With{Debug|Info|Warn|Error}Color

Overrides the default color for the log level.

#### WithShortLevels

Prints 3 character log levels instead of the full name. In text mode, this helps keep the log lines visually straight.

#### WithPid

Adds the process ID to the log message.

#### WithGroupRightJustify

Right justifies the log group name. This is useful for visually grouping log messages.
Will attempt to calculate terminal width; if an error occurs, it will default to 80 characters.
Overrides WithGroupTextOutputFormat

## Examples

```go
logger = slog.New(shandler.NewHandler(
 shandler.WithLogLevel(slog.LevelDebug),
 shandler.WithTimeFormat(time.RFC822),
 shandler.WithTextOutputFormat("%s | %s | %s\n"),
 shandler.WithStdErr(os.Stdout),
))
logger.With(slog.String("app", "myapp")).Debug("test")
```

#### Trace Log Level

Library includes an easier way to log trace messages. This is useful for debugging chatty logs.

```go
logger = slog.New(shandler.NewHandler(
 shandler.WithLogLevel(shandler.LevelTrace),
))
logger.Log(context.Background(), shandler.LevelTrace, "trace test")
```

## Benchmarks if you're into that sort of thing

```shell
goos: linux
goarch: amd64
pkg: disorder.dev/shandler
cpu: 13th Gen Intel(R) Core(TM) i9-13900H

BenchmarkHandlers/handler_text_log-20            3211198               351.4 ns/op
BenchmarkHandlers/stdlib_text_log-20             3325005               356.1 ns/op
BenchmarkHandlers/handler_json_log-20            2108128               592.3 ns/op
BenchmarkHandlers/stdlib_json_log-20             3445792               379.6 ns/op

PASS
ok      disorder.dev/shandler   6.604s
```

> The JSON Handler here is much slower as it uses the JSON library under the covers. The stdlib implementation builds the string manually, so its faster. 🤷🏼‍♀️
