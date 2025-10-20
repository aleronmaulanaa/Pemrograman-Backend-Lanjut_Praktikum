package repository

import (
	"fmt"
	"strings"

	"praktikum4-crud/app/model"
	"praktikum4-crud/database"
)

// GetAllAlumni (existing)
func GetAllAlumni() ([]model.Alumni, error) {
	rows, err := database.DB.Query("SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at FROM alumni")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Alumni
	for rows.Next() {
		var a model.Alumni
		if err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

// GetAlumniByID (existing)
func GetAlumniByID(id int) (*model.Alumni, error) {
	var a model.Alumni
	err := database.DB.QueryRow("SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at FROM alumni WHERE id=$1", id).
		Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// CreateAlumni (existing)
func CreateAlumni(a *model.Alumni) (int, error) {
	var id int
	err := database.DB.QueryRow("INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id",
		a.NIM, a.Nama, a.Jurusan, a.Angkatan, a.TahunLulus, a.Email, a.NoTelepon, a.Alamat).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateAlumni (existing) — note: nim/email not updated here (preserve uniqueness)
func UpdateAlumni(id int, a *model.Alumni) error {
	_, err := database.DB.Exec("UPDATE alumni SET nama=$1, jurusan=$2, angkatan=$3, tahun_lulus=$4, no_telepon=$5, alamat=$6, updated_at=NOW() WHERE id=$7",
		a.Nama, a.Jurusan, a.Angkatan, a.TahunLulus, a.NoTelepon, a.Alamat, id)
	return err
}

// DeleteAlumni (existing)
func DeleteAlumni(id int) error {
	_, err := database.DB.Exec("DELETE FROM alumni WHERE id=$1", id)
	return err
}

// --------------------- New: GetAlumniWithFilter ---------------------
// GetAlumniWithFilter mengambil data alumni dengan pagination, sorting, dan search.
// - page: 1-based
// - limit: jumlah per halaman
// - sortBy: kolom yang diizinkan untuk sorting (divalidasi)
// - order: "asc" atau "desc"
// - search: kata kunci search pada kolom nama atau jurusan (ILIKE)
func GetAlumniWithFilter(page, limit int, sortBy, order, search string) ([]model.Alumni, int, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// whitelist kolom untuk sorting (hindari SQL injection)
	allowed := map[string]bool{
		"id":          true,
		"nim":         true,
		"nama":        true,
		"jurusan":     true,
		"angkatan":    true,
		"tahun_lulus": true,
		"email":       true,
	}

	sortBy = strings.ToLower(sortBy)
	if !allowed[sortBy] {
		sortBy = "id"
	}

	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	// Build WHERE + args
	conditions := []string{}
	args := []interface{}{}

	if search != "" {
		// dua placeholder: nama dan jurusan
		args = append(args, "%"+search+"%")
		args = append(args, "%"+search+"%")
		// positions are 1-indexed; after append len(args) == 2 => placeholders $1 and $2
		n1 := len(args) - 1
		n2 := len(args)
		conditions = append(conditions, fmt.Sprintf("(nama ILIKE $%d OR jurusan ILIKE $%d)", n1, n2))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 1) Hitung total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM alumni %s", whereClause)
	var total int
	if err := database.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 2) Ambil data (tambahkan limit & offset ke args)
	countArgs := len(args) // posisi sebelum append
	args = append(args, limit, offset)
	limitPos := countArgs + 1
	offsetPos := countArgs + 2

	selectQuery := fmt.Sprintf(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortBy, order, limitPos, offsetPos)

	rows, err := database.DB.Query(selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []model.Alumni
	for rows.Next() {
		var a model.Alumni
		if err := rows.Scan(
			&a.ID,
			&a.NIM,
			&a.Nama,
			&a.Jurusan,
			&a.Angkatan,
			&a.TahunLulus,
			&a.Email,
			&a.NoTelepon,
			&a.Alamat,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, a)
	}

	return list, total, nil
}
