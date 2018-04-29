package gui

import (
	"github.com/felixangell/phi/cfg"
	"github.com/felixangell/strife"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"unicode"
)

// View is an array of buffers basically.
type View struct {
	BaseComponent

	conf           *cfg.TomlConfig
	buffers        map[int]*Buffer
	focusedBuff    int
	commandPalette *CommandPalette
}

func NewView(width, height int, conf *cfg.TomlConfig) *View {
	view := &View{
		conf:    conf,
		buffers: map[int]*Buffer{},
	}

	view.Translate(width, height)
	view.Resize(width, height)

	view.commandPalette = NewCommandPalette(*conf, view)
	view.UnfocusBuffers()

	return view
}

func (n *View) hidePalette() {
	p := n.commandPalette
	p.clearInput()
	p.SetFocus(false)

	// set focus to the buffer
	// that invoked the cmd palette
	if p.parentBuff != nil {
		p.parentBuff.SetFocus(true)
		n.focusedBuff = p.parentBuff.index
	}

	// remove focus from palette
	p.buff.SetFocus(false)
}

func (n *View) focusPalette(buff *Buffer) {
	p := n.commandPalette
	p.SetFocus(true)

	// focus the command palette
	p.buff.SetFocus(true)

	// remove focus from the buffer
	// that invoked the command palette
	p.parentBuff = buff
}

func (n *View) UnfocusBuffers() {
	// clear focus from buffers
	for _, buff := range n.buffers {
		buff.SetFocus(false)
	}
}

func sign(dir int) int {
	if dir > 0 {
		return 1
	} else if dir < 0 {
		return -1
	}
	return 0
}

func (n *View) removeBuffer(index int) {
	log.Println("Removing buffer index:", index)
	delete(n.buffers, index)

	// only resize the buffers if we have
	// some remaining in the window
	if len(n.buffers) > 0 {
		bufferWidth := n.w / len(n.buffers)

		// translate all the components accordingly.
		for i, buff := range n.buffers {
			buff.Resize(bufferWidth, n.h)
			buff.SetPosition(bufferWidth*i, 0)
		}
	}

}

func (n *View) ChangeFocus(dir int) {
	prevBuff, _ := n.buffers[n.focusedBuff]

	if dir == -1 {
		n.focusedBuff--
	} else if dir == 1 {
		n.focusedBuff++
	}

	if n.focusedBuff < 0 {
		n.focusedBuff = len(n.buffers) - 1
	} else if n.focusedBuff >= len(n.buffers) {
		n.focusedBuff = 0
	}

	prevBuff.SetFocus(false)
	if buff, ok := n.buffers[n.focusedBuff]; ok {
		buff.SetFocus(true)
	}
}

func (n *View) getCurrentBuff() *Buffer {
	if buff, ok := n.buffers[n.focusedBuff]; ok {
		return buff
	}
	return nil
}

func (n *View) OnInit() {
}

func (n *View) OnUpdate() bool {
	dirty := false

	CONTROL_DOWN = strife.KeyPressed(sdl.K_LCTRL) || strife.KeyPressed(sdl.K_RCTRL)
	if CONTROL_DOWN && strife.PollKeys() {
		r := rune(strife.PopKey())

		actionName, actionExists := cfg.Shortcuts.Controls[string(unicode.ToLower(r))]
		if actionExists {
			if action, ok := actions[actionName]; ok {
				log.Println("Executing action '" + actionName + "'")
				return action.proc(n, []string{})
			}
		} else {
			log.Println("warning, unimplemented shortcut ctrl +", string(unicode.ToLower(r)), actionName)
		}
	}

	if buff, ok := n.buffers[n.focusedBuff]; ok {
		buff.processInput(nil)
		buff.OnUpdate()
	}

	n.commandPalette.OnUpdate()

	return dirty
}

func (n *View) OnRender(ctx *strife.Renderer) {
	for _, buffer := range n.buffers {
		buffer.OnRender(ctx)
	}

	n.commandPalette.OnRender(ctx)
}

func (n *View) OnDispose() {}

func (n *View) AddBuffer() *Buffer {
	n.UnfocusBuffers()

	cfg := n.conf
	c := NewBuffer(cfg, BufferConfig{
		cfg.Theme.Background,
		cfg.Theme.Foreground,
		cfg.Theme.Cursor,
		cfg.Theme.Cursor_Invert,
		cfg.Theme.Gutter_Background,
		cfg.Theme.Gutter_Foreground,
	}, n, len(n.buffers))

	c.SetFocus(true)

	// work out the size of the buffer and set it
	// note that we +1 the components because
	// we haven't yet added the panel
	var bufferWidth int
	bufferWidth = n.w / (c.index + 1)

	n.buffers[c.index] = c
	n.focusedBuff = c.index

	// translate all the components accordingly.
	for i, buff := range n.buffers {
		buff.Resize(bufferWidth, n.h)
		buff.SetPosition(bufferWidth*i, 0)
	}

	return c
}
