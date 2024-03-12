package sl

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type HandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type Handler struct {
	slog.Handler
	l *log.Logger
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.HiMagentaString(level)
	case slog.LevelInfo:
		level = color.HiBlueString(level)
	case slog.LevelWarn:
		level = color.HiYellowString(level)
	case slog.LevelError:
		level = color.HiRedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := color.GreenString(r.Time.Format("[06-01-02 15:04:05.000]"))
	msg := color.CyanString(r.Message)
	h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))

	return nil
}

func NewHandler(out io.Writer, opts HandlerOptions) *Handler {
	return &Handler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
}
