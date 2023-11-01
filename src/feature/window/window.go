package window

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Kazemir/discord-soundboard/feature/sound"
	"github.com/Kazemir/discord-soundboard/util/files"
	"github.com/bwmarrin/discordgo"

	log "github.com/sirupsen/logrus"
)

func CreateWindow(discord *discordgo.Session, guildID, channelID string) fyne.Window {
	var voiceConnection *discordgo.VoiceConnection
	var err error

	a := app.New()
	w := a.NewWindow("Саундбар профессионала")

	var buttons = []*widget.Button{
		widget.NewButton("Войти", func() {
			voiceConnection, err = discord.ChannelVoiceJoin(guildID, channelID, false, true)
			if err != nil {
				log.Error(err)
				return
			}
		}),
		widget.NewButton("Выйти", func() {
			if voiceConnection != nil {
				voiceConnection.Disconnect()
			}
		}),
	}

	for _, s := range files.FindSounds(".ogg") {
		soundName := s
		soundPath := "../sounds/" + soundName
		button := widget.NewButton(soundName, func() {
			if voiceConnection != nil {
				println("play " + soundName)
				err2 := sound.PlaySound(voiceConnection, soundPath)
				if err2 != nil {
					log.Error(err2)
				}
			}
		})
		buttons = append(buttons, button)
	}

	content := container.NewGridWithColumns(3)
	scrollContainer := container.NewVScroll(content)
	for _, b := range buttons {
		content.Add(b)
	}

	scrollContainer.SetMinSize(fyne.NewSize(300, 800))
	w.SetContent(scrollContainer)

	return w
}
