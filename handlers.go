package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

type handlers struct {
	als *AccessLogStore
}

func (h *handlers) HelloHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(label.String("server", "handling this..."))

	_, err := h.als.Insert(ctx, &AccessLog{
		ID: uuid.New().String(),
	})
	if err != nil {
		_, _ = fmt.Fprint(w, err.Error())
	}

	_, _ = io.WriteString(w, "Hello, otel world!")
}

func (h *handlers) Hello2Handler(w http.ResponseWriter, req *http.Request) {
	ctx := SpanContextFromHttpRequest(req.Context(), req)
	ctx, _ = StartSpan(ctx, "hello2")
	var err error
	defer EndSpan(ctx, err)

	_, err = h.als.Insert(ctx, &AccessLog{
		ID: uuid.New().String(),
	})
	if err != nil {
		_, _ = fmt.Fprint(w, err.Error())
	}

	msg := req.FormValue("message")
	SetAttributesKV(ctx, map[string]interface{}{"message": msg, "time.Unix": time.Now().Unix()})

	if msg == "" {
		err = fmt.Errorf("message is required")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message is required"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Ack : %s : %s", time.Now(), msg)))
}
