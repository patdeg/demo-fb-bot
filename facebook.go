package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"net/http"
)

// Facebook Callback Handler
func FacebookCallbackHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	log.Debugf(c, ">>> Facebook Callback Handler")

	if r.Method == "GET" {
		FacebookCallbackGETHandler(w, r)
	} else if r.Method == "POST" {
		FacebookCallbackPOSTHandler(w, r)
	} else {
		log.Errorf(c, "Error, unkown method: %v", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

}

// Facebook Callback Handler for a GET request
// Usually from a page subscription in the developer dashboard
func FacebookCallbackGETHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	log.Debugf(c, ">>>> FacebookCallbackGETHandler")

	mode := r.FormValue("hub.mode")
	log.Debugf(c, "Hub Mode: %v", mode)

	challenge := r.FormValue("hub.challenge")
	log.Debugf(c, "Hub Challenge: %v", challenge)

	verify_token := r.FormValue("hub.verify_token")
	log.Debugf(c, "Hub Verify Token: %v", verify_token)

	if verify_token != VERIFY_TOKEN {
		log.Errorf(c, "Error, bad verification token: %v", verify_token)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if mode != "subscribe" {
		log.Errorf(c, "Error, bad mode: %v", mode)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%v", challenge)
}

// Facebook Callback Handler for a POST request
// Usually from an incoming user message, or a message delivery confirmation
func FacebookCallbackPOSTHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	log.Debugf(c, ">>>> FacebookCallbackPOSTHandler")

	var payload FacebookPayload
	err := UnmarshalRequest(c, r, &payload)
	if err != nil {
		log.Errorf(c, "Error reading JSON: %v", err)
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if payload.Object == "page" {
		log.Debugf(c, "Received %v entries", len(payload.Entry))

		for i, e := range payload.Entry {
			log.Debugf(c, "Entry #%v", i)
			log.Debugf(c, "Message on page %v", e.Id)
			log.Debugf(c, "Time: %v", e.Time)

			if len(e.Messaging) > 0 {
				log.Debugf(c, "Received %v messages", len(e.Messaging))

				for _, m := range e.Messaging {
					log.Debugf(c, "Sender %v Recipient %v", m.Sender.Id, m.Recipient.Id)

					if m.Message != nil {
						log.Debugf(c, " --- Message Received")
						log.Debugf(c, "Message Id: %v", m.Message.Mid)
						log.Debugf(c, "Message Seq: %v", m.Message.Seq)
						log.Debugf(c, "Message: %v", m.Message.Text)

						response := GetResponse(c, m.Sender.Id, m.Message.Text)

						err = SendFacebookMessage(c, m.Sender.Id, response)
						if err != nil {
							log.Errorf(c, "Error, sending message: %v", err)
						}
					}

					if m.Delivery != nil {
						log.Debugf(c, " --- Message Delivered")
						log.Debugf(c, "Message Ids: %v", m.Delivery.Mids)
						log.Debugf(c, "Message Seq: %v", m.Delivery.Seq)
						log.Debugf(c, "Watermark: %v", m.Delivery.Watermark)
					}

				}
			} else {
				log.Warningf(c, "No message in payload - feature not developed!")
			}
		}

	} else {
		log.Errorf(c, "Error, unkown Object type: %v", payload.Object)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "")
}

// Provide the reponse of the bot from the message of a Facebook user
// In this demo example, we just add the word "Hello" in front of the message
func GetResponse(c context.Context, facebookUser string, message string) string {
	response := "Hello " + message
	return response
}

// Send a message to a Facebook User
// Per Facebook Policy, this CAN NOT be marketing/promotional contents and should only
// send organic content
func SendFacebookMessage(c context.Context, facebookUser string, message string) error {

	messaging := FacebookMessaging{
		Recipient: &FacebookPerson{
			Id: facebookUser,
		},
		Message: &FacebookMessage{
			Text: message,
		},
	}

	jsonData, err := json.Marshal(messaging)
	if err != nil {
		log.Errorf(c, "Error converting to JSON: %v", err)
		return err
	}

	URL := "https://graph.facebook.com/v2.6/me/messages?access_token=" + PAGE_ACCESS_TOKEN
	log.Debugf(c, "Calling %v", URL)
	resp, err := urlfetch.Client(c).Post(URL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		log.Errorf(c, "Error posting message: %v", err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf(c, "Error reading response: %v", err)
		return err
	}

	log.Debugf(c, "Facebook response: %s", body)

	return nil

}
