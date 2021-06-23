package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vvakame/sdlog/aelog"
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
	if err := h.hello2Handler(w, req); err != nil {
		// noop
	}
}

func (h *handlers) hello2Handler(w http.ResponseWriter, req *http.Request) (err error) {
	ctx := SpanContextFromHttpRequest(req.Context(), req)
	ctx, _ = StartSpan(ctx, "hello2")
	ctx = aelog.WithHTTPRequest(ctx, req)
	defer EndSpan(ctx, err)

	_, err = h.als.Insert(ctx, &AccessLog{
		ID: uuid.New().String(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		aelog.Errorf(ctx, "failed AccessLogStore.Insert : %s", err)
		_, _ = fmt.Fprint(w, err.Error())
		return
	}

	msg := req.FormValue("message")
	SetAttributesKV(ctx, map[string]interface{}{"message": msg, "time.Unix": time.Now().Unix()})

	if msg == "" {
		err = fmt.Errorf("message is required")
		aelog.Errorf(ctx, "message is required")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message is required"))
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Ack : %s : %s", time.Now(), msg)))

	return nil
}
