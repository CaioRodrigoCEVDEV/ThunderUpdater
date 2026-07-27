//go:build !windows

package odbc

func ListDSNs() []string {
	return nil
}
