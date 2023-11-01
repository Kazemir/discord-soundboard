package sound

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca/v2"
	"github.com/jonas747/ogg"
	log "github.com/sirupsen/logrus"
)

func PlaySound(vc *discordgo.VoiceConnection, file string) (err error) {
	encodeSession, err := dca.EncodeFile(file, dca.StdEncodeOptions)
	if err != nil {
		return err
	}
	defer encodeSession.Cleanup()

	err = vc.Speaking(true)
	if err != nil {
		return err
	}

	for _, buff := range makeBuffer(file) {
		vc.OpusSend <- buff
	}

	err = vc.Speaking(false)
	return err
}

func makeBuffer(path string) (output [][]byte) {
	reader, err := os.Open(path)
	if err != nil {
		log.Error("Error opening file:", err)
		return
	}

	// Setup our ogg and packet decoders
	oggdecoder := ogg.NewDecoder(reader)
	packetdecoder := ogg.NewPacketDecoder(oggdecoder)

	// Run through the packet decoder appending the bytes to our output [][]byte
	for {
		packet, _, err := packetdecoder.Decode()
		if err != nil {
			log.Error(err.Error())
			return output
		}
		output = append(output, packet)
	}
}
