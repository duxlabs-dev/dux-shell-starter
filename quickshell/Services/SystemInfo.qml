pragma Singleton
pragma ComponentBehavior: Bound

import QtQuick
import Quickshell
import qs.Services

Singleton {
    id: root

    property real load: 0
    property int memoryUsedPercent: 0

    Component.onCompleted: Backend.subscribe("sysinfo", info => {
        root.load = info.load;
        root.memoryUsedPercent = info.memoryUsedPercent;
    })
}
