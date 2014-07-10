package main

import(
  "database/sql"
  "fmt"
  "github.com/go-martini/martini"
  _ "github.com/lib/pq"
  "github.com/martini-contrib/render"
  "net/http"
)

var (
  createTable = `CREATE TABLE IF NOT EXISTS users(
		name character varying(100) NOT NULL,
		email character varying(100) NOT NULL,
		description character varying(500) NOT NULL
    );`
)

func SetupDB() *sql.DB {
  db, err := sql.Open("postgres", "host=localhost user=postgres password=ankit1234 dbname=App42PaaS-Rails-PostgreSQL-Sample_development sslmode=disable")
  fmt.Println(db)
  PanicIf(err)
  ctble, err := db.Query(createTable)
	PanicIf(err)
	fmt.Println("Table created successfully", ctble)
  return db	
}

type User struct {
	Name       string
	Email      string
	Description string
	Id          int
}

func PanicIf(err error) {
  if err != nil {
    panic(err)  	
  }
}

func main(){
  m := martini.Classic()
  m.Map(SetupDB())

  m.Use(render.Renderer(render.Options{
	Layout: "layout",
  }))
  
  m.Post("/users", func(ren render.Render, r *http.Request, db *sql.DB) {

    fmt.Println(r.FormValue("name"))
    fmt.Println(r.FormValue("email"))
    fmt.Println(r.FormValue("description"))

    _, err := db.Query("INSERT INTO users (name, email, description) VALUES ($1, $2, $3)",
      r.FormValue("name"),
      r.FormValue("email"),
      r.FormValue("description"))

    PanicIf(err)
    
    rows, err := db.Query("SELECT * FROM users")
		PanicIf(err)
		defer rows.Close()

		fmt.Println(rows)
		fmt.Println(err)

		users := []User{}
		for rows.Next() {

      user := User{}
			err := rows.Scan(&user.Name, &user.Email, &user.Description)
			PanicIf(err)
			users = append(users, user)

			fmt.Println(users);
		}
		ren.HTML(200, "users", users)
	})
  
  m.Run()
}
