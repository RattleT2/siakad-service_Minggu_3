package siakad

import (
    "errors"
    "net/http"
)

var (
    ErrNIMTidakValid      = errors.New("NIM wajib 8 digit angka")
    ErrNilaiTidakValid    = errors.New("mutu harus di antara 0.0 - 4.0")
    ErrSKSTidakValid      = errors.New("SKS harus lebih dari 0")
    ErrMahasiswaTidakAda  = errors.New("mahasiswa tidak ditemukan")
    ErrMahasiswaSudahAda  = errors.New("mahasiswa dengan NIM tersebut sudah terdaftar")
)

func MapErrorToHTTP(err error) int {
    switch {
    case errors.Is(err, ErrNIMTidakValid),
        errors.Is(err, ErrNilaiTidakValid),
        errors.Is(err, ErrSKSTidakValid):
        return http.StatusBadRequest
    case errors.Is(err, ErrMahasiswaTidakAda):
        return http.StatusNotFound
    case errors.Is(err, ErrMahasiswaSudahAda):
        return http.StatusConflict
    default:
        return http.StatusInternalServerError
    }
}
