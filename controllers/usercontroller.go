package controllers

import (
	"context"
	"encoding/json"
	"log"
	"mux-mongo-api/configs"
	"mux-mongo-api/helpers"
	"mux-mongo-api/models"
	"mux-mongo-api/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user models.User
	defer cancel()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&user)
	if err != nil {
		newUser := models.User{
			Id:        primitive.NewObjectID(),
			Name:      user.Name,
			Email:     user.Email,
			Company:   user.Company,
			Password:  helpers.GenerateHash(user.Password),
			Role:      user.Role,
			IsActive:  false,
			TsCreated: time.Now(),
			TsUpdated: time.Now(),
		}
		err := userCollection.FindOne(ctx, bson.M{"role": user.Role}).Decode(&user)
		if err != nil {
			result, err := userCollection.InsertOne(ctx, newUser)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}
			w.WriteHeader(http.StatusCreated)
			response := responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusConflict)
		response := responses.UserResponse{Status: http.StatusConflict, Message: "error", Data: map[string]interface{}{"data": "Cannot create more superadmins."}}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusConflict)
	response := responses.UserResponse{Status: http.StatusConflict, Message: "fail", Data: map[string]interface{}{"data": "email already exists"}}
	json.NewEncoder(w).Encode(response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	params := mux.Vars(r)
	userId := params["userId"]

	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)
	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}}
	json.NewEncoder(w).Encode(response)

}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
		}
		users = append(users, singleUser)

	}
	w.WriteHeader(http.StatusOK)
	response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}}
	json.NewEncoder(w).Encode(response)

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	params := mux.Vars(r)
	userId := params["userId"]
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})
	log.Println(result.DeletedCount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}
	if result.DeletedCount < 1 {
		w.WriteHeader(http.StatusNotFound)
		response := responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "user with given id not found"}}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "user deleted successfully"}}
	json.NewEncoder(w).Encode(response)
}

func ActivateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	params := mux.Vars(r)
	userId := params["userId"]
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(userId)
	log.Println(objId)
	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}
	status := helpers.CheckEmail(user.Email, user.Name)
	if status == true {
		w.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	response := responses.UserResponse{Status: http.StatusBadRequest, Message: "failed", Data: map[string]interface{}{"data": "email id is invalid"}}
	json.NewEncoder(w).Encode(response)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	email := r.Header.Get("email")
	password := r.Header.Get("password")
	defer cancel()
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}
	token := helpers.GenerateToken(user.Email)
	status := helpers.ValidateHash(user.Password, password)
	if status {
		w.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"details": user, "token": token}}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	response := responses.UserResponse{Status: http.StatusUnauthorized, Message: "failed", Data: map[string]interface{}{"data": "invalid credentials"}}
	json.NewEncoder(w).Encode(response)

}
