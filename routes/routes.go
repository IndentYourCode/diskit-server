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
	Name    string `json:"name"`
	Address string `json:"address"`
}

type CoursesModel struct {
	Courses *mongo.Collection
	Logger  *log.Logger
}

func (m *CoursesModel) GetCourse(w http.ResponseWriter, req *http.Request) error {
	cid := httprouter.Param(req, "id")

	objId, _ := primitive.ObjectIDFromHex(cid)
	filter := bson.D{{"_id", objId}}

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

	var c Course

	err := json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	m.Logger.Printf("Body: %+v", c)

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
