# ThunderUpdaterGO

Atualizador da aplicação Thunder, escrito em Go.

## Objetivo

Substituir o atualizador legado por uma versão moderna em Go, mantendo as mesmas regras de negócio porém com código mais simples, fácil de manter e com dependências mínimas.

## Estrutura

```
ThunderUpdaterGO/
├── cmd/thunderupdater/   # entry point
├── internal/             # todo o código da aplicação
│   ├── app/             # inicialização e orquestração
│   ├── architecture/    # definição de arquitetura (x86/x64)
│   ├── backup/          # backup antes de atualizações
│   ├── config/          # configurações
│   ├── download/        # download de arquivos
│   ├── logger/          # logger central (slog)
│   ├── odbc/            # conexão ODBC
│   ├── release/         # gerenciamento de releases
│   ├── repository/      # acesso a dados
│   ├── thunder/         # lógica da aplicação Thunder
│   ├── ui/              # interface gráfica (Fyne)
│   └── updater/         # processo de atualização
├── assets/
├── build/
├── docs/
├── scripts/
├── README.md
└── go.mod
