package cmd

import (
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/data/binding"
)

type boundValidatedString struct {
	base

	val *string
	// >>> added code
	validator fyne.StringValidator
	// <<< added code
}

func validatedString(fn fyne.StringValidator) *boundValidatedString {
	empty := ""

	return &boundValidatedString{
		val: &empty,
		// >>> added code
		validator: fn,
		// <<< added code
	}
}

func (b *boundValidatedString) Get() (string, error) {
	if b.val == nil {
		return "", nil
	}
	return *b.val, nil
}

func (b *boundValidatedString) Set(val string) error {
	if *b.val == val {
		return nil
	}
	// >>> added code
	if b.validator != nil {
		if err := b.validator(val); err != nil {
			return err
		}
	}
	// <<< added code
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
	return nil
}

// ---

var itemQueue = make(chan itemData, 1024)

type itemData struct {
	fn   func()
	done chan interface{}
}

func queueItem(f func()) {
	itemQueue <- itemData{fn: f}
}

func init() {
	go processItems()
}

func processItems() {
	for {
		i := <-itemQueue
		if i.fn != nil {
			i.fn()
		}
		if i.done != nil {
			i.done <- struct{}{}
		}
	}
}

type base struct {
	listeners []binding.DataListener
	lock      sync.RWMutex
}

// AddListener allows a data listener to be informed of changes to this item.
func (b *base) AddListener(l binding.DataListener) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.listeners = append(b.listeners, l)
	queueItem(l.DataChanged)
}

// RemoveListener should be called if the listener is no longer interested in being informed of data change events.
func (b *base) RemoveListener(l binding.DataListener) {
	b.lock.Lock()
	defer b.lock.Unlock()

	for i, listen := range b.listeners {
		if listen != l {
			continue
		}

		if i == len(b.listeners)-1 {
			b.listeners = b.listeners[:len(b.listeners)-1]
		} else {
			b.listeners = append(b.listeners[:i], b.listeners[i+1:]...)
		}
	}
}

func (b *base) trigger() {
	b.lock.RLock()
	defer b.lock.RUnlock()

	for _, listen := range b.listeners {
		queueItem(listen.DataChanged)
	}
}
