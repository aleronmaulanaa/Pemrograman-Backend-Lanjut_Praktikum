package service

import (
	"context"
	"time"

	"praktikum4-crud/app/model"
	"praktikum4-crud/app/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// ===========================
// CREATE (Admin Only)
// ===========================
// @Summary Create Pekerjaan (MongoDB)
// @Description Create a new job record in MongoDB (Admin only)
// @Tags Pekerjaan-MongoDB
// @Security Bearer
// @Accept json
// @Produce json
// @Param pekerjaan body model.PekerjaanMongo true "Pekerjaan data"
// @Success 201 {object} model.PekerjaanMongo
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 403 {object} map[string]interface{} "Only admin can create"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /pekerjaan-mongo [post]
func CreatePekerjaanMongo(c *fiber.Ctx) error {
    role, _ := c.Locals("role").(string)
    if role != "admin" {
        return c.Status(403).JSON(fiber.Map{"error": "Hanya admin yang bisa menambah pekerjaan"})
    }

    var p model.PekerjaanMongo
    if err := c.BodyParser(&p); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    created, err := repository.NewPekerjaanMongoRepository().Create(ctx, &p)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.Status(201).JSON(created)
}

// ===========================
// GET ALL (RBAC)
// ===========================
// @Summary Get All Pekerjaan (MongoDB)
// @Description Get list of active jobs (Admin sees all, User sees only their own)
// @Tags Pekerjaan-MongoDB
// @Security Bearer
// @Success 200 {array} model.PekerjaanMongo
// @Failure 404 {object} map[string]interface{} "No active jobs found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /pekerjaan-mongo [get]
func GetAllPekerjaanMongo(c *fiber.Ctx) error {
    role, _ := c.Locals("role").(string)
    userIDFloat := c.Locals("user_id").(float64)
    userID := int(userIDFloat)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    repo := repository.NewPekerjaanMongoRepository()
    var filter bson.M

    // Admin melihat semua data, user hanya miliknya sendiri
    if role == "admin" {
        filter = bson.M{"is_deleted": nil}
    } else {
        filter = bson.M{"is_deleted": nil, "alumni_id": userID}
    }

    data, err := repo.FindAll(ctx, filter)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    if len(data) == 0 {
        return c.Status(404).JSON(fiber.Map{"error": "Tidak ada pekerjaan aktif"})
    }

    return c.JSON(data)
}

// ===========================
// GET BY ID (Hanya data aktif)
// ===========================
// @Summary Get Pekerjaan by ID (MongoDB)
// @Description Get details of an active job
// @Tags Pekerjaan-MongoDB
// @Security Bearer
// @Param id path string true "Pekerjaan ID"
// @Success 200 {object} model.PekerjaanMongo
// @Failure 400 {object} map[string]interface{} "Job is deleted"
// @Failure 404 {object} map[string]interface{} "Job not found"
// @Router /pekerjaan-mongo/{id} [get]
func GetPekerjaanByIDMongo(c *fiber.Ctx) error {
    id := c.Params("id")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    repo := repository.NewPekerjaanMongoRepository()
    data, err := repo.FindByID(ctx, id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
    }

    if data.IsDeleted != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Pekerjaan sudah dihapus"})
    }

    return c.JSON(data)
}

// ===========================
// SOFT DELETE (Admin/User Owner)
// ===========================
// @Summary Soft Delete Pekerjaan (MongoDB)
// @Description Soft delete a job (Admin or job owner)
// @Tags Pekerjaan-MongoDB
// @Security Bearer
// @Param id path string true "Pekerjaan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Job not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /pekerjaan-mongo/{id} [delete]
func SoftDeletePekerjaanMongo(c *fiber.Ctx) error {
    id := c.Params("id")
    role, _ := c.Locals("role").(string)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    repo := repository.NewPekerjaanMongoRepository()

    data, err := repo.FindByID(ctx, id)
    if err != nil || data == nil {
        return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
    }

    // Cegah soft delete jika sudah dihapus
    if data.IsDeleted != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Data sudah dihapus sebelumnya"})
    }

    // Validasi kepemilikan jika bukan admin
    if role != "admin" {
        userID := int(c.Locals("user_id").(float64))
        if data.AlumniID != userID {
            return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Anda hanya bisa menghapus pekerjaan Anda sendiri"})
        }
    }

    if err := repo.SoftDelete(ctx, id); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Pekerjaan berhasil dihapus (soft delete)"})
}

// ===========================
// RESTORE (Admin/User Owner)
// ===========================
// @Summary Restore Pekerjaan (MongoDB)
// @Description Restore a soft-deleted job (Admin or job owner)
// @Tags Pekerjaan-MongoDB
// @Security Bearer
// @Param id path string true "Pekerjaan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Not soft-deleted"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Job not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /pekerjaan-mongo/restore/{id} [put]
func RestorePekerjaanMongo(c *fiber.Ctx) error {
    id := c.Params("id")
    role, _ := c.Locals("role").(string)
    userID := int(c.Locals("user_id").(float64))

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    repo := repository.NewPekerjaanMongoRepository()

    data, err := repo.FindByID(ctx, id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
    }

    // Pastikan pekerjaan memang sudah dihapus
    if data.IsDeleted == nil {
        return c.Status(400).JSON(fiber.Map{"error": "Data belum dihapus (soft delete) sehingga tidak bisa direstore"})
    }

    // ===== Logika baru (disamakan dengan PostgreSQL) =====
    if role == "admin" {
        // Admin boleh restore meskipun tidak ada relasi ke alumni
        if err := repo.Restore(ctx, id); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(fiber.Map{"message": "Pekerjaan berhasil direstore oleh admin"})
    }

    // Jika user biasa dan tidak ada relasi
    if data.AlumniID == 0 {
        return c.Status(403).JSON(fiber.Map{"error": "Data ini tidak memiliki relasi ke akun mana pun"})
    }

    // Jika user bukan pemilik
    if data.AlumniID != userID {
        return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Anda hanya bisa me-restore pekerjaan Anda sendiri"})
    }

    // Restore oleh user pemilik
    if err := repo.Restore(ctx, id); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Pekerjaan berhasil direstore oleh user"})
}

// ===========================
// HARD DELETE (Admin/User Owner)
// ===========================
// @Summary Hard Delete Pekerjaan (MongoDB)
// @Description Permanently delete a soft-deleted job (Admin or job owner)
// @Tags Pekerjaan-MongoDB
// @Security Bearer
// @Param id path string true "Pekerjaan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Not soft-deleted"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Job not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /pekerjaan-mongo/hard/{id} [delete]
func HardDeletePekerjaanMongo(c *fiber.Ctx) error {
    id := c.Params("id")
    role, _ := c.Locals("role").(string)
    userID := int(c.Locals("user_id").(float64))

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    repo := repository.NewPekerjaanMongoRepository()

    data, err := repo.FindByID(ctx, id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
    }

    // Pastikan sudah di-soft-delete dulu
    if data.IsDeleted == nil {
        return c.Status(400).JSON(fiber.Map{"error": "Data belum dihapus (soft delete) sehingga tidak bisa dihapus permanen"})
    }

    // ===== Logika baru (disamakan dengan PostgreSQL) =====
    if role != "admin" {
        // Jika pekerjaan tidak punya relasi alumni_id
        if data.AlumniID == 0 {
            return c.Status(403).JSON(fiber.Map{"error": "Data ini tidak memiliki relasi ke akun mana pun"})
        }

        // Jika user bukan pemilik
        if data.AlumniID != userID {
            return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Anda hanya bisa menghapus pekerjaan Anda sendiri"})
        }
    }

    if err := repo.HardDelete(ctx, id); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Pekerjaan berhasil dihapus permanen"})
}
