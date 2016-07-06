package main

type EmailConfig struct {
	Sender     Sender
	Recipients []string
}

type Sender struct {
	Username string
	Password string
}