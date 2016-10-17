package main

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
