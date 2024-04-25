package routes

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type Status struct {
	Status map[string]int
}

type StatusModel struct {
	StatusCollection *mongo.Collection
	Logger           *log.Logger
	S                Status
}

func (m *StatusModel) addStatus(key string, val int) error {
	m.S.Status[key] = val
	m.Logger.Println(m.S.Status[key])
	return nil
}
