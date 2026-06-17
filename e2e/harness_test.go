//go:build e2e

// Package e2e drives the compiled payments binary over a real HTTP socket,
// exercising whole request flows against the default fake gateway and an
// in-memory database. It needs no network and no provider credentials.
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

var baseURL string

func TestMain(m *testing.M) { os.Exit(run(m)) }

func run(m *testing.M) int {
	dir, err := os.MkdirTemp("", "payments-e2e")
	if err != nil {
		fmt.Fprintln(os.Stderr, "e2e: temp dir:", err)
		return 1
	}
	defer os.RemoveAll(dir)

	bin := filepath.Join(dir, "payments")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = ".."
	build.Stdout, build.Stderr = os.Stderr, os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "e2e: build:", err)
		return 1
	}

	port, err := freePort()
	if err != nil {
		fmt.Fprintln(os.Stderr, "e2e: free port:", err)
		return 1
	}
	baseURL = fmt.Sprintf("http://127.0.0.1:%d", port)

	srv := exec.Command(bin)
	srv.Env = append(os.Environ(), "PORT="+strconv.Itoa(port), "DB_PATH=:memory:")
	srv.Stdout, srv.Stderr = os.Stderr, os.Stderr
	if err := srv.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "e2e: start:", err)
		return 1
	}
	defer func() {
		_ = srv.Process.Kill()
		_ = srv.Wait()
	}()

	if err := waitReady(baseURL, 10*time.Second); err != nil {
		fmt.Fprintln(os.Stderr, "e2e: not ready:", err)
		return 1
	}
	return m.Run()
}

func freePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func waitReady(base string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(base + "/")
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("server did not become ready within %s", timeout)
}

// --- request helpers ---

var emailSeq atomic.Int64

func uniqueEmail() string {
	return fmt.Sprintf("e2e-%d@example.com", emailSeq.Add(1))
}

func postJSON(t *testing.T, path, body string) *http.Response {
	t.Helper()
	resp, err := http.Post(baseURL+path, "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	return resp
}

func getJSON(t *testing.T, path string) *http.Response {
	t.Helper()
	resp, err := http.Get(baseURL + path)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	return resp
}

func decode(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("decode body: %v", err)
	}
}

func wantStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode != want {
		t.Fatalf("status = %d, want %d", resp.StatusCode, want)
	}
}

// createCustomer creates a customer with a fresh email and returns its id.
func createCustomer(t *testing.T) int64 {
	t.Helper()
	resp := postJSON(t, "/v1/customers", fmt.Sprintf(`{"email":%q}`, uniqueEmail()))
	wantStatus(t, resp, http.StatusCreated)
	var c struct {
		ID int64 `json:"id"`
	}
	decode(t, resp, &c)
	if c.ID == 0 {
		t.Fatal("created customer has zero id")
	}
	return c.ID
}
