package models

type LocationInfo struct {
	PlayerColumn int32
	ScoreColumn  int32
}

type ProtocolImage struct {
	Path string
	Info LocationInfo
}

type PlayerStat struct {
	Name    string
	Surname string
	Score   int32
}

var ProtocolTypes = map[string]LocationInfo{
	"fiba":    {1, 21},
	"classic": {1, 2},
}
