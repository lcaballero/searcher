package console

import (
	"os"
	"bufio"
	"fmt"
	"sync"
)


type Prompt struct {
	onRead chan string
	onClose chan *sync.WaitGroup
}

func NewPrompt() *Prompt {
	return &Prompt{
		onRead: make(chan string, 1),
		onClose: make(chan *sync.WaitGroup, 1),
	}
}

func (p *Prompt) Start() {
	reader := bufio.NewReader(os.Stdin)
	onPrompt := make(chan bool, 1)
	prompt := "Search> "

	go func() {
		onPrompt <- true
		for {
			select {
			case <-onPrompt:
				fmt.Print(prompt)
			case clean := <-p.onClose:
				clean.Done()
				return
			}
		}
	}()
}

func (p *Prompt) Close() error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	os.Stdin.Write([]byte{byte(1)})
	p.onClose <- wg
	wg.Wait()
	return nil
}
