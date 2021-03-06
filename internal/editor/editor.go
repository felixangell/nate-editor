package editor

import (
	"github.com/felixangell/phi/internal/buff"
	"github.com/felixangell/phi/internal/cfg"
	"github.com/felixangell/phi/internal/command_handler"
	"github.com/felixangell/phi/internal/gui"
	"github.com/felixangell/strife"
	"io/ioutil"
	"log"
	"os"
)

type PhiEditor struct {
	running     bool
	mainView    *buff.BufferView
}

func NewPhiEditor() *PhiEditor {
	return &PhiEditor{running: true}
}

func (n *PhiEditor) Resize(w, h int) {
	n.mainView.Resize(w, h)
}

func (n *PhiEditor) HandleEvent(_ strife.StrifeEvent) {}

func (n *PhiEditor) ApplyConfig(conf *cfg.PhiEditorConfig) {
	gui.LoadDefaultFont(conf)
	command_handler.SetupCommandHandler(conf)

	mainView := buff.NewView(int(1280.0*cfg.ScaleFactor), int(720.0*cfg.ScaleFactor), conf)

	args := os.Args
	if len(args) > 1 {
		// TODO check these are files
		// that actually exist here?
		for _, arg := range args[1:] {
			mainView.AddBuffer().OpenFile(arg)
		}
	} else {
		// we have no args, open up a scratch file
		tempFile, err := ioutil.TempFile("", "phi-editor-")
		if err != nil {
			log.Println("Failed to create temp file", err.Error())
			os.Exit(1)
		}

		mainView.AddBuffer().OpenFile(tempFile.Name())
	}

	n.mainView = mainView
}

func (n *PhiEditor) Update() bool {
	return n.mainView.OnUpdate()
}

func (n *PhiEditor) Render(ctx *strife.Renderer) {
	ctx.SetFont(gui.GetDefaultFont())
	n.mainView.OnRender(ctx)
}
