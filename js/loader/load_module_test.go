package loader

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadModule_Local(t *testing.T) {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "loader_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// main.js
	mainContent := `import "./foo.js"; console.log("main");`
	err = os.WriteFile(filepath.Join(tmpDir, "main.js"), []byte(mainContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// foo.js
	fooContent := `import "./bar/baz"; console.log("foo");`
	err = os.WriteFile(filepath.Join(tmpDir, "foo.js"), []byte(fooContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// bar directory
	err = os.Mkdir(filepath.Join(tmpDir, "bar"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	// bar/baz.js
	bazContent := `console.log("baz");`
	err = os.WriteFile(filepath.Join(tmpDir, "bar", "baz.js"), []byte(bazContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Test LoadModule
	codes, err := LoadModule(tmpDir, tmpDir, "main.js", false)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	// Check keys
	// Keys should be "/main.js", "/foo.js", "/bar/baz.js"
	expectedKeys := []string{"/main.js", "/foo.js", "/bar/baz.js"}
	for _, key := range expectedKeys {
		if GetCode(codes, key) == nil {
			t.Errorf("Missing key: %s", key)
		}
	}

	// Check rewrites in main.js
	// import "./foo.js" should become import "/foo.js"
	mainCode := GetCode(codes, "/main.js").Code
	if !strings.Contains(mainCode, `import "/foo.js"`) {
		t.Errorf("Import rewrite failed in main.js. Got: %s", mainCode)
	}

	// Check rewrites in foo.js
	// import "./bar/baz" -> from foo.js (at root) -> bar/baz.js -> key /bar/baz.js
	fooCode := GetCode(codes, "/foo.js").Code
	if !strings.Contains(fooCode, `import "/bar/baz.js"`) {
		t.Errorf("Import rewrite failed in foo.js. Got: %s", fooCode)
	}
}

func TestLoadModule_URL(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/remote.js":
			// Imports a relative file on the server
			fmt.Fprint(w, `import "./remote_dep.js"; console.log("remote");`)
		case "/remote_dep.js":
			fmt.Fprint(w, `console.log("remote dependency");`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Create a temporary directory for the local entry point
	tmpDir, err := os.MkdirTemp("", "loader_url_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// main.js imports the remote file
	remoteMapUrl := server.URL + "/remote.js"
	mainContent := fmt.Sprintf(`import "%s";`, remoteMapUrl)
	err = os.WriteFile(filepath.Join(tmpDir, "main.js"), []byte(mainContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Test LoadModule with fetch=true
	codes, err := LoadModule(tmpDir, tmpDir, "main.js", true)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	// Expected keys:
	// 1. /main.js
	// 2. server.URL + "/remote.js"
	// 3. server.URL + "/remote_dep.js"

	if GetCode(codes, "/main.js") == nil {
		t.Error("Missing key: /main.js")
	}

	remoteKey := remoteMapUrl
	if GetCode(codes, remoteKey) == nil {
		t.Errorf("Missing key: %s", remoteKey)
	}

	remoteDepKey := server.URL + "/remote_dep.js"
	if GetCode(codes, remoteDepKey) == nil {
		t.Errorf("Missing key: %s", remoteDepKey)
	}

	// Check rewrites

	// main.js should still have the full URL import because the key IS the full URL
	// import "http://.../remote.js" -> import "http://.../remote.js"
	mainCode := GetCode(codes, "/main.js").Code
	if !strings.Contains(mainCode, fmt.Sprintf(`import "%s"`, remoteKey)) {
		t.Errorf("Import in main.js incorrect. Got: %s", mainCode)
	}

	// remote.js import "./remote_dep.js" should be rewritten to "http://.../remote_dep.js"
	remoteCode := GetCode(codes, remoteKey).Code
	if !strings.Contains(remoteCode, fmt.Sprintf(`import "%s"`, remoteDepKey)) {
		t.Errorf("Import rewrite in remote.js failed. Got: %s", remoteCode)
	}
}
