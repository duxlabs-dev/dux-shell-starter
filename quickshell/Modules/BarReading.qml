import QtQuick
import qs.DCommon.Common
import qs.DCommon.Widgets

Row {
    id: root

    property string iconName: ""
    property string text: ""

    spacing: Style.spacingXS

    DIcon {
        anchors.verticalCenter: parent.verticalCenter
        name: root.iconName
        size: Style.iconSizeSmall
        color: Style.primary
    }

    StyledText {
        anchors.verticalCenter: parent.verticalCenter
        text: root.text
    }
}
