package repositories

import (
    "context"
    "fmt"
    "log"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    cursosDAO "proyecto_arqui_soft_2/cursos-api/dao"
)

type MongoConfig struct {
    Host       string
    Port       string
    Username   string
    Password   string
    Database   string
    Collection string
}

type Mongo struct {
    client     *mongo.Client
    database   string
    collection string
}

const (
    connectionURI = "mongodb://%s:%s"
)

func NewMongo(config MongoConfig) Mongo {
    credentials := options.Credential{
        Username: config.Username,
        Password: config.Password,
    }

    ctx := context.Background()
    uri := fmt.Sprintf(connectionURI, config.Host, config.Port)
    cfg := options.Client().ApplyURI(uri).SetAuth(credentials)

    client, err := mongo.Connect(ctx, cfg)
    if err != nil {
        log.Panicf("error connecting to MongoDB: %v", err)
    }

    return Mongo{
        client:     client,
        database:   config.Database,
        collection: config.Collection,
    }
}

func (repository Mongo) GetCursoByID(ctx context.Context, id string) (cursosDAO.Curso, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return cursosDAO.Curso{}, fmt.Errorf("error converting id to mongo ID: %w", err)
    }
    result := repository.client.Database(repository.database).Collection(repository.collection).FindOne(ctx, bson.M{"_id": objectID})
    if result.Err() != nil {
        return cursosDAO.Curso{}, fmt.Errorf("error finding document: %w", result.Err())
    }

    var curso cursosDAO.Curso
    if err := result.Decode(&curso); err != nil {
        return cursosDAO.Curso{}, fmt.Errorf("error decoding result: %w", err)
    }
    return curso, nil
}

func (repository Mongo) Create(ctx context.Context, curso cursosDAO.Curso) (string, error) {
    result, err := repository.client.Database(repository.database).Collection(repository.collection).InsertOne(ctx, curso)
    if err != nil {
        return "", fmt.Errorf("error creating document: %w", err)
    }

    objectID, ok := result.InsertedID.(primitive.ObjectID)
    if !ok {
        return "", fmt.Errorf("error converting mongo ID to object ID")
    }
    return objectID.Hex(), nil
}

func (repository Mongo) Update(ctx context.Context, cursoID string, updateData bson.M) error {
    objectID, err := primitive.ObjectIDFromHex(cursoID)
    if err != nil {
        return fmt.Errorf("error converting id to mongo ID: %w", err)
    }

    if len(updateData) == 0 {
        return fmt.Errorf("no fields to update for course ID %s", cursoID)
    }

    filter := bson.M{"_id": objectID}
    result, err := repository.client.Database(repository.database).Collection(repository.collection).UpdateOne(ctx, filter, bson.M{"$set": updateData})
    if err != nil {
        return fmt.Errorf("error updating document: %w", err)
    }
    if result.MatchedCount == 0 {
        return fmt.Errorf("no document found with ID %s", cursoID)
    }

    return nil
}

func (repository Mongo) Delete(ctx context.Context, id string) error {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return fmt.Errorf("error converting id to mongo ID: %w", err)
    }

    filter := bson.M{"_id": objectID}
    result, err := repository.client.Database(repository.database).Collection(repository.collection).DeleteOne(ctx, filter)
    if err != nil {
        return fmt.Errorf("error deleting document: %w", err)
    }
    if result.DeletedCount == 0 {
        return fmt.Errorf("no document found with ID %s", id)
    }

    return nil
}


func (repository Mongo) TestConnection(ctx context.Context) error {
    if err := repository.client.Ping(ctx, nil); err != nil {
        return fmt.Errorf("error al hacer ping a MongoDB: %w", err)
    }
    return nil
}

