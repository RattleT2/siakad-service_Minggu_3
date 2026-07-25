package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AkademikClient interface {
	DaftarMahasiswa(ctx context.Context) ([]Mahasiswa, error)
	Ringkasan(ctx context.Context, nim string) (Ringkasan, error)
}

type akademikHTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewAkademikClient(baseURL string) AkademikClient {
	return &akademikHTTPClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

type daftarResponse struct {
	Sukses bool        `json:"sukses"`
	Data   []Mahasiswa `json:"data"`
	Error  string      `json:"error,omitempty"`
}

type ringkasanResponse struct {
	Sukses bool      `json:"sukses"`
	Data   Ringkasan `json:"data"`
	Error  string    `json:"error,omitempty"`
}

type ringkasanErrorResponse struct {
	Sukses bool   `json:"sukses"`
	Error  string `json:"error"`
}

func (a *akademikHTTPClient) DaftarMahasiswa(ctx context.Context) ([]Mahasiswa, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.baseURL+"/api/v1/mahasiswa", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAkademikTidakTersedia, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("%w: %v", ErrAkademikTimeout, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrAkademikTidakTersedia, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: gagal membaca respons", ErrAkademikTidakTersedia)
	}

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("%w: status %d", ErrAkademikTidakTersedia, resp.StatusCode)
	}

	var result daftarResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("%w: gagal mendekode JSON", ErrAkademikTidakTersedia)
	}

	if !result.Sukses {
		return nil, fmt.Errorf("%w: %s", ErrAkademikTidakTersedia, result.Error)
	}

	return result.Data, nil
}

func (a *akademikHTTPClient) Ringkasan(ctx context.Context, nim string) (Ringkasan, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.baseURL+"/api/v1/mahasiswa/"+nim+"/ringkasan", nil)
	if err != nil {
		return Ringkasan{}, fmt.Errorf("%w: %v", ErrAkademikTidakTersedia, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return Ringkasan{}, fmt.Errorf("%w: %v", ErrAkademikTimeout, err)
		}
		return Ringkasan{}, fmt.Errorf("%w: %v", ErrAkademikTidakTersedia, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Ringkasan{}, fmt.Errorf("%w: gagal membaca respons", ErrAkademikTidakTersedia)
	}

	if resp.StatusCode == http.StatusNotFound {
		var errResp ringkasanErrorResponse
		json.Unmarshal(body, &errResp)
		return Ringkasan{}, fmt.Errorf("%w: %s", ErrMahasiswaTidakAda, errResp.Error)
	}

	if resp.StatusCode >= 500 {
		return Ringkasan{}, fmt.Errorf("%w: status %d", ErrAkademikTidakTersedia, resp.StatusCode)
	}

	var result ringkasanResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return Ringkasan{}, fmt.Errorf("%w: gagal mendekode JSON", ErrAkademikTidakTersedia)
	}

	if !result.Sukses {
		return Ringkasan{}, fmt.Errorf("%w: %s", ErrAkademikTidakTersedia, result.Error)
	}

	return result.Data, nil
}
