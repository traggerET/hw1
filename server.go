package main

import (
 "encoding/json"
 "fmt"
 "net/http"
)

var dataStore = make(map[string]string)

func putHandler(w http.ResponseWriter, r *http.Request) {
  key := r.URL.Query().Get("key")
  value := r.URL.Query().Get("value")
  if key == "" || value == "" {
    http.Error(w, "Key and value parameters are required", http.StatusBadRequest)
    return
  }
  dataStore[key] = value
  w.WriteHeader(http.StatusCreated)
  fmt.Fprintf(w, "Value added successfully")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
 key := r.URL.Query().Get("key")
 if key == "" {
  http.Error(w, "Key parameter is missing", http.StatusBadRequest)
  return
 }

 value, ok := dataStore[key]
 if !ok {
  http.Error(w, "Key not found", http.StatusNotFound)
  return
 }

 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(map[string]string{"value": value})
}


func main() {
 http.HandleFunc("/put", putHandler)
 http.HandleFunc("/get", getHandler)

 fmt.Println("Server listening on :8080")
 http.ListenAndServe(":8080", nil)
}
