package conf

var (
	// log
	LogLevel string
	LogPath  string
	LogFlag  int

	// console
	ConsolePort   int
	ConsolePrompt string
	ProfilePath   string

	// cluster
	ServerName        string
	ListenAddr        string
	ConnAddrs         map[string]string
	PendingWriteNum   int
	HeartBeatInterval int

	MaxMsgLen    uint32
	LenMsgLen    int
	LittleEndian bool
)
