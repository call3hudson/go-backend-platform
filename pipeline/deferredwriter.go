package pipeline

import (
	"net/http"
	"strings"
)

type DeferredResponseWriter struct {
	http.ResponseWriter
	strings.Builder
	statusCode int
}

func (dw *DeferredResponseWriter) Write(data []byte) (int, error) {
	return dw.Builder.Write(data)
}

func (dw *DeferredResponseWriter) FlushData() {
	if dw.statusCode == 0 {
		dw.statusCode = http.StatusOK
	}
	dw.ResponseWriter.WriteHeader(dw.statusCode)
	dw.ResponseWriter.Write([]byte(dw.Builder.String()))
}

func (dw *DeferredResponseWriter) writeHeader(statusCode int) {
	dw.statusCode = statusCode
}
