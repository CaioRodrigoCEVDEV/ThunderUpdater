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
│   ├── architecture/    # definição da arquitetura do executável
│   ├── config/          # configurações
│   ├── download/        # download de arquivos
│   ├── logger/          # logger central (slog)
│   ├── odbc/            # conexão ODBC
│   ├── release/         # gerenciamento de releases
│   ├── repository/      # acesso a dados
│   ├── consoleui/       # interface de console
│   ├── ui/              # componentes de interface
│   └── updater/         # processo de atualização
├── assets/
├── build/
├── scripts/
├── README.md
└── go.mod
```

## Executando e compilando o projeto

### Pré-requisitos

- Go 1.24, conforme definido em `go.mod`.
- Windows para executar o atualizador.
- DSN ODBC `Thunder` configurado no Windows.
- Driver ODBC compatível com a arquitetura do executável: 32 bits para x86 ou 64 bits para x64, com o DSN `Thunder` configurado para a mesma arquitetura.
- Para a compilação oficial, um compilador C compatível com CGO e com o alvo Windows da arquitetura compilada, pois os scripts definem `CGO_ENABLED=1`.
- Git Bash no Windows, ou outro terminal Bash compatível. Em WSL/Linux, é necessário também um toolchain C capaz de compilar para Windows x86.

O projeto usa o driver ODBC `github.com/alexbrainman/odbc`. O executável e o driver ODBC devem usar a mesma arquitetura. Os scripts de build instalam `goversioninfo` automaticamente caso a ferramenta não esteja disponível, portanto o primeiro build também exige Go disponível no `PATH`, acesso à rede e permissão para executar `go install`.

O comando `bc` é opcional e serve apenas para exibir o tamanho do executável em megabytes; sem ele, o script exibe o tamanho em bytes. O `stat` é usado pelo script para obter o tamanho do arquivo.

### Executar em desenvolvimento

Na raiz do repositório, o comando de execução direta pelo Go é:

```bash
go run .
```

O comando exige Go 1.24 e o ambiente Windows com o DSN `Thunder` e seu driver ODBC de 32 bits configurados. Como o entry point atual está em `cmd/thunderupdater`, se a raiz retornar `no Go files`, execute o pacote real:

```bash
go run ./cmd/thunderupdater
```

### Gerar o executável x86

Na raiz do projeto, usando Git Bash ou outro terminal compatível com Bash, execute o processo oficial:

```bash
./build/build-x86.sh
```

O script executa estas etapas:

- garante a existência de `dist/`;
- instala `goversioninfo` se necessário;
- gera temporariamente o recurso Windows x86 com `goversioninfo -64=false` e `assets/icon.ico`;
- compila com `CGO_ENABLED=1`, `GOOS=windows`, `GOARCH=386` e `go build -ldflags="-s -w"`;
- remove o arquivo temporário de recursos após a compilação.

O executável gerado é:

```text
dist/ThunderUpdater-x86.exe
```

Esse executável x86 é o build oficial necessário para funcionar com o driver ODBC e o DSN `Thunder` de 32 bits. Em Linux ou WSL, caso o script não tenha permissão de execução, use:

```bash
chmod +x ./build/build-x86.sh
./build/build-x86.sh
```

### Gerar o executável x64

Na raiz do projeto, usando Git Bash ou outro terminal compatível com Bash, execute:

```bash
./build/build-x64.sh
```

Esse script compila com `CGO_ENABLED=1`, `GOOS=windows`, `GOARCH=amd64` e `go build -ldflags="-s -w"`. O executável gerado é:

```text
dist/ThunderUpdater-x64.exe
```

O executável x64 requer o driver ODBC e o DSN `Thunder` de 64 bits. Em Linux ou WSL, conceda permissão de execução se necessário:

```bash
chmod +x ./build/build-x64.sh
./build/build-x64.sh
```

### Gerar os executáveis x86 e x64

Para executar os dois builds em sequência, na raiz do projeto, use:

```bash
./build/build-all.sh
```

O script `build-all.sh` chama `build-x86.sh` e, depois, `build-x64.sh`. Ao concluir, os arquivos estarão em `dist/`:

```text
dist/ThunderUpdater-x86.exe
dist/ThunderUpdater-x64.exe
```

Em Linux ou WSL, conceda permissão ao script se necessário:

```bash
chmod +x ./build/build-all.sh
./build/build-all.sh
```

## Downloads

Os executáveis prontos para Windows podem ser baixados na página de [Releases](../../releases) do repositório.

## Publicar uma nova versão

Crie e envie uma tag de versão para iniciar automaticamente o build e a publicação dos executáveis:

```bash
git tag v1.0.0
git push origin v1.0.0
```

Também é possível executar o workflow manualmente pela aba **Actions**, informando a versão no formato `v1.0.0`.

### Fluxo da atualização

Após iniciar, o atualizador:

- consulta o campo `confrelease` na tabela `conf` pelo ODBC;
- baixa exatamente o arquivo `Thunder_<versão>.zip` correspondente à versão consultada;
- extrai o conteúdo em `C:\Thunder`, criando a pasta se necessário;
- sobrescreve arquivos extraídos que já tenham o mesmo nome;
- preserva os demais arquivos existentes em `C:\Thunder`, pois não limpa a pasta nem remove arquivos que não estejam no ZIP;
- não realiza backup dos arquivos;
- exibe uma barra de progresso para o download e outra para a extração, em etapas separadas;
- encerra automaticamente após uma atualização concluída com sucesso;
- mantém o terminal aberto em caso de erro, aguardando ENTER para que a mensagem possa ser lida.
