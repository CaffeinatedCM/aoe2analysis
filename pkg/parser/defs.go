package parser

import "syscall"

type RecordedGame struct {
	Header Header
}

/*
Header information found in aoe2record.

Note: Specifically uses int32 for integers to be consistent with the file format since int appears to be more flexible
in go with it's definition of "at least 32 bits"
*/
type Header struct {
	length      int32
	Version     string
	SaveVersion float32
	DE          DEHeader
	AI          AIHeader
	Replay		ReplayHeader
}

type DEHeader struct {
	Version            float32
	IntervalVersion    uint32
	GameOptionsVersion uint32
	DLCCount           uint32
	DLCIds             []uint32
	DatasetRef         uint32
	Difficulty         Difficulty
	SelectedMapId      uint32
	ResolvedMapId      uint32
	RevealMap          uint32
	VictoryType        VictoryType
	StartingResources  ResourceLevel
	StartingAge        Age
	EndingAge          Age
	GameType           uint32
	Speed              float32
	TreatyLength       uint32
	PopulationLimit    uint32
	NumPlayers         uint32
	UnusedPlayerColor  uint32
	VictoryAmount      uint32
	TradeEnabled       bool
	TeamBonusDisabled  bool
	RandomPositions    bool
	AllTechs           bool
	NumStartingUnits   byte
	LockTeams          bool
	LockSpeed          bool
	Multiplayer        bool
	Cheats             bool
	RecordGame         bool
	AnimalsEnabled     bool
	PredatorsEnabled   bool
	TurboEnabled       bool
	SharedExploration  bool
	TeamPositions      bool
	Players            [8]Player
	FogOfWar           bool
	CheatNotifications bool
	ColoredChat        bool
	Strings            [23]StringsThing
	StrategicNumbers   [59]int32
	NumAiFiles uint64
	AiFiles []AIFile
	GUID syscall.GUID
	LobbyName DEString
	ModdedDataSet DEString
	RandomString DEString
}

type StringsThing struct {
	String DEString
}

type AIHeader struct {
	HasAi bool
}

type ReplayHeader struct {
	OldTime uint32
	WorldTime uint32
	OldWorldTime uint32
	GameSpeedId uint32 // world_time_delta
	WorldTimeDeltaSeconds uint32
	Timer float32
	GameSpeedFloat float32
	TempPause byte
	NextObjectId uint32
	NextReusableObjectId int32
	RandomSeed uint32
	RecPlayer uint16
	NumPlayers uint8 // includes gaia

}


type AIFile struct {
	unknown []byte
	name DEString
	unknown2 []byte
}

type Player struct {
	DLCId	uint32
	ColorId uint32
	SelectedColor byte
	SelectedTeamId byte
	ResolvedTeamId byte
	DatCrc []byte
	MpGameVersion byte
	CivId byte
	blank []byte
	AIType DEString
	AiCivNameIndex byte
	AiName DEString
	Name DEString
	Type PlayerType
	ProfileId uint32
	blank2 []byte
	PlayerNumber int32
	HdRmElo uint32
	HdDmElo uint32
	AnimatedDestructionEnabled bool
	CustomAI bool
}

type DEString struct {
	Length uint16
	Value []byte
}

func (deString DEString) String() string {
	return string(deString.Value[:])
}

type PlayerType uint32

type Age int32

const (
	AgeUnknown        = -2
	AgeUnset          = -1
	AgeStandard       = 0
	AgeFeudal         = 1
	AgeCastle         = 2
	AgeImperial       = 3
	AgePostImperial   = 4
	AgeDMPostImperial = 6
)

type ResourceLevel int32

const (
	ResourceLevelNone     ResourceLevel = -1
	ResourceLevelStandard               = 0
	ResourceLevelLow                    = 1
	ResourceLevelMedium                 = 2
	ResourceLevelHigh                   = 3
	ResourceLevelUnknown1               = 4
	ResourceLevelUnknown2               = 5
)

type VictoryType uint32

const (
	VictoryStandard VictoryType = iota
	VictoryConquest
	VictoryExploration
	VictoryRuins
	VictoryArtifacts
	VictoryDiscoveries
	VictoryGold
	VictoryTimeLimit
	VictoryScore
	VictoryStandard2
	VictoryRegicide
	VictoryLastMan
)

func (e VictoryType) String() string {
	return [...]string{"Standard", "Conquest", "Exploration", "Ruins", "Artifacts", "Discoveries", "Gold", "TimeLimit", "Score", "Standard2", "Regicide", "LastMan"}[e]
}

type Difficulty uint32

const (
	DifficultyHardest Difficulty = iota
	DifficultyHard
	DifficultyMedium
	DifficultyStandard
	DifficultyEasiest
	DifficultyExtreme
	DifficultyUnknown
)

func (e Difficulty) String() string {
	switch e {
	case DifficultyHardest:
		return "Hardest"
	case DifficultyHard:
		return "Hard"
	case DifficultyMedium:
		return "Moderate"
	case DifficultyStandard:
		return "Standard"
	case DifficultyEasiest:
		return "Easiest"
	case DifficultyExtreme:
		return "Extreme"
	default:
		return "Unknown"
	}
}
