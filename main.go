package main

import (
    "encoding/json"
    "strconv"

    //    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "movies"
    "net/http"
)

var print = fmt.Printf

func getMovies(w http.ResponseWriter , r* http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(movies.GetAll())
}

func getMovieById(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    id, parseIdError := strToUint(params["id"])

    if parseIdError != nil {
        json.NewEncoder(w).Encode( struct { Error string } { parseIdError.Error() } )
    }

    movie, err := movies.GetById(uint(id))

    if err != nil {
        json.NewEncoder(w).Encode( struct { Error string } { err.Error() } )
        return
    }

    json.NewEncoder(w).Encode(movie)
}

func deleteMovieById(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    id, parseIdError := strToUint(params["id"])

    if parseIdError != nil {
        json.NewEncoder(w).Encode( struct { Error string } { parseIdError.Error() } )
    }

    err := movies.DeleteById(uint(id))

    if err != nil {
        json.NewEncoder(w).Encode( struct { Error string } { err.Error() } )
        return
    }
}

func createMovie(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var movie movies.Movie
    _ = json.NewDecoder(r.Body).Decode(&movie)
    _, err := movies.Add(movie)
    if err != nil {
        json.NewEncoder(w).Encode( struct { Error string } { err.Error() } )
    }
    json.NewEncoder(w).Encode(movie)
}


func updateMovieById(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var movie movies.Movie
    _ = json.NewDecoder(r.Body).Decode(&movie)
    fmt.Println(movie)

    params := mux.Vars(r)
    id, parseIdError := strToUint(params["id"])

    if parseIdError != nil {
        json.NewEncoder(w).Encode( struct { Error string } { parseIdError.Error() } )
    }

    err := movies.UpdateById(uint(id), movie)

    if err != nil {
        json.NewEncoder(w).Encode( struct { Error string } { err.Error() } )
        return
    }
}

func strToUint(str string) (uint64, error){
    res, err := strconv.ParseUint(str, 10, 32)
    if err != nil {
        return 0, err
    }
    return res, nil
}

func main(){
    r := mux.NewRouter()

    movies.Init()

    r.HandleFunc("/movies", getMovies).Methods("GET")
    r.HandleFunc("/movies/{id}", getMovieById).Methods("GET")
    r.HandleFunc("/movies", createMovie).Methods("POST")
    r.HandleFunc("/movies/{id}", deleteMovieById).Methods("DELETE")
    r.HandleFunc("/movies/{id}", updateMovieById).Methods("PUT")

    print("Starting server at port 8000\n")


    http.ListenAndServe(":8000", r)
}
