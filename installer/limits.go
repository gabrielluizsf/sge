package installer

// Limits bound what an installer copies, against a payload that would
// fill the disk: Files how many files, Bytes how many bytes in all.
type Limits struct {
	Files int
	Bytes int64
}

// DefaultLimits are taken when a Setup's Limits are zero.
var DefaultLimits = Limits{Files: 200_000, Bytes: 64 << 30}

func (l Limits) orDefault() Limits {
	if l.Files <= 0 {
		l.Files = DefaultLimits.Files
	}
	if l.Bytes <= 0 {
		l.Bytes = DefaultLimits.Bytes
	}
	return l
}
