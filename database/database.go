package database

import (
	"context"
	"log"
	"time"

	"github.com/forzyz/crudql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var connectionString string = "mongodb+srv://forzyz:qazqazqaz2006@cluster0.4xoqljs.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &DB{
		client: client,
	}
}

func (db *DB) GetJob(id string) *model.JobListing {
	jobCollec := db.client.Database("grapql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	var jobListing model.JobListing
	err := jobCollec.FindOne(ctx, filter).Decode(&jobListing)
	if err != nil {
		log.Fatal(err)
	}

	return &jobListing
}

func (db *DB) GetJobs() []*model.JobListing {
	jobCollec := db.client.Database("grapql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var jobListings []*model.JobListing
	cursor, err := jobCollec.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(context.TODO(), &jobListings); err != nil {
		panic(err)
	}

	return jobListings
}

func (db *DB) CreateJobListing(jobInfo model.CreateJobListingInput) *model.JobListing {
	jobCollec := db.client.Database("grapql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	inserted, err := jobCollec.InsertOne(ctx, bson.M{
		"title":       jobInfo.Title,
		"description": jobInfo.Description,
		"url":         jobInfo.URL,
		"company":     jobInfo.Company,
	})
	if err != nil {
		log.Fatal(err)
	}
	insertedID := inserted.InsertedID.(primitive.ObjectID).Hex()
	returnJobListing := model.JobListing{
		ID:          insertedID,
		Title:       jobInfo.Title,
		Description: jobInfo.Description,
		URL:         jobInfo.URL,
		Company:     jobInfo.Company,
	}

	return &returnJobListing
}

func (db *DB) UpdateJobListing(jobId string, jobInfo model.UpdateJobListingInput) *model.JobListing {
	jobCollec := db.client.Database("grapql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateJobInfo := bson.M{}

	if jobInfo.Title != nil {
		updateJobInfo["title"] = jobInfo.Title
	}

	if jobInfo.Description != nil {
		updateJobInfo["description"] = jobInfo.Description
	}

	if jobInfo.URL != nil {
		updateJobInfo["url"] = jobInfo.URL
	}

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateJobInfo}

	results := jobCollec.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var jobListing model.JobListing

	if err := results.Decode(&jobListing); err != nil {
		log.Fatal()
	}
	return &jobListing
}

func (db *DB) DeleteJobListing(jobId string) *model.DeleteJobResponse {
	jobCollec := db.client.Database("grapql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	_, err := jobCollec.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	return &model.DeleteJobResponse{DeleteJobID: jobId}
}
