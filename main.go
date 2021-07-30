package main

import (
	"fmt"
	"strings"
)

//Notifier ...
type Notifier interface {
	notify() bool
	setStatus(status bool)
	printStatus()
}

//SMSNotification ...
type SMSNotification struct {
	Phone   string
	Message string
	Status  string
}

func (not *SMSNotification) notify() bool {
	fmt.Printf("\nSending SMS to %s with text [%s]", not.Phone, not.Message)
	return true
}

func (not *SMSNotification) setStatus(sent bool) {
	status := "not delivered"
	if sent {
		status = "delivered"
	}
	not.Status = status
}

func (not *SMSNotification) printStatus() {
	fmt.Printf("\nSMS to %s status [%s]", not.Phone, not.Status)
}

//EmailNotification ...
type EmailNotification struct {
	From    string
	To      []string
	Subject string
	Body    string
	Status  string
}

func (not *EmailNotification) notify() bool {
	fmt.Println("\nSending Email with follwing detais")
	fmt.Println("Subject:", not.Subject)
	fmt.Println("From:", not.From)
	fmt.Println("To:", strings.Join(not.To, ","))
	fmt.Println("Body:", not.Body)
	return true
}

func (not *EmailNotification) setStatus(sent bool) {
	status := "not delivered"
	if sent {
		status = "delivered"
	}
	not.Status = status
}

func (not *EmailNotification) printStatus() {
	fmt.Printf("\nEmail to %s status [%s]", strings.Join(not.To, ","), not.Status)
}

func sendNotification(notifications <-chan Notifier) chan Notifier {
	output := make(chan Notifier, 0)
	go func() {
		for not := range notifications {
			sent := not.notify()
			not.setStatus(sent)

			output <- not
		}
		close(output)
	}()
	return output

}

func notificationGenerator(notifications []Notifier) chan Notifier {
	output := make(chan Notifier, 0)
	go func() {
		for _, not := range notifications {
			output <- not
		}
		close(output)
	}()
	return output
}

func main() {
	notifications := make([]Notifier, 0)

	sms := &SMSNotification{
		Phone:   "+1234567890",
		Message: "Hi",
	}
	email := &EmailNotification{
		From:    "from@test.com",
		To:      []string{"to@test.com"},
		Subject: "Welcome",
		Body:    "Good Morning",
	}

	notifications = append(notifications, sms, email)
	notChannels := notificationGenerator(notifications)
	statusList := sendNotification(notChannels)

	for stat := range statusList {
		stat.printStatus()
	}

}
