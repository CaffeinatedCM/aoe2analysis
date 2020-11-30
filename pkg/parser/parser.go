package parser

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"math"
	"strings"
	"syscall"
)

func readChunk(reader io.Reader, chunkLen int) []byte {
	buf := make([]byte, chunkLen)
	read, err := reader.Read(buf)
	if err != nil || read != chunkLen {
		log.Fatalf("could not read chunk: %Ev", err)
	}

	return buf
}

func readCString(reader io.Reader) string {
	stringBuilder := strings.Builder{}
	for {
		buf := make([]byte, 1)
		_, err := reader.Read(buf)
		if err != nil {
			log.Fatalf("Failed to read CString: %Ev", err)
		}
		if buf[0] == 0x00 {
			break
		}
		stringBuilder.WriteByte(buf[0])
	}
	return stringBuilder.String()
}

func readFloat32(reader io.Reader) float32 {
	return math.Float32frombits(readUInt32(reader))
}

func readInt32(reader io.Reader) int32 {
	return (int32)(readUInt32(reader))
}

func readUInt64(reader io.Reader) uint64 {
	intBytes := readChunk(reader, 8)
	return binary.LittleEndian.Uint64(intBytes)
}

func readUInt32(reader io.Reader) uint32 {
	intBytes := readChunk(reader, 4)
	return binary.LittleEndian.Uint32(intBytes)
}

func readUInt16(reader io.Reader) uint16 {
	intBytes := readChunk(reader, 2)
	return binary.LittleEndian.Uint16(intBytes)
}

func readUInt8(reader io.Reader) uint8 {
	intBytes := readChunk(reader, 1)
	return intBytes[0]
}

func readBool(reader io.Reader) bool {
	boolByte := readChunk(reader, 1)
	if boolByte[0] == 0x01 {
		return true
	}
	return false
}

func readDEString(reader io.Reader) DEString {
	deStringStart := readChunk(reader, 2)
	if !bytes.Equal(deStringStart, []byte{0x60, 0x0A}) {
		log.Fatalf("Unexpected start of DE String: %+#v", deStringStart)
	}
	length := readUInt16(reader)
	return DEString{
		Length: length,
		Value: readChunk(reader, (int)(length)),
	}
}

func skip(reader io.Reader, bytesToSkip int) {
	readChunk(reader, bytesToSkip)
}

func skipSeperator(reader io.Reader) {
	seperator := readChunk(reader, 4)
	if !bytes.Equal(seperator, []byte{0xa3, 0x5F, 0x02, 0x00}) {
		log.Fatalf("Unexpected seperator bytes: %+#v", seperator)
	}
}

func contains(arr []uint32, val uint32) bool {
	for _, b := range arr {
		if b == val {
			return true
		}
	}
	return false
}

func readHeader(reader io.Reader) (*Header, error) {
	header := Header{}
	header.length = readInt32(reader)
	skip(reader, 4)

	compressedHeaderBytes := readChunk(reader, (int)(header.length)-8)

	headerReader := flate.NewReader(bytes.NewReader(compressedHeaderBytes))
	defer func() {
		err := headerReader.Close()
		if err != nil {
			log.Fatalf("could not close header reader")
		}
	}()

	header.Version = readCString(headerReader)
	header.SaveVersion = readFloat32(headerReader)
	header.DE.Version = readFloat32(headerReader)
	header.DE.IntervalVersion = readUInt32(headerReader)
	header.DE.GameOptionsVersion = readUInt32(headerReader)
	header.DE.DLCCount = readUInt32(headerReader)
	header.DE.DLCIds = make([]uint32, header.DE.DLCCount)
	for i := range header.DE.DLCIds {
		header.DE.DLCIds[i] = readUInt32(headerReader)
	}
	header.DE.DatasetRef = readUInt32(headerReader)
	difficulty := readUInt32(headerReader)
	header.DE.Difficulty = Difficulty(difficulty)
	header.DE.SelectedMapId = readUInt32(headerReader)
	header.DE.ResolvedMapId = readUInt32(headerReader)
	header.DE.RevealMap = readUInt32(headerReader)
	header.DE.VictoryType = VictoryType(readUInt32(headerReader))
	header.DE.StartingResources = ResourceLevel(readUInt32(headerReader))
	header.DE.StartingAge = Age(readInt32(headerReader))
	header.DE.EndingAge = Age(readInt32(headerReader))
	header.DE.GameType = readUInt32(headerReader)

	skipSeperator(headerReader)
	skipSeperator(headerReader)

	header.DE.Speed = readFloat32(headerReader)
	header.DE.TreatyLength = readUInt32(headerReader)
	header.DE.PopulationLimit = readUInt32(headerReader)
	header.DE.NumPlayers = readUInt32(headerReader)
	header.DE.UnusedPlayerColor = readUInt32(headerReader)
	header.DE.VictoryAmount = readUInt32(headerReader)

	skipSeperator(headerReader)

	header.DE.TradeEnabled = readBool(headerReader)
	header.DE.TeamBonusDisabled = readBool(headerReader)
	header.DE.RandomPositions = readBool(headerReader)
	header.DE.AllTechs = readBool(headerReader)
	header.DE.NumStartingUnits = readChunk(headerReader, 1)[0]
	header.DE.LockTeams = readBool(headerReader)
	header.DE.LockSpeed = readBool(headerReader)
	header.DE.Multiplayer = readBool(headerReader)
	header.DE.Cheats = readBool(headerReader)
	header.DE.RecordGame = readBool(headerReader)
	header.DE.AnimalsEnabled = readBool(headerReader)
	header.DE.PredatorsEnabled = readBool(headerReader)
	header.DE.TurboEnabled = readBool(headerReader)
	header.DE.SharedExploration = readBool(headerReader)
	header.DE.TeamPositions = readBool(headerReader)

	if header.SaveVersion >= 13.34 {
		skip(headerReader, 8) // TODO: Wonder if these mean anything?
	}

	skipSeperator(headerReader)

	for i := range header.DE.Players {
		header.DE.Players[i] = Player{
			DLCId: readUInt32(headerReader),
			ColorId: readUInt32(headerReader),
			SelectedColor: readChunk(headerReader, 1)[0],
			SelectedTeamId: readChunk(headerReader, 1)[0],
			ResolvedTeamId: readChunk(headerReader, 1)[0],
			DatCrc: readChunk(headerReader, 8),
			MpGameVersion: readChunk(headerReader,1)[0],
			CivId: readChunk(headerReader, 1)[0],
			blank: readChunk(headerReader, 3),
			AIType: readDEString(headerReader),
			AiCivNameIndex: readChunk(headerReader, 1)[0],
			AiName: readDEString(headerReader),
			Name: readDEString(headerReader),
			Type: PlayerType(readUInt32(headerReader)),
			ProfileId: readUInt32(headerReader),
			blank2: readChunk(headerReader,4),
			PlayerNumber: readInt32(headerReader),
			HdRmElo: readUInt32(headerReader),
			HdDmElo: readUInt32(headerReader),
			AnimatedDestructionEnabled: readBool(headerReader),
			CustomAI: readBool(headerReader),
		}
	}

	header.DE.FogOfWar = readBool(headerReader)
	header.DE.CheatNotifications = readBool(headerReader)
	header.DE.ColoredChat = readBool(headerReader)

	skip(headerReader, 9)
	skipSeperator(headerReader)
	skip(headerReader, 12)

	if header.SaveVersion >= 13.13 {
		skip(headerReader, 5)
	}

	for i := range header.DE.Strings {
		header.DE.Strings[i] = StringsThing{
			String: readDEString(headerReader),
		}
		lst := readUInt32(headerReader)
		for {
			ignores := []uint32{3, 21, 23, 42, 44, 45}
			if !contains(ignores, lst) {
				break
			}
			lst = readUInt32(headerReader)
		}
	}

	for i := range header.DE.StrategicNumbers {
		header.DE.StrategicNumbers[i] = readInt32(headerReader)
	}

	header.DE.NumAiFiles = readUInt64(headerReader)

	aiFiles := make([]AIFile, header.DE.NumAiFiles)
	for i := range header.DE.AiFiles {
		header.DE.AiFiles[i] = AIFile{
			unknown: readChunk(headerReader, 4),
			name: readDEString(headerReader),
			unknown2: readChunk(headerReader, 4),
		}
	}

	header.DE.AiFiles = aiFiles

	header.DE.GUID = syscall.GUID{
		Data1: readUInt32(headerReader),
		Data2: readUInt16(headerReader),
		Data3: readUInt16(headerReader),
	}

	copy(header.DE.GUID.Data4[:], readChunk(headerReader, 8))

	header.DE.LobbyName = readDEString(headerReader)
	header.DE.ModdedDataSet = readDEString(headerReader)

	skip(headerReader, 19)

	if header.SaveVersion >= 13.13 {
		skip(headerReader, 5)
	}

	if header.SaveVersion >= 13.17 {
		skip(headerReader, 9)
	}

	header.DE.RandomString = readDEString(headerReader)

	skip(headerReader, 5)

	if header.SaveVersion >= 13.13 {
		skip(headerReader, 1)
	}

	if header.SaveVersion < 13.17 {
		readDEString(headerReader)
		readUInt32(headerReader)
		skip(headerReader, 4)
	}

	if header.SaveVersion >= 13.17 {
		skip(headerReader, 2)
	}

	hasAiInt := readUInt32(headerReader)
	if hasAiInt == 1 {
		header.AI.HasAi = true
		skip(headerReader, 4096) // skip the actual ai structure for now
	} else {
		header.AI.HasAi = false
	}

	// TODO: This isn't ending up right... we're probably off by some bytes
	header.Replay = ReplayHeader{
		OldTime: readUInt32(headerReader),
		WorldTime: readUInt32(headerReader),
		OldWorldTime: readUInt32(headerReader),
		GameSpeedId: readUInt32(headerReader),
		WorldTimeDeltaSeconds: readUInt32(headerReader),
		Timer: readFloat32(headerReader),
		GameSpeedFloat: readFloat32(headerReader),
		TempPause: readUInt8(headerReader),
		NextObjectId: readUInt32(headerReader),
		NextReusableObjectId: readInt32(headerReader),
		RandomSeed: readUInt32(headerReader),
		RecPlayer: readUInt16(headerReader),
		NumPlayers: readUInt8(headerReader),
	}

	return &header, nil
}

func Parse(reader io.Reader) (*RecordedGame, error) {
	game := RecordedGame{}

	header, err := readHeader(reader)
	if err != nil {
		return nil, errors.New("could not read header: " + err.Error())
	}
	game.Header = *header

	return &game, nil
}
