package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/antoniofmoliveira/apis/internal/dto"
	"github.com/antoniofmoliveira/apis/internal/entity"
	"github.com/antoniofmoliveira/apis/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDB database.UserRepositoryInterface
}

func NewUserHandler(userDB database.UserRepositoryInterface) *UserHandler {
	return &UserHandler{UserDB: userDB}
}

// Get Jwt godoc
// @Summary      Get Jwt
// @Description  Get Jwt
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input  body      dto.GetJWTInput  true  "user request"
// @Success      200     {object}  dto.AccessToken
// @Failure      400     {object}  Error
// @Failure      500     {object}  Error
// @Router       /users/generate_token [post]
func (h *UserHandler) GetJwt(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)

	var userdto dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&userdto)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var entityUser *entity.User
	entityUser, err = h.UserDB.FindByEmail(userdto.Email)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !entityUser.ValidatePassword(userdto.Password) {
		json.NewEncoder(w).Encode(Error{Message: "Invalid credentials"})
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": entityUser.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	accessToken := dto.AccessToken{
		AccessToken: tokenString,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)

}

// Create User godoc
// @Summary      Create a new user
// @Description  Create a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input  body      dto.CreateUserInput  true  "user request"
// @Success      201
// @Failure      400     {object}  Error
// @Failure      500     {object}  Error
// @Router       /users [post]
// @Security     ApiKeyAuth
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userdto dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&userdto)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var entityUser *entity.User
	entityUser, err = entity.NewUser(userdto.Name, userdto.Email, userdto.Password)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.UserDB.Create(entityUser)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// @Summary      Find user by email
// @Description  Find user by email
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        email   query     string  true  "User email"
// @Success      200  {object}  entity.User
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /users [get]
// @Security     ApiKeyAuth
func (h *UserHandler) FindByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		json.NewEncoder(w).Encode(Error{Message: "Email is required"})
		http.Error(w, "Email is required", http.StatusBadRequest)
	}
	user, err := h.UserDB.FindByEmail(email)
	if err != nil {
		if err.Error() == "record not found" {
			json.NewEncoder(w).Encode(Error{Message: "User not found"})
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
