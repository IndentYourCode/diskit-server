package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	Status  map[string]int     `bson:"status"`
}

type CoursesModel struct {
	Courses *mongo.Collection
	Logger  *log.Logger
	Jobs    chan string
}

func CourseModel(c *mongo.Collection, l *log.Logger) *CoursesModel {
	model_jobs := make(chan string, 256)
	cm := CoursesModel{
		Courses: c,
		Logger:  l,
		Jobs:    model_jobs,
	}
	go func(jobs chan string) {
		for val := range jobs {
			l.Printf("Grabbed Value: %s\n", val)
			time.Sleep(10 * time.Second)
			key := "status." + val
			filter := bson.D{{Key: key, Value: bson.D{{Key: "$gt", Value: 0}}}}
			update := bson.D{{Key: "$inc", Value: bson.D{{Key: "status.crowded", Value: -1}}}}
			resp, _ := c.UpdateMany(context.TODO(), filter, update)
			l.Printf("Matched Count: %d\n", resp.MatchedCount)
			if resp.ModifiedCount > 0 {
				l.Printf("Modified Count: %d\n", resp.ModifiedCount)
				jobs <- val
				l.Printf("added val to jobs")
			}
		}
	}(model_jobs)
	model_jobs <- "crowded"
	model_jobs <- "empty"
	model_jobs <- "wind"
	model_jobs <- "rain"
	return &cm
}
func (m *CoursesModel) IncrementRain(w http.ResponseWriter, req *http.Request) error {
	cid := httprouter.Param(req, "id")
	objId, _ := primitive.ObjectIDFromHex(cid)
	filter := bson.D{{Key: "_id", Value: objId}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "status.rain", Value: 1}}}}
	resp, err := m.Courses.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	m.Logger.Printf("Update Rain: %+v\n", resp)
	m.Jobs <- "rain"
	return nil
}

func (m *CoursesModel) IncrementEmpty(w http.ResponseWriter, req *http.Request) error {
	cid := httprouter.Param(req, "id")
	objId, _ := primitive.ObjectIDFromHex(cid)
	filter := bson.D{{Key: "_id", Value: objId}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "status.empty", Value: 1}}}}
	resp, err := m.Courses.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	m.Logger.Printf("Update Empty: %+v\n", resp)
	m.Jobs <- "empty"
	return nil
}

func (m *CoursesModel) IncrementWindy(w http.ResponseWriter, req *http.Request) error {
	cid := httprouter.Param(req, "id")
	objId, _ := primitive.ObjectIDFromHex(cid)
	filter := bson.D{{Key: "_id", Value: objId}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "status.wind", Value: 1}}}}
	resp, err := m.Courses.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	m.Logger.Printf("Update Wind: %+v\n", resp)
	m.Jobs <- "wind"
	return nil
}

func (m *CoursesModel) IncrementCrowd(w http.ResponseWriter, req *http.Request) error {
	cid := httprouter.Param(req, "id")
	objId, _ := primitive.ObjectIDFromHex(cid)
	filter := bson.D{{Key: "_id", Value: objId}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "status.crowded", Value: 1}}}}
	resp, err := m.Courses.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	m.Logger.Printf("Update Crowd: %+v\n", resp)
	m.Jobs <- "crowded"
	return nil
}

func (m *CoursesModel) GetPopulation(w http.ResponseWriter, req *http.Request) error {
	cid := httprouter.Param(req, "id")

	objId, _ := primitive.ObjectIDFromHex(cid)
	filter := bson.D{{Key: "_id", Value: objId}}

	var course Course

	err := m.Courses.FindOne(context.TODO(), filter).Decode(&course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	empty := course.Status["empty"]
	crowded := course.Status["crowded"]
	result := crowded - empty

	prop := struct {
		Property string
		Value    int
	}{
		Property: "Population",
		Value:    result,
	}

	resp, _ := json.Marshal(prop)
	m.Logger.Print(resp)
	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
	return nil
}

func (m *CoursesModel) GetCourseStats(w http.ResponseWriter, req *http.Request) error {
	cid := httprouter.Param(req, "id")

	objId, _ := primitive.ObjectIDFromHex(cid)
	filter := bson.D{{Key: "_id", Value: objId}}

	var course Course

	err := m.Courses.FindOne(context.TODO(), filter).Decode(&course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	m.Logger.Print(course.Status)

	resp, _ := json.Marshal(course.Status)
	m.Logger.Print(resp)
	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)

	return nil
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
