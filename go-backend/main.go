package main

// May be provided automatically by your editor, so they may be removed if you hit save right after typing this bit.
import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

type Item struct {
	Id       int
	Price    float32
	Name     string
	ImageURL string
}

func returnResponse(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/response.html"))
	response := r.Header.Get("Hx-Prompt")
	tmpl.Execute(w, response)
}

func spitOutItems(w http.ResponseWriter, r *http.Request) {
	db := getDBConnection()
	defer db.Close()
	items := loadNextItems(db, -1, 5)
	// items = []Item{
	// 	{Id: 1, Name: "test", Price: 1.0, ImageURL: "test"},
	// 	{Id: 2, Name: "test2", Price: 2.0, ImageURL: "test2"},
	// 	{Id: 3, Name: "test3", Price: 3.0, ImageURL: "test3"},
	// }
	tmpl := template.Must(template.ParseFiles("templates/shop-item.html"))
	tmpl.Execute(w, items)
}

// Loads next N items from the database
func loadNextItems(db *sql.DB, lastIndex int, amountToLoad int) []Item {
	query := fmt.Sprintf("SELECT * FROM items WHERE id > %d LIMIT %d", lastIndex, amountToLoad)
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	items := make([]Item, 0)
	for rows.Next() {
		item := Item{}
		err := rows.Scan(&item.Id, &item.Name, &item.Price, &item.ImageURL)
		if err != nil {
			panic(err)
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	for _, item := range items {
		fmt.Println(item)
	}
	return items
}

func getDBConnection() *sql.DB {
	// connect to postgres
	sqlConnectionString := "host=siwenna port=5432 user=toyuser password=88jelena! dbname=toydb sslmode=disable search_path=shoppe"

	// connect to the database
	db, err := sql.Open("postgres", sqlConnectionString)
	if err != nil {
		panic(err)
	}

	return db
}

func main() {
	body := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	}
	http.HandleFunc("/", body)
	http.HandleFunc("/getItems", spitOutItems)
	http.HandleFunc("/submit", returnResponse)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe("127.0.0.1:3000", nil)

}
