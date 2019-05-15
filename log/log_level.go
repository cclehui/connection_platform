package log

// 日志层级类型
type Level int

func (ll Level) String() string {
	if ll < LEVEL_DEBUG || ll > LEVEL_FATAL {
		return LEVEL_NAME_UNKNOWN
	}

	return levelName[ll]
}

// 日志级别，重要性从低到高一次递增，LEVEL_UNKNOWN不是一种可用类型
const (
	LEVEL_UNKNOWN = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

var LEVEL_NAME_UNKNOWN = "UNKNOWN"

var levelName = [...]string{
	LEVEL_NAME_UNKNOWN,
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}
