package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/nahojer/httprouter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Course struct {
	Id      primitive.ObjectID `bson:"_id"`
	Name    string             `bson:"name"`
	Address string             `bson:"address"`
	City    string             `bson:"city"`
	State   string             `bson:"state"`
	ZipCode int                `bson:"zip"`
}

type CoursesModel struct {
	Courses *mongo.Collection
	Logger  *log.Logger
}

func CourseModel(c *mongo.Collection, l *log.Logger) *CoursesModel {
	cm := CoursesModel{
		Courses: c,
		Logger:  l,
	}
	return &cm
}

func (m *CoursesModel) GetCoursesByRegion(w http.ResponseWriter, req *http.Request) error {
	rid := httprouter.Param(req, "city")

	filter := bson.D{{Key: "city", Value: rid}}

	cursor, err := m.Courses.Find(context.TODO(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	var results []Course

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		m.Logger.Printf("%+v\n", result)
	}

	resp, _ := json.Marshal(results)
	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
	return nil
}

func (m *CoursesModel) GetCourse(w http.ResponseWriter, req *http.Request) error {
	cid := httprouter.Param(req, "id")

	objId, _ := primitive.ObjectIDFromHex(cid)
	filter := bson.D{{Key: "_id", Value: objId}}

	var course Course

	err := m.Courses.FindOne(context.TODO(), filter).Decode(&course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	resp, _ := json.Marshal(course)
	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)

	return nil
}

func (m *CoursesModel) PostCourse(w http.ResponseWriter, req *http.Request) error {

	type PostBody struct {
		Name    string `bson:"name"`
		Address string `bson:"address"`
		City    string `bson:"city"`
		State   string `bson:"state"`
		ZipCode int    `bson:"zip"`
	}

	var c PostBody

	err := json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	resp, err := m.Courses.InsertOne(context.TODO(), c)
	if err != nil {
		m.Logger.Println("Insertion Failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	} else {
		m.Logger.Println("Insertion Succeeded")
		w.Write([]byte(resp.InsertedID.(primitive.ObjectID).Hex()))
	}
	return nil
}
