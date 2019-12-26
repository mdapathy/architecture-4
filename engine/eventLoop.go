package engine

import (
	"sync"
)

// Command represents actions that can be performed in a single event loop iteration.
type Command interface {
	Execute(handler Handler)
}

// Handler allows to send commands to an event loop it's associated with.
type Handler interface {
	Post(cmd Command)
}

type CommandFunc func(handler Handler)

func (c CommandFunc) Execute(handler Handler) {
	c(handler)
}

type messageQueue struct {
	data struct {
		sync.Mutex
		arr     []Command
		waiting bool
	}

	receiveSignal chan struct{}
}

func (queue *messageQueue) waitForSignal() {
	<-queue.receiveSignal
}

func (queue *messageQueue) size() int {
	return len(queue.data.arr)
}

type EventLoop struct {
	queue            *messageQueue
	receivingStopped bool
	stopSignal       chan struct{}
}

func (queue *messageQueue) push(cmd Command) {
	queue.data.Lock()
	defer queue.data.Unlock()
	queue.data.arr = append(queue.data.arr, cmd)

	if queue.data.waiting {
		queue.data.waiting = false
		queue.receiveSignal <- struct{}{}

	}

}

func (queue *messageQueue) pull() Command {
	queue.data.Lock()
	defer queue.data.Unlock()

	if len(queue.data.arr) == 0 {
		queue.data.waiting = true
		queue.data.Unlock()
		<-queue.receiveSignal
		queue.data.Lock()
	}

	res := queue.data.arr[0]
	queue.data.arr[0] = nil
	queue.data.arr = queue.data.arr[1:]

	return res
}

func (loop *EventLoop) Start() {
	loop.queue = new(messageQueue)
	loop.queue.receiveSignal = make(chan struct{})
	loop.stopSignal = make(chan struct{})

	go func() {
		for !loop.receivingStopped || loop.queue.size() != 0 {
			cmd := loop.queue.pull()
			cmd.Execute(loop)
		}
		loop.stopSignal <- struct{}{}
	}()
}

func (loop *EventLoop) Post(cmd Command) {
	loop.queue.push(cmd)
}

func (loop *EventLoop) AwaitFinish() {
	loop.Post(CommandFunc(func(h Handler) { h.(*EventLoop).receivingStopped = true }))
	<-loop.stopSignal
}
