package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/raziel2244/geckosite/s3"
)

// S3 handles minio operations.
func S3(w http.ResponseWriter, r *http.Request) {
	bucket := mux.Vars(r)["bucket"]

	obj, err := s3.Client.GetObject(
		context.Background(), bucket,
		r.URL.Path[len(bucket)+1:],
		minio.GetObjectOptions{},
	)

	if err != nil {
		NotFound(w, r)
		return
	}

	io.Copy(w, obj)
}
