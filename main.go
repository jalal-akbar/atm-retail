package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type Transaction struct {
	UserName string
	Nominal  int
	Action   string
}

var users []*User
var transactions []Transaction

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/transactions", getTransactions).Methods("GET")
	router.HandleFunc("/transactions/{name}", getTransactionsByName).Methods("GET")
	router.HandleFunc("/transactions", doTransaction).Methods("POST")

	fmt.Println("running")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(transactions)
}

func getTransactionsByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	user := findUser(name)

	if user != nil {
		userTransaction := make([]Transaction, 0)

		for _, tx := range transactions {
			if tx.UserName == user.Name {
				userTransaction = append(userTransaction, tx)
			}
		}
		json.NewEncoder(w).Encode(userTransaction)
		return
	}
	http.Error(w, "name not found", http.StatusBadRequest)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func doTransaction(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	name := params.Get("name")
	nominal, _ := strconv.Atoi(params.Get("nominal"))
	action := params.Get("action")

	user := findUser(name)

	if user == nil {
		user = &User{
			Name: name,
		}
		users = append(users, user)
	}
	if action == "TRANSFER" {
		recipientParam := params.Get("recipient")
		recipient := findUser(recipientParam)
		if recipient == nil {
			http.Error(w, "Recipient not found", http.StatusBadRequest)
			return
		}
		if user.Balance < nominal {
			http.Error(w, "Insufficient Balance", http.StatusBadRequest)
			return
		}
		user.Balance -= nominal
		recipient.Balance += nominal
	} else if action == "SAVING" {
		user.Balance += nominal
	} else if action == "WITHDRAW" {
		if user.Balance < nominal {
			http.Error(w, "Insufficient Balance", http.StatusBadRequest)
			return
		}
		user.Balance -= nominal
	} else {
		http.Error(w, "Invalid Action", http.StatusBadRequest)
	}
	//executeTransaction(user, nominal, action)
	transactions = append(transactions, Transaction{
		UserName: user.Name,
		Nominal:  nominal,
		Action:   action,
	})

	json.NewEncoder(w).Encode(nil)
}

func findUser(name string) *User {
	for _, user := range users {
		if user.Name == name {
			return user
		}
	}
	return nil
}

// func executeTransaction(user *User, nominal int, action string) {
// 	switch action {
// 	case "SAVING":
// 		user.Balance += nominal
// 	case "TRANSFER":
// 		// Implement transfer logic here
// 	case "WITHDRAW":
// 		user.Balance -= nominal
// 	default:
// 	}

// 	transactions = append(transactions, Transaction{
// 		UserName: user.Name,
// 		Nominal:  nominal,
// 		Action:   action,
// 	})
// }
