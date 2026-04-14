# terminal-interface

An SSH-accessible TUI portfolio. Connect with:

```
ssh tui.arkadiuszjuszczyk.com
```

Built in Go with [Wish](https://github.com/charmbracelet/wish), [Bubble Tea](https://github.com/charmbracelet/bubbletea), and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

## Quick start

```bash
go run . --port 2222
```

Then in another terminal:

```bash
ssh -p 2222 localhost
```

A host key is generated on first run under `.ssh/id_ed25519`.

## Flags

```
--host        host to bind to (default: ::, accepts IPv6 and IPv4)
--port        port to listen on (default: 2222)
--host-key    path to SSH host key (default: .ssh/id_ed25519)
```

## Project layout

```
.
├── main.go              # entry point, Wish server, tea.Model
├── boot.go              # animated boot sequence
├── styles.go            # Lip Gloss styles and color palette
├── content.go           # bio, experience, stack, projects, contacts
├── deploy.sh            # cross-compile + deploy script
└── deploy/
    └── terminal-interface.service  # systemd unit template
```

## Deployment

See [DEPLOY.md](./DEPLOY.md) for one-time server setup and deploy instructions.

## License

[MIT](./LICENSE)
