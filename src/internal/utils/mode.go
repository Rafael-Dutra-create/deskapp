package utils

type MODE uint8

func (m MODE) String() string {
	switch m {
	case DEBUG: return "DEBUG"
	case RELEASE: return "RELEASE"
	default: return "UNKOWN"
	}
}



const (
	DEBUG MODE = iota
	RELEASE
)