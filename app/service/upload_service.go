package service

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "praktikum4-crud/app/model"
    "praktikum4-crud/app/repository"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "go.mongodb.org/mongo-driver/bson/primitive"
)


// ==================================
// UPLOAD FILE (Admin/User)
// ==================================
func UploadFile(c *fiber.Ctx) error {
    role := c.Locals("role").(string)
    userID := int(c.Locals("user_id").(float64))
    paramUserID := c.Query("user_id") // admin bisa upload untuk user lain

    var targetUserID int
    if role == "admin" && paramUserID != "" {
        fmt.Sscanf(paramUserID, "%d", &targetUserID)
    } else {
        targetUserID = userID
    }

    // Ambil kategori: foto / sertifikat
    category := c.Query("category")
    if category != "foto" && category != "sertifikat" {
        return c.Status(400).JSON(fiber.Map{"error": "Kategori harus 'foto' atau 'sertifikat'"})
    }

    // Ambil file
    fileHeader, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "File tidak ditemukan"})
    }

    // Validasi ukuran file
    maxSize := int64(1 * 1024 * 1024)
    if category == "sertifikat" {
        maxSize = 2 * 1024 * 1024
    }
    if fileHeader.Size > maxSize {
        return c.Status(400).JSON(fiber.Map{"error": "Ukuran file melebihi batas"})
    }

    // // Validasi tipe file
    // allowed := map[string]bool{}
    // if category == "foto" {
    //     allowed = map[string]bool{
    //         "image/jpeg": true, "image/jpg": true, "image/png": true,
    //     }
    // } else {
    //     allowed = map[string]bool{"application/pdf": true}
    // }

    // contentType := fileHeader.Header.Get("Content-Type")
    // if !allowed[contentType] {
    //     return c.Status(400).JSON(fiber.Map{"error": "Format file tidak diperbolehkan"})
    // }
    
    // Validasi tipe file berdasarkan kategori
contentType := fileHeader.Header.Get("Content-Type")

if category == "foto" {
    if contentType != "image/jpeg" && contentType != "image/jpg" && contentType != "image/png" {
        return c.Status(400).JSON(fiber.Map{"error": "Format foto harus JPG atau PNG"})
    }
} else if category == "sertifikat" {
    if contentType != "application/pdf" {
        return c.Status(400).JSON(fiber.Map{"error": "Format sertifikat harus PDF"})
    }
}


    // Simpan file ke folder uploads/<kategori>
    folder := filepath.Join("uploads", category)
    os.MkdirAll(folder, os.ModePerm)
    ext := filepath.Ext(fileHeader.Filename)
    newName := uuid.New().String() + ext
    path := filepath.Join(folder, newName)

    if err := c.SaveFile(fileHeader, path); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan file"})
    }

    // Simpan metadata ke MongoDB
    repo := repository.NewUploadRepository()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    upload := model.Upload{
        UserID:       targetUserID,
        FileName:     newName,
        OriginalName: fileHeader.Filename,
        FilePath:     path,
        FileSize:     fileHeader.Size,
        FileType:     contentType,
        Category:     category,
    }

    if err := repo.Create(ctx, &upload); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan metadata"})
    }

    return c.Status(201).JSON(fiber.Map{
        "message": "Upload berhasil",
        "data":    upload,
    })
}



// ==================================
// GET ALL FILES
// ==================================
func GetAllUploads(c *fiber.Ctx) error {
    role := c.Locals("role").(string)
    userID := int(c.Locals("user_id").(float64))

    repo := repository.NewUploadRepository()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var uploads []model.Upload
    var err error

    if role == "admin" {
        uploads, err = repo.FindAll(ctx)
    } else {
        uploads, err = repo.FindByUser(ctx, userID)
    }

    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    if len(uploads) == 0 {
        return c.Status(404).JSON(fiber.Map{"error": "Tidak ada file ditemukan"})
    }

    return c.JSON(uploads)
}



// ==================================
// GET FILE BY ID
// ==================================
func GetUploadByID(c *fiber.Ctx) error {
    idStr := c.Params("id")
    objID, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
    }

    role := c.Locals("role").(string)
    userID := int(c.Locals("user_id").(float64))

    repo := repository.NewUploadRepository()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    upload, err := repo.FindByID(ctx, objID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "File tidak ditemukan"})
    }

    if role != "admin" && upload.UserID != userID {
        return c.Status(403).JSON(fiber.Map{"error": "Anda tidak memiliki akses ke file ini"})
    }

    return c.JSON(upload)
}



// ==================================
// DELETE FILE (Admin/User Owner)
// ==================================
func DeleteUpload(c *fiber.Ctx) error {
    idStr := c.Params("id")
    objID, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
    }

    role := c.Locals("role").(string)
    userID := int(c.Locals("user_id").(float64))

    repo := repository.NewUploadRepository()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    upload, err := repo.FindByID(ctx, objID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "File tidak ditemukan"})
    }

    if role != "admin" && upload.UserID != userID {
        return c.Status(403).JSON(fiber.Map{"error": "Anda tidak memiliki izin menghapus file ini"})
    }

    // Hapus file dari sistem
    if err := os.Remove(filepath.Join(upload.FilePath)); err != nil {
        fmt.Println("Peringatan: Gagal menghapus file dari disk:", err)
    }

    // Hapus metadata dari MongoDB
    if err := repo.Delete(ctx, objID); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus metadata dari database"})
    }

    return c.JSON(fiber.Map{"message": "File berhasil dihapus"})
}
