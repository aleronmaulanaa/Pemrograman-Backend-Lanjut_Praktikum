// package service

// import (
// 	"strconv"

// 	"praktikum4-crud/app/model"
// 	"praktikum4-crud/app/repository"

// 	"github.com/gofiber/fiber/v2"
// )

// // ---------------- existing handlers (CRUD) using repository ----------------

// func GetAllAlumni(c *fiber.Ctx) error {
// 	list, err := repository.GetAllAlumni()
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data", "detail": err.Error()})
// 	}
// 	return c.JSON(list)
// }

// func GetAlumniByID(c *fiber.Ctx) error {
// 	idStr := c.Params("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
// 	}
// 	a, err := repository.GetAlumniByID(id)
// 	if err != nil {
// 		return c.Status(404).JSON(fiber.Map{"error": "Alumni tidak ditemukan"})
// 	}
// 	return c.JSON(a)
// }

// func CreateAlumni(c *fiber.Ctx) error {
// 	var a model.Alumni
// 	if err := c.BodyParser(&a); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Invalid input", "detail": err.Error()})
// 	}

// 	id, err := repository.CreateAlumni(&a)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data", "detail": err.Error()})
// 	}
// 	a.ID = id
// 	return c.Status(201).JSON(a)
// }

// func UpdateAlumni(c *fiber.Ctx) error {
// 	idStr := c.Params("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
// 	}
// 	var a model.Alumni
// 	if err := c.BodyParser(&a); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Invalid input", "detail": err.Error()})
// 	}

// 	if err := repository.UpdateAlumni(id, &a); err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "Gagal update data", "detail": err.Error()})
// 	}
// 	return c.JSON(fiber.Map{"message": "Alumni berhasil diupdate"})
// }

// func DeleteAlumni(c *fiber.Ctx) error {
// 	idStr := c.Params("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
// 	}
// 	if err := repository.DeleteAlumni(id); err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus data", "detail": err.Error()})
// 	}
// 	return c.JSON(fiber.Map{"message": "Alumni berhasil dihapus"})
// }

// // ---------------- New handler: Pagination / Sorting / Search ----------------

// // GetAlumniWithPagination -> handler GET /api/alumni
// // Query params:
// //   page (int, default 1), limit (int, default 10), sort (kolom), order (asc|desc), search (string)
// func GetAlumniWithPagination(c *fiber.Ctx) error {
// 	// parse query params
// 	pageStr := c.Query("page", "1")
// 	limitStr := c.Query("limit", "10")
// 	sortBy := c.Query("sort", "id")
// 	order := c.Query("order", "desc")
// 	search := c.Query("search", "")

// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil || page < 1 {
// 		page = 1
// 	}
// 	limit, err := strconv.Atoi(limitStr)
// 	if err != nil || limit < 1 {
// 		limit = 10
// 	}
// 	if limit > 100 {
// 		limit = 100
// 	}

// 	items, total, err := repository.GetAlumniWithFilter(page, limit, sortBy, order, search)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data", "detail": err.Error()})
// 	}

// 	totalPages := 0
// 	if total > 0 {
// 		totalPages = (total + limit - 1) / limit
// 	}

// 	return c.JSON(fiber.Map{
// 		"page":        page,
// 		"limit":       limit,
// 		"total":       total,
// 		"total_pages": totalPages,
// 		"data":        items,
// 	})
// }

// // GetJumlahByAngkatan -> hitung jumlah alumni berdasarkan angkatan
// func GetJumlahByAngkatan(c *fiber.Ctx) error {
//     rows, err := database.DB.Query("SELECT angkatan, COUNT(*) as jumlah FROM alumni GROUP BY angkatan")
//     if err != nil {
//         return c.Status(500).JSON(fiber.Map{
//             "error": "Gagal mengambil data",
//         })
//     }
//     defer rows.Close()

//     type Result struct {
//         Angkatan int `json:"angkatan"`
//         Jumlah   int `json:"jumlah"`
//     }

//     var results []Result

//     for rows.Next() {
//         var r Result
//         if err := rows.Scan(&r.Angkatan, &r.Jumlah); err != nil {
//             return c.Status(500).JSON(fiber.Map{"error": "Gagal baca data"})
//         }
//         results = append(results, r)
//     }

//     return c.JSON(results)
// }

package service

import (
	"strconv"

	"praktikum4-crud/app/model"
	"praktikum4-crud/app/repository"
	"praktikum4-crud/database" // ✅ Tambahkan ini agar bisa akses database.DB

	"github.com/gofiber/fiber/v2"
)

// ---------------- existing handlers (CRUD) using repository ----------------

func GetAllAlumni(c *fiber.Ctx) error {
	list, err := repository.GetAllAlumni()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data", "detail": err.Error()})
	}
	return c.JSON(list)
}

func GetAlumniByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}
	a, err := repository.GetAlumniByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Alumni tidak ditemukan"})
	}
	return c.JSON(a)
}

func CreateAlumni(c *fiber.Ctx) error {
	var a model.Alumni
	if err := c.BodyParser(&a); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input", "detail": err.Error()})
	}

	id, err := repository.CreateAlumni(&a)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data", "detail": err.Error()})
	}
	a.ID = id
	return c.Status(201).JSON(a)
}

func UpdateAlumni(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}
	var a model.Alumni
	if err := c.BodyParser(&a); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input", "detail": err.Error()})
	}

	if err := repository.UpdateAlumni(id, &a); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update data", "detail": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Alumni berhasil diupdate"})
}

func DeleteAlumni(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}
	if err := repository.DeleteAlumni(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus data", "detail": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Alumni berhasil dihapus"})
}

// ---------------- New handler: Pagination / Sorting / Search ----------------

func GetAlumniWithPagination(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")
	sortBy := c.Query("sort", "id")
	order := c.Query("order", "desc")
	search := c.Query("search", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	items, total, err := repository.GetAlumniWithFilter(page, limit, sortBy, order, search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data", "detail": err.Error()})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return c.JSON(fiber.Map{
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
		"data":        items,
	})
}

// ---------------- New handler: Jumlah Alumni per Angkatan ----------------

func GetJumlahByAngkatan(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT angkatan, COUNT(*) as jumlah FROM alumni GROUP BY angkatan")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil data",
			"detail": err.Error(),
		})
	}
	defer rows.Close()

	type Result struct {
		Angkatan int `json:"angkatan"`
		Jumlah   int `json:"jumlah"`
	}

	var results []Result

	for rows.Next() {
		var r Result
		if err := rows.Scan(&r.Angkatan, &r.Jumlah); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal baca data", "detail": err.Error()})
		}
		results = append(results, r)
	}

	return c.JSON(results)
}

// GetAlumniDenganDuaPekerjaan -> menampilkan alumni dengan >= 2 pekerjaan
func GetAlumniDenganDuaPekerjaan(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`
		SELECT a.nama, COUNT(p.id) AS jumlah_pekerjaan
		FROM alumni a
		JOIN pekerjaan_alumni p ON a.id = p.alumni_id
		GROUP BY a.id, a.nama
		HAVING COUNT(p.id) >= 2
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}
	defer rows.Close()

	type Result struct {
		Nama             string `json:"nama"`
		JumlahPekerjaan  int    `json:"jumlah_pekerjaan"`
	}

	var results []Result

	for rows.Next() {
		var r Result
		if err := rows.Scan(&r.Nama, &r.JumlahPekerjaan); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal baca data"})
		}
		results = append(results, r)
	}

	return c.JSON(results)
}
