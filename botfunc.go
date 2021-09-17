package modules

import (
	"database/sql"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type slice struct {
	field map[string]string
}

var Database, _ = sql.Open("mysql", "username:password@tcp(URL:PORT)/DBNAME")

func CreateWebhook() *discordgo.WebhookParams {

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "NuggieACO",
		},
		Title:       "Welcome to NuggieACO!",
		Description: "Create a ticket by reacting", //can change upon different releases
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Property of Nuggie Labs",
			IconURL: "<discord webhook>",
		},
		Color: 14902784,
	}
	var id []*discordgo.MessageEmbed
	id = append(id, embed)
	webhook := discordgo.WebhookParams{
		Username: "NuggieACO",
		Embeds:   id,
	}

	return &webhook
}

func CustomWebhook(title string, description string, fields slice, dg *discordgo.Session, webhook *discordgo.Webhook) (*discordgo.Message, error) {
	var field []*discordgo.MessageEmbedField
	for titl, desc := range fields.field {
		temp := &discordgo.MessageEmbedField{
			Name:  titl,
			Value: desc,
		}
		field = append(field, temp)
	}
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "NuggieACO",
		},
		Title:       title,
		Description: description, //can change upon different releases
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Property of Nuggie Labs",
			IconURL: "<discord webhook>",
		},
		Color: 14902784,
	}
	var id []*discordgo.MessageEmbed
	id = append(id, embed)
	web := &discordgo.WebhookParams{
		Username: "NuggieACO",
		Embeds:   id,
	}

	result, err := dg.WebhookExecute(webhook.ID, webhook.Token, true, web)

	return result, err
}

func InfoWebhook(title string, password string, fields map[string]string, file []*discordgo.MessageAttachment, dg *discordgo.Session) {
	var field []*discordgo.MessageEmbedField
	for titl, desc := range fields {
		temp := &discordgo.MessageEmbedField{
			Name:  titl,
			Value: desc,
		}
		field = append(field, temp)
	}
	temp := &discordgo.MessageEmbedField{
		Name:  "Profiles",
		Value: "[Link](" + file[0].URL + ")",
	}
	field = append(field, temp)

	temp = &discordgo.MessageEmbedField{
		Name:  "Password For Profiles",
		Value: "```" + password + "```",
	}
	field = append(field, temp)

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "NuggieACO",
		},
		Title:  title,
		Fields: field,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Property of Nuggie Labs",
			IconURL: "",
		},
		Color: 14902784,
	}
	var id []*discordgo.MessageEmbed
	id = append(id, embed)
	web := &discordgo.WebhookParams{
		Username: "NuggieACO",
		Embeds:   id,
	}

	_, err := dg.WebhookExecute("<discord webhook>", "<discord webhook>", true, web)
	if err != nil {
		log.Fatal(err)
	}

}

func SendWebhook(webhook *discordgo.WebhookParams, dg *discordgo.Session) *discordgo.Message {
	res, _ := dg.WebhookExecute("<discord webhook>", "<discord webhook>", true, webhook)
	if res != nil {
		println("success")
	}

	err := dg.MessageReactionAdd(res.ChannelID, res.ID, "🎟️")
	if err != nil {
		println("error")
		log.Fatal(err)
	}

	return res
}

func HandReaction(dg *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.MessageID == "880917833503100969" && r.Emoji.Name == "🎟️" && r.UserID != "844677813285945345" {
		println(r.UserID)
		err := dg.MessageReactionRemove("880831881829031956", "880917833503100969", "🎟️", r.UserID)
		if err != nil {
			println(err)
			time.Sleep(400 * time.Millisecond)
			return
		}

		perms := []*discordgo.PermissionOverwrite{}
		temp := &discordgo.PermissionOverwrite{ //USER CAN SEE CHANNEL
			ID:    r.UserID,
			Type:  1,
			Allow: 0x0000000400,
		}
		perms = append(perms, temp)
		temp = &discordgo.PermissionOverwrite{ //EVERYONE ELSE CANNOT SEE
			ID:   "439590627429056512",
			Type: 0,
			Deny: 0x0000000400,
		}
		perms = append(perms, temp)
		//CREATE PERM FOR STAFF TO SEE
		user, err := dg.GuildMember("439590627429056512", r.UserID)
		if err != nil {
			println("user not present?")
			return
		}

		channel := discordgo.GuildChannelCreateData{
			Name:                 user.User.Username,
			Type:                 0,
			PermissionOverwrites: perms,
			ParentID:             "880831831648403576",
		}

		result, err := dg.GuildChannelCreateComplex("439590627429056512", channel)
		if err != nil {
			log.Fatal(err)
		}

		hook, err := dg.WebhookCreate(result.ID, "NuggieAco", "https://cdn.discordapp.com/emojis/781803541060386826.png?v=1")
		dg.ChannelMessageSend(result.ID, "<@"+r.UserID+">")

		//start steps
		GetInfo(result.ID, hook, dg)
	}
}

func HandleClose(dg *discordgo.Session, m *discordgo.MessageCreate) bool {
	//add staff role only
	if strings.Compare(m.Content, "-close") == 0 {
		channel, err := dg.Channel(m.ChannelID)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Compare(channel.ParentID, "880831831648403576") == 0 { //check if it is from correct channel
			_, err = dg.ChannelDelete(channel.ID)
			if err != nil {
				log.Fatal(err)
			}
			return true
		}
	}
	return false
}

func AwaitClose(curr *discordgo.Message, channelid string, dg *discordgo.Session) bool {
	var messages []*discordgo.Message
	messages, _ = dg.ChannelMessages(channelid, 100, "", curr.ID, "")
	if messages == nil {
		return false
	}
	if len(messages) < 1 {
		return false
	} else {
		current := messages[0]
		if strings.Compare(current.Content, "-close") == 0 && strings.Compare(current.ChannelID, channelid) == 0 {
			return true
		}
		time.Sleep(time.Second)
	}
	return false
}

func AwaitMessage(curr *discordgo.Message, channelid string, dg *discordgo.Session) *discordgo.Message {
	var messages []*discordgo.Message
	for {
		messages, _ = dg.ChannelMessages(channelid, 100, "", curr.ID, "")
		if messages == nil {
			time.Sleep(time.Millisecond * 750)
			continue
		}
		if len(messages) < 1 {
			time.Sleep(time.Millisecond * 750)
			continue
		} else {
			break
		}
	}
	println(messages[0].Content) //content
	current := messages[0]

	return current
}

func AwaitEnd(curr *discordgo.Message, channelid string, dg *discordgo.Session) {
	var messages []*discordgo.Message
	for {
		messages, _ = dg.ChannelMessages(channelid, 100, "", curr.ID, "")
		if messages == nil {
			time.Sleep(time.Second)
			continue
		} else if len(messages) < 1 {
			time.Sleep(time.Second)
			continue
		}
		if strings.Compare(messages[0].Content, "-endrelease") == 0 {
			return
		}
		time.Sleep(time.Second)
		continue
	}
	println(messages[0].Content)
}

//ratelimiter? look discord
func GetInfo(channelid string, webhook *discordgo.Webhook, dg *discordgo.Session) {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "NuggieACO",
		},
		Title:       "Thank You For Choosing NuggieACO!",
		Description: "Which Release Would You like?", //can change upon different releases
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Property of Nuggie Labs",
			IconURL: "",
		},
		Color: 14902784,
	}
	var id []*discordgo.MessageEmbed
	id = append(id, embed)
	web := &discordgo.WebhookParams{
		Username: "NuggieACO",
		Embeds:   id,
	}

	temp, _ := dg.WebhookExecute(webhook.ID, webhook.Token, true, web)

	// start to wait for -close
	//go AwaitClose(temp, channelid, dg)
	go func() {
		for {
			if AwaitClose(temp, channelid, dg) {
				_, err := dg.ChannelDelete(channelid)
				if err != nil {
					log.Fatal(err)
				}
				// maybe delete from database
				println("success deleting channel")
				return
			}
			time.Sleep(time.Second)
		}
	}()

	result := AwaitMessage(temp, channelid, dg)

	release := result.Content
	_, err := dg.ChannelEdit(channelid, result.Author.Username+release)
	if err != nil {
		log.Fatal(err)
	}

	var nothing slice
	msg, err := CustomWebhook("Thank You For Choosing NuggieACO!", "What Size Range Would You Like?\nExamples: **random** / **only GS** / **7-13** / **6,7,10,11**", nothing, dg, webhook)
	if err != nil {
		println("close detected?")
		return
	}

	current := msg
	result = AwaitMessage(current, channelid, dg)
	sizes := result.Content

	msg, err = CustomWebhook("Thank You For Choosing NuggieACO!", "What Is the Max Checkouts You want?\nExamples: **No Max** / **10 Pairs**", nothing, dg, webhook)
	if err != nil {
		println("close detected?")
		return
	}

	current = msg
	result = AwaitMessage(current, channelid, dg)
	maxcheckouts := result.Content

	msg, err = CustomWebhook("Thank You For Choosing NuggieACO!", "Please List Any Additional Details\nExample: **AMEX/ENO/PRIVACY Limit**", nothing, dg, webhook)
	if err != nil {
		println("close detected?")
		return
	}

	current = msg
	result = AwaitMessage(current, channelid, dg)
	adddescription := result.Content

	msg, err = CustomWebhook("Thank You For Choosing NuggieACO!", "Please Provide Your Profiles in a **ZIPPED AYCD FILE** With the Password Set to "+"`"+result.Author.Username+"-"+release+"`\n **__Please Name All Your Profiles "+result.Author.Username+"-"+release+"__**", nothing, dg, webhook)
	if err != nil {
		println("close detected?")
		return
	}
	current = msg
	result = AwaitMessage(current, channelid, dg)
	file := result.Attachments

	msg, err = CustomWebhook("Thank You For Choosing NuggieACO!", "Your ACO Request Has Been Recorded, Let's Cook!", nothing, dg, webhook)

	info := map[string]string{
		"Release":         release,
		"Size Range":      sizes,
		"Max Checkouts":   maxcheckouts,
		"Additional Info": adddescription,
	}

	InfoWebhook(result.Author.Username+"-"+release, result.Author.Username+"-"+release, info, file, dg)

	statement, err := Database.Prepare("drop table Orders")
	if err != nil {
		println("error deleting")
		log.Fatal(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}

	statement, err = Database.Prepare("CREATE TABLE IF NOT EXISTS Customers(user_release text, webhookid text, token text)")
	if err != nil {
		println("error creating table")
		log.Fatal(err)
	}
	statement.Exec()

	statement, err = Database.Prepare("CREATE TABLE IF NOT EXISTS Orders(user_release text, profile text, orderinfo text, size text)")
	if err != nil {
		println("error creating orders table")
		log.Fatal(err)
	}
	statement.Exec()

	statement, err = Database.Prepare("INSERT INTO Customers(user_release, webhookid, token) VALUES (?, ?, ?)")
	if err != nil {
		println("error inputting into table")
		log.Fatal(err)
	}
	_, err = statement.Exec(result.Author.Username+"-"+release, webhook.ID, webhook.Token)
	if err != nil {
		println("error insertting into table")
		log.Fatal(err)
	}

	var webid string
	err = Database.QueryRow("select webhookid from Customers where user_release = ?", result.Author.Username+"-"+release).Scan(&webid)
	if err != nil {
		println("here")
		log.Fatal(err)
	}

	println("found in database: " + webid)

	var user_release string
	var webhookid string
	var token string
	rows, err := Database.Query("SELECT user_release, webhookid, token FROM Customers")

	var total int
	println()
	for rows.Next() {
		total++
		rows.Scan(&user_release, &webhookid, &token)
		println(user_release + ": " + webhookid + " / " + token)
	}
	AwaitEnd(msg, channelid, dg)
	println("done with release\nDeleting user from database and retrieving checkouts...")

	// delete webhook + retrieve total checkouts
	statement, err = Database.Prepare("DELETE from Customers where webhookid = " + webhook.ID)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}

	println("successfully webhookid: " + webhook.ID + " deleted from database")
	defer Database.Close()

	total = 0
	var profile string
	var orderinfo string
	var size string
	var list string
	rows, err = Database.Query("SELECT profile, orderinfo, size FROM Orders WHERE user_release = ?", result.Author.Username+"-"+release)
	if err != nil {
		log.Fatal(err)
	} else {
		total = 0
		println()
		for rows.Next() {
			total++
			rows.Scan(&profile, &orderinfo, &size)
			list = list + profile + ": " + orderinfo + " ~ " + size + "\n"
			println(profile + ": " + orderinfo + " ~ " + size)
		}
	}
	dg.ChannelMessageSend(channelid, "You have "+strconv.Itoa(total)+" checkouts!\n"+list)
	println("You have " + strconv.Itoa(total) + " checkouts")
	println("done")

	statement, err = Database.Prepare("DELETE FROM Orders WHERE user_release = ?")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(user_release)
	println("deleted data from database")

}

func HandleReaction(channelid string, msgid string, dg *discordgo.Session) *discordgo.Channel {
	users, err := dg.MessageReactions(channelid, msgid, "🎟️", 2, "", "")
	if err != nil {
		println(err)
		time.Sleep(400 * time.Millisecond)
		return nil
	}
	println(users[0].Username)
	if len(users) != 2 {
		println("NULL")
		time.Sleep(400 * time.Millisecond)
		return nil
	} else {
		err = dg.MessageReactionRemove(channelid, msgid, "🎟️", users[0].ID)
		if err != nil {
			println(err)
			time.Sleep(400 * time.Millisecond)
			return nil
		}

		perms := []*discordgo.PermissionOverwrite{}
		temp := &discordgo.PermissionOverwrite{ //USER CAN SEE CHANNEL
			ID:    users[0].ID,
			Type:  1,
			Allow: 0x0000000400,
		}
		perms = append(perms, temp)
		temp = &discordgo.PermissionOverwrite{ //EVERYONE ELSE CANNOT SEE
			ID:   "439590627429056512",
			Type: 0,
			Deny: 0x0000000400,
		}
		perms = append(perms, temp)

		channel := discordgo.GuildChannelCreateData{
			Name:                 users[0].Username + "aco",
			Type:                 0,
			PermissionOverwrites: perms,
			ParentID:             "880831831648403576",
		}

		result, err := dg.GuildChannelCreateComplex("439590627429056512", channel)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(400 * time.Millisecond)

		return result
	}

	time.Sleep(400 * time.Millisecond)
	return nil

}
