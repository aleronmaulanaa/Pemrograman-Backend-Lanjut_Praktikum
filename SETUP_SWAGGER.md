// Swagger Route Setup Guide
// ===========================

// Jika belum ada, tambahkan di file main.go atau config/app.go:

package config

import (
	_ "praktikum4-crud/docs" // Import the generated docs
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupSwagger(app *fiber.App) {
	// Swagger endpoint
	app.Get("/swagger/*", swagger.HandlerDefault)
	// Alternative custom path
	// app.Get("/api/docs/*", swagger.New(swagger.Config{
	// 	URL: "http://localhost:3000/swagger/doc.json",
	// }))
}

// ===========================
// LANGKAH-LANGKAH IMPLEMENTASI
// ===========================

// 1. INSTALL SWAG CLI
//    go install github.com/swaggo/swag/cmd/swag@latest

// 2. INSTALL SWAG FIBER
//    go get -u github.com/gofiber/swagger
//    go get -u github.com/swaggo/files
//    go get -u github.com/swaggo/gin-swagger

// 3. GENERATE SWAGGER DOCUMENTATION
//    swag init
//    
//    Perintah ini akan membuat folder 'docs' dengan file:
//    - docs.go
//    - swagger.json
//    - swagger.yaml

// 4. IMPORT DOCS DI MAIN.GO
//    Tambahkan di atas main.go:
//    _ "praktikum4-crud/docs"

// 5. SETUP ROUTE SWAGGER
//    Tambahkan di config.go atau main.go:
//    app.Get("/swagger/*", swagger.HandlerDefault)

// 6. JALANKAN PROJECT
//    go run main.go

// 7. AKSES SWAGGER UI
//    http://localhost:3000/swagger/index.html

// ===========================
// STRUKTUR ANOTASI SWAGGER
// ===========================

// Basic Endpoint Example:
// @Summary [Deskripsi singkat]
// @Description [Deskripsi detail]
// @Tags [Kategori endpoint]
// @Security Bearer [atau security scheme lainnya]
// @Accept json
// @Produce json
// @Param [nama] [in] [tipe] [required/optional] "[deskripsi]"
// @Success [kode] {[tipe]} [model] "[deskripsi]"
// @Failure [kode] {object} map[string]interface{} "[deskripsi]"
// @Router /path [method]

// ===========================
// CONTOH ANOTASI LENGKAP
// ===========================

/*
// @Summary Create Alumni
// @Description Create a new alumni record
// @Tags Alumni
// @Security Bearer
// @Accept json
// @Produce json
// @Param alumni body model.Alumni true "Alumni data"
// @Success 201 {object} model.Alumni
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Server Error"
// @Router /alumni [post]
func CreateAlumni(c *fiber.Ctx) error {
    // implementation...
}
*/

// ===========================
// REFERENSI PARAMETER
// ===========================

// @Param Parameter Locations:
// - path: URL path parameter (e.g., /users/{id})
// - query: URL query string (e.g., ?page=1&limit=10)
// - header: Request header
// - formData: Form data (untuk file upload)
// - body: Request body

// @Success / @Failure Formats:
// {object} model.StructName - Single object
// {array} model.StructName - Array of objects
// map[string]interface{} - Generic response object

// Security Schemes:
// @Security Bearer - JWT Token
// @Security ApiKey - API Key
// @Security OAuth2 - OAuth 2.0

// ===========================
// TIPS DAN BEST PRACTICES
// ===========================

// 1. Selalu gunakan tag yang konsisten untuk grouping endpoints
// 2. Pastikan router path dan @Router anotasi sama
// 3. Gunakan deskripsi yang jelas dan detail
// 4. Tentukan failure responses dengan status code yang tepat
// 5. Update docs setiap kali menambah/mengubah endpoint dengan: swag init
// 6. Sertakan contoh request/response di dokumentasi
// 7. Gunakan model struct untuk request/response yang kompleks

// ===========================
// TROUBLESHOOTING
// ===========================

// Problem: Swagger docs tidak muncul
// Solution: 
// - Pastikan sudah menjalankan: swag init
// - Periksa import: _ "praktikum4-crud/docs"
// - Pastikan route /swagger/* terdaftar

// Problem: Model tidak muncul di Swagger
// Solution:
// - Pastikan model struct memiliki JSON tags
// - Contoh: type Alumni struct {
//            ID int `json:"id"`
//            Nama string `json:"nama"`
//          }

// Problem: Parameter tidak terdeteksi
// Solution:
// - Gunakan format yang tepat: @Param name [in] [type] [required] "description"
// - Contoh: @Param id path int true "Alumni ID"

// ===========================
// SWAGGER UI FEATURES
// ===========================

// Di Swagger UI, Anda dapat:
// 1. Lihat semua endpoint yang tersedia
// 2. Coba request endpoint langsung dari UI
// 3. Lihat response contoh
// 4. Download spec dalam format JSON/YAML
// 5. Share dokumentasi dengan tim
