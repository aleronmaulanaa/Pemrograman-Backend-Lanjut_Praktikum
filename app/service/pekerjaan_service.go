package service

import (
	"praktikum4-crud/app/model"
	"praktikum4-crud/database"
	"praktikum4-crud/app/repository"
	"time"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

// GET ALL
func GetAllPekerjaan(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, 
               tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at 
        FROM pekerjaan_alumni
        WHERE is_deleted IS NULL`) // perbarui query untuk menampilkan pekerjaan
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		var tglMulai time.Time
		var tglSelesai *time.Time

		err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
			&tglMulai, &tglSelesai, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		p.TanggalMulaiKerja = tglMulai.Format("2006-01-02")
		if tglSelesai != nil {
			p.TanggalSelesaiKerja = tglSelesai.Format("2006-01-02")
		} else {
			p.TanggalSelesaiKerja = ""
		}

		pekerjaanList = append(pekerjaanList, p)
	}

	return c.JSON(pekerjaanList)
}

// GET BY ID
func GetPekerjaanByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var p model.PekerjaanAlumni
	var tglMulai time.Time
	var tglSelesai *time.Time

	err := database.DB.QueryRow(`
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, 
               tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at 
        FROM pekerjaan_alumni WHERE id=$1`, id).
		Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
			&tglMulai, &tglSelesai, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
	}

	p.TanggalMulaiKerja = tglMulai.Format("2006-01-02")
	if tglSelesai != nil {
		p.TanggalSelesaiKerja = tglSelesai.Format("2006-01-02")
	} else {
		p.TanggalSelesaiKerja = ""
	}

	return c.JSON(p)
}

// GET BY ALUMNI ID → diperbarui agar otomatis ambil dari JWT
func GetPekerjaanByAlumniID(c *fiber.Ctx) error {
	// Ambil user_id dari JWT
	userVal := c.Locals("user_id")
	userID, ok := userVal.(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user_id in token"})
	}

	rows, err := database.DB.Query(`
		SELECT p.id, p.alumni_id, p.nama_perusahaan, p.posisi_jabatan, p.bidang_industri, 
		       p.lokasi_kerja, p.gaji_range, p.tanggal_mulai_kerja, p.tanggal_selesai_kerja, 
		       p.status_pekerjaan, p.deskripsi_pekerjaan, p.created_at, p.updated_at
		FROM pekerjaan_alumni p
		JOIN alumni a ON p.alumni_id = a.id
		WHERE a.user_id=$1 AND p.is_deleted IS NULL`, int(userID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		var tglMulai time.Time
		var tglSelesai *time.Time

		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &tglMulai, &tglSelesai,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		p.TanggalMulaiKerja = tglMulai.Format("2006-01-02")
		if tglSelesai != nil {
			p.TanggalSelesaiKerja = tglSelesai.Format("2006-01-02")
		} else {
			p.TanggalSelesaiKerja = ""
		}

		pekerjaanList = append(pekerjaanList, p)
	}

	return c.JSON(pekerjaanList)
}


// CREATE
func CreatePekerjaan(c *fiber.Ctx) error {
	var p model.PekerjaanAlumni
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Parse tanggal string → time.Time
	tglMulai, err := time.Parse("2006-01-02", p.TanggalMulaiKerja)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tanggal_mulai_kerja salah, gunakan YYYY-MM-DD"})
	}

	var tglSelesai *time.Time
	if p.TanggalSelesaiKerja != "" {
		ts, err := time.Parse("2006-01-02", p.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format tanggal_selesai_kerja salah, gunakan YYYY-MM-DD"})
		}
		tglSelesai = &ts
	}

	err = database.DB.QueryRow(
		`INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan) 
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		p.AlumniID, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange, tglMulai, tglSelesai, p.StatusPekerjaan, p.DeskripsiPekerjaan,
	).Scan(&p.ID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(p)
}

// UPDATE
func UpdatePekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")
	var p model.PekerjaanAlumni
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Parse tanggal string → time.Time
	tglMulai, err := time.Parse("2006-01-02", p.TanggalMulaiKerja)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tanggal_mulai_kerja salah, gunakan YYYY-MM-DD"})
	}

	var tglSelesai *time.Time
	if p.TanggalSelesaiKerja != "" {
		ts, err := time.Parse("2006-01-02", p.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format tanggal_selesai_kerja salah, gunakan YYYY-MM-DD"})
		}
		tglSelesai = &ts
	}

	_, err = database.DB.Exec(
		`UPDATE pekerjaan_alumni SET alumni_id=$1, nama_perusahaan=$2, posisi_jabatan=$3, bidang_industri=$4, lokasi_kerja=$5, gaji_range=$6, tanggal_mulai_kerja=$7, tanggal_selesai_kerja=$8, status_pekerjaan=$9, deskripsi_pekerjaan=$10, updated_at=NOW() WHERE id=$11`,
		p.AlumniID, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange, tglMulai, tglSelesai, p.StatusPekerjaan, p.DeskripsiPekerjaan, id,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil diupdate"})
}

// DELETE
func DeletePekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := database.DB.Exec("DELETE FROM pekerjaan_alumni WHERE id=$1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil dihapus"})
}

// DELETE dengan soft delete + RBAC
func DeletePekerjaanRBAC(c *fiber.Ctx) error {
    idStr := c.Params("id")
    idInt, err := strconv.Atoi(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
    }

    // Ambil role & user_id dari JWT
    roleVal := c.Locals("role")
    userVal := c.Locals("user_id")

    role, _ := roleVal.(string)
    userID, ok := userVal.(float64)
    if !ok {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid user_id in token"})
    }

    // Jika admin → boleh hapus apapun
    if role == "admin" {
        if err := repository.SoftDeletePekerjaan(idInt); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(fiber.Map{"message": "Pekerjaan berhasil dihapus oleh admin"})
    }

    // Jika user → cek apakah pekerjaan milik alumni dia
    var alumniUserID int
    err = database.DB.QueryRow(`
        SELECT a.user_id
        FROM pekerjaan_alumni p
        JOIN alumni a ON p.alumni_id = a.id
        WHERE p.id=$1 AND p.is_deleted IS NULL`, idInt).Scan(&alumniUserID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Data pekerjaan tidak ditemukan"})
    }

    if int(userID) != alumniUserID {
        return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Anda hanya bisa menghapus pekerjaan Anda sendiri"})
    }

    if err := repository.SoftDeletePekerjaan(idInt); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Pekerjaan berhasil dihapus oleh user"})
}

// GET /api/pekerjaan/trash
func GetTrashPekerjaanRBAC(c *fiber.Ctx) error {
	roleVal := c.Locals("role")
	userVal := c.Locals("user_id")

	role, _ := roleVal.(string)
	userID, ok := userVal.(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user_id in token"})
	}

	data, err := repository.GetTrashPekerjaanRBACRepo(role, int(userID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if len(data) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tidak ada pekerjaan yang dihapus"})
	}

	return c.JSON(data)
}

// PUT /api/pekerjaan/restore/:id
func RestorePekerjaanRBAC(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	roleVal := c.Locals("role")
	userVal := c.Locals("user_id")

	role, _ := roleVal.(string)
	userID, ok := userVal.(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user_id in token"})
	}

	// Cek siapa pemiliknya
	alumniUserID, _, err := repository.GetOwnerAndDeleteStatus(idInt)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
	}

	// Hanya admin atau pemilik
	if role != "admin" && int(userID) != alumniUserID {
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Anda hanya bisa me-restore pekerjaan Anda sendiri"})
	}

	// Restore
	err = repository.RestorePekerjaanRepo(idInt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil direstore"})
}

// DELETE /api/pekerjaan/hard/:id
func HardDeletePekerjaanRBAC(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	roleVal := c.Locals("role")
	userVal := c.Locals("user_id")

	role, _ := roleVal.(string)
	userID, ok := userVal.(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user_id in token"})
	}

	alumniUserID, isDeleted, err := repository.GetOwnerAndDeleteStatus(idInt)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
	}

	if isDeleted == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data belum dihapus (soft delete) sehingga tidak bisa dihapus permanen"})
	}

	if role != "admin" && int(userID) != alumniUserID {
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Anda hanya bisa menghapus pekerjaan Anda sendiri"})
	}

	err = repository.HardDeletePekerjaanRepo(idInt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil dihapus permanen"})
}
