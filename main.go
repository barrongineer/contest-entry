package main

import (
	"github.com/sclevine/agouti"
	"fmt"
	"time"
	"io/ioutil"
	"encoding/json"
	"github.com/go-gomail/gomail"
	"crypto/tls"
)

func main() {
	fmt.Println("Running automated contest entry...")

	var emails []string

	driver := agouti.ChromeDriver()
	defer killTheBrowser(driver)

	driver.Start()

	page, err := driver.NewPage()
	if (err != nil) {
		handleError(err)
		return
	}

	for _, email := range emails {
		enterSweepstakes(page, "http://www.hgtv.com/design/hgtv-smart-home/sweepstakes", "ngxFrame48405", email)
		enterSweepstakes(page, "http://www.diynetwork.com/hgtv-smart-home", "ngxFrame48409", email)
	}

	sendSuccessEmail()

	fmt.Println("Completed automated contest entry.")
}

func sendSuccessEmail() {
	fmt.Println("Sending success emails...")

	emailConfig := getEmailConfig()
	message := gomail.NewMessage()
	message.SetHeader("From", emailConfig.Sender.Username)
	message.SetHeader("To", emailConfig.Recipients...)
	message.SetHeader("Subject", "Smart Home Contest Entry Success")
	message.SetBody("text/html", "Successfully entered the smart home contest.")

	mailClient := gomail.NewDialer("smtp.gmail.com", 587, emailConfig.Sender.Username, emailConfig.Sender.Password)
	mailClient.TLSConfig = &tls.Config{InsecureSkipVerify:true}
	if err := mailClient.DialAndSend(message); err != nil {
		panic(err)
	}

	fmt.Println("Successfully sent success emails.")
}

func getEmailConfig() (emailConfig EmailConfig) {
	emailConfigFile, err := ioutil.ReadFile("./email_config.json")
	if err != nil {
		fmt.Println("Error reading the file.")
		handleError(err)
		return
	}

	json.Unmarshal(emailConfigFile, &emailConfig)

	return emailConfig
}

func enterSweepstakes(page *agouti.Page, url string, frame string, email string) {
	var err error
	err = page.Navigate(url)
	if err != nil {
		handleError(err)
		return
	}

	err = page.FindByID(frame).SwitchToFrame()
	if err != nil {
		handleError(err)
		return
	}

	err = page.FindByID("xReturningUserEmail").Fill(email)
	if err != nil {
		handleError(err)
		return
	}

	err = page.FindByID("xCheckUser").Click()
	if err != nil {
		handleError(err)
		return
	}

	time.Sleep(5 * time.Second)

	err = page.SwitchToRootFrame()
	if err != nil {
		handleError(err)
		return
	}
	err = page.FindByID(frame).SwitchToFrame()
	if err != nil {
		handleError(err)
		return
	}

	form := page.FindByID("xSecondaryForm")
	fmt.Println(form.Count())

	container := form.FindByID("xSubmitContainer")
	fmt.Println(container.Count())

	button := container.Find("button")
	fmt.Println(button.Count())

	err = button.Click()
	if err != nil {
		handleError(err)
		return
	}

	time.Sleep(5 * time.Second)
}

func handleError(err error) {
	fmt.Println("There was an error:")
	fmt.Println(err.Error())
}

func killTheBrowser(driver *agouti.WebDriver) {
	err := driver.Stop()
	if err != nil {
		handleError(err)
		return
	}
}