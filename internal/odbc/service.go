package odbc

import (
	"log/slog"
	"strconv"
)

type DatabaseService interface {
	Connect(dsn string) error
	Disconnect()
	IsConnected() bool
	GetInstalledVersion() (int, error)
	ExecuteScalar(query string) (string, error)
}

type databaseService struct {
	conn *odbcConnection
	log  *slog.Logger
}

func NewDatabaseService(log *slog.Logger) DatabaseService {
	if log == nil {
		log = slog.Default()
	}
	return &databaseService{
		conn: newConnection(log),
		log:  log,
	}
}

func (s *databaseService) Connect(dsn string) error {
	s.log.Info("Conectando...", "dsn", dsn)

	if err := s.conn.open(dsn); err != nil {
		return err
	}

	s.log.Info("Conectado.")
	return nil
}

func (s *databaseService) Disconnect() {
	s.conn.close()
}

func (s *databaseService) IsConnected() bool {
	return s.conn.isConnected()
}

func (s *databaseService) ExecuteScalar(query string) (string, error) {
	s.log.Debug("Executando consulta...", "query", query)
	return s.conn.queryRowScalar(query)
}

func (s *databaseService) GetInstalledVersion() (int, error) {
	s.log.Info("Executando consulta...")

	val, err := s.ExecuteScalar("SELECT confrelease FROM conf LIMIT 1")
	if err != nil {
		return 0, err
	}

	version, err := strconv.Atoi(val)
	if err != nil {
		return 0, ErrInvalidQuery
	}

	s.log.Info("Versão encontrada", "version", version)
	return version, nil
}
