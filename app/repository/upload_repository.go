package repository

import (
    "context"
    "praktikum4-crud/app/model"
    "praktikum4-crud/database"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type UploadRepository struct {
    collection *mongo.Collection
}

func NewUploadRepository() *UploadRepository {
    return &UploadRepository{
        collection: database.MongoClient.Database("uts_mongo").Collection("uploads"), // ganti sesuai nama DB-mu
    }
}

func (r *UploadRepository) Create(ctx context.Context, file *model.Upload) error {
    file.UploadedAt = time.Now()
    result, err := r.collection.InsertOne(ctx, file)
    if err != nil {
        return err
    }
    file.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}

func (r *UploadRepository) FindByUser(ctx context.Context, userID int) ([]model.Upload, error) {
    var uploads []model.Upload
    cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    if err := cursor.All(ctx, &uploads); err != nil {
        return nil, err
    }
    return uploads, nil
}

func (r *UploadRepository) FindAll(ctx context.Context) ([]model.Upload, error) {
    var uploads []model.Upload
    cursor, err := r.collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    if err := cursor.All(ctx, &uploads); err != nil {
        return nil, err
    }
    return uploads, nil
}

func (r *UploadRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Upload, error) {
    var upload model.Upload
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&upload)
    if err != nil {
        return nil, err
    }
    return &upload, nil
}

func (r *UploadRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
    _, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
    return err
}
