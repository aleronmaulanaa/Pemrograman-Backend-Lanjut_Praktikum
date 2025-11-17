package service

import (
	"strconv"

	"praktikum4-crud/app/model"
	"praktikum4-crud/app/repository"
	"praktikum4-crud/database" // ✅ Tambahkan ini agar bisa akses database.DB

	"github.com/gofiber/fiber/v2"
)

// ---------------- existing handlers (CRUD) using repository ----------------

// @Summary Get All Alumni
// @Description Get list of all alumni with their details
// @Tags Alumni
// @Security Bearer
// @Success 200 {array} model.Alumni
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /alumni [get]
func GetAllAlumni(c *fiber.Ctx) error {
	list, err := repository.GetAllAlumni()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data", "detail": err.Error()})
	}
	return c.JSON(list)
}

// @Summary Get Alumni by ID
// @Description Get alumni details by their ID
// @Tags Alumni
// @Security Bearer
// @Param id path int true "Alumni ID"
// @Success 200 {object} model.Alumni
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 404 {object} map[string]interface{} "Alumni not found"
// @Router /alumni/{id} [get]
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

// @Summary Create Alumni
// @Description Create a new alumni record
// @Tags Alumni
// @Security Bearer
// @Accept json
// @Produce json
// @Param alumni body model.Alumni true "Alumni data"
// @Success 201 {object} model.Alumni
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /alumni [post]
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

// @Summary Update Alumni
// @Description Update an existing alumni record
// @Tags Alumni
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Alumni ID"
// @Param alumni body model.Alumni true "Alumni data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /alumni/{id} [put]
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

// @Summary Delete Alumni
// @Description Delete an alumni record
// @Tags Alumni
// @Security Bearer
// @Param id path int true "Alumni ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /alumni/{id} [delete]
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

// @Summary Get Alumni with Pagination, Sort, and Search
// @Description Get list of alumni with pagination, sorting, and search capabilities
// @Tags Alumni
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param sort query string false "Sort by field" default(id)
// @Param order query string false "Sort order (asc/desc)" default(desc)
// @Param search query string false "Search by name"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /alumni [get]
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

// @Summary Get Jumlah Alumni by Angkatan
// @Description Get the count of alumni grouped by year (angkatan)
// @Tags Alumni
// @Security Bearer
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /alumni/jumlah-angkatan [get]
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

// @Summary Get Alumni dengan Dua atau Lebih Pekerjaan
// @Description Get list of alumni who have two or more jobs
// @Tags Alumni
// @Security Bearer
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /alumni/jumlah-pekerjaan [get]
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
