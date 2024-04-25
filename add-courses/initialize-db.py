import csv
import json
from pymongo import MongoClient # type: ignore

with open('courses.csv', newline='') as csvfile:
    client = MongoClient("mongodb://localhost:27017")
    dbname = client["diskit-db"]
    collection = dbname["courses"]
    # print(collection)
    courses = csv.reader(csvfile, delimiter=',')
    fields = {
        0: "name",
        1: "address",
        2: "city",
        3: "state",
        4: "zip"
    }
    i = 0
    for course in courses:
        obj = {}
        for field in course:
            if(field.isdigit()):
                field = int(field)
            else:
                field = field.upper()
            obj[fields[i]]=field
            i=i+1
        i=0
        print(obj)
        x = collection.insert_one(obj)
        print(x)