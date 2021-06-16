package middleware

import (
	// models package where User schema is defined
	"crud-backend/models" // models package where User schema is defined
	"database/sql"        // package to encode and decode the json into struct and vice versa
	"encoding/json"
	"fmt"
	"log" // used to access the request and response object of the api
	"net/http"
	"os" // used to read the environment variable
	"strconv"

	// package used to covert string into int type
	"github.com/gorilla/mux"   // used to get the params from the route
	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {

	// load .env file

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check connection

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Success Connected!")

	return db
}

// create bbm

func CreateBbm(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var bbm models.Bbm

	// decoce the json request to bbm

	err := json.NewDecoder(r.Body).Decode(&bbm)

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	// call insert bbm function and pass the bbm

	insertID := insertBbm(bbm)

	// format a response object

	res := response{
		ID:      insertID,
		Message: "Created Succesfully",
	}

	json.NewEncoder(w).Encode(res)
}

// will return a single bbm by its id

func GetBbm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get the bbmid from the request params, key is "id"

	params := mux.Vars(r)

	//convert the id type from string to int
	id, err := strconv.Atoi((params["id"]))

	if err != nil {
		log.Fatalf("unable to convert the string into int. %v", err)
	}

	// call the getBbm function with bbm id to retrieve a single bbm

	bbm, err := getBbm(int64(id))

	if err != nil {
		log.Fatalf("Unable to get bbm. %v", err)
	}

	// send the response

	json.NewEncoder(w).Encode(bbm)
}

func GetAllBbm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get all the bbm in the db

	bbms, err := getAllBbms()

	if err != nil {
		log.Fatalf("unable to get all bbm. %v", err)
	}

	// send all the bbms as response
	json.NewEncoder(w).Encode(bbms)
}

func UpdateBbm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the bbmid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	// create an empty bbm of the models.Bbm

	var bbm models.Bbm

	// decode the json request to bbm

	err = json.NewDecoder(r.Body).Decode(&bbm)

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	// call update bbm to update the bbm

	updateRows := updateBbm(int64(id), bbm)

	// format the message string

	msg := fmt.Sprintf("BBM updated usccessfully. total rows/records affected %v", updateRows)

	// format the response message

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response

	json.NewEncoder(w).Encode(res)

}

func DeleteBbm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	deletedRows := deleteBbm(int64(id))

	msg := fmt.Sprintf("BBM updated succesfully. total rows/record affected %v", deletedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// ======================= Handler Function ====================

// insert one bbm in the DB

func insertBbm(bbm models.Bbm) int64 {
	// create the postgres db connection

	db := createConnection()

	// close the db connection

	defer db.Close()

	// create the insert sql query
	// returning bbmid will return the id of the inserted bbm

	sqlStatement := `INSERT INTO bbms (jumlah_liter, premium, pertalite) VALUES ($1, $2, $3) RETURNING bbmid`

	// the inserted id will store in this id
	var id int64

	//execute the sql statement
	//scan function will save the insert id in the id

	err := db.QueryRow(sqlStatement, bbm.Jumlah_liter, bbm.Premium, bbm.Pertalite).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// get one bbm from the DB by its bbmid

func getBbm(id int64) (models.Bbm, error) {
	// create the postgres fb connection
	db := createConnection()

	//close the db connection
	defer db.Close()

	// create a bbm of models.Bbm type

	var bbm models.Bbm

	// create the selet sql query

	sqlStatement := `SELECT * FROM bbms WHERE bbmid=$1`

	//execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to bbm

	err := row.Scan(&bbm.ID, &bbm.Jumlah_liter, &bbm.Pertalite, &bbm.Premium)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return bbm, nil
	case nil:
		return bbm, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty bbm on error

	return bbm, err
}

func getAllBbms() ([]models.Bbm, error) {
	//create the postgress db connection

	db := createConnection()

	// close conenction

	defer db.Close()

	var bbms []models.Bbm

	//create the select sql query
	sqlStatement := `SELECT * from bbms`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement

	defer rows.Close()

	// iterate over the rows

	for rows.Next() {
		var bbm models.Bbm

		//unmarshal the row object to bbm
		err = rows.Scan(&bbm.ID, &bbm.Jumlah_liter, &bbm.Pertalite, &bbm.Premium)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		bbms = append(bbms, bbm)
	}

	// return empty user on error

	return bbms, err
}

func updateBbm(id int64, bbm models.Bbm) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE bbms SET jumlah_liter=$2, premium=$3, pertalite=$4 WHERE bbmid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, bbm.Jumlah_liter, bbm.Premium, bbm.Pertalite)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete user in the DB
func deleteBbm(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM bbms WHERE bbmid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
