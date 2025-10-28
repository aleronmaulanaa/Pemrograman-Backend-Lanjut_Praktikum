package service

import (
	"context"
	"time"

	"praktikum4-crud/app/model"
	"praktikum4-crud/app/repository"

	"github.com/gofiber/fiber/v2"
)

// ===========================
// CREATE
// ===========================

// func CreatePekerjaanMongo(c *fiber.Ctx) error {
// 	var p model.PekerjaanMongo
// 	if err := c.BodyParser(&p); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	repo := repository.NewPekerjaanMongoRepository() // inisialisasi repository di sini

// 	created, err := repo.Create(ctx, &p)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return c.Status(201).JSON(created)
// }
func CreatePekerjaanMongo(c *fiber.Ctx) error {
    var p model.PekerjaanMongo
    if err := c.BodyParser(&p); err != nil {
        // Tambahkan log error agar tahu kesalahannya
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid input",
            "detail": err.Error(),
        })
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
// GET ALL
// ===========================
func GetAllPekerjaanMongo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repo := repository.NewPekerjaanMongoRepository()

	data, err := repo.FindAll(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if len(data) == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "Belum ada data pekerjaan"})
	}

	return c.JSON(data)
}

// ===========================
// GET BY ID
// ===========================
func GetPekerjaanByIDMongo(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repo := repository.NewPekerjaanMongoRepository()

	data, err := repo.FindByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
	}

	return c.JSON(data)
}

// ===========================
// SOFT DELETE
// ===========================
func SoftDeletePekerjaanMongo(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repo := repository.NewPekerjaanMongoRepository()

	if err := repo.SoftDelete(ctx, id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil dihapus (soft delete)"})
}

// ===========================
// RESTORE
// ===========================
func RestorePekerjaanMongo(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repo := repository.NewPekerjaanMongoRepository()

	if err := repo.Restore(ctx, id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil direstore"})
}

// ===========================
// HARD DELETE
// ===========================
func HardDeletePekerjaanMongo(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repo := repository.NewPekerjaanMongoRepository()

	if err := repo.HardDelete(ctx, id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Pekerjaan dihapus permanen"})
}
