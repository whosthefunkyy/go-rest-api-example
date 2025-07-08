package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"database/sql"
	"context"
	"errors"
	

    "myapi/utils"
	"myapi/models"
	 "myapi/db"
	 "myapi/hateoas"
)


	// GET ALL USERS
func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.QueryContext(r.Context(),"SELECT id, name, age FROM users")
	if err != nil {
	// Timeout
	 if errors.Is(err, context.DeadlineExceeded) ||
	    errors.Is(r.Context().Err(), context.DeadlineExceeded) {
		utils.SendError(w, "Request timeout", http.StatusGatewayTimeout)
		return
}
	// User not found
		utils.SendError(w, "User not found", http.StatusNotFound)
		return

}
	defer rows.Close()
	var users []models.User
	
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			utils.SendError(w, err.Error(), http.StatusInternalServerError)
			return
	}
	users = append(users, u)
}
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(users)
}
	// GET 1 USER by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var u models.User
	err = db.DB.QueryRowContext(r.Context(),"SELECT id,name,age FROM users WHERE id = $1",id).Scan(&u.ID,  &u.Name,&u.Age)
		if err != nil {
	// Timeout
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(r.Context().Err(), context.DeadlineExceeded) {
			http.Error(w, "Request timeout", http.StatusGatewayTimeout)
			return
		}
	// User not Found
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	// Any others error
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hateoas.CreateUserResponse(u))
}  

	// CREATE A NEW USER
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User

	json.NewDecoder(r.Body).Decode(&u)
	err := db.DB.QueryRowContext(r.Context(),"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id", u.Name,u.Age).Scan(&u.ID)
		if err != nil {
	// Timeout
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(r.Context().Err(), context.DeadlineExceeded) {
			http.Error(w, "Request timeout", http.StatusGatewayTimeout)
			return
		}
	// User not Found
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	// Any others error
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(hateoas.CreateUserResponse(u))
}

	// UPDATE INFO into USER
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err :=strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var u models.User
	json.NewDecoder(r.Body).Decode(&u)

	res, err := db.DB.ExecContext(r.Context(),"UPDATE users SET name= $1, age= $2 WHERE id = $3", u.Name, u.Age, id)
	
	if err != nil {
	// TimeOut
	if errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(r.Context().Err(), context.DeadlineExceeded) {
		http.Error(w, "Request timeout", http.StatusGatewayTimeout)
		return
	}
	// Any others error
	utils.SendError(w, err.Error(), http.StatusInternalServerError)
	return
}
	// Any Affected rows
	rows, err := res.RowsAffected()
	if err != nil {
	utils.SendError(w, "RowsAffected error: "+err.Error(), http.StatusInternalServerError)
	return
}
	if rows == 0 {	
		utils.SendError(w, "User not found", 404)
		return 
	}
	
	u.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hateoas.CreateUserResponse(u))
}

// DELETE USER
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res,err := db.DB.ExecContext(r.Context(),"DELETE FROM users WHERE id = $1",id)
	 if err != nil {
	// Timeout
	 if errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(r.Context().Err(), context.DeadlineExceeded) {
		http.Error(w, "Request timeout", http.StatusGatewayTimeout)
		return
	}
	// Any others error
	 utils.SendError(w, err.Error(), http.StatusInternalServerError)
	 return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		utils.SendError(w, "RowsAffected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	
	if rows == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

	// CHANGE USER 
func PatchUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendError(w, "Invalid ID", http.StatusBadRequest)
	}
	
	var u models.User
	err = db.DB.QueryRowContext(r.Context(),
					"SELECT id,name,age FROM users WHERE id = $1",id).
					Scan(&u.ID,&u.Name,&u.Age)
	// Timeout
		if errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(r.Context().Err(), context.DeadlineExceeded) {
		http.Error(w, "Request timeout", http.StatusGatewayTimeout)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var patchData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&patchData); err != nil {
			utils.SendError(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		for k := range patchData {
		if k != "name" && k != "age" {
			utils.SendError(w, "Unknown field: "+k, http.StatusBadRequest)
			return
		}
	}
		
		if name, ok := patchData["name"].(string); ok {
			u.Name = name
		}
		if age, ok := patchData["age"].(float64); ok {
			u.Age = int(age)
		}
		
		if u.Name == " " {
			utils.SendError(w, "Name can't be empty", http.StatusBadRequest)
			return
		}
		
	res, err := db.DB.ExecContext(r.Context(),
		`UPDATE users SET name=$1, age=$2 WHERE id=$3`,
		u.Name, u.Age, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(r.Context().Err(), context.DeadlineExceeded) {
			http.Error(w, "Request timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
				
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(hateoas.CreateUserResponse(u))
	}




