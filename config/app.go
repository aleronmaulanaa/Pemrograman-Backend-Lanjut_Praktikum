// package config

// import (
// 	"praktikum4-crud/middleware"
// 	"praktikum4-crud/app/service"

// 	"github.com/gofiber/fiber/v2"
// )

// func NewApp() *fiber.App {
// 	app := fiber.New()

// 	// ---------- PUBLIC ROUTES ----------
// 	app.Post("/login", service.Login)

// 	// ---------- PROTECTED ROUTES ----------
// 	api := app.Group("/api")

// 	// ===== ALUMNI =====
// 	alumni := api.Group("/alumni", middleware.JWTMiddleware)

// 	// Tambahkan route khusus jumlah angkatan DULU, sebelum ":id"
// 	alumni.Get("/jumlah-angkatan", middleware.RoleMiddleware("admin", "user"), service.GetJumlahByAngkatan)

// 	alumni.Get("/jumlah-pekerjaan", middleware.RoleMiddleware("admin", "user"), service.GetAlumniDenganDuaPekerjaan)


// 	// List alumni dengan pagination, sort, dan search
// 	alumni.Get("/", middleware.RoleMiddleware("admin", "user"), service.GetAlumniWithPagination)

// 	// Get 1 alumni by ID
// 	alumni.Get("/:id", middleware.RoleMiddleware("admin", "user"), service.GetAlumniByID)

// 	// Create alumni (hanya admin)
// 	alumni.Post("/", middleware.RoleMiddleware("admin"), service.CreateAlumni)

// 	// Update alumni (hanya admin)
// 	alumni.Put("/:id", middleware.RoleMiddleware("admin"), service.UpdateAlumni)

// 	// Delete alumni (hanya admin)
// 	alumni.Delete("/:id", middleware.RoleMiddleware("admin"), service.DeleteAlumni)

// 	// ===== PEKERJAAN =====
// 	pekerjaan := api.Group("/pekerjaan", middleware.JWTMiddleware)

// 	// List dan detail pekerjaan
// 	pekerjaan.Get("/", middleware.RoleMiddleware("admin", "user"), service.GetAllPekerjaan)
// 	pekerjaan.Get("/trash", middleware.RoleMiddleware("admin", "user"), service.GetTrashPekerjaanRBAC)
// 	pekerjaan.Get("/:id", middleware.RoleMiddleware("admin", "user"), service.GetPekerjaanByID)
// 	pekerjaan.Get("/alumni/:alumni_id", middleware.RoleMiddleware("admin", "user"), service.GetPekerjaanByAlumniID)

// 	// Create / Update / Delete pekerjaan
// 	pekerjaan.Post("/", middleware.RoleMiddleware("admin"), service.CreatePekerjaan)
// 	pekerjaan.Put("/restore/:id", middleware.RoleMiddleware("admin", "user"), service.RestorePekerjaanRBAC)
// 	pekerjaan.Put("/:id", middleware.RoleMiddleware("admin"), service.UpdatePekerjaan)
// 	pekerjaan.Delete("/hard/:id", middleware.RoleMiddleware("admin", "user"), service.HardDeletePekerjaanRBAC)
// 	pekerjaan.Delete("/:id", middleware.RoleMiddleware("admin", "user"), service.DeletePekerjaanRBAC)

	
// // ===== PEKERJAAN (MongoDB) =====
// pekerjaanMongo := api.Group("/pekerjaan-mongo", middleware.JWTMiddleware)

// pekerjaanMongo.Get("/", middleware.RoleMiddleware("admin", "user"), service.GetAllPekerjaanMongo)
// pekerjaanMongo.Get("/:id", middleware.RoleMiddleware("admin", "user"), service.GetPekerjaanByIDMongo)
// pekerjaanMongo.Post("/", middleware.RoleMiddleware("admin"), service.CreatePekerjaanMongo)
// pekerjaanMongo.Delete("/:id", middleware.RoleMiddleware("admin"), service.SoftDeletePekerjaanMongo)
// pekerjaanMongo.Put("/restore/:id", middleware.RoleMiddleware("admin"), service.RestorePekerjaanMongo)
// pekerjaanMongo.Delete("/hard/:id", middleware.RoleMiddleware("admin"), service.HardDeletePekerjaanMongo)

// 	return app
// }



package config

import (
	"praktikum4-crud/middleware"
	"praktikum4-crud/app/service"

	"github.com/gofiber/fiber/v2"
)

func NewApp() *fiber.App {
	app := fiber.New()

	// ---------- PUBLIC ROUTES ----------
	app.Post("/login", service.Login)

	// ---------- PROTECTED ROUTES ----------
	api := app.Group("/api")

	// ===== ALUMNI =====
	alumni := api.Group("/alumni", middleware.JWTMiddleware)

	// Statistik
	alumni.Get("/jumlah-angkatan", middleware.RoleMiddleware("admin", "user"), service.GetJumlahByAngkatan)
	alumni.Get("/jumlah-pekerjaan", middleware.RoleMiddleware("admin", "user"), service.GetAlumniDenganDuaPekerjaan)

	// CRUD
	alumni.Get("/", middleware.RoleMiddleware("admin", "user"), service.GetAlumniWithPagination)
	alumni.Get("/:id", middleware.RoleMiddleware("admin", "user"), service.GetAlumniByID)
	alumni.Post("/", middleware.RoleMiddleware("admin"), service.CreateAlumni)
	alumni.Put("/:id", middleware.RoleMiddleware("admin"), service.UpdateAlumni)
	alumni.Delete("/:id", middleware.RoleMiddleware("admin"), service.DeleteAlumni)

	// ===== PEKERJAAN (PostgreSQL) =====
	pekerjaan := api.Group("/pekerjaan", middleware.JWTMiddleware)
	pekerjaan.Get("/", middleware.RoleMiddleware("admin", "user"), service.GetAllPekerjaan)
	pekerjaan.Get("/trash", middleware.RoleMiddleware("admin", "user"), service.GetTrashPekerjaanRBAC)
	pekerjaan.Get("/:id", middleware.RoleMiddleware("admin", "user"), service.GetPekerjaanByID)
	pekerjaan.Get("/alumni/:alumni_id", middleware.RoleMiddleware("admin", "user"), service.GetPekerjaanByAlumniID)
	pekerjaan.Post("/", middleware.RoleMiddleware("admin"), service.CreatePekerjaan)
	pekerjaan.Put("/restore/:id", middleware.RoleMiddleware("admin", "user"), service.RestorePekerjaanRBAC)
	pekerjaan.Put("/:id", middleware.RoleMiddleware("admin"), service.UpdatePekerjaan)
	pekerjaan.Delete("/hard/:id", middleware.RoleMiddleware("admin", "user"), service.HardDeletePekerjaanRBAC)
	pekerjaan.Delete("/:id", middleware.RoleMiddleware("admin", "user"), service.DeletePekerjaanRBAC)

	// ===== PEKERJAAN (MongoDB) =====
	pekerjaanMongo := api.Group("/pekerjaan-mongo", middleware.JWTMiddleware)

	// Hanya admin bisa create
	pekerjaanMongo.Post("/", middleware.RoleMiddleware("admin"), service.CreatePekerjaanMongo)

	// Admin & User bisa melihat pekerjaan aktif (non-deleted)
	pekerjaanMongo.Get("/", middleware.RoleMiddleware("admin", "user"), service.GetAllPekerjaanMongo)
	pekerjaanMongo.Get("/:id", middleware.RoleMiddleware("admin", "user"), service.GetPekerjaanByIDMongo)

	// Soft Delete: admin & owner
	pekerjaanMongo.Delete("/:id", middleware.RoleMiddleware("admin", "user"), service.SoftDeletePekerjaanMongo)

	// Restore: admin & owner, hanya jika is_deleted != null
	pekerjaanMongo.Put("/restore/:id", middleware.RoleMiddleware("admin", "user"), service.RestorePekerjaanMongo)

	// Hard Delete: admin & owner, hanya jika is_deleted != null
	pekerjaanMongo.Delete("/hard/:id", middleware.RoleMiddleware("admin", "user"), service.HardDeletePekerjaanMongo)


	// ===== UPLOAD FILE =====
	upload := api.Group("/upload", middleware.JWTMiddleware)
	upload.Post("/", middleware.RoleMiddleware("admin", "user"), service.UploadFile)
	upload.Get("/", middleware.RoleMiddleware("admin", "user"), service.GetAllUploads)
	upload.Get("/:id", middleware.RoleMiddleware("admin", "user"), service.GetUploadByID)
	upload.Delete("/:id", middleware.RoleMiddleware("admin", "user"), service.DeleteUpload)
	
	return app
}