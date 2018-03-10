package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"encoding/base64"
	"crypto/rand"
	"time"
	"os"
	// "errors"
	// "regexp"
	"github.com/mattsykesj/hello/views"
	"log"
	"net/http"
)

const (
	host           = "localhost"
	port           = 5432
	user           = "matt"
	password       = "stunl0ck"
	dbname         = "matt"
	httpPortString = ":8080"
)

type Session struct {
	UserId int64
	UserName string
	Id string
}

var db *sql.DB
var session *Session

func checkErr(err error) {
	if err != nil {
		panic(err)
	}			
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login/", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		template.RenderLogin("login", false, w)
	}
	
	if r.Method == "POST" {
		r.ParseForm()

		inputUserName := r.PostFormValue("username")

		var userName string
		var userId int64

		fmt.Println("DB#LOGIN#Connecting")
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlInfo)
		checkErr(err)
		defer db.Close()

		fmt.Println("DB#LOGIN#Querying users")
		err = db.QueryRow("SELECT user_name, user_id FROM users where user_name=$1", inputUserName).Scan(&userName, &userId)

		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("Couldn't find user")
				template.RenderLogin("login", true, w)
			} else { 
				panic(err)
			}
		}

		fmt.Println("Successfully found user")
		fmt.Println("Creating Session")	

		bytes := make([]byte, 32)
		_ , err = rand.Read(bytes)
		checkErr(err)

		sessionId := base64.StdEncoding.EncodeToString(bytes)

		session.Id = sessionId
		session.UserName = userName
		session.UserId = userId 

		fmt.Printf("Session Id: %v\n", session.Id)
		fmt.Printf("User Id: %v\n", session.UserId)
		fmt.Printf("User Name: %v\n", session.UserName)

		cookie := &http.Cookie{Name: "test_id", Value: sessionId, Path: "/", HttpOnly:true, MaxAge:0}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/food/", http.StatusFound)
		return
	}
}

func addFoodHandler(w http.ResponseWriter, r *http.Request) {
	if(r.Method == "GET") {
		template.RenderAddFood("Add Food", "", w)
	}

	if(r.Method == "POST") {
		//TODO(matt) make this better
		cookie, err := r.Cookie("test_id")
		checkErr(err)

		if session.Id == cookie.Value {

			r.ParseForm()

			var insertedId int64

			fmt.Println("DB#ADDFOOD#Connecting")
			psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
			db, err := sql.Open("postgres", psqlInfo)
			checkErr(err)
			defer db.Close()

			fmt.Println("DB#ADDFOOD#Inserting foods")
			err = db.QueryRow("INSERT INTO foods(name, protein, carbohydrate, fat, calories) VALUES($1, $2, $3, $4, $5) RETURNING food_id;", r.PostFormValue("name"), r.PostFormValue("protein"), r.PostFormValue("carbohydrate"), r.PostFormValue("fat"), r.PostFormValue("calories")).Scan(&insertedId)
			checkErr(err)

			fmt.Printf("DB#ADDFOOD#Inserted 1 row Id: %v\n", insertedId)

			fmt.Println("DB#ADDFOOD#Preparing")
			stmt, err := db.Prepare("INSERT INTO user_foods(user_id, food_id) VALUES($1, $2)")
			checkErr(err)

			fmt.Println("DB#ADDFOOD#Inserting user foods")
			_ , err = stmt.Exec(session.UserId, insertedId)
			checkErr(err)

			template.RenderAddFood("Add Food", r.PostFormValue("name"), w)
		} else { 
			//TODO(matt): no user redirect login
		}
	}
}

func foodHandler(w http.ResponseWriter, r *http.Request) {
	
	//TODO(matt)Replace this with proper checking - check cookie vs session every req?
	cookie, err := r.Cookie("test_id")
	checkErr(err)

	if session.Id == cookie.Value {
		username := session.UserName
		u := &template.UserVM{}

		fmt.Println("DB#FOOD#Connecting")
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlInfo)
		checkErr(err)
		defer db.Close()

		fmt.Println("DB#FOOD#Querying users")
		userRows, err := db.Query("SELECT user_id, user_name, target_calories, target_protein, target_carbohydrate, target_fat FROM users where user_name=$1", username)
		checkErr(err)
		defer userRows.Close()

		for userRows.Next() {
			err := userRows.Scan(&u.UserId, &u.UserName, &u.TargetCalories, &u.TargetProtein, &u.TargetCarbohydrate, &u.TargetFat)
			checkErr(err)
		}

		fmt.Println("DB#FOOD#Querying user foods")
		userFoodRows, err := db.Query("SELECT foods.name FROM foods, user_foods WHERE foods.food_id = user_foods.food_id AND user_foods.user_id = $1", session.UserId)
		checkErr(err)
		defer userFoodRows.Close()

		foods := []string{}

		for userFoodRows.Next() {
			var foodName string
			err = userFoodRows.Scan(&foodName)
			foods = append(foods, foodName)
			checkErr(err)
		}

		template.RenderFood("food", u, foods, w)
	
	} else { 
		
		//Cant find user in session - redirect to login?	
		fmt.Println("Cant find user in session")	
	}

}

func initDbCon()  {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	checkErr(err)
	defer db.Close()
	err = db.Ping()

	if err != nil {
		log.Println("Failed to connect to db")
		log.Println(psqlInfo)
		panic(err)
	}

	log.Println("Successfully connected to db")
}

func resourceHandler(w http.ResponseWriter, r *http.Request) {
	resource, err := os.Open("." + r.URL.Path)
	checkErr(err)
    http.ServeContent(w, r, r.URL.Path, time.Now(), resource)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/food/", foodHandler)
	http.HandleFunc("/food/add/", addFoodHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/content/css/", resourceHandler)
	http.HandleFunc("/content/scripts/", resourceHandler)

	initDbCon()
	session = &Session{}

	log.Println("App started listening on ", httpPortString)
	log.Fatal(http.ListenAndServe(httpPortString, nil))
}

