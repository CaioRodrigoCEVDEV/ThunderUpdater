package consoleui

import (
	"fmt"
	"math"

	"github.com/caio/ThunderUpdaterGO/internal/download"
	"github.com/caio/ThunderUpdaterGO/internal/ui"
	"github.com/caio/ThunderUpdaterGO/internal/updater"
)

type ConsoleUI struct{}

func New() *ConsoleUI {
	return &ConsoleUI{}
}

func (c *ConsoleUI) PrintConnecting() {
	fmt.Println("==============================================")
	fmt.Println("          Thunder Updater")
	fmt.Println("==============================================")
	fmt.Println()
	fmt.Print("  Conectando ao banco...")
}

func (c *ConsoleUI) PrintHeader(version int, releaseName string) {
	fmt.Print("\033[2K\r")
	fmt.Println("==============================================")
	fmt.Println("          Thunder Updater")
	fmt.Println("==============================================")
	fmt.Println()
	fmt.Printf("  Versão instalada: %d\n", version)
	fmt.Printf("  Release: %s\n", releaseName)
	fmt.Println()
}

func (c *ConsoleUI) DownloadCallback() download.ProgressCallback {
	bar := ui.NewProgress(4)

	return func(info download.ProgressInfo) {
		if info.Percent == 0 && info.BytesTotal > 0 {
			bar.Start("Baixando atualização...")
			return
		}

		downloadedMB := float64(info.BytesDownloaded) / (1024 * 1024)
		totalMB := float64(info.BytesTotal) / (1024 * 1024)

		var speedLine string
		if info.SpeedMBps > 0 {
			speedLine = fmt.Sprintf("Velocidade: %.1f MB/s", info.SpeedMBps)
		} else {
			speedLine = "Velocidade: -- MB/s"
		}

		var etaLine string
		if info.ETA > 0 {
			etaSecs := int(math.Round(info.ETA.Seconds()))
			if etaSecs >= 60 {
				etaLine = fmt.Sprintf("Tempo restante: %d m %d s", etaSecs/60, etaSecs%60)
			} else {
				etaLine = fmt.Sprintf("Tempo restante: %d s", etaSecs)
			}
		} else {
			etaLine = "Tempo restante: --"
		}

		if info.Percent >= 100 {
			bar.Done("\033[32m\u2714 Download conclu\u00eddo.\033[0m")
		} else {
			bar.Update(info.Percent,
				fmt.Sprintf("%.0f MB / %.0f MB", downloadedMB, totalMB),
				speedLine,
				etaLine,
			)
		}
	}
}

func (c *ConsoleUI) ExtractCallback() func(info updater.ExtractProgressInfo) {
	bar := ui.NewProgress(4)

	return func(info updater.ExtractProgressInfo) {
		if info.Percent == 0 && info.BytesTotal > 0 {
			bar.Start("Extraindo arquivos")
			return
		}

		if info.Percent >= 100 {
			bar.Done("\033[32m\u2714 Extra\u00e7\u00e3o conclu\u00edda.\033[0m")
		} else {
			bar.Update(info.Percent,
				fmt.Sprintf("%d / %d arquivos", info.FilesExtracted, info.FilesTotal),
				"Arquivo atual:",
				info.CurrentFile,
			)
		}
	}
}

func (c *ConsoleUI) PrintUpdateComplete() {
	fmt.Println("  \033[32m\u2714 Atualiza\u00e7\u00e3o conclu\u00edda com sucesso!\033[0m")
	fmt.Println()
}

func (c *ConsoleUI) PrintFatalError(msg string) {
	fmt.Println()
	fmt.Println("=========================================")
	fmt.Println("                ERRO")
	fmt.Println()
	fmt.Printf("  %s\n", msg)
	fmt.Println()
	fmt.Println("=========================================")
	fmt.Print("  Pressione ENTER para sair...")
	fmt.Println()
	fmt.Println("=========================================")
}
