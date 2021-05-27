package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"strconv"  // package used to covert string into int type

	"github.com/TrondSpjelakvik/golang-backend/models"
	"github.com/gorilla/mux" // used to get the params from the route

	// "golang.org/x/crypto/bcrypt"

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response

type response struct {
	ID int64 `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// connect to postgres db
func createConnection() *sql.DB {

	// get .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file :(")
	}

	// Open connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// Check connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Success! You are connected")

	// Return connection
	return db

}


// Create a user in DB
func CreateUser(w http.ResponseWriter, r * http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")


	// empty user of the type User
	
	var user models.User

	// hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 10)

	// user := models.User{
	// 	ID: u.ID,
	// 	Name: u.Name,
	// 	Location: u.Location,
	// 	Age: u.Age,
	// 	Password: string(hash),
	// 	Email: u.Email,
	// 	Username: u.Username,

	// }
	
	// Json request to user
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Not able to decode the request: %v", err)
	}

	// Insert the user and pass the user
	insertID := insertUser(user)

	// format a response object
	res := response{
		ID: insertID,
		Message: "User created",
	}

	// Send the response
	json.NewEncoder(w).Encode(res)

}

// Return a single user by id

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Get the user id from the request
	params := mux.Vars(r)

	// Convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Cant convert the string to int. %v", err)
	}

	// Call getUser with user id and retrive a single user

	user, err := getUser(int64(id))

	if err != nil {
		log.Fatalf("Cant get user. %v", err)
	}

	// Send response

	json.NewEncoder(w).Encode(user)
}

	// Get all users
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get all users from the db
	users, err := getAllUsers()

	if err != nil {
		log.Fatalf("Cant get all users. %v", err)
	}

	// all the users response
	json.NewEncoder(w).Encode(users)

}

// Update users details
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Get the userid from request

	params := mux.Vars(r)

	// Convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Cant convert the string to int. %v", err)
	}

	// Create an empty user of type models.User
	var user models.User

	// Decode the json request
	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Cant decode the request body. %v", err)
	}

	// Update the user
	updateRows := updateUser(int64(id), user)

	// Format the message string
	msg := fmt.Sprintf("User updated. Total rows affected %v", updateRows)

	// Format the response message
	res := response{
		ID: int64(id),
		Message: msg,
	}

	// Send the response
	json.NewEncoder(w).Encode(res)

}

	// Delete users details 
	func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")	

	// Get the user id from the request
	params := mux.Vars(r)

	// Convert string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Cant convert string to int. %v", err)
	}

	// Call the deleteUser. int => int64
	deletedRows := deleteUser(int64(id))

	// Format the message string
	msg := fmt.Sprintf("User updated. Total rows/record affected %v", deletedRows)

	// Format the response message
	res := response{
		ID: int64(id),
		Message: msg,
	}

	// Send the response
	json.NewEncoder(w).Encode(res)

	}


	// Handler functions

	// Insert one user in the DB

	func insertUser(user models.User) int64 {

		// Create the db connection
		db := createConnection()

		// Close the db connection
		defer db.Close()

		


		
		// Insert sql query
		// Returning user id will return the id of the inserted user
		sqlStatement := `INSERT INTO users (name, location, age, password, email, username) VALUES ($1, $2, $3, $4, $5, $6) RETURNING userid`

		// Inserted id will store in this id
		var id int64

		

		// Execute the sql statement
		// Scan function will save the insert id in the id
		err := db.QueryRow(sqlStatement, user.Name, user.Location, user.Age, user.Password, user.Email, user.Username).Scan(&id)

		
	


		if err != nil {
			log.Fatalf("Cant execute the query. %v", err)
		}

		fmt.Printf("Inserted a single record %v", id)

		return id

	}

	// Get one user from the DB by id 
	func getUser(id int64) (models.User, error) {
		// Create db connection
		db := createConnection()

		// Close the connection
		defer db.Close()

		// Create a user of models.User
		var user models.User

		// Create the SELECT query
		sqlStatement := `SELECT * FROM users WHERE userid=$1`

		// Execute swl statement
		row := db.QueryRow(sqlStatement, id)

		// Unmarshal the row object to user
		err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Location, &user.Email, &user.Username, &user.Password)

		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows was returned")
			return user, nil
		case nil:
			return user, nil
		default:
			log.Fatalf("Unable to scan the row. %v", err)
		}


		// Return empty user in error
		return user, err

	}

	// Get one user from the DB by ID

	func getAllUsers() ([]models.User, error) {

		// Create connection
		db := createConnection()

		// Close connection
		defer db.Close()

		var users []models.User

		// Create the SELECT query
		sqlStatement := `SELECT * FROM users`

		// Execute statement
		rows, err := db.Query(sqlStatement)

		if err != nil {
			log.Fatalf("Cant execute the query. %v", err)
		}

		// Close statement
		defer rows.Close()

		// Iterate over the rows
		for rows.Next() {
			var user models.User

			// Unmarshal the row object to user
			err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.Location, &user.Email, &user.Password, &user.Username)

			if err != nil {
				log.Fatalf("Cant scan the row. %v", err)
			}

			// Append the user in the users slice
			users = append(users, user)

		}

		// Return empty user on error
		return users, err

	}

	// Update user in the DB
	func updateUser(id int64, user models.User) int64 {

		// Create db connection
		db := createConnection()

		// Close the db connection
		defer db.Close()

		// Create the update sql query

		sqlStatement := `UPDATE users SET name=$2, location=$3, age=$4, password=$5, username=$6, email=$7 WHERE userid=$1`

		res, err := db.Exec(sqlStatement, id, user.Name, user.Location, user.Age, user.Password, user.Username, user.Email)

		if err != nil {
			log.Fatalf("Cant execute the query. %v", err)
		}

		// Check how many rows are affected

		rowsAffected, err := res.RowsAffected()

		if err != nil {
			log.Fatalf("Error while checking the affected rows. %v", err)
		}

		fmt.Printf("Total rows/records affected %v", rowsAffected)

		return rowsAffected

	}

	// Delete user in DB
	func deleteUser(id int64) int64 {

		// Create db connection
		db := createConnection()

		// Close db connection
		defer db.Close()

		// Create the DELETE query
		sqlStatement := `DELETE FROM users WHERE userid=$1`

		// Execute the statement
		res, err := db.Exec(sqlStatement, id)

		if err != nil {
			log.Fatalf("Cant execute the query. %v", err)
		}

		// Check how many rows affected
		rowsAffected, err := res.RowsAffected()

		if err != nil {
			log.Fatalf("Error while cheking the affected rows. %v", err)
		}

		fmt.Printf("Total rows/record affected %v", rowsAffected)

		return rowsAffected

	}
