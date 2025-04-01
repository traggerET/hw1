package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "gorm.io/driver/postgres" 
  "gorm.io/gorm" 
)

var db *gorm.DB

type Item struct {
  ID string `gorm:"primaryKey"`
  Value string `gorm:"not null"`
}

func putHandler(w http.ResponseWriter, r *http.Request) {
  key := r.URL.Query().Get("key") 
  value := r.URL.Query().Get("value") 
  if key == "" || value == "" { 
    http.Error(w, "Key and value parameters are required", http.StatusBadRequest) 

    return 
  } 

  result := db.Save(&Item{ ID: key, Value: value }) 
  if result.Error != nil { 
    http.Error(w, "failed to save to db", http.StatusInternalServerError) 
  }
  w.WriteHeader(http.StatusCreated)

  fmt.Fprintf(w, "Value added successfully") 
}

func getHandler(w http.ResponseWriter, r *http.Request) {
  key := r.URL.Query().Get("key") 
  if key == "" { 
    http.Error(w, "Key parameter is missing", http.StatusBadRequest) 
    
    return 
  } 
  var item Item 
  result := db.First(&item, key) 
  if result.Error != nil { 
    http.Error(w, "Key not found", http.StatusNotFound) 

    return 
  }
  w.Header().Set("Content-Type", "application/json") 
  
  json.NewEncoder(w).Encode(item) 
}

func main() {
  connstr := "host=localhost user=postgres password=postgres dbname=testdb port=5432 sslmode=disable"
  db, _ = gorm.Open(postgres.Open(connstr), &gorm.Config{})
  db.AutoMigrate(&Item{}) 

  http.HandleFunc("/put", putHandler)
  http.HandleFunc("/get", getHandler)

  fmt.Println("Server listening on :8080")
  http.ListenAndServe(":8080", nil)
}
