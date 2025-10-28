package repository

import (
	"database/sql"
	"praktikum4-crud/app/model"
	"praktikum4-crud/database"
	"time"
)

func GetAllPekerjaan() ([]model.PekerjaanAlumni, error) {
	rows, err := database.DB.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		       tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni
		WHERE is_deleted IS NULL`) // perbarui query untuk menampilkan pekerjaan
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		var tMulai time.Time
		var tSelesai *time.Time
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
			&tMulai, &tSelesai, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.TanggalMulaiKerja = tMulai.Format("2006-01-02")
		if tSelesai != nil {
			p.TanggalSelesaiKerja = tSelesai.Format("2006-01-02")
		} else {
			p.TanggalSelesaiKerja = ""
		}
		list = append(list, p)
	}
	return list, nil
}

func GetPekerjaanByID(id int) (*model.PekerjaanAlumni, error) {
	var p model.PekerjaanAlumni
	var tMulai time.Time
	var tSelesai *time.Time
	err := database.DB.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		       tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni WHERE id=$1`, id).
		Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
			&tMulai, &tSelesai, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	p.TanggalMulaiKerja = tMulai.Format("2006-01-02")
	if tSelesai != nil {
		p.TanggalSelesaiKerja = tSelesai.Format("2006-01-02")
	} else {
		p.TanggalSelesaiKerja = ""
	}
	return &p, nil
}

func GetPekerjaanByAlumniID(alumniID int) ([]model.PekerjaanAlumni, error) {
	rows, err := database.DB.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		       tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni WHERE alumni_id=$1`, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		var tMulai time.Time
		var tSelesai *time.Time
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
			&tMulai, &tSelesai, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.TanggalMulaiKerja = tMulai.Format("2006-01-02")
		if tSelesai != nil {
			p.TanggalSelesaiKerja = tSelesai.Format("2006-01-02")
		} else {
			p.TanggalSelesaiKerja = ""
		}
		list = append(list, p)
	}
	return list, nil
}

func CreatePekerjaan(p *model.PekerjaanAlumni) (int, error) {
	tMulai, err := time.Parse("2006-01-02", p.TanggalMulaiKerja)
	if err != nil {
		return 0, err
	}
	var tSelesai *time.Time
	if p.TanggalSelesaiKerja != "" {
		ts, err := time.Parse("2006-01-02", p.TanggalSelesaiKerja)
		if err != nil {
			return 0, err
		}
		tSelesai = &ts
	}

	var id int
	err = database.DB.QueryRow(`
		INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		p.AlumniID, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange, tMulai, tSelesai, p.StatusPekerjaan, p.DeskripsiPekerjaan).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdatePekerjaan(id int, p *model.PekerjaanAlumni) error {
	tMulai, err := time.Parse("2006-01-02", p.TanggalMulaiKerja)
	if err != nil {
		return err
	}
	var tSelesai *time.Time
	if p.TanggalSelesaiKerja != "" {
		ts, err := time.Parse("2006-01-02", p.TanggalSelesaiKerja)
		if err != nil {
			return err
		}
		tSelesai = &ts
	}

	_, err = database.DB.Exec(`
		UPDATE pekerjaan_alumni SET alumni_id=$1, nama_perusahaan=$2, posisi_jabatan=$3, bidang_industri=$4, lokasi_kerja=$5, gaji_range=$6, tanggal_mulai_kerja=$7, tanggal_selesai_kerja=$8, status_pekerjaan=$9, deskripsi_pekerjaan=$10, updated_at=NOW()
		WHERE id=$11`,
		p.AlumniID, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange, tMulai, tSelesai, p.StatusPekerjaan, p.DeskripsiPekerjaan, id)
	return err
}

func DeletePekerjaan(id int) error {
	_, err := database.DB.Exec("DELETE FROM pekerjaan_alumni WHERE id=$1", id)
	return err
}

// SoftDeletePekerjaan hanya menandai pekerjaan sebagai terhapus
func SoftDeletePekerjaan(id int) error {
	_, err := database.DB.Exec(`UPDATE pekerjaan_alumni SET is_deleted=NOW() WHERE id=$1`, id)
	return err
}

func GetTrashPekerjaanRBACRepo(role string, userID int) ([]model.TrashPekerjaan, error) {
	var rows *sql.Rows
	var err error

	if role == "admin" {
		rows, err = database.DB.Query(`
			SELECT p.id, p.alumni_id, a.user_id, a.nama AS nama_alumni,
				   p.nama_perusahaan, p.posisi_jabatan, p.bidang_industri,
				   p.lokasi_kerja, p.is_deleted
			FROM pekerjaan_alumni p
			JOIN alumni a ON p.alumni_id = a.id
			WHERE p.is_deleted IS NOT NULL
			ORDER BY p.is_deleted DESC`)
	} else {
		rows, err = database.DB.Query(`
			SELECT p.id, p.alumni_id, a.user_id, a.nama AS nama_alumni,
				   p.nama_perusahaan, p.posisi_jabatan, p.bidang_industri,
				   p.lokasi_kerja, p.is_deleted
			FROM pekerjaan_alumni p
			JOIN alumni a ON p.alumni_id = a.id
			WHERE a.user_id=$1 AND p.is_deleted IS NOT NULL
			ORDER BY p.is_deleted DESC`, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.TrashPekerjaan
	for rows.Next() {
		var d model.TrashPekerjaan
		err := rows.Scan(
			&d.ID, &d.AlumniID, &d.UserID, &d.NamaAlumni,
			&d.NamaPerusahaan, &d.PosisiJabatan, &d.BidangIndustri,
			&d.LokasiKerja, &d.IsDeleted,
		)
		if err != nil {
			return nil, err
		}

		// Jika user_id NULL → ganti 0 agar tetap aman di JSON
		if d.UserID == nil {
			tmp := 0
			d.UserID = &tmp
		}

		list = append(list, d)
	}
	return list, nil
}


// Restore pekerjaan yang dihapus (set is_deleted = NULL)
func RestorePekerjaanRepo(id int) error {
	_, err := database.DB.Exec(`UPDATE pekerjaan_alumni SET is_deleted=NULL WHERE id=$1`, id)
	return err
}

// Hapus permanen pekerjaan
func HardDeletePekerjaanRepo(id int) error {
	_, err := database.DB.Exec(`DELETE FROM pekerjaan_alumni WHERE id=$1`, id)
	return err
}

// Cek pemilik pekerjaan dan status is_deleted
func GetOwnerAndDeleteStatus(id int) (alumniUserID int, isDeleted *time.Time, err error) {
	err = database.DB.QueryRow(`
		SELECT a.user_id, p.is_deleted
		FROM pekerjaan_alumni p
		JOIN alumni a ON p.alumni_id = a.id
		WHERE p.id=$1`, id).Scan(&alumniUserID, &isDeleted)
	return
}
