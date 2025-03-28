package configs

import (
	"encoding/json"
	"fmt"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioClient struct {
	*twilio.RestClient
	PhoneNumber string
}

func NewTwilioClient(t Twilio) *TwilioClient {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: t.AccountSid,
		Password: t.AuthToken,
	})

	return &TwilioClient{
		client,
		t.PhoneNum,
	}
}

func (client *TwilioClient) SendSMS(to, message string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetFrom(client.PhoneNumber)
	params.SetTo(to)
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}

	return err
}
