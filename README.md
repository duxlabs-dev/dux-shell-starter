# dux-shell-starter

A template for a desktop shell built like [Dux Shell](https://github.com/duxlabs-dev/dux-shell): a Go binary that carries its own QML UI and runs it on [dux-qt-base](https://github.com/duxlabs-dev/dux-qt-base).

The demo is a top bar on every screen. The clock is plain QML, the memory and load readings are pushed from the Go backend every two seconds, and `myshell ipc bar.toggle` hides and shows the bar.

## What's in here

| Path | What it is |
|---|---|
| `core/` | The Go side. [dux-go](https://github.com/duxlabs-dev/dux-go) provides the `run`, `restart` and `kill` commands, supervises the UI process, and serves the unix socket the UI talks to. |
| `quickshell/` | The QML side. `shell.qml` is the entrypoint and `Services/Backend.qml` is the socket client. |
| `dux-qml-common/` | Submodule with the shared widgets and theme tokens, symlinked in as `quickshell/DCommon`. |
| `dux-qt-base/` | Submodule with the runtime. Only built with `QTBASE=1`. |
| `assets/systemd/` | User unit that starts the shell with the session. |
| `scripts/format-staged.py` | Formats staged QML with qmlls. |

## Build and run

You need Go, dux-qt-base (the `dux-qt-base` package, `libdux-qt-base0` on Debian and Ubuntu) and a Wayland compositor with layer shell support, like niri, Hyprland or Sway.

```sh
git clone --recurse-submodules https://github.com/you/yourshell.git
cd yourshell
make run
```

`make run` makes a build without the embedded UI and loads `quickshell/` straight from the repo, with hot reload.

```sh
make
sudo make install
```

`make` builds `core/bin/myshell` with the UI embedded. `make install` puts the binary and the systemd user unit under `/usr/local`.

`make QTBASE=1` also builds dux-qt-base from the submodule and embeds it, so the binary runs without the package installed. The build dependencies are listed in [BUILD.md](https://github.com/duxlabs-dev/dux-qt-base/blob/master/BUILD.md).

## Using it

```sh
myshell run               # run in the foreground
myshell run -d            # run in the background
myshell ipc bar.toggle    # bind this to a key in your compositor
myshell ipc sysinfo.get
myshell restart
myshell kill
```

To start it with the session, run `myshell run` from your compositor config or enable the unit:

```sh
systemctl --user enable --now myshell
```

## How it works

1. `myshell run` starts the backend, which listens on a unix socket in `$XDG_RUNTIME_DIR`.
2. It then starts itself again as `myshell qs`. That process loads `libdux-qt-base.so` and runs `shell.qml`, from the copy embedded in the binary or from the dir given with `-c`.
3. The UI gets the socket path in `MYSHELL_SOCKET`. `Services/Backend.qml` connects to it, sends requests and subscribes to topics.
4. Messages are one JSON object per line. `{"id":1,"method":"bar.toggle"}` is answered with `{"id":1,"result":{"visible":false}}`, and every subscriber of `bar` gets `{"event":"bar","data":{"visible":false}}`.
5. `myshell ipc` uses the same socket, so anything the UI can call you can call from a terminal or a keybind.

## Adding a feature

The bar shows the three pieces a feature needs:

- **Go handlers** in `core/internal/backend`, registered in `routes()` in `backend.go`. `bar.go` changes state on request, `sysinfo.go` publishes on a timer. Add the topic to `replayState` so a UI that subscribes later gets the current state.
- **A QML service** in `quickshell/Services` that subscribes to the topic and exposes it as properties, like `BarState.qml` and `SystemInfo.qml`.
- **A view** that binds to the service, like `Modules/Bar.qml`.

## Theme

The widgets take colors, spacing and fonts from `Style` in dux-qml-common. With nothing set it uses its stock theme. To use your own, hand it your singletons at startup as described in [the contract](https://github.com/duxlabs-dev/dux-qml-common#the-contract).

## Renaming

Four names identify the shell. From the repo root, with your own values in place of `mybar`, `MYBAR`, `org.you.mybar` and `you/mybar`:

```sh
git grep -lE 'myshell|MYSHELL|dux-shell-starter' | xargs sed -i \
  -e 's|duxlabs-dev/dux-shell-starter|you/mybar|g' \
  -e 's|com\.example\.myshell|org.you.mybar|g' \
  -e 's|myshell|mybar|g' \
  -e 's|MYSHELL|MYBAR|g'
git mv core/cmd/myshell core/cmd/mybar
git mv assets/systemd/myshell.service assets/systemd/mybar.service
```

| Name | Used for |
|---|---|
| `myshell` | Binary, socket and pid file names, cache dir, systemd unit, layer shell namespace |
| `MYSHELL` | Prefix of the env vars: `MYSHELL_SOCKET`, `MYSHELL_SHELL_DIR`, `MYSHELL_LOG_LEVEL` |
| `com.example.myshell` | Wayland app id |
| `github.com/duxlabs-dev/dux-shell-starter` | Go module path |

## License

MIT
