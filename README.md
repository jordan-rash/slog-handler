# SLOG Handler with more knobs

As great as `log/slog` is, the provided handlers don't have enough customization knobs.  This tries to 
provide more flexibility to the user.

## Features

#### WithJSON
Enables JSON output for the log message.  This is useful for structured logging.

#### WithLogLevel
Controls the log level for the message.  This is useful for filtering messages.

#### WithTimeFormat
Controls the time format for the messages.

#### WithTextOutputFormat
This is a format string that gets used in text based logs.  It takes 3 strings: time, level, and message (in that order).  Include a newline at the end of your string.

#### WithStdOut
Controls which `io.Writer` is used for non-error log messages.

#### WithStdErr 
Controls which `io.Writer` is used for error messages.

#### WithColor
Adds color to the log levels in text mode

#### With{Debug|Info|Warn|Error}Color
Overrides the default color for the log level.

#### WithShortLevels
Prints 3 character log levels instead of the full name.  In text mode, this helps keep the log lines visually straight.

## Example

```go 
logger = slog.New(handler.NewHandler(
	handler.WithLogLevel(slog.LevelDebug),
	handler.WithTimeFormat(time.RFC822),
	handler.WithTextOutputFormat("%s | %s | %s\n"),
	handler.WithStdErr(os.Stdout),
))
logger.With(slog.String("app", "myapp")).Debug("test")
```

### Benchmarks if you're into that sort of thing

##### Text Log Handler
```shell
BenchmarkTextLog/handler_text_log-20            33507366               350.1 ns/op           102 B/op          3 allocs/op
BenchmarkTextLogStdLib/handler_text_log-20      35334312               334.1 ns/op           117 B/op          0 allocs/op
```

##### JSON Log Handler
```shell
BenchmarkJSONLog/handler_json_log-20            21419242               546.6 ns/op           474 B/op          6 allocs/op
BenchmarkJSONLogStdLib/stdlib_json_log-20       38057089               340.3 ns/op           163 B/op          0 allocs/op
```

> The JSON Handler here is much slower as it uses the JSON library under the covers.  The stdlib implementation builds the string manually, so its faster.  🤷🏼‍♀️
