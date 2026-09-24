package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
)

type FormatHandler struct {
	mu    *sync.Mutex
	w     *os.File
	attrs []slog.Attr
}

func NewFormatHandler(w *os.File) *FormatHandler {
	return &FormatHandler{
		mu: &sync.Mutex{},
		w:  w,
	}
}

func (h *FormatHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *FormatHandler) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder

	b.WriteString(r.Time.Format("2006-01-02 15:04:05"))
	b.WriteByte('\t')
	b.WriteString(r.Level.String())
	b.WriteByte('\t')
	b.WriteString(r.Message)

	for _, a := range h.attrs {
		b.WriteByte('\t')
		writeAttr(&b, a)
	}

	r.Attrs(func(a slog.Attr) bool {
		b.WriteByte('\t')
		writeAttr(&b, a)
		return true
	})

	b.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.WriteString(b.String())
	if err != nil {
		return fmt.Errorf("write log: %w", err)
	}

	return nil
}

func (h *FormatHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	mergedAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	mergedAttrs = append(mergedAttrs, h.attrs...)
	mergedAttrs = append(mergedAttrs, attrs...)

	return &FormatHandler{
		mu:    h.mu,
		w:     h.w,
		attrs: mergedAttrs,
	}
}

func (h *FormatHandler) WithGroup(_ string) slog.Handler {
	return h
}

func writeAttr(b *strings.Builder, a slog.Attr) {
	b.WriteString(a.Key)
	b.WriteByte('=')
	b.WriteString(fmt.Sprint(a.Value.Any()))
}
