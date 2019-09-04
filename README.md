# app
Reusable framework for Go apps & CLIs

## Examples
### App
### [Short](https://github.com/byliuyang/short)
![Short screenshots](example/short.png)

### CLI
![CLI screenshots](example/cli.png)
```go
import (
	"fmt"
	"os"

	"github.com/byliuyang/app/tool/cli"

	"github.com/byliuyang/app/tool/terminal"
	"github.com/byliuyang/app/tool/ui"
	"github.com/byliuyang/eventbus"
	"github.com/spf13/cobra"
)

type ExampleTool struct {
	term            terminal.Terminal
	exitChannel     eventbus.DataChannel
	keyUpChannel    eventbus.DataChannel
	keyDownChannel  eventbus.DataChannel
	keyEnterChannel eventbus.DataChannel
	cli             cli.CommandLineTool
	rootCmd         *cobra.Command
	radio           ui.Radio
	languages       []string
}

func (e ExampleTool) Execute() {
	if err := e.rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (e ExampleTool) bindKeys() {
	e.term.OnKeyPress(terminal.CtrlEName, e.exitChannel)
	e.term.OnKeyPress(terminal.CursorUpName, e.keyUpChannel)
	e.term.OnKeyPress(terminal.CursorDownName, e.keyDownChannel)
	e.term.OnKeyPress(terminal.EnterName, e.keyEnterChannel)
	fmt.Println("To exit, press Ctrl + E")
	fmt.Println("To select an item, press Enter")
}

func (e ExampleTool) handleEvents() {
	e.cli.EnterMainLoop(func() {
		select {
		case <-e.exitChannel:
			e.radio.Remove()
			fmt.Println("Terminating process...")
			e.cli.Exit()
		case <-e.keyUpChannel:
			e.radio.Prev()
		case <-e.keyDownChannel:
			e.radio.Next()
		case <-e.keyEnterChannel:
			e.radio.Remove()
			selectedItem := e.languages[e.radio.SelectedIdx()]
			fmt.Printf("Selected %s\n", selectedItem)
			e.cli.Exit()
		}
	})
}

func NewExampleTool() ExampleTool {
	term := terminal.NewTerminal()
	languages := []string{
		"Go",
		"C++",
		"Java",
		"Python",
		"JavaScript",
		"Rust",
	}

	exampleTool := ExampleTool{
		term:            term,
		cli:             cli.NewCommandLineTool(term),
		exitChannel:     make(eventbus.DataChannel),
		keyUpChannel:    make(eventbus.DataChannel),
		keyDownChannel:  make(eventbus.DataChannel),
		keyEnterChannel: make(eventbus.DataChannel),
		radio:           ui.NewRadio(languages, term),
		languages:       languages,
	}
	rootCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			exampleTool.bindKeys()
			exampleTool.radio.Render()
			exampleTool.handleEvents()
		},
	}
	exampleTool.rootCmd = rootCmd
	return exampleTool
}
```

## Author
Harry Liu - [byliuyang](https://github.com/byliuyang)

## License
This project is maintained under MIT license