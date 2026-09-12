package main

func JumpToEndOfLine(m model) model {
    lineLength := len([]rune(m.lines[m.cursor[0]]))

    m.cursor[1] = lineLength
    m.desiredCol = m.cursor[1]
    m.endLine = true

    visualCol := GetVisualCol(m)
    if visualCol > m.horizontalScrollOffset+GetMainPanelWidth(m)-ScrollMargin {
        m.horizontalScrollOffset = visualCol - GetMainPanelWidth(m) + ScrollMargin
    }

    return m
}
