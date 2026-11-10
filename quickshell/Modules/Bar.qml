pragma ComponentBehavior: Bound

import QtQuick
import Quickshell
import Quickshell.Wayland
import qs.DCommon.Common
import qs.DCommon.Widgets
import qs.Services

PanelWindow {
    id: bar

    anchors.top: true
    anchors.left: true
    anchors.right: true
    implicitHeight: 36
    color: Style.surface
    visible: BarState.visible
    WlrLayershell.namespace: "myshell:bar"

    SystemClock {
        id: clock
        precision: SystemClock.Minutes
    }

    StyledText {
        anchors.left: parent.left
        anchors.leftMargin: Style.spacingL
        anchors.verticalCenter: parent.verticalCenter
        text: bar.screen?.name ?? ""
        color: Style.surfaceVariantText
    }

    StyledText {
        anchors.centerIn: parent
        text: Qt.formatDateTime(clock.date, "ddd d MMM  hh:mm")
        font.weight: Font.Medium
    }

    Row {
        anchors.right: parent.right
        anchors.rightMargin: Style.spacingL
        anchors.verticalCenter: parent.verticalCenter
        spacing: Style.spacingM

        BarReading {
            iconName: "memory"
            text: SystemInfo.memoryUsedPercent + "%"
        }

        BarReading {
            iconName: "speed"
            text: SystemInfo.load.toFixed(2)
        }
    }
}
