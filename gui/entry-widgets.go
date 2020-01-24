package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// ExtendedEntry is used to make an entry that reacts to key presses.
type ExtendedEntry struct {
	widget.Entry
	*Action
}

// Action handles the Button press action.
type Action struct {
	widget.Button
}

// TypedKey handles the key presses inside our UsernameEntry and uses Action to press the linked button.
func (e *ExtendedEntry) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyReturn:
		e.Action.Button.OnTapped()
	default:
		e.Entry.TypedKey(ev)
	}
}

// NewExtendedEntry creates an ExtendedEntry button.
func NewExtendedEntry(placeholder string, password bool) *ExtendedEntry {
	entry := &ExtendedEntry{}

	// Extend the base widget.
	entry.ExtendBaseWidget(entry)

	// Set placeholder for the entry.
	entry.SetPlaceHolder(placeholder)

	// Check if we are creating a password entry.
	if password {
		entry.Password = true
	}

	return entry
}

// NewEntryWithPlaceholder makes it easy to create entry widgets with placeholders.
func NewEntryWithPlaceholder(text string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(text)

	return entry
}