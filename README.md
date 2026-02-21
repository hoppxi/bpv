# bpv

Terminal and Browser based local music player

**BPV** - a lightweight music player made using go. BPV is a daemon based player with multiple clients: terminal based TUI client and Browser based web client. It uses bubbletea for the tui and vuejs for web clients.

---

## features

- multiple clients syncing perfectly
- large music library handling daemon
- beautiful design bot for tui and web clients
- multiple visualizer for the web client
- lightweight by cacheing things

---

## Developing

1. start the dameon

```bash
go run ./cmd/bpvd
```

2. run your client

```bash
# for web client
cd web/
npm install
npm run dev # then open the provided link

# or after running npm run build, Run
go run ./cmd/bpv --client web

# for tui client
go run ./cmd/bpv
```

## Installing and building

You can use the flake, makefile or manual building

```bash
# with nix
nix build .#default # for daemon use .#bpvd, for bpv main binary use .#bpv and for web build use .#web

# with makefile
make build

# manual
go build ./cmd/bpv
go build ./cmd/bpvd
cd web && npm install & npm run build
```

to install, use home-manager or nixos module, manual

```bash
nix profile add github:hoppxi/bpv

# or manually after building move the binaries to $PATH and set BPV_WEB_DIR to the path where your web build is
```

```nix
inputs = {
  bpv.url = "github:hoppxi/bpv"
};
# ...
# in your home-manager setup
services.bpvd.enable = true; # runs bvpd as a systemd service
programs.bpv.enable = true;
```

---

## License

MIT License
