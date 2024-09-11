package main

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type productTypes struct {
	St    int `json:"st"`
	Pt    int `json:"pt"`
	Units int `json:"units"`
}

type Order struct {
	OrderID  string         `json:"orderID"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Phone    string         `json:"phone"`
	Pincode  string         `json:"pincode"`
	Time     string         `json:"time"`
	Address  string         `json:"address"`
	Products []productTypes `json:"pt"`
}

func (app *application) FormData(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "", 400)
		app.errLog.Println(err)
		return
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	number := r.PostForm.Get("number")
	pincode := r.PostForm.Get("pincode")
	time := r.PostForm.Get("time")
	address := r.PostForm.Get("address")

	units := r.Form["units"]
	serviceTypes := r.Form["serviceType"]
	productTypes := r.Form["productType"]

	_, err = app.db.Exec("INSERT INTO legorders(name, email, number, pincode, time, address) VALUES(?,?,?,?,?,?);", name, email, number, pincode, time, address)

	if err != nil {
		http.Error(w, "", 500)
		app.errLog.Println(err)
		return
	}

	var order_id int
	row := app.db.QueryRow("SELECT LAST_INSERT_ID();")

	err = row.Scan(&order_id)

	if err != nil {
		app.errLog.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(units); i++ {

		_, err = app.db.Exec("INSERT INTO legorderTypes VALUES(?,?,?,?);", order_id, serviceTypes[i], productTypes[i], units[i])

		if err != nil {
			app.errLog.Println(err)
			http.Error(w, "", 500)
		}
	}

}

func (app *application) getForm(w http.ResponseWriter, r *http.Request) {

	rows, err := app.db.Query(`
	SELECT 
		*
	FROM 
		legorders`)
	if err != nil {
		http.Error(w, "Error querying orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []Order

	// Iterate over all orders
	for rows.Next() {
		var order Order
		if err := rows.Scan(
			&order.OrderID, &order.Name, &order.Email, &order.Phone,
			&order.Pincode, &order.Time, &order.Address,
		); err != nil {
			http.Error(w, "Error scanning orders", http.StatusInternalServerError)
			return
		}

		// Query for product types associated with the current order
		ptRows, err := app.db.Query(`
		SELECT 
			st, pt, u 
		FROM 
			legorderTypes 
		WHERE 
			orderid = ?`, order.OrderID)
		if err != nil {
			app.errLog.Println(err)
			http.Error(w, "Error querying product types", http.StatusInternalServerError)
			return
		}
		defer ptRows.Close()

		var products []productTypes
		for ptRows.Next() {
			var pt productTypes
			if err := ptRows.Scan(&pt.St, &pt.Pt, &pt.Units); err != nil {
				http.Error(w, "Error scanning product types", http.StatusInternalServerError)
				return
			}
			products = append(products, pt)
		}
		if err := ptRows.Err(); err != nil {
			http.Error(w, "Error iterating over product types", http.StatusInternalServerError)
			return
		}

		order.Products = products
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over orders", http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}
func (app *application) showOrdersPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("show.html")

	if err != nil {
		app.errLog.Println(err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = tmpl.ExecuteTemplate(w, "show.html", nil)
	if err != nil {
		app.errLog.Println(err)
		return
	}
}
