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
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var userSessionCollection *mongo.Collection = configs.GetCollection(configs.DB, "usersession")

var validate = validator.New()

func Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	w.Header().Set("Content-Type", "application/json")
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

func CreateAdmin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	w.Header().Set("Content-Type", "application/json")
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
			Role:      "admin",
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
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
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
	response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"users": users}}
	json.NewEncoder(w).Encode(response)

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	w.Header().Set("Content-Type", "application/json")
	requestEmail := r.Context().Value("user-id")
	log.Println(requestEmail)

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
	w.Header().Set("Content-Type", "application/json")
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
	if status {
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
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	var session models.UserSession
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

	status := helpers.ValidateHash(user.Password, password)
	if status {
		token := helpers.GenerateToken(user.Email)
		refresh := helpers.GenerateRefreshToken(user.Email)
		session.AccessToken = token
		session.RefreshToken = refresh
		session.UserAgent = r.Header.Get("User-Agent")
		session.TsCreated = time.Now()

		_, err = userSessionCollection.InsertOne(ctx, session)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"details": user, "access-token": token, "refresh-token": refresh}}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	response := responses.UserResponse{Status: http.StatusUnauthorized, Message: "failed", Data: map[string]interface{}{"data": "invalid credentials"}}
	json.NewEncoder(w).Encode(response)

}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	w.Header().Set("Content-Type", "application/json")
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	requestEmail := r.Context().Value("user-id")
	var user models.User
	var session models.UserSession

	defer cancel()
	err := userCollection.FindOne(ctx, bson.M{"email": requestEmail}).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Println(authHeader[1])
	err = userSessionCollection.FindOne(ctx, bson.M{"refreshtoken": authHeader[1]}).Decode(&session)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}
	var token = helpers.GenerateToken(user.Email)

	result, err := userSessionCollection.UpdateOne(
		ctx,
		bson.M{"_id": session.Id},
		bson.D{
			{"$set", bson.D{{"accesstoken", token}}},
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Println(result)

	w.WriteHeader(http.StatusOK)
	response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"details": user, "access-token": token, "refresh-token": session.RefreshToken}}
	json.NewEncoder(w).Encode(response)
}
