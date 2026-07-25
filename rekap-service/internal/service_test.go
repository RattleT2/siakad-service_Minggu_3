package internal

import (
	"context"
	"testing"
)

type mockAkademikClient struct {
	mahasiswas []Mahasiswa
	ringkasans map[string]Ringkasan
	ringkErr   map[string]error
}

func newMockAkademikClient() *mockAkademikClient {
	return &mockAkademikClient{
		mahasiswas: []Mahasiswa{},
		ringkasans: make(map[string]Ringkasan),
		ringkErr:   make(map[string]error),
	}
}

func (m *mockAkademikClient) DaftarMahasiswa(ctx context.Context) ([]Mahasiswa, error) {
	return m.mahasiswas, nil
}

func (m *mockAkademikClient) Ringkasan(ctx context.Context, nim string) (Ringkasan, error) {
	if err, ok := m.ringkErr[nim]; ok {
		return Ringkasan{}, err
	}
	return m.ringkasans[nim], nil
}

func TestTopIPK(t *testing.T) {
	client := newMockAkademikClient()
	client.mahasiswas = []Mahasiswa{
		{NIM: "23010001", Nama: "Bunga", Jurusan: "Teknik Informatika", Status: "Aktif"},
		{NIM: "23010002", Nama: "Mawar", Jurusan: "Sistem Informasi", Status: "Aktif"},
		{NIM: "23010003", Nama: "Melati", Jurusan: "Teknik Informatika", Status: "Aktif"},
		{NIM: "23010004", Nama: "Anggrek", Jurusan: "Manajemen", Status: "Cuti"},
	}
	client.ringkasans = map[string]Ringkasan{
		"23010001": {NIM: "23010001", Nama: "Bunga", Jurusan: "Teknik Informatika", Status: "Aktif", TotalSKS: 20, IPK: 3.80},
		"23010002": {NIM: "23010002", Nama: "Mawar", Jurusan: "Sistem Informasi", Status: "Aktif", TotalSKS: 18, IPK: 3.50},
		"23010003": {NIM: "23010003", Nama: "Melati", Jurusan: "Teknik Informatika", Status: "Aktif", TotalSKS: 16, IPK: 3.20},
		"23010004": {NIM: "23010004", Nama: "Anggrek", Jurusan: "Manajemen", Status: "Cuti", TotalSKS: 10, IPK: 2.90},
	}

	svc := NewService(client)
	result, err := svc.TopIPK(context.Background(), 3)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result))
	}
	if result[0].NIM != "23010001" {
		t.Errorf("expected top NIM 23010001, got %s", result[0].NIM)
	}
	if result[0].IPK != 3.80 {
		t.Errorf("expected top IPK 3.80, got %v", result[0].IPK)
	}
	if result[1].NIM != "23010002" {
		t.Errorf("expected second NIM 23010002, got %s", result[1].NIM)
	}
	if result[2].NIM != "23010003" {
		t.Errorf("expected third NIM 23010003, got %s", result[2].NIM)
	}
}

func TestTopIPKDefaultN(t *testing.T) {
	client := newMockAkademikClient()
	client.mahasiswas = []Mahasiswa{
		{NIM: "23010001", Nama: "Bunga", Jurusan: "TI", Status: "Aktif"},
		{NIM: "23010002", Nama: "Mawar", Jurusan: "SI", Status: "Aktif"},
		{NIM: "23010003", Nama: "Melati", Jurusan: "TI", Status: "Aktif"},
		{NIM: "23010004", Nama: "Anggrek", Jurusan: "MNJ", Status: "Cuti"},
	}
	client.ringkasans["23010001"] = Ringkasan{NIM: "23010001", IPK: 3.80}
	client.ringkasans["23010002"] = Ringkasan{NIM: "23010002", IPK: 3.50}
	client.ringkasans["23010003"] = Ringkasan{NIM: "23010003", IPK: 3.20}
	client.ringkasans["23010004"] = Ringkasan{NIM: "23010004", IPK: 2.90}

	svc := NewService(client)
	result, err := svc.TopIPK(context.Background(), 0)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 results (default), got %d", len(result))
	}
}

func TestPerJurusan(t *testing.T) {
	client := newMockAkademikClient()
	client.mahasiswas = []Mahasiswa{
		{NIM: "23010001", Nama: "Bunga", Jurusan: "Teknik Informatika", Status: "Aktif"},
		{NIM: "23010002", Nama: "Mawar", Jurusan: "Sistem Informasi", Status: "Aktif"},
		{NIM: "23010003", Nama: "Melati", Jurusan: "Teknik Informatika", Status: "Aktif"},
	}
	client.ringkasans["23010001"] = Ringkasan{NIM: "23010001", IPK: 3.80}
	client.ringkasans["23010002"] = Ringkasan{NIM: "23010002", IPK: 3.50}
	client.ringkasans["23010003"] = Ringkasan{NIM: "23010003", IPK: 3.20}

	svc := NewService(client)
	result, err := svc.PerJurusan(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	ti, ok := result["Teknik Informatika"]
	if !ok {
		t.Fatal("expected Teknik Informatika jurusan")
	}
	if len(ti.DaftarMahasiswa) != 2 {
		t.Errorf("expected 2 mahasiswa in TI, got %d", len(ti.DaftarMahasiswa))
	}
	expectedAvg := (3.80 + 3.20) / 2.0
	if ti.RataRataIPK != expectedAvg {
		t.Errorf("expected rata-rata %v, got %v", expectedAvg, ti.RataRataIPK)
	}

	si, ok := result["Sistem Informasi"]
	if !ok {
		t.Fatal("expected Sistem Informasi jurusan")
	}
	if len(si.DaftarMahasiswa) != 1 {
		t.Errorf("expected 1 mahasiswa in SI, got %d", len(si.DaftarMahasiswa))
	}
	if si.RataRataIPK != 3.50 {
		t.Errorf("expected rata-rata 3.50, got %v", si.RataRataIPK)
	}
}

func TestRingkasanMahasiswa(t *testing.T) {
	client := newMockAkademikClient()
	client.mahasiswas = []Mahasiswa{
		{NIM: "23010001", Nama: "Bunga", Jurusan: "Teknik Informatika", Status: "Aktif"},
	}
	client.ringkasans["23010001"] = Ringkasan{
		NIM:      "23010001",
		Nama:     "Bunga",
		Jurusan:  "Teknik Informatika",
		Status:   "Aktif",
		TotalSKS: 20,
		IPK:      3.80,
	}

	svc := NewService(client)
	result, err := svc.RingkasanMahasiswa(context.Background(), "23010001")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.NIM != "23010001" {
		t.Errorf("expected NIM 23010001, got %s", result.NIM)
	}
	if result.IPK != 3.80 {
		t.Errorf("expected IPK 3.80, got %v", result.IPK)
	}
}

func TestRingkasanMahasiswaNotFound(t *testing.T) {
	client := newMockAkademikClient()
	client.ringkErr["99999999"] = ErrMahasiswaTidakAda

	svc := NewService(client)
	_, err := svc.RingkasanMahasiswa(context.Background(), "99999999")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !containsError(err, ErrMahasiswaTidakAda) {
		t.Errorf("expected ErrMahasiswaTidakAda, got %v", err)
	}
}

func containsError(err, target error) bool {
	if err == nil || target == nil {
		return false
	}
	return err.Error() == target.Error() ||
		len(err.Error()) >= len(target.Error()) &&
			err.Error()[:len(target.Error())] == target.Error()
}
