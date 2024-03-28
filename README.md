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
