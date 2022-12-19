package movies

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

type Person struct {
    gorm.Model
    Name    string `json:"name"`
    Surname string `json:"surname"`
}

type Movie struct {
    gorm.Model
    EIDR      string `json:"eidr" gorm:"unique;column:eidr"`
    Title     string `json:"title"`
    Directors []Person `json:"directors" gorm:"many2many:movie_directors"`
    Actors    []Person `json:"actors" gorm:"many2many:movie_actors"`
}

var DB *gorm.DB

func Init() {
    dsn := "host=localhost user=postgres password=root dbname=movies port=5432"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    DB = db
    DB.AutoMigrate(&Movie{})

    if err != nil {
        log.Fatal(err)
        return
    }

}

func GetAll() []Movie {
    var movies []Movie
    DB.Model(&Movie{}).Preload("Actors").Preload("Directors").Find(&movies).Find(&movies)
    return movies
}

func GetById(id uint) (Movie, error) {
    var movie = Movie{Model: gorm.Model{ID: id}}
    res := DB.Model(&Movie{}).Preload("Actors").Preload("Directors").First(&movie)
    if res.Error != nil {
        return Movie{}, res.Error
    }
    return movie, nil
}

func DeleteById(id uint) error {
    var movie = Movie{Model: gorm.Model{ID: id}}
    res := DB.Delete(&movie)
    if res.Error != nil {
        return res.Error
    }
    return res.Error
}

func UpdateById(id uint, newMovie Movie) error {
    var movie = Movie{Model: gorm.Model{ID: id}}
    res := DB.First(&movie)
    if res.Error != nil {
        return res.Error
    }
    movie.EIDR = newMovie.EIDR
    movie.Title = newMovie.Title
    movie.Directors = newMovie.Directors
    movie.Actors = newMovie.Actors

    res2 := DB.Save(&movie)
    if res2.Error != nil {
        return res2.Error
    }

    return nil
}

func Add(newMovie Movie) (Movie, error) {
    res := DB.Save(&newMovie)
    if res.Error != nil {
        return Movie{}, res.Error
    }
    return newMovie, nil
}
