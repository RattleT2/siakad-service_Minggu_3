package internal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRingkasan200OK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"sukses": true,
			"data": {
				"nim": "23010001",
				"nama": "Bunga",
				"jurusan": "Teknik Informatika",
				"status": "Aktif",
				"total_sks": 5,
				"ipk": 3.60
			}
		}`))
	}))
	defer server.Close()

	client := NewAkademikClient(server.URL)
	result, err := client.Ringkasan(context.Background(), "23010001")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.NIM != "23010001" {
		t.Errorf("expected NIM 23010001, got %s", result.NIM)
	}
	if result.IPK != 3.60 {
		t.Errorf("expected IPK 3.60, got %v", result.IPK)
	}
}

func TestRingkasan404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{
			"sukses": false,
			"error": "mahasiswa tidak ditemukan"
		}`))
	}))
	defer server.Close()

	client := NewAkademikClient(server.URL)
	_, err := client.Ringkasan(context.Background(), "99999999")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), ErrMahasiswaTidakAda.Error()) {
		t.Errorf("expected ErrMahasiswaTidakAda, got %v", err)
	}
}

func TestRingkasan500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client := NewAkademikClient(server.URL)
	_, err := client.Ringkasan(context.Background(), "23010001")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), ErrAkademikTidakTersedia.Error()) {
		t.Errorf("expected ErrAkademikTidakTersedia, got %v", err)
	}
}

func TestRingkasanTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sukses": true, "data": {}}`))
	}))
	defer server.Close()

	client := NewAkademikClient(server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, err := client.Ringkasan(ctx, "23010001")

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), ErrAkademikTimeout.Error()) {
		t.Errorf("expected ErrAkademikTimeout, got %v", err)
	}
}

func TestRingkasanJSONRusak(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{broken json`))
	}))
	defer server.Close()

	client := NewAkademikClient(server.URL)
	_, err := client.Ringkasan(context.Background(), "23010001")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), ErrAkademikTidakTersedia.Error()) {
		t.Errorf("expected ErrAkademikTidakTersedia, got %v", err)
	}
}

func TestDaftarMahasiswa200OK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"sukses": true,
			"data": [
				{"nim": "23010001", "nama": "Bunga", "jurusan": "Teknik Informatika", "status": "Aktif"},
				{"nim": "23010002", "nama": "Mawar", "jurusan": "Sistem Informasi", "status": "Aktif"}
			]
		}`))
	}))
	defer server.Close()

	client := NewAkademikClient(server.URL)
	result, err := client.DaftarMahasiswa(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 mahasiswa, got %d", len(result))
	}
	if result[0].NIM != "23010001" {
		t.Errorf("expected NIM 23010001, got %s", result[0].NIM)
	}
}
