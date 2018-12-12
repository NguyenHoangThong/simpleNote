package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var db *sql.DB

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

func dbConfig () map[string]string {
	//conf := make(map[string]string)
	conf, err := godotenv.Read("./config/development.env")
	if err != nil {
		log.Fatal("Can't load config file")
	}
	return conf
}

func initDb () {
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
}


func main() {
	initDb()
	defer db.Close()
	http.HandleFunc("/api/index", indexHandler)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type repositorySummary struct {
	ID         int
	Name       string
	Owner      string
	TotalStars int
}

type repositories struct {
	Repositories []repositorySummary
}

// indexHandler calls `queryRepos()` and marshals the result as JSON
func indexHandler(w http.ResponseWriter, req *http.Request) {
	repos := repositories{}

	err := queryRepos(&repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

// queryRepos first fetches the repositories data from the db
func queryRepos(repos *repositories) error {
	rows, err := db.Query(`
		SELECT
			id,
			owner,
			name,
			total_stars
		FROM repositories
		ORDER BY total_stars DESC`)

	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		repo := repositorySummary{}
		err = rows.Scan(
			&repo.ID,
			&repo.Owner,
			&repo.Name,
			&repo.TotalStars,
		)
		if err != nil {
			return err
		}
		repos.Repositories = append(repos.Repositories, repo)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

//func indexHandler(w http.ResponseWriter, r *http.Request) {
//	//...
//	fmt.Fprintf(w, string("hello"))
//
//}
//
//func repoHandler(w http.ResponseWriter, r *http.Request) {
//	//...
//}