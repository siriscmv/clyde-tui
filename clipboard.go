package main

import "golang.design/x/clipboard"

func InitClipboard() {
	err := clipboard.Init()

	if err != nil {
		tui.Send(logMsg{Msg: "Unable to init clipboard", Type: Error})
	}
}

func ReadClipboard() string {
	return string(clipboard.Read(clipboard.FmtText))
}

func WriteClipboard(content string) {
	clipboard.Write(clipboard.FmtText, []byte(content))
}
