package window

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Kazemir/discord-soundboard/feature/sound"
	"github.com/Kazemir/discord-soundboard/util/files"
	"github.com/bwmarrin/discordgo"

	log "github.com/sirupsen/logrus"
)

func CreateWindow(discord *discordgo.Session) fyne.Window {
	var voiceConnection *discordgo.VoiceConnection

	a := app.New()
	w := a.NewWindow("Саундбар профессионала")

	var joinButton *widget.Button
	var leaveButton *widget.Button

	joinButton = widget.NewButton("Войти", func() {
		guilds := make(map[string]*discordgo.Guild)
		var channels map[string]*discordgo.Channel
		for _, g := range discord.State.Guilds {
			guilds[g.Name] = g
		}
		if len(guilds) == 0 {
			dialog.ShowInformation("Ошибка", "Бот не добавлен ни на один сервер", w)
			return
		}

		var selectedGuild *discordgo.Guild
		var selectedChannel *discordgo.Channel

		guildNames := make([]string, 0, len(guilds))
		for k := range guilds {
			guildNames = append(guildNames, k)
		}

		channelIDCombo := widget.NewSelect([]string{}, func(value string) {
			channel, isFound := channels[value]
			if !isFound {
				return
			}
			selectedChannel = channel
		})

		guildIDCombo := widget.NewSelect(guildNames, func(value string) {
			guild, isFound := guilds[value]
			if !isFound {
				return
			}
			selectedGuild = guild
			channels = make(map[string]*discordgo.Channel)
			channelNames := make([]string, 0)
			for _, c := range selectedGuild.Channels {
				if c.Type == discordgo.ChannelTypeGuildVoice {
					channelNames = append(channelNames, c.Name)
					channels[c.Name] = c
				}
			}
			channelIDCombo.SetOptions(channelNames)
			channelIDCombo.SetSelectedIndex(0)
		})
		guildIDCombo.SetSelectedIndex(0)

		dialog.ShowForm("Вход на канал", "Войти", "Отмена", []*widget.FormItem{
			widget.NewFormItem("Гильдия", guildIDCombo),
			widget.NewFormItem("Канал", channelIDCombo),
		}, func(response bool) {
			if response {
				var err error
				voiceConnection, err = discord.ChannelVoiceJoin(selectedGuild.ID, selectedChannel.ID, false, true)
				if err != nil {
					dialog.ShowInformation("Ошибка", "Произошла ошибка при подключении к каналу... "+err.Error(), w)
					log.Error(err)
					return
				}
				joinButton.Disable()
				leaveButton.Enable()
			}
		}, w)
	})

	leaveButton = widget.NewButton("Выйти", func() {
		if voiceConnection != nil {
			voiceConnection.Disconnect()
			voiceConnection = nil
			joinButton.Enable()
			leaveButton.Disable()
		}
	})
	leaveButton.Disable()

	var buttons = []*widget.Button{joinButton, leaveButton}

	for _, s := range files.FindSounds(".ogg") {
		soundName := s
		soundPath := "../sounds/" + soundName
		button := widget.NewButton(soundName, func() {
			if voiceConnection != nil {
				println("play " + soundName)
				err := sound.PlaySound(voiceConnection, soundPath)
				if err != nil {
					dialog.ShowInformation("Ошибка", "Произошла ошибка при отправке звука... "+err.Error(), w)
					log.Error(err)
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
