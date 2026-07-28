package server

import (
	"encoding/json"
	"os"
	"time"
)

type ResponseLog struct {
	ResponseWriter
	size   int
	status int
}

func (w *ResponseLog) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *ResponseLog) Write(p []byte) (int, error) {
	w.size += len(p)
	return w.ResponseWriter.Write(p)
}

type Middleware func(HandlerFunc) HandlerFunc

func RequestMiddleware(hand HandlerFunc) HandlerFunc {
	return func(w ResponseWriter, r *Request) {

		start := time.Now()

		logFields := []any{"method", r.Method, "path", r.Path}

		if r.Body != nil && len(r.Body) > 0 {
			bodyStr := r.Body
			logFields = append(logFields, "body", bodyStr)
		}

		if len(r.Headers) > 0 {
			logFields = append(logFields, "headers", r.Headers)
		}

		writeLog(logFields)

		hand(w, r)

		rw := &ResponseLog{ResponseWriter: w,
			size:   0,
			status: 200}

		hand(rw, r)

		responseFields := []any{
			"type", "response",
			"method", r.Method,
			"path", r.Path,
			"status", rw.status,
			"size", rw.size,
			"duration_ms", time.Since(start).Milliseconds(),
		}

		writeLog(responseFields)
	}
}

func writeLog(fields []any) {
	logMap := make(map[string]interface{})
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			if key, ok := fields[i].(string); ok {
				logMap[key] = fields[i+1]
			}
		}
	}
	logMap["time"] = time.Now().Format(time.RFC3339)

	jsonData, err := json.MarshalIndent(logMap, "", " ")
	if err != nil {
		return
	}

	if err := os.MkdirAll("logs", 0755); err != nil {
		return
	}

	file, err := os.OpenFile("logs/requests.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	file.Write(append(jsonData, '\n'))
}
