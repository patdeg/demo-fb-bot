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

type FacebookAttachment struct {
	/* image, audio, video, file or location */
	Type string `json:"type,omitempty"`

	/* multimedia or location payload */
	Payload string `json:"payload,omitempty"`
}

type FacebookQuickReply struct {
	/* Custom data provided by the app */
	Payload string `json:"payload,omitempty"`
}

type FacebookMessage struct {

	/* Indicates the message sent from the page itself */
	IsEcho bool `json:"is_echo,omitempty"`

	/* ID of the app from which the message was sent */
	AppId string `json:"app_id,omitempty"`

	/* Custom string passed to the Send API as the metadata field */
	Metadata string `json:"metadata,omitempty"`

	/* Message ID */
	Mid string `json:"mid,omitempty"`

	/* Message sequence number */
	Seq int64 `json:"seq,omitempty"`

	/* Text of message */
	Text string `json:"text,omitempty"`

	/* Array containing attachment data */
	Attachments *[]FacebookAttachment `json:"attachment,omitempty"`

	/* Optional custom data provided by the sending app */
	QuickReply *FacebookQuickReply `json:"quick_reply,omitempty"`
}

type FacebookPerson struct {
	Id string `json:"id,omitempty"`
}

type FacebookPostback struct {
	/* payload parameter that was defined with the button */
	Payload string `json:"payload,omitempty"`
}

type FacebookOptIn struct {
	/* data-ref parameter that was defined with the entry point */
	Ref string `json:"ref,omitempty"`
}

type FacebookAccountLinking struct {
	/* linked or unlinked */
	status string `json:"status,omitempty"`

	/* Value of pass-through authorization_code provided in the Linking Account flow */
	AuthorizationCode string `json:"authorization_code,omitempty"`
}

type FacebookDelivery struct {
	/* Array containing message IDs of messages that were delivered. Field may not be present. */
	Mids []string `json:"mids,omitempty"`

	/* All messages that were sent before this timestamp were delivered */
	Watermark int64 `json:"watermark,omitempty"`

	/* Sequence number */
	Seq int64 `json:"seq,omitempty"`
}

type FacebookRead struct {
	/* All messages that were sent before this timestamp were read */
	Watermark int64 `json:"watermark,omitempty"`

	/* Sequence number */
	Seq int64 `json:"seq,omitempty"`
}

type FacebookMessaging struct {
	/* Sender user ID */
	Sender *FacebookPerson `json:"sender,omitempty"`

	/* Recipient user ID */
	Recipient *FacebookPerson `json:"recipient,omitempty"`

	/* Time of messaging (epoch time in milliseconds) */
	Timestamp int64 `json:"timestamp,omitempty"`

	/* Message Received Payload */
	Message *FacebookMessage `json:"message,omitempty"`

	/* Postback Received Payload */
	Postback *FacebookPostback `json:"postback,omitempty"`

	/* Authentification Payload */
	OptIn *FacebookOptIn `json:"optin,omitempty"`

	/* Account Linking Payload */
	AccountLinking *FacebookAccountLinking `json:"account_linking,omitempty"`

	/* Message Delivered Payload */
	Delivery *FacebookDelivery `json:"delivery,omitempty"`

	/* Message Read Payload */
	Read *FacebookRead `json:"read,omitempty"`
}

type FacebookEntry struct {
	/* Page ID of page */
	Id string `json:"id,omitempty"`

	/* Time of update (epoch time in milliseconds) */
	Time int64 `json:"time,omitempty"`

	/* Array containing objects related to messaging */
	Messaging []FacebookMessaging `json:"messaging,omitempty"`
}

type FacebookPayload struct {
	Object string `json:"object,omitempty"`

	/* Array containing event data */
	Entry []FacebookEntry `json:"entry,omitempty"`
}

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

						err = TreatMessage(c, m.Sender.Id, m.Message.Text)
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

func TreatMessage(c context.Context, recipientId string, message string) error {
	response := "Hello " + message
	return SendFacebookMessage(c, recipientId, response)
}

func SendFacebookMessage(c context.Context, recipientId string, message string) error {

	messaging := FacebookMessaging{
		Recipient: &FacebookPerson{
			Id: recipientId,
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
