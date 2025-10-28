package repository

import (
    "context"
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

// CREATE
// func (r *PekerjaanMongoRepository) Create(ctx context.Context, p *model.PekerjaanMongo) (*model.PekerjaanMongo, error) {
//     now := time.Now()
//     p.CreatedAt = now
//     p.UpdatedAt = now
//     res, err := r.coll.InsertOne(ctx, p)
//     if err != nil {
//         return nil, err
//     }
//     p.ID = res.InsertedID.(primitive.ObjectID)
//     return p, nil
// }

// CREATE
func (r *PekerjaanMongoRepository) Create(ctx context.Context, p *model.PekerjaanMongo) (*model.PekerjaanMongo, error) {
    // Inisialisasi ID dan waktu otomatis
    p.ID = primitive.NewObjectID()
    now := time.Now()
    p.CreatedAt = now
    p.UpdatedAt = now
    p.IsDeleted = nil

    // Insert ke MongoDB
    _, err := r.coll.InsertOne(ctx, p)
    if err != nil {
        return nil, err
    }

    return p, nil
}

// GET ALL (hanya yang belum dihapus)
func (r *PekerjaanMongoRepository) FindAll(ctx context.Context) ([]model.PekerjaanMongo, error) {
    cursor, err := r.coll.Find(ctx, bson.M{"is_deleted": bson.M{"$eq": nil}}, options.Find())
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

// FIND BY ID
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

// UPDATE
func (r *PekerjaanMongoRepository) Update(ctx context.Context, id string, update bson.M) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    update["updated_at"] = time.Now()
    _, err = r.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": update})
    return err
}

// SOFT DELETE
func (r *PekerjaanMongoRepository) SoftDelete(ctx context.Context, id string) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    now := time.Now()
    _, err = r.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"is_deleted": now}})
    return err
}

// RESTORE
func (r *PekerjaanMongoRepository) Restore(ctx context.Context, id string) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    _, err = r.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"is_deleted": nil}})
    return err
}

// HARD DELETE
func (r *PekerjaanMongoRepository) HardDelete(ctx context.Context, id string) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    _, err = r.coll.DeleteOne(ctx, bson.M{"_id": oid})
    return err
}
