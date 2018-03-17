// Package view provides command
package view

import (
	"fmt"
	"log"

	"github.com/koolay/console-chat/rethink"
	"github.com/marcusolsson/tui-go"
)

type ChatUI struct {
	rawTUI   tui.UI
	Messages []rethink.MessageBase
}

func NewChatUI() *ChatUI {
	sidebar := tui.NewVBox(
		tui.NewLabel("CHANNELS"),
		tui.NewLabel(""),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		err := rethink.RethinkActor.SendPublicMessage(e.Text())
		if err != nil {
			panic(err)
		}
		//history.Append(tui.NewHBox(
		//tui.NewLabel(time.Now().Format("15:04")),
		//tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", "john"))),
		//tui.NewLabel(e.Text()),
		//tui.NewSpacer(),
		//))
		input.SetText("")
	})
	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("Up", func() { historyScroll.Scroll(0, -1) })
	ui.SetKeybinding("Down", func() { historyScroll.Scroll(0, 1) })
	ui.SetKeybinding("Left", func() { historyScroll.Scroll(-1, 0) })
	ui.SetKeybinding("Right", func() { historyScroll.Scroll(1, 0) })

	go func() {
		for {
			record := <-rethink.ChatRecordChan
			history.Append(tui.NewHBox(
				tui.NewLabel(record.TimeSince()),
				tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", record.Sender))),
				tui.NewLabel(record.Content),
				tui.NewSpacer(),
			))
			rethink.ChatRecordPoll.Put(record)
		}
	}()
	return &ChatUI{
		rawTUI: ui,
	}
}

func (p *ChatUI) Run() error {
	return p.rawTUI.Run()
}
