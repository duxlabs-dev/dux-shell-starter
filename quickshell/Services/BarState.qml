pragma Singleton
pragma ComponentBehavior: Bound

import QtQuick
import Quickshell
import qs.Services

Singleton {
    id: root

    property bool visible: true

    function toggle() {
        Backend.call("bar.toggle", null, null);
    }

    Component.onCompleted: Backend.subscribe("bar", state => root.visible = state.visible)
}
