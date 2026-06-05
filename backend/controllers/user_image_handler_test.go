package controllers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"github.com/fabiokaelin/f-oauth/config"
	user_pkg "github.com/fabiokaelin/f-oauth/pkg/user"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func makeMultipartUploadRequest(t *testing.T, route string, fieldName string, fileName string, contentType string, content []byte) *http.Request {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", `form-data; name="`+fieldName+`"; filename="`+fileName+`"`)
	if contentType != "" {
		partHeader.Set("Content-Type", contentType)
	}
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		t.Fatalf("failed to create multipart part: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("failed to write multipart content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, route, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestUserPostMeImage_SuccessForwardsToImageService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	currentUserID := uuid.NewV4()
	forwarded := false
	var forwardedPath string

	imageService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forwarded = true
		forwardedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer imageService.Close()

	config.InternalImageService = imageService.URL

	r := gin.New()
	r.POST("/upload", func(c *gin.Context) {
		c.Set("currentUser", user_pkg.User{ID: currentUserID})
		userPostMeImage(c)
	})

	req := makeMultipartUploadRequest(t, "/upload", "image", "avatar.jpg", "image/jpeg", []byte("fake-jpeg-content"))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	if !forwarded {
		t.Fatalf("expected request to be forwarded to image service")
	}

	expectedPath := "/api/profile/" + currentUserID.String()
	if forwardedPath != expectedPath {
		t.Fatalf("expected forwarded path %s, got %s", expectedPath, forwardedPath)
	}

	var payload map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected json response, got error: %v", err)
	}

	worked, ok := payload["worked"].(bool)
	if !ok || !worked {
		t.Fatalf("expected response worked=true, got: %v", payload)
	}
}

func TestUserPostMeImage_RejectsNonImageContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	currentUserID := uuid.NewV4()
	forwarded := false

	imageService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forwarded = true
		w.WriteHeader(http.StatusOK)
	}))
	defer imageService.Close()

	config.InternalImageService = imageService.URL

	r := gin.New()
	r.POST("/upload", func(c *gin.Context) {
		c.Set("currentUser", user_pkg.User{ID: currentUserID})
		userPostMeImage(c)
	})

	req := makeMultipartUploadRequest(t, "/upload", "image", "notes.txt", "text/plain", []byte("not-an-image"))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}

	if forwarded {
		t.Fatalf("expected no forwarding to image service for invalid content type")
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("invalid file type")) {
		t.Fatalf("expected invalid file type message, got: %s", w.Body.String())
	}
}

func TestUserPostMeImage_FallsBackToLocalStorageWhenImageServiceIsDown(t *testing.T) {
	gin.SetMode(gin.TestMode)

	currentUserID := uuid.NewV4()
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		if chdirErr := os.Chdir(oldWd); chdirErr != nil {
			t.Fatalf("failed to restore working directory: %v", chdirErr)
		}
	}()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}

	config.InternalImageService = "http://127.0.0.1:1"

	r := gin.New()
	r.POST("/upload", func(c *gin.Context) {
		c.Set("currentUser", user_pkg.User{ID: currentUserID})
		userPostMeImage(c)
	})

	req := makeMultipartUploadRequest(t, "/upload", "image", "avatar.jpg", "image/jpeg", []byte("fake-jpeg-content"))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("fallback")) {
		t.Fatalf("expected fallback response, got: %s", w.Body.String())
	}

	localFilePattern := filepath.Join("public", "images", "profileimage-"+currentUserID.String()+".*")
	matches, err := filepath.Glob(localFilePattern)
	if err != nil {
		t.Fatalf("failed to glob local image: %v", err)
	}
	if len(matches) == 0 {
		t.Fatalf("expected local image file to be created with pattern %s", localFilePattern)
	}
}
