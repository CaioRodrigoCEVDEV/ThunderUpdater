package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/caio/ThunderUpdaterGO/internal/architecture"
	"github.com/caio/ThunderUpdaterGO/internal/config"
	"github.com/caio/ThunderUpdaterGO/internal/consoleui"
	"github.com/caio/ThunderUpdaterGO/internal/download"
	"github.com/caio/ThunderUpdaterGO/internal/odbc"
	"github.com/caio/ThunderUpdaterGO/internal/repository"
	"github.com/caio/ThunderUpdaterGO/internal/updater"
)

func Run() error {
	log := slog.Default()
	ui := consoleui.New()

	arch := executableArchitecture()
	_ = arch

	dsns := odbc.ListDSNs()
	if len(dsns) == 0 {
		return fmt.Errorf("Nenhum DSN disponível")
	}

	ui.PrintConnecting()
	db := odbc.NewDatabaseService(log)
	if err := db.Connect(dsns[0]); err != nil {
		return fmt.Errorf("Falha na conexão: %v", err)
	}

	version, err := db.GetInstalledVersion()
	if err != nil {
		db.Disconnect()
		return fmt.Errorf("Falha ao obter versão: %v", err)
	}
	db.Disconnect()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("Erro ao carregar configuração: %v", err)
	}

	repo := repository.NewRepositoryService(cfg.RepositoryURL(), log)
	release, err := repo.FindRelease(version)
	if err != nil {
		return fmt.Errorf("Release não encontrada: %v", err)
	}

	ui.PrintHeader(version, release.FileName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	dl := download.NewDownloadManager(log)
	dl.SetCallback(ui.DownloadCallback())
	zipPath, err := dl.Download(ctx, release.FileName, release.DownloadURL)
	if err != nil {
		return fmt.Errorf("Falha no download: %v", err)
	}

	if err := os.MkdirAll(config.ThunderPath, 0755); err != nil {
		return fmt.Errorf("Falha ao criar diretório Thunder: %v", err)
	}

	updateMgr := updater.NewUpdateManager(log, version)
	if _, err := updateMgr.Update(zipPath, ui.ExtractCallback()); err != nil {
		return fmt.Errorf("Falha na atualização: %v", err)
	}

	ui.PrintUpdateComplete()

	return nil
}

func executableArchitecture() architecture.Architecture {
	switch runtime.GOARCH {
	case "386":
		return architecture.X86
	default:
		return architecture.X64
	}
}
