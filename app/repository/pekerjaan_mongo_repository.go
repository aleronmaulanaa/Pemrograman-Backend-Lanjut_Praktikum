// package repository

// import (
//     "context"
//     "time"

//     "praktikum4-crud/app/model"
//     "praktikum4-crud/database"

//     "go.mongodb.org/mongo-driver/bson"
//     "go.mongodb.org/mongo-driver/bson/primitive"
//     "go.mongodb.org/mongo-driver/mongo"
//     "go.mongodb.org/mongo-driver/mongo/options"
// )

// type PekerjaanMongoRepository struct {
//     coll *mongo.Collection
// }

// func NewPekerjaanMongoRepository() *PekerjaanMongoRepository {
//     return &PekerjaanMongoRepository{
//         coll: database.MongoDB.Collection("pekerjaan_alumni"),
//     }
// }

// // CREATE
// func (r *PekerjaanMongoRepository) Create(ctx context.Context, p *model.PekerjaanMongo) (*model.PekerjaanMongo, error) {
//     // Inisialisasi ID dan waktu otomatis
//     p.ID = primitive.NewObjectID()
//     now := time.Now()
//     p.CreatedAt = now
//     p.UpdatedAt = now
//     p.IsDeleted = nil

//     // Insert ke MongoDB
//     _, err := r.coll.InsertOne(ctx, p)
//     if err != nil {
//         return nil, err
//     }

//     return p, nil
// }

// // GET ALL (hanya yang belum dihapus)
// func (r *PekerjaanMongoRepository) FindAll(ctx context.Context) ([]model.PekerjaanMongo, error) {
//     cursor, err := r.coll.Find(ctx, bson.M{"is_deleted": bson.M{"$eq": nil}}, options.Find())
//     if err != nil {
//         return nil, err
//     }
//     defer cursor.Close(ctx)

//     var list []model.PekerjaanMongo
//     if err := cursor.All(ctx, &list); err != nil {
//         return nil, err
//     }
//     return list, nil
// }

// // FIND BY ID
// func (r *PekerjaanMongoRepository) FindByID(ctx context.Context, id string) (*model.PekerjaanMongo, error) {
//     oid, err := primitive.ObjectIDFromHex(id)
//     if err != nil {
//         return nil, err
//     }

//     var p model.PekerjaanMongo
//     if err := r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&p); err != nil {
//         return nil, err
//     }
//     return &p, nil
// }

// // UPDATE
// func (r *PekerjaanMongoRepository) Update(ctx context.Context, id string, update bson.M) error {
//     oid, err := primitive.ObjectIDFromHex(id)
//     if err != nil {
//         return err
//     }
//     update["updated_at"] = time.Now()
//     _, err = r.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": update})
//     return err
// }

// // SOFT DELETE
// func (r *PekerjaanMongoRepository) SoftDelete(ctx context.Context, id string) error {
//     oid, err := primitive.ObjectIDFromHex(id)
//     if err != nil {
//         return err
//     }
//     now := time.Now()
//     _, err = r.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"is_deleted": now}})
//     return err
// }

// // RESTORE
// func (r *PekerjaanMongoRepository) Restore(ctx context.Context, id string) error {
//     oid, err := primitive.ObjectIDFromHex(id)
//     if err != nil {
//         return err
//     }
//     _, err = r.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"is_deleted": nil}})
//     return err
// }

// // HARD DELETE
// func (r *PekerjaanMongoRepository) HardDelete(ctx context.Context, id string) error {
//     oid, err := primitive.ObjectIDFromHex(id)
//     if err != nil {
//         return err
//     }
//     _, err = r.coll.DeleteOne(ctx, bson.M{"_id": oid})
//     return err
// }



package repository

import (
    "context"
    "errors"
    "time"

    "praktikum4-crud/app/model"
    "praktikum4-crud/database"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type PekerjaanMongoRepository struct {
    coll *mongo.Collection
}

func NewPekerjaanMongoRepository() *PekerjaanMongoRepository {
    return &PekerjaanMongoRepository{
        coll: database.MongoDB.Collection("pekerjaan_alumni"),
    }
}

// ==========================
// CREATE (admin only)
// ==========================
func (r *PekerjaanMongoRepository) Create(ctx context.Context, p *model.PekerjaanMongo) (*model.PekerjaanMongo, error) {
    p.ID = primitive.NewObjectID()
    now := time.Now()
    p.CreatedAt = now
    p.UpdatedAt = now
    p.IsDeleted = nil

    _, err := r.coll.InsertOne(ctx, p)
    if err != nil {
        return nil, err
    }
    return p, nil
}

// ==========================
// FIND ALL
// ==========================
func (r *PekerjaanMongoRepository) FindAll(ctx context.Context, filter bson.M) ([]model.PekerjaanMongo, error) {
    opts := options.Find().SetSort(bson.M{"created_at": -1})
    cursor, err := r.coll.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var list []model.PekerjaanMongo
    if err := cursor.All(ctx, &list); err != nil {
        return nil, err
    }
    return list, nil
}

// ==========================
// FIND BY ID
// ==========================
func (r *PekerjaanMongoRepository) FindByID(ctx context.Context, id string) (*model.PekerjaanMongo, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var p model.PekerjaanMongo
    if err := r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&p); err != nil {
        return nil, err
    }
    return &p, nil
}

// ==========================
// SOFT DELETE
// ==========================
func (r *PekerjaanMongoRepository) SoftDelete(ctx context.Context, id string) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    now := time.Now()
    res, err := r.coll.UpdateOne(ctx, bson.M{"_id": oid, "is_deleted": bson.M{"$eq": nil}}, bson.M{"$set": bson.M{"is_deleted": now}})
    if err != nil {
        return err
    }
    if res.MatchedCount == 0 {
        return errors.New("data tidak ditemukan atau sudah dihapus")
    }
    return nil
}

// ==========================
// RESTORE
// ==========================
func (r *PekerjaanMongoRepository) Restore(ctx context.Context, id string) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    // Pastikan hanya bisa restore jika sudah dihapus
    res, err := r.coll.UpdateOne(ctx,
        bson.M{"_id": oid, "is_deleted": bson.M{"$ne": nil}},
        bson.M{"$set": bson.M{"is_deleted": nil, "updated_at": time.Now()}},
    )
    if err != nil {
        return err
    }
    if res.MatchedCount == 0 {
        return errors.New("data belum dihapus atau tidak ditemukan")
    }
    return nil
}

// ==========================
// HARD DELETE
// ==========================
func (r *PekerjaanMongoRepository) HardDelete(ctx context.Context, id string) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    res, err := r.coll.DeleteOne(ctx, bson.M{"_id": oid, "is_deleted": bson.M{"$ne": nil}})
    if err != nil {
        return err
    }
    if res.DeletedCount == 0 {
        return errors.New("data belum dihapus (soft delete) atau tidak ditemukan")
    }
    return nil
}
