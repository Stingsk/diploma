package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Stingsk/diploma/internal/repository/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func UserRegisterHandler(userStore users.Store, auth *jwtauth.JWTAuth) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", userRegisterHandler(userStore, auth))
	}
}

func UserLoginHandler(userStore users.Store, auth *jwtauth.JWTAuth) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", userLoginHandler(userStore, auth))
	}
}

func userRegisterHandler(store users.Store, auth *jwtauth.JWTAuth) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
		defer requestCancel()

		var user users.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot decode provided data: %q", err), http.StatusBadRequest)

			return
		}

		err = user.Validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		err = store.CreateUser(requestContext, user.Login, user.Password)
		if errors.Is(err, users.ErrUserExists) {
			http.Error(
				w,
				err.Error(),
				http.StatusConflict,
			)

			return
		}
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Something went wrong during user create: %q", err),
				http.StatusInternalServerError,
			)

			return
		}

		userToken, err := getUserToken(user.Login, auth)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Something went wrong during user create: %q", err),
				http.StatusInternalServerError,
			)

			return
		}
		w.Header().Set("Authorization", "Bearer "+userToken)
		w.WriteHeader(http.StatusOK)
	}
}

func userLoginHandler(userStore users.Store, authToken *jwtauth.JWTAuth) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestContext, requestCancel := context.WithTimeout(r.Context(), requestTimeout)
		defer requestCancel()

		var user users.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot decode provided data: %q", err), http.StatusBadRequest)

			return
		}

		err = user.Validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		err = userStore.ValidateUser(requestContext, user.Login, user.Password)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Unauthorized: %q", err),
				http.StatusUnauthorized,
			)

			return
		}

		userToken, err := getUserToken(user.Login, authToken)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Something went wrong during user login: %q", err),
				http.StatusInternalServerError,
			)

			return
		}
		w.Header().Set("Authorization", "Bearer "+userToken)
		w.WriteHeader(http.StatusOK)
	}
}

func getUserToken(login string, auth *jwtauth.JWTAuth) (string, error) {
	_, tokenString, err := auth.Encode(map[string]interface{}{"login": login})

	return tokenString, err
}
