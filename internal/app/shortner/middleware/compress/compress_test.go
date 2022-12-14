package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/service"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
)

func TestGzipHandler_GzipMiddleware(t *testing.T) {

	config.BuildConfig()
	usm := storage.NewURLStorage()
	us := service.NewURLService(usm)
	uh := handlers.NewURLHandler(us)

	tests := []struct {
		name                     string
		url                      string
		headerAcceptEncoding     string
		headerContentEncoding    string
		expectedResponseBody     []byte
		expectedResponseCEHeader string
	}{
		// TODO: Add more test cases.
		{
			name:                     "decoded",
			url:                      "http://www.booking.com/",
			headerAcceptEncoding:     "gzip",
			headerContentEncoding:    "gzip",
			expectedResponseBody:     []byte(`http://localhost:8080/4cd89a20`),
			expectedResponseCEHeader: "gzip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Need to gzip-encode the URL before passing to request
			var buf bytes.Buffer       // Will write to buf
			gz := gzip.NewWriter(&buf) // Create gzip.Writer

			_, err := gz.Write([]byte(tt.url))
			if err != nil {
				log.Fatalf("failed write data to compress temporary buffer: %v", err)
			}

			err = gz.Close()
			if err != nil {
				log.Fatalf("failed compress data: %v", err)
			}

			r := httptest.NewRequest(http.MethodPost, "http://localhost:8080", &buf)
			w := httptest.NewRecorder()

			r.Header.Set("Content-Encoding", tt.headerAcceptEncoding)
			r.Header.Set("Accept-Encoding", tt.headerAcceptEncoding)

			h := GzipMiddleware(http.HandlerFunc(uh.Create))

			h.ServeHTTP(w, r)
			log.Println("compress_test: statuscode", w.Code)

			zr, err := gzip.NewReader(w.Body)
			if err != nil {
				log.Fatal(err)
			}

			body, err := io.ReadAll(zr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Println("compress_test: w.Body =", string(body))

			assert.Equal(t, tt.expectedResponseBody, body)
			assert.Equal(t, tt.expectedResponseCEHeader, w.Header().Get("Content-Encoding"))
			if err := zr.Close(); err != nil {
				log.Fatal(err)
			}
		})
	}
}
