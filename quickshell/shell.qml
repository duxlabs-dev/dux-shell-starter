//@ pragma Env QT_WAYLAND_DISABLE_WINDOWDECORATION=1
//@ pragma Env QSG_RENDER_LOOP=threaded
//@ pragma Env QT_QUICK_CONTROLS_STYLE=Material
//@ pragma UseQApplication
//@ pragma AppId com.example.myshell
//@ pragma RenderControl
//@ pragma DiskCache

pragma ComponentBehavior: Bound

import QtQuick
import Quickshell
import qs.Modules

ShellRoot {
    Component.onCompleted: Quickshell.watchFiles = Quickshell.env("MYSHELL_DISABLE_HOT_RELOAD") !== "1"

    Variants {
        model: Quickshell.screens

        Bar {
            required property var modelData

            screen: modelData
        }
    }
}
