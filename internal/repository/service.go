package repository

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

const httpTimeout = 30 * time.Second

type RepositoryService interface {
	LoadCatalog() (*ReleaseCatalog, error)
	GetLatestRelease() (*Release, error)
	FindRelease(version int) (*Release, error)
}

type repositoryService struct {
	baseURL string
	log     *slog.Logger
	client  *http.Client
	catalog *ReleaseCatalog
	mu      sync.Mutex
}

func NewRepositoryService(baseURL string, log *slog.Logger) RepositoryService {
	if log == nil {
		log = slog.Default()
	}
	return &repositoryService{
		baseURL: baseURL,
		log:     log,
		client: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

func (s *repositoryService) LoadCatalog() (*ReleaseCatalog, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.catalog != nil {
		return s.catalog, nil
	}

	s.log.Info("Consultando repositório", "url", s.baseURL)

	req, err := http.NewRequest(http.MethodGet, s.baseURL, nil)
	if err != nil {
		s.log.Error("Erro ao criar requisição HTTP", "error", err)
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	startReq := time.Now()
	s.log.Info("Iniciando requisição HTTP", "url", s.baseURL)

	resp, err := s.client.Do(req)
	elapsedReq := time.Since(startReq)
	if err != nil {
		s.log.Error("Erro na requisição HTTP",
			"error", err,
			"elapsed", elapsedReq.String(),
		)
		return nil, fmt.Errorf("erro ao acessar repositório: %w", err)
	}
	defer resp.Body.Close()

	s.log.Info("Resposta recebida",
		"status", resp.StatusCode,
		"elapsed", elapsedReq.String(),
	)

	startBody := time.Now()
	s.log.Info("Iniciando leitura do body")

	body, err := io.ReadAll(resp.Body)
	elapsedBody := time.Since(startBody)
	if err != nil {
		s.log.Error("Erro ao ler body da resposta",
			"error", err,
			"elapsed", elapsedBody.String(),
		)
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	s.log.Info("Fim da leitura do body",
		"tamanho", len(body),
		"elapsed", elapsedBody.String(),
	)

	s.log.Info("Iniciando parser das releases")
	releases, err := parseReleases(string(body), s.baseURL)
	if err != nil {
		s.log.Error("Erro ao fazer parser das releases",
			"error", err,
			"bodySize", len(body),
		)
		return nil, fmt.Errorf("erro ao processar releases: %w", err)
	}

	s.log.Info("Parser concluído",
		"quantidade", len(releases),
	)

	latest := releases[len(releases)-1].Version
	s.log.Info("Releases encontradas", "quantidade", len(releases), "ultima", latest)

	s.catalog = &ReleaseCatalog{
		Releases:      releases,
		LatestVersion: latest,
	}

	return s.catalog, nil
}

func (s *repositoryService) GetLatestRelease() (*Release, error) {
	catalog, err := s.LoadCatalog()
	if err != nil {
		return nil, err
	}

	release := catalog.LatestRelease()
	if release == nil {
		return nil, fmt.Errorf("nenhuma release disponível")
	}

	return release, nil
}

func (s *repositoryService) FindRelease(version int) (*Release, error) {
	catalog, err := s.LoadCatalog()
	if err != nil {
		return nil, err
	}

	release := catalog.FindByVersion(version)
	if release == nil {
		return nil, fmt.Errorf("release %d não encontrada", version)
	}

	return release, nil
}
