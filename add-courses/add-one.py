from pymongo import MongoClient

client = MongoClient("mongodb://localhost:27017")
dbname = client["diskit-db"]
collection = dbname["courses"]

json = {
    "name" : "foobartoo",
    "address" : "1234 bar st",
    "city" : "vancouver",
    "state" : "WA",
    "zip" : 98683,
    "status" : {
        "crowded" : 10
    }
}

x = collection.insert_one(json)
print(x)