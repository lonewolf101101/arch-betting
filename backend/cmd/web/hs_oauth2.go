package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/lonewolf101101/Architect-betting/backend/cmd/web/app"
	"github.com/lonewolf101101/Architect-betting/backend/common/easyOAuth2"
	"github.com/lonewolf101101/Architect-betting/backend/common/oapi"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/customerman"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/supersetman"
)

func oauthLogin(oauthClient *easyOAuth2.EasyOAuthClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := oauthClient.RedirectToLogin(w, r); err != nil {
			handleOAuthError(w, r, fmt.Sprintf("%v %v %v", oauthClient.Name, "oauth2 login error:", err))
			return
		}
	}
}

func oauthCallback(oauthClient *easyOAuth2.EasyOAuthClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := oauthClient.HandleCallback(w, r)
		if err != nil {
			handleOAuthError(w, r, fmt.Sprintf("%v %v %v", oauthClient.Name, "oauth2 callback error:", err))
			return
		}

		data, err := oauthClient.GetUserInfo(token.AccessToken)
		if err != nil {
			handleOAuthError(w, r, fmt.Sprintf("%v %v %v", oauthClient.Name, "oauth2 callback error:", err))
			return
		}

		var customer *customerman.Customer
		switch oauthClient.Name {
		case "google":
			var userinfo *customerman.GoogleUserInfo
			if err := json.Unmarshal(data, &userinfo); err != nil {
				handleOAuthError(w, r, fmt.Sprintf("google unmarshal error: %v data: %v", err, string(data)))
				return
			}
			if userinfo.ID == "" {
				handleOAuthError(w, r, fmt.Sprintf("google userinfo had empty ID. Data: %v", string(data)))
				return
			}

			if !userinfo.VerifiedEmail {
				handleOAuthError(w, r, fmt.Sprintf("google userinfo email not verified. Data: %v", string(data)))
				return
			}

			if userinfo.Email == "" {
				handleOAuthError(w, r, fmt.Sprintf("google userinfo had empty email. Data: %v", string(data)))
				return
			}

			customer, err = app.Customers.GetWithEmail(userinfo.Email)
			if err != nil {
				if !errors.Is(err, customerman.ErrNotFound) {
					handleOAuthError(w, r, fmt.Sprintf("%v %v %v", oauthClient.Name, "oauth2 callback error:", err))
					return
				}
				customer = &customerman.Customer{
					GoogleID:       userinfo.ID,
					Name:           userinfo.Name,
					ProfilePicture: userinfo.Picture,
					Email:          userinfo.Email,
				}
				if _, err := app.Customers.Save(customer); err != nil {
					handleOAuthError(w, r, fmt.Sprintf("%v %v %v", oauthClient.Name, "oauth2 callback error:", err))
					return
				}
				customerCreated := supersetman.CustomerCreated{
					Action:   "customer_created",
					Customer: *customer,
				}

				go func() {
					if err := app.KafkaServer.Push(app.Config.Topic, customerCreated, app.Config.Topic, nil); err != nil {
						app.ErrorLog.Println("Failed to produce customer event:", err)
					}
				}()
			}
		default:
			oapi.ServerError(w, fmt.Errorf("invalid oauth2 provider: %v", oauthClient.Name))
			return
		}
		go func() {
			if err := app.KafkaServer.Push(app.Config.Topic, supersetman.CustomerLogin{Action: "customer_login", Customer: *customer}, app.Config.Topic, nil); err != nil {
				app.ErrorLog.Println("Failed to produce customer event:", err)
			}
		}()

		app.Session.Put(r, "email", customer.Email)
		app.Session.Put(r, "auth_user_id", customer.ID)
		app.Session.Put(r, "oauth2_provider_name", oauthClient.Name)
		http.Redirect(w, r, "http://localhost:3000/", http.StatusTemporaryRedirect)
	}
}

func handleOAuthError(w http.ResponseWriter, r *http.Request, errorStr string) {
	app.ErrorLog.Println(errorStr)
	Error := &supersetman.Error{
		Action:  "oauth2_error",
		Message: errorStr,
	}

	go func() {
		if err := app.KafkaServer.Push(app.Config.Topic, Error, app.Config.Topic, nil); err != nil {
			app.ErrorLog.Println("failed to push document to kafka", err)
		}
	}()

	http.Redirect(w, r, "http://localhost:3000/", http.StatusTemporaryRedirect)
}
