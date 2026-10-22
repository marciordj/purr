# purr

purr converte arquivos QuickTime `.mov` em `.mp4` pelo terminal. Ele encapsula o `ffmpeg` com um seletor de arquivos interativo e uma tela de progresso, e recodifica gravações para HEVC usando o encoder em hardware VideoToolbox do macOS.

English: [README.md](README.md)

## Requisitos

- **macOS** para o caminho padrão de compressão (encoder da Apple VideoToolbox, que já vem no `ffmpeg` da Homebrew). O caminho de remux via `--fast` funciona em qualquer sistema com `ffmpeg`/`ffprobe`.
- `ffmpeg` e `ffprobe` no `PATH` — `brew install ffmpeg` no macOS. O purr sai imediatamente se algum dos dois não existir.
- Go 1.26.3 ou mais recente (versão declarada no `go.mod`), apenas para compilar o binário.

## Início rápido

```sh
git clone https://github.com/marciordj/purr
cd purr
make build          # gera o binário ./purr
./purr              # abre o seletor no diretório atual
```

O seletor lista apenas os arquivos `.mov` do diretório atual. Para converter um vídeo específico direto, passe-o como argumento:

```sh
./purr ~/Desktop/recording.mov
```

Se não tiver nenhum `.mov` em mãos, crie uma gravação de tela (QuickTime Player → Arquivo → Nova Gravação de Tela) e passe o caminho do arquivo salvo para `./purr`, ainda no diretório do projeto.

## Instalação opcional

```sh
make install        # compila e copia o binário para ~/.local/bin/purr
make uninstall      # remove o binário
```

Verifique se `~/.local/bin` está no seu `PATH`; se não estiver, adicione ao `~/.zshrc` (ou `~/.bashrc`):

```sh
export PATH="$HOME/.local/bin:$PATH"
```

Depois de alterar o código, rode `make install` novamente para atualizar o binário instalado.

## Uso

```sh
purr                    # seletor de arquivos (só o diretório atual, sem subpastas)
purr video.mov          # converte um arquivo direto, sem seletor
purr clip1.mov clip2.mov  # converte vários arquivos, em ordem
purr --fast video.mov   # remux sem recodificar (-f equivale a --fast)
```

- Os argumentos precisam ser arquivos existentes com extensão `.mov` (maiúsculas/minúsculas tanto faz); qualquer coisa diferente sai com erro.
- O shell expande `purr *.mov` antes da execução, então o comando recebe os arquivos `.mov` que o glob encontra no disco.
- A saída fica junto da entrada, com o mesmo nome base: `video.mov` → `video.mp4`.
- Se o arquivo de saída já existir, o trabalho é pulado (exibido como "output exists, skipped"); o purr nunca sobrescreve.
- Os trabalhos são executados um por vez, na ordem da lista.
- O purr inspeciona cada arquivo com o `ffprobe` antes de converter com o `ffmpeg`.

### Padrão: compressão

Sem a opção `--fast`, os arquivos são recodificados para **HEVC (H.265)** para reduzir o tamanho; a redução depende do material e não é garantida.

- Vídeo: `hevc_videotoolbox` (VideoToolbox, macOS) na qualidade 65, marcado como `hvc1` para ser reproduzido pelo QuickTime Player, Safari e iOS; a codificação por software é permitida quando o hardware não está disponível.
- Áudio: copiado como está quando AAC, recodificado para AAC 192k caso contrário.
- Como ele depende do encoder VideoToolbox, o caminho padrão funciona apenas no macOS; não há encoder de fallback multiplataforma configurado.

### `--fast`: remux

`--fast` (ou `-f`) inspeciona o arquivo e, quando seus codecs são considerados compatíveis com MP4, o copia para um contêiner MP4 sem recodificar. Isso costuma ser mais rápido que a compressão, preserva a qualidade codificada e mantém o tamanho praticamente igual.

- Vídeo compatível: H.264, HEVC, VP9, AV1, MPEG-4.
- Áudio compatível: AAC, MP3, Opus, AC3, E-AC3, FLAC.
- Com codecs incompatíveis (como ProRes), o purr recorre à compressão completa, que exige macOS. A conversão ainda pode falhar se o FFmpeg não conseguir processar a entrada ou colocar seus codecs no contêiner MP4.

## Teclas

**Seletor de arquivos**

| Tecla               | Ação                                           |
|---------------------|------------------------------------------------|
| `↑`/`↓`, `k`/`j`    | Mover cursor                                   |
| `space`             | Marcar/desmarcar                               |
| `enter`             | Iniciar conversão (exige pelo menos um selecionado) |
| `q`, `esc`, `ctrl+c`| Sair                                           |

**Conversão**

| Tecla     | Ação                                                                        |
|-----------|------------------------------------------------------------------------------|
| `esc`    | Prompt "Abort remaining jobs? [y/N]": `y` cancela o trabalho em execução e pula o restante; qualquer outra tecla retoma |
| `ctrl+c` | Interrompe a conversão em execução imediatamente e sai                        |
| `q`      | Sair — disponível apenas quando nenhum trabalho está em execução               |

Cada trabalho mostra ○ na fila / ▸ em execução (com barra de progresso, velocidade e ETA) / ✓ concluído / – pulado / ✗ falhou (saída de erro do ffmpeg abaixo da linha) / ◌ cancelado. Um trabalho cancelado ou interrompido remove o arquivo parcial que estava sendo gerado.

## Limitações conhecidas

- Argumentos diretos e seletor aceitam apenas `.mov`; não há varredura recursiva de diretórios.
- Apenas uma conversão por vez; a fila é sequencial.
- Erros aparecem no stderr do ffmpeg abaixo da linha correspondente na fila.
- Não há arquivo de configuração, escolha de encoder ou substituição do diretório de saída.
