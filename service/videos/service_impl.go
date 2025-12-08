package videos

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type videoServiceImpl struct {
	mongo *mongo.Client
}

func newVideoService(mongo *mongo.Client) videoService {
	return &videoServiceImpl{mongo: mongo}
}

func (s *videoServiceImpl) ListVideos() ([]video, error) {
	collection := s.mongo.Database("videoDB").Collection("videos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		slog.Error("failed to create DB cursor for videos", "error", err)
		return nil, err
	}

	var videos []video
	if err = cursor.All(ctx, &videos); err != nil {
		slog.Error("failed to fetch videos", "error", err)
		return nil, err
	}

	return videos, nil
}

func (s *videoServiceImpl) GetVideoByID(ID string) (video, error) {
	collection := s.mongo.Database("videoDB").Collection("videos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var v video
	objectID, err := bson.ObjectIDFromHex(ID)
	if err != nil {
		slog.Error("failed to create objectID from hex", "error", err, "videoID", ID)
		return video{}, err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	err = collection.FindOne(ctx, filter).Decode(&v)
	if err != nil {
		slog.Error("failed to find video by ID", "error", err, "videoID", ID)
		return video{}, err
	}

	return v, nil
}

func (s *videoServiceImpl) CreateVideo(path string) (string, error) {
	collection := s.mongo.Database("videoDB").Collection("videos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newVideo := video{
		Path: path,
	}

	result, err := collection.InsertOne(ctx, newVideo)
	if err != nil {
		slog.Error("failed to insert new video", "error", err, "videoPath", path)
		return "", err
	}

	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

func (s *videoServiceImpl) DeleteVideo(ID string) error {
	collection := s.mongo.Database("videoDB").Collection("videos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := bson.ObjectIDFromHex(ID)
	if err != nil {
		slog.Error("failed to create objectID from hex", "error", err, "videoID", ID)
		return err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		slog.Error("failed to delete video by ID", "error", err, "videoID", ID)
		return err
	}

	if result.DeletedCount == 0 {
		slog.Error("no video found to delete", "videoID", ID)
		return mongo.ErrNoDocuments
	}
	return nil
}

func (s *videoServiceImpl) UpdateVideo(ID string, newPath string, isUpsert bool) error {
	collection := s.mongo.Database("videoDB").Collection("videos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := bson.ObjectIDFromHex(ID)
	if err != nil {
		slog.Error("failed to create objectID from hex", "error", err, "videoID", ID)
		return err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "path", Value: newPath}}}}
	opts := options.UpdateOne().SetUpsert(isUpsert)

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		slog.Error("failed to update video by ID", "error", err, "videoID", ID)
		return err
	}

	if result.ModifiedCount == 0 {
		slog.Warn("no video found to update", "videoID", ID)
		return mongo.ErrNoDocuments
	}

	return nil
}
