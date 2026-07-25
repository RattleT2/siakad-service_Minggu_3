package internal

type Mahasiswa struct {
	NIM     string `json:"nim"`
	Nama    string `json:"nama"`
	Jurusan string `json:"jurusan"`
	Status  string `json:"status"`
}

type MahasiswaDetail struct {
	NIM     string  `json:"nim"`
	Nama    string  `json:"nama"`
	Jurusan string  `json:"jurusan"`
	Status  string  `json:"status"`
	IPK     float64 `json:"ipk"`
}

type Ringkasan struct {
	NIM      string  `json:"nim"`
	Nama     string  `json:"nama"`
	Jurusan  string  `json:"jurusan"`
	Status   string  `json:"status"`
	TotalSKS int     `json:"total_sks"`
	IPK      float64 `json:"ipk"`
}

type JurusanRekap struct {
	DaftarMahasiswa []MahasiswaDetail `json:"daftar_mahasiswa"`
	RataRataIPK     float64           `json:"rata_rata_ipk"`
}
