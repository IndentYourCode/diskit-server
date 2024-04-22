package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cs455_server/routes"

	"github.com/nahojer/httprouter"
	"github.com/nahojer/httprouter/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Course struct {
	name    string
	address string
}

func healthCheck(w http.ResponseWriter, req *http.Request) error {
	w.Write([]byte("healthy"))
	return nil
}

func getCourses(w http.ResponseWriter, req *http.Request) error {
	w.Write([]byte("This provides a list of all the nearby courses"))
	return nil
}

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	logger.Println("This is an info message")
	uri := os.Getenv("MONGO_URL")
	logger.Printf("MONGODB_URL = %s", uri)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	opts.SetDirect(true)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		logger.Println("Error when Connecting to MongoDB")
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			logger.Fatal("Error when Disconnecting")
			panic(err)
		}
	}()
	var result bson.M

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	logger.Println("Connected to MongoDB!")

	coll := client.Database("diskit-db").Collection("courses")

	courses := routes.CoursesModel{
		Courses: coll,
		Logger:  logger,
	}

	r := httprouter.New() // new router
	r.Use(middleware.RecoverPanics())
	r.Handle(http.MethodGet, "/health", healthCheck)
	r.Handle(http.MethodGet, "/region/:region", courses.GetCoursesByRegion)
	r.Handle(http.MethodGet, "/course/:id", courses.GetCourse)
	r.Handle(http.MethodPost, "/courses", courses.PostCourse)
	go log.Fatal(http.ListenAndServe(":3000", r))
}
