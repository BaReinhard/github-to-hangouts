package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/api/chat/v1"
	"google.golang.org/appengine" // Required external App Engine library
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type Message struct {
	Text        string       `json:"text"`
	User        UserObject   `json:"sender"`
	Annotations []Annotation `json:"annotations"`
}
type Annotation struct {
	Usermention UserMention `json:"userMention"`
}
type UserMention struct {
	User UserObject `json:"user"`
}
type UserObject struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
}

// Payload send from GChat
type Payload struct {
	Message Message   `json:"message"`
	User    chat.User `json:"user"`
	Space   Space     `json:"space"`
	Type    string    `json:"type"`
}

// Space Struct for Unmarshalling
type Space struct {
	Name string `json:"name"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Set Headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Set Context to appengine context
	ctx := appengine.NewContext(r)

	// Read Body into Bytes Array
	b, e := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if e != nil {
		json.NewEncoder(w).Encode(e)
	}
	log.Infof(ctx, "Body: %+v", string(b))
	var message Payload
	// Unmarshall Byte Array (b) into a Message Struct (message) to interact easier
	json.Unmarshal(b, &message)

	switch message.Space.Name {
	case "spaces/AAAA0c_TyMI":
		// Switch Case Space : spaces/AAAA0c_TyMI . Sends to JavaScript Cloud Function
		// Setup as a beginning codebase for beginners
		// Repo Here: https://github.com/BaReinhard/help-bot

		log.Infof(ctx, "Sending to Bot Dev Room")
		msg, err := postToRoom(ctx, "https://us-central1-uplifted-elixir-203119.cloudfunctions.net/helpBot", bytes.NewReader(b))
		if err != nil {
			// Log Error and Return An Error Message in a Chat Friendly Format
			log.Errorf(ctx, "An Error Occurred: ", err)
			json.NewEncoder(w).Encode(chat.Message{Text: "An error has occurred"})
		}
		log.Infof(ctx, "Returned from Bot Dev Room: %+v", msg)
		json.NewEncoder(w).Encode(msg)
	case "spaces/AAAAifGFyYk":
		log.Infof(ctx, "Sending to Python Room")
		msg, err := postToRoom(ctx, "https://python-bot-dot-uplifted-elixir-203119.appspot.com", bytes.NewReader(b))
		if err != nil {
			// Log Error and Return An Error Message in a Chat Friendly Format
			log.Errorf(ctx, "An Error Occurred: ", err)
			json.NewEncoder(w).Encode(chat.Message{Text: "An error has occurred"})
		}
		log.Infof(ctx, "Returned from Bot Dev Room: %+v", msg)
		json.NewEncoder(w).Encode(msg)

	case "spaces/AAAAcZ1PlGk":
		log.Infof(ctx, "Sending to Tacos Room")
		msg, err := postToRoom(ctx, "https://basic-taco-bot-dot-uplifted-elixir-203119.appspot.com", bytes.NewReader(b))
		if err != nil {
			// Log Error and Return An Error Message in a Chat Friendly Format
			log.Errorf(ctx, "An Error Occurred: ", err)
			json.NewEncoder(w).Encode(chat.Message{Text: "An error has occurred"})
		}
		log.Infof(ctx, "Returned from Tacos Room: %+v", msg)
		json.NewEncoder(w).Encode(msg)
	case "spaces/AAAAyXeUgAM", "spaces/AAAALPK7rTg":
		log.Infof(ctx, "Sending to Spotify Room")
		msg, err := postToRoom(ctx, "https://spotify-chat-dot-uplifted-elixir-203119.appspot.com/bot", bytes.NewReader(b))
		if err != nil {
			// Log Error and Return An Error Message in a Chat Friendly Format
			log.Errorf(ctx, "An Error Occurred: ", err)
			json.NewEncoder(w).Encode(chat.Message{Text: "An error has occurred"})
		}

		log.Infof(ctx, "Returned from Spotify Room: %+v", msg)
		json.NewEncoder(w).Encode(msg)
	// case "spaces/AAAA3dgDXKM":
	// 	log.Infof(ctx, "Sending to Hubot Room")
	// 	msg, err := postToRoom(ctx, "https://")
	default:
		// Default Switch Function, sends to Go Bot

		log.Infof(ctx, "Sending to Bot Development")
		msg, err := postToRoom(ctx, "https://bitmoji-bot-dot-uplifted-elixir-203119.appspot.com", bytes.NewReader(b))
		if err != nil {
			// Log Error and Return An Error Message in a Chat Friendly Format
			log.Errorf(ctx, "An Error Occurred: ", err)
			json.NewEncoder(w).Encode(chat.Message{Text: "An error has occurred"})
		}
		log.Infof(ctx, "Returned from Bot Development: %+v", msg)
		json.NewEncoder(w).Encode(msg)
	}

}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main() // Starts the server to receive requests
}

// Helper Function to cut down on code redundancy
func postToRoom(ctx context.Context, url string, body io.Reader) (chat.Message, error) {
	var br chat.Message

	// Use urlfetch in App Engine
	client := urlfetch.Client(ctx)
	resp, err := client.Post(url, "application/json; charset=utf-8", body)
	if err != nil {
		log.Infof(ctx, "Error In Post to Room %+v", err)
		return br, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	log.Infof(ctx, "Byte to String %v", string(b))
	if err != nil {
		return br, err
	}
	err = json.Unmarshal(b, &br)
	if err != nil {
		return br, err
	}
	return br, nil

}
