package Emp

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Emp struct {
	Id    string `json:"Id"`
	Name  string `json:"Name"`
	Phone string `json:"Phone"`
	Dept  Dept   `json:"Dept"`
}

type Dept struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

var DB *sql.DB

func PostEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("mysql", "root:Saikiran@18@tcp(127.0.0.1:3306)/employeedata")
	if err != nil {
		w.WriteHeader(404)
		return
	}
	var e Emp
	emp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
	}
	json.Unmarshal(emp, &e)
	_, er := db.Exec("insert into Employee values (?,?,?,?)", e.Id, e.Name, e.Dept.Id, e.Phone)
	if er != nil {
		log.Println(er)
		return
	}
	w.WriteHeader(200)
}

func GetEmployees(db *sql.DB) ([]Emp, error) {
	rows, err := db.Query("select e.id,e.name,e.phone,e.deptid,d.name from Employee e join Department d on e.deptid = d.id")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var employees []Emp

	for rows.Next() {
		var e Emp
		err = rows.Scan(&e.Id, &e.Name, &e.Phone, &e.Dept.Id, &e.Dept.Name)
		if err != nil {
			return nil, err
		}

		employees = append(employees, e)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func GetEmployee(db *sql.DB, id string) (Emp, error) {
	var e Emp
	row := db.QueryRow("select e.id,e.name,e.phone,e.deptid,d.name from Employee e join Department d on e.deptid = d.id where e.id = ?", id)
	err := row.Scan(&e.Id, &e.Name, &e.Phone, &e.Dept.Id, &e.Dept.Name)
	if err != nil {
		return Emp{}, err
	}
	return e, nil
}

func GetEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	db, err := sql.Open("mysql",
		"root:Saikiran@18@tcp(127.0.0.1:3306)/employeedata")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	emp, err := GetEmployee(db, id)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(emp)
	respBody, _ := json.Marshal(emp)

	w.Write(respBody)
}

func GetAllEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("mysql", "root:Saikiran@18@tcp(127.0.01)/employeedata")
	if err != nil {
		log.Println(err)
		return
	}
	emp, er := GetEmployees(db)
	if er != nil {
		w.WriteHeader(400)
		log.Println(err)
	}
	resp, _ := json.Marshal(emp)
	w.Write(resp)
}
