package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	gitlab "github.com/xanzy/go-gitlab"

	authentication "k8s.io/client-go/pkg/apis/authentication/v1beta1"
)

func main() {
	//https://gitlab.com/api/v4
	gitlabUrl := os.Getenv("GITLAB_URL")
	http.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var tr authentication.TokenReview
		err := decoder.Decode(&tr)
		if err != nil {
			log.Println("[Error]", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"apiVersion": "authentication.k8s.io/v1beta1",
				"kind":       "TokenReview",
				"status": authentication.TokenReviewStatus{
					Authenticated: false,
				},
			})
			return
		}
		// Check User
		git := gitlab.NewClient(nil, tr.Spec.Token)
		git.SetBaseURL(gitlabUrl)
		user, _, err := git.Users.CurrentUser(nil)
		if err != nil {
			log.Println("[Error]", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"apiVersion": "authentication.k8s.io/v1beta1",
				"kind":       "TokenReview",
				"status": authentication.TokenReviewStatus{
					Authenticated: false,
				},
			})
			return
		}
		projects, _, err := git.Groups.ListGroups(&gitlab.ListGroupsOptions{})
		var groups []string
		for _, g := range projects {
			groups = append(groups, g.Name)
		}

		log.Printf("[Success] login as %s, groups: %v", user.Username, groups)
		w.WriteHeader(http.StatusOK)
		trs := authentication.TokenReviewStatus{
			Authenticated: true,
			User: authentication.UserInfo{
				Username: user.Username,
				UID:      user.Username,
				Groups:   groups,
			},
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"apiVersion": "authentication.k8s.io/v1beta1",
			"kind":       "TokenReview",
			"status":     trs,
		})
	})
	log.Fatal(http.ListenAndServe(":3000", nil))
}
