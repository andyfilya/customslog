package customslog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"sync"

	"github.com/fatih/color"
)

var (
	timeStamp           = "15:04:05"
	eRR_INNER_HANDLER   = errors.New("error when call unwrap Handler")
	eRR_INNER_UNMARSHAL = errors.New("error unmarshal result of unwrap Handler")
)

type CustomSlogHandler struct {
	handler slog.Handler
	buff    *bytes.Buffer
	mu      *sync.Mutex
}

func (h *CustomSlogHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.handler.Enabled(ctx, lvl)
}

func (h *CustomSlogHandler) WithAttrs(attr []slog.Attr) slog.Handler {
	return &CustomSlogHandler{
		handler: h.handler.WithAttrs(attr),
		buff:    h.buff,
		mu:      h.mu,
	}
}

func (h *CustomSlogHandler) WithGroup(group string) slog.Handler {
	return &CustomSlogHandler{
		handler: h.handler.WithGroup(group),
		buff:    h.buff,
		mu:      h.mu,
	}
}

func (h *CustomSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	lvl := r.Level.String()
	switch r.Level {
	case slog.LevelInfo:
		lvl = color.HiGreenString("%s", lvl)
	case slog.LevelDebug:
		lvl = color.BlueString("%s", lvl)
	case slog.LevelError:
		lvl = color.RedString("%s", lvl)
	case slog.LevelWarn:
		lvl = color.HiRedString("%s", lvl)
	}
	lvl += ":"

	attr, err := h.attrMake(ctx, r)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(attr, "", " ")
	if err != nil {
		return err
	}
	out := "[" + color.YellowString(r.Time.Format(timeStamp)) + "] " + color.WhiteString(r.Message) + " " + color.BlackString(string(bytes))
	os.Stdout.WriteString(out)

	return nil
}

func (h *CustomSlogHandler) attrMake(ctx context.Context, r slog.Record) (map[string]any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.handler.Handle(ctx, r); err != nil {
		return nil, eRR_INNER_HANDLER
	}
	attr := map[string]any{}
	err := json.Unmarshal(h.buff.Bytes(), &attr)
	if err != nil {
		return nil, eRR_INNER_UNMARSHAL
	}
	return attr, nil
}

func NewHandler(opts *slog.HandlerOptions) *CustomSlogHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	buff := &bytes.Buffer{}
	return &CustomSlogHandler{
		handler: slog.NewJSONHandler(buff, &slog.HandlerOptions{
			Level:     opts.Level,
			AddSource: opts.AddSource,
		}),
		buff: buff,
		mu:   &sync.Mutex{},
	}
}
