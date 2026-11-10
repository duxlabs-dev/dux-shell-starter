BINARY_NAME=myshell
CORE_DIR=core
BUILD_DIR=$(CORE_DIR)/bin
SHELL_DIR=quickshell
ASSETS_DIR=assets
PREFIX ?= /usr/local
DESTDIR ?=
INSTALL_DIR=$(PREFIX)/bin
SYSTEMD_USER_DIR=$(PREFIX)/lib/systemd/user

.PHONY: all build dev run clean test fmt vet lint update-common update-qtbase install uninstall help

all: build

build:
	@$(MAKE) -C $(CORE_DIR) build

dev:
	@$(MAKE) -C $(CORE_DIR) dev

run: dev
	@$(BUILD_DIR)/$(BINARY_NAME) run -c $(CURDIR)/$(SHELL_DIR)

clean:
	@$(MAKE) -C $(CORE_DIR) clean

test:
	@$(MAKE) -C $(CORE_DIR) test

fmt:
	@$(MAKE) -C $(CORE_DIR) fmt

vet:
	@$(MAKE) -C $(CORE_DIR) vet

lint:
	@$(MAKE) -C $(CORE_DIR) lint

update-common:
	git submodule update --remote --merge dux-qml-common

update-qtbase:
	git submodule update --remote --merge dux-qt-base

install:
	@test -f $(BUILD_DIR)/$(BINARY_NAME) || { echo "$(BUILD_DIR)/$(BINARY_NAME) not found; run 'make' first"; exit 1; }
	@install -D -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(DESTDIR)$(INSTALL_DIR)/$(BINARY_NAME)
	@install -d -m 755 $(DESTDIR)$(SYSTEMD_USER_DIR)
	@sed 's|/usr/bin/$(BINARY_NAME)|$(INSTALL_DIR)/$(BINARY_NAME)|g' $(ASSETS_DIR)/systemd/$(BINARY_NAME).service > $(DESTDIR)$(SYSTEMD_USER_DIR)/$(BINARY_NAME).service
	@chmod 644 $(DESTDIR)$(SYSTEMD_USER_DIR)/$(BINARY_NAME).service
	@echo "Installed. Start it with '$(BINARY_NAME) run' or 'systemctl --user enable --now $(BINARY_NAME)'."

uninstall:
	@rm -f $(DESTDIR)$(INSTALL_DIR)/$(BINARY_NAME)
	@rm -f $(DESTDIR)$(SYSTEMD_USER_DIR)/$(BINARY_NAME).service

help:
	@echo "build          - Release build, UI embedded (QTBASE=1 also embeds the dux-qt-base runtime)"
	@echo "dev            - Fast build without the embedded UI"
	@echo "run            - dev build, then run against quickshell/ with hot reload"
	@echo "test fmt vet lint clean"
	@echo "update-common  - Bump the dux-qml-common submodule"
	@echo "update-qtbase  - Bump the dux-qt-base submodule"
	@echo "install        - Binary and systemd user unit under PREFIX=$(PREFIX)"
	@echo "uninstall      - Remove them"
