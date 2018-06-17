package main

import (
	"net/url"
	"errors"
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

var (
	VOLUNTEERO_AUTH = "https://volunteero-auth.herokuapp.com/auth"
)

type Role struct {
	// role name
	Title string `json:"title"`
	// access level - e.g admin or so
	Level string `json:"level"`
	OrganizationName string `json:"organizationName"`
	OrganizationID string `json:"organizationID"`
}

type Roles []Role

type ErrorResponse struct {
	Message string `json:"message"`
}

type Response struct {
	Roles Roles `json:"roles"`
}

func infoEndpoint(w http.ResponseWriter, r *http.Request){
	message := "Info Endpoint Hit"
	log.Print(message)
	fmt.Fprintf(w, message)
}

func resolveAccessToken(r *http.Request) string{
	tokens, ok := r.URL.Query()["accessToken"]
	if !ok || len(tokens) < 1 {
		log.Println("Url Param 'accessToken' is missing")
		return ""
	}
	token := tokens[0]

	log.Println("Got the token: " + token)
	return token
}

func resolveRoles(token string) Response {
	// mock stuff
	roles := Roles{
		Role{Title:"cleric", Level:"master", OrganizationName:"Dungeons", OrganizationID:"123"},
	}
	return Response{Roles:roles}
}

func handleNoToken(w http.ResponseWriter) {
	errorResponse := ErrorResponse{
		Message:errors.New("please, provide a token").Error(),
	}
	resp, _ := json.Marshal(errorResponse)
	http.Error(w, string(resp), http.StatusUnauthorized)
}

func getRolesFromAuth(token string) (*http.Response, error) {
	payload := url.Values{}
	payload.Add("accessToken",token)
	req, err := http.Get(VOLUNTEERO_AUTH +"/roles?" + payload.Encode())
	log.Println("got roles response")
	log.Print(req)
	return req, err
}


func getRoles(w http.ResponseWriter, r *http.Request){
	message := "Get Roles Endpoint Hit"
	log.Print(message)
	
	token := resolveAccessToken(r)

	if(token == ""){
		handleNoToken(w)
		return 
	}

	getRolesFromAuth(token)
	response := resolveRoles(token)
	fmt.Println(message)
	json.NewEncoder(w).Encode(response)
}



func handleRequests() {
	// Map a callback `infoEndpoint` to the `/` route
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", infoEndpoint)
	myRouter.HandleFunc("/roles", getRoles).Methods("GET")
	// Log any errros that happen when we serve and start listening
	log.Print("Starting the server!")
	log.Fatal(http.ListenAndServe(":8001", myRouter))
}

func main() { 
	handleRequests()
}