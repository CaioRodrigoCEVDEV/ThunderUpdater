package architecture

// Architecture representa a arquitetura de build (x86 ou x64)
type Architecture string

const (
	X86 Architecture = "x86"
	X64 Architecture = "x64"
)

func (a Architecture) String() string {
	return string(a)
}
