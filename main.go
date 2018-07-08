package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/oauth2/google"

	"golang.org/x/net/context"
	"google.golang.org/api/chat/v1"
	"google.golang.org/appengine" // Required external App Engine library
	"google.golang.org/appengine/log"
)

type GithubPayload struct {
	Action      string      `json:"action"`
	Number      int         `json:"number"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
}
type PullRequest struct {
	URL    string `json:"url"`
	ID     int    `json:"id"`
	User   User   `json:"user"`
	Body   string `json:"body"`
	Merged bool   `json:"merged"`
}
type User struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
}

type Repository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	// Set Headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	log.Infof(ctx, "Endpoint reached "+r.URL.Path)
	// Check Endpoint for Secure Endpoint
	if r.URL.Path != "/"+os.Getenv("SECURE_ENDPOINT") {
		http.Error(w, "Bad Request", http.StatusForbidden)
		return
	}

	// Check Key
	if r.URL.Query().Get("key") != os.Getenv("SECURE_KEY") {
		http.Error(w, "Bad Shared Key", http.StatusForbidden)
		return
	}

	// Set Context to appengine context

	// Read Body into Bytes Array
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Errorf(ctx, "Error Reading Body "+err.Error())
		http.Error(w, "Error Reading Body", http.StatusInternalServerError)
		return
	}
	log.Infof(ctx, "Body: %+v", string(b))
	var gp GithubPayload
	err = json.Unmarshal(b, &gp)
	if err != nil {
		log.Errorf(ctx, "Error Unmarshalling Github Payload", err)
		http.Error(w, "Error Reading Payload", http.StatusInternalServerError)
		return
	}
	log.Infof(ctx, "Check payload", gp)
	var action string
	if gp.PullRequest.Merged {
		action = "merged"
	} else {
		action = gp.Action
	}
	err = postToRoom(ctx, chat.Message{Text: gp.PullRequest.User.Login + " " + action + " a Pull Request on repo: " + gp.Repository.FullName + "\n" + gp.PullRequest.URL}, "AAAAV2Ons90", strconv.Itoa(gp.Number))
	if err != nil {
		log.Errorf(ctx, "Error Posting to Room", err)
		http.Error(w, "Error Sending Alert", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Success"))
}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main() // Starts the server to receive requests
}

// Helper Function to cut down on code redundancy
func postToRoom(ctx context.Context, payload chat.Message, space string, threadKey string) error {
	url := "https://chat.googleapis.com/v1/spaces/" + space + "/messages?threadKey=" + threadKey
	client, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/chat.bot")
	if err != nil {
		log.Errorf(ctx, "Error Getting Default Token Source", err)
		return err
	}
	body, err := json.Marshal(payload)
	resp, err := client.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(body))
	if err != nil {
		log.Infof(ctx, "Error In Post to Room %+v", err)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	log.Infof(ctx, "Byte to String %v", string(b))
	if err != nil {
		return err
	}

	return nil

}
