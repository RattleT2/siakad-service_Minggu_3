package internal

import (
	"errors"
	"net/http"
)

var (
	ErrAkademikTimeout       = errors.New("waktu pemanggilan ke akademik service habis")
	ErrAkademikTidakTersedia = errors.New("akademik service tidak tersedia")
	ErrMahasiswaTidakAda     = errors.New("mahasiswa tidak ditemukan")
)

func MapErrorToHTTP(err error) int {
	switch {
	case errors.Is(err, ErrMahasiswaTidakAda):
		return http.StatusNotFound
	case errors.Is(err, ErrAkademikTimeout):
		return http.StatusGatewayTimeout
	case errors.Is(err, ErrAkademikTidakTersedia):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
