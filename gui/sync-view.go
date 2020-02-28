package gui

import (
	"sparta/gui/widgets"
	"sparta/sync"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// SyncView displays the tab page for syncing data between devices.
func SyncView(window fyne.Window, user *user) fyne.CanvasObject {

	// Create the channel that we will use to share the code needed for receiving the data.
	syncCodeChan := make(chan string)

	// aboutGroup is a small group to display information about file sharing.
	aboutGroup := widget.NewGroup("About Exercise Synchronization", widget.NewLabelWithStyle("The share support in Sparta enables the user to syncronize their exercises.\nThe file is served encrypted and will be automatically decrypted on a second device using the same login credentials.\nIt is as simple as starting the share on one device and then starting receiving on another device on the same network.\nThe receiving computer will then get all the new activities added locally.", fyne.TextAlignCenter, fyne.TextStyle{}))

	// recieveCode makes it possible to type in receive code.
	recieveCodeEntry := widgets.NewEntryWithPlaceholder("Receive Code")

	// recieveDataButton starts looking for shared data on the local network.
	recieveDataButton := widget.NewButtonWithIcon("Start Receiving Exercises", theme.MailComposeIcon(), func() {
		go sync.Retrieve(&user.Data, user.ReorderExercises, user.FirstExercise, user.Errors, &user.EncryptionKey, recieveCodeEntry.Text)
	})

	// sendDataButton starts the network server and shares the file over the local internet.
	startSendingDataButton := widget.NewButtonWithIcon("Start Syncing Exercises", theme.ViewRefreshIcon(), func() {
		if len(user.Data.Exercise) == 0 {
			dialog.ShowInformation("No exercises to syncronize", "You need to add exercises to be able to sync between devices.", window)
			return
		}

		// Create a channel for comunicating when we are done.
		finished := make(chan struct{})

		// Disable the receive button to hinder people from receiving on same computer.
		recieveDataButton.Disable()

		// Start the sharing.
		go sync.StartSync(syncCodeChan, user.Errors, finished)

		// Listen in on the channels and clear stuff up when we are finished.
		go func() {
			for {
				select {
				case code := <-syncCodeChan:
					recieveCodeEntry.SetText(code)
				case <-finished:
					recieveDataButton.Enable()
					dialog.ShowInformation("Syncronization sucessfull", "The syncronization of exercises finsished sucessfully.", window)
					return
				case err := <-user.Errors:
					dialog.ShowError(err, window)
					recieveDataButton.Enable()
					return
				}
			}
		}()
	})

	// shareGroup is a group containing all the options for sharing data.
	shareGroup := widget.NewGroup("Syncronizing Data", startSendingDataButton)

	// recieveGroup is a group containing all options related to receiving.
	recieveGroup := widget.NewGroup("Receiving Data", fyne.NewContainerWithLayout(layout.NewGridLayout(2), recieveCodeEntry, recieveDataButton))

	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(), aboutGroup, layout.NewSpacer(), shareGroup, layout.NewSpacer(), recieveGroup)
}
