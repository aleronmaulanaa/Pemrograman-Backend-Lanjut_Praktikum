// package config

// import (
// 	"praktikum4-crud/middleware"
// 	"praktikum4-crud/app/service"

// 	"github.com/gofiber/fiber/v2"
// )

// func NewApp() *fiber.App {
// 	app := fiber.New()

// 	// public
// 	app.Post("/login", service.Login)

// 	api := app.Group("/api")

// 	// Alumni (protected) -> GET /api/alumni sekarang pakai pagination/sort/search
// 	alumni := api.Group("/alumni", middleware.JWTMiddleware)
// 	// daftar route:
// 	// GET /api/alumni            -> list (pagination, sort, search) (admin & user)
// 	alumni.Get("/", middleware.RoleMiddleware("admin", "user"), service.GetAlumniWithPagination)
// 	// GET /api/alumni/:id       -> get single (admin & user)
// 	alumni.Get("/:id", middleware.RoleMiddleware("admin", "user"), service.GetAlumniByID)
// 	// POST /api/alumni          -> create (admin)
// 	alumni.Post("/", middleware.RoleMiddleware("admin"), service.CreateAlumni)
// 	// PUT /api/alumni/:id       -> update (admin)
// 	alumni.Put("/:id", middleware.RoleMiddleware("admin"), service.UpdateAlumni)
// 	// DELETE /api/alumni/:id    -> delete (admin)
// 	alumni.Delete("/:id", middleware.RoleMiddleware("admin"), service.DeleteAlumni)

// 	// Pekerjaan (protected) — tetap seperti sebelumnya
// 	pekerjaan := api.Group("/pekerjaan", middleware.JWTMiddleware)
// 	pekerjaan.Get("/", middleware.RoleMiddleware("admin", "user"), service.GetAllPekerjaan)
// 	pekerjaan.Get("/:id", middleware.RoleMiddleware("admin", "user"), service.GetPekerjaanByID)
// 	// pekerjaan.Get("/alumni/:alumni_id", middleware.RoleMiddleware("admin"), service.GetPekerjaanByAlumniID)
// 	pekerjaan.Get("/alumni/:alumni_id", middleware.RoleMiddleware("admin", "user"), service.GetPekerjaanByAlumniID)


// 	pekerjaan.Post("/", middleware.RoleMiddleware("admin"), service.CreatePekerjaan)
// 	pekerjaan.Put("/:id", middleware.RoleMiddleware("admin"), service.UpdatePekerjaan)
// 	// pekerjaan.Delete("/:id", middleware.RoleMiddleware("admin"), service.DeletePekerjaan)

// 	// app.Get("/api/alumni/jumlah-angkatan", service.GetJumlahByAngkatan) ini ditaruh atas

// 	pekerjaan.Delete("/:id", middleware.RoleMiddleware("admin", "user"), service.DeletePekerjaanRBAC)

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

	// Tambahkan route khusus jumlah angkatan DULU, sebelum ":id"
	alumni.Get("/jumlah-angkatan", middleware.RoleMiddleware("admin", "user"), service.GetJumlahByAngkatan)

	alumni.Get("/jumlah-pekerjaan", middleware.RoleMiddleware("admin", "user"), service.GetAlumniDenganDuaPekerjaan)


	// List alumni dengan pagination, sort, dan search
	alumni.Get("/", middleware.RoleMiddleware("admin", "user"), service.GetAlumniWithPagination)

	// Get 1 alumni by ID
	alumni.Get("/:id", middleware.RoleMiddleware("admin", "user"), service.GetAlumniByID)

	// Create alumni (hanya admin)
	alumni.Post("/", middleware.RoleMiddleware("admin"), service.CreateAlumni)

	// Update alumni (hanya admin)
	alumni.Put("/:id", middleware.RoleMiddleware("admin"), service.UpdateAlumni)

	// Delete alumni (hanya admin)
	alumni.Delete("/:id", middleware.RoleMiddleware("admin"), service.DeleteAlumni)

	// ===== PEKERJAAN =====
	pekerjaan := api.Group("/pekerjaan", middleware.JWTMiddleware)

	// List dan detail pekerjaan
	pekerjaan.Get("/", middleware.RoleMiddleware("admin", "user"), service.GetAllPekerjaan)
	pekerjaan.Get("/trash", middleware.RoleMiddleware("admin", "user"), service.GetTrashPekerjaanRBAC)
	pekerjaan.Get("/:id", middleware.RoleMiddleware("admin", "user"), service.GetPekerjaanByID)
	pekerjaan.Get("/alumni/:alumni_id", middleware.RoleMiddleware("admin", "user"), service.GetPekerjaanByAlumniID)

	// Create / Update / Delete pekerjaan
	pekerjaan.Post("/", middleware.RoleMiddleware("admin"), service.CreatePekerjaan)
	pekerjaan.Put("/restore/:id", middleware.RoleMiddleware("admin", "user"), service.RestorePekerjaanRBAC)
	pekerjaan.Put("/:id", middleware.RoleMiddleware("admin"), service.UpdatePekerjaan)
	pekerjaan.Delete("/hard/:id", middleware.RoleMiddleware("admin", "user"), service.HardDeletePekerjaanRBAC)
	pekerjaan.Delete("/:id", middleware.RoleMiddleware("admin", "user"), service.DeletePekerjaanRBAC)

	return app
}
