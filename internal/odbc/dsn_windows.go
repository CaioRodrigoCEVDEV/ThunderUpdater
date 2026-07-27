//go:build windows

package odbc

import (
	"log/slog"

	"golang.org/x/sys/windows/registry"
)

const dsnKeyPath = `SOFTWARE\ODBC\ODBC.INI\ODBC Data Sources`

func ListDSNs() []string {
	dsns := listFromRoot(registry.LOCAL_MACHINE, dsnKeyPath)
	dsns = append(dsns, listFromRoot(registry.CURRENT_USER, dsnKeyPath)...)
	return dsns
}

func listFromRoot(root registry.Key, keyPath string) []string {
	k, err := registry.OpenKey(root, keyPath, registry.READ)
	if err != nil {
		slog.Debug("dsn key not found", "root", root, "path", keyPath)
		return nil
	}
	defer k.Close()

	names, err := k.ReadValueNames(0)
	if err != nil {
		slog.Debug("failed to read dsn value names", "error", err)
		return nil
	}

	return names
}
