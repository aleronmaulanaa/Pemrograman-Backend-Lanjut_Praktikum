// package model

// import "time"

// type TrashPekerjaan struct {
// 	ID             int        `json:"id"`
// 	AlumniID       int        `json:"alumni_id"`
// 	UserID         int        `json:"user_id"`
// 	NamaAlumni     string     `json:"nama_alumni"`
// 	NamaPerusahaan string     `json:"nama_perusahaan"`
// 	PosisiJabatan  string     `json:"posisi_jabatan"`
// 	BidangIndustri string     `json:"bidang_industri"`
// 	LokasiKerja    string     `json:"lokasi_kerja"`
// 	IsDeleted      *time.Time `json:"is_deleted"`
// }

package model

import "time"

type TrashPekerjaan struct {
	ID             int        `json:"id"`
	AlumniID       int        `json:"alumni_id"`
	UserID         *int       `json:"user_id"`
	NamaAlumni     string     `json:"nama_alumni"`
	NamaPerusahaan string     `json:"nama_perusahaan"`
	PosisiJabatan  string     `json:"posisi_jabatan"`
	BidangIndustri string     `json:"bidang_industri"`
	LokasiKerja    string     `json:"lokasi_kerja"`
	IsDeleted      *time.Time `json:"is_deleted"`
}
