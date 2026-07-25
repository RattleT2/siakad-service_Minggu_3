package internal

import (
	"context"
	"fmt"
	"sort"
)

type Service interface {
	PerJurusan(ctx context.Context) (map[string]JurusanRekap, error)
	TopIPK(ctx context.Context, n int) ([]MahasiswaDetail, error)
	RingkasanMahasiswa(ctx context.Context, nim string) (MahasiswaDetail, error)
}

type rekapService struct {
	client AkademikClient
}

func NewService(client AkademikClient) Service {
	return &rekapService{client: client}
}

func (s *rekapService) PerJurusan(ctx context.Context) (map[string]JurusanRekap, error) {
	mahasiswas, err := s.client.DaftarMahasiswa(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar mahasiswa: %w", err)
	}

	jurusanMap := make(map[string][]MahasiswaDetail)

	for _, m := range mahasiswas {
		ringkasan, err := s.client.Ringkasan(ctx, m.NIM)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil ringkasan %s: %w", m.NIM, err)
		}

		detail := MahasiswaDetail{
			NIM:     m.NIM,
			Nama:    m.Nama,
			Jurusan: m.Jurusan,
			Status:  m.Status,
			IPK:     ringkasan.IPK,
		}

		jurusanMap[m.Jurusan] = append(jurusanMap[m.Jurusan], detail)
	}

	result := make(map[string]JurusanRekap)
	for jurusan, details := range jurusanMap {
		var totalIPK float64
		for _, d := range details {
			totalIPK += d.IPK
		}
		rataRata := totalIPK / float64(len(details))

		result[jurusan] = JurusanRekap{
			DaftarMahasiswa: details,
			RataRataIPK:     rataRata,
		}
	}

	return result, nil
}

func (s *rekapService) TopIPK(ctx context.Context, n int) ([]MahasiswaDetail, error) {
	if n <= 0 {
		n = 3
	}

	mahasiswas, err := s.client.DaftarMahasiswa(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar mahasiswa: %w", err)
	}

	var details []MahasiswaDetail
	for _, m := range mahasiswas {
		ringkasan, err := s.client.Ringkasan(ctx, m.NIM)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil ringkasan %s: %w", m.NIM, err)
		}

		details = append(details, MahasiswaDetail{
			NIM:     m.NIM,
			Nama:    m.Nama,
			Jurusan: m.Jurusan,
			Status:  m.Status,
			IPK:     ringkasan.IPK,
		})
	}

	sort.Slice(details, func(i, j int) bool {
		return details[i].IPK > details[j].IPK
	})

	if n > len(details) {
		n = len(details)
	}
	return details[:n], nil
}

func (s *rekapService) RingkasanMahasiswa(ctx context.Context, nim string) (MahasiswaDetail, error) {
	ringkasan, err := s.client.Ringkasan(ctx, nim)
	if err != nil {
		return MahasiswaDetail{}, err
	}

	return MahasiswaDetail{
		NIM:     ringkasan.NIM,
		Nama:    ringkasan.Nama,
		Jurusan: ringkasan.Jurusan,
		Status:  ringkasan.Status,
		IPK:     ringkasan.IPK,
	}, nil
}
