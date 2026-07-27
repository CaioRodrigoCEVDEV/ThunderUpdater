package ui

// UI gerencia a interface gráfica (Fyne)
type UI struct{}

func New() *UI {
	return &UI{}
}

func (u *UI) Run() error {
	return nil
}
