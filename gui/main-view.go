package gui

import (
	"sparta/file"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// ShowMainDataView shows the main view after we are logged in.
func ShowMainDataView(window fyne.Window, app fyne.App, XMLData *file.Data, newAddedExercise chan string) {
	// Create a label for displaing some info for the user. Default to showing nothing.
	dataLabel := widget.NewLabel("")

	go func() {

		// Handle an empty data file.
		if file.Empty() {
			// Start by inorming  the user that no data is avaliable.
			dataLabel.SetText("No exercieses have been created yet.")

			// Then wait for more data to come running down the pipe.
			dataLabel.SetText(<-newAddedExercise)
		} else {
			// We loop through the imported file and add the formated info before the previous info (new information comes out on top).
			for i := range XMLData.Exercise {
				dataLabel.SetText(XMLData.Format(i) + dataLabel.Text)
				// TODO: When #607 is fixed in fyne we can optimize this by just setting the text and refreshing the widget when we are done.
			}
		}

		// We then block the channel while waiting for an update on the channel.
		for {
			dataLabel.SetText(<-newAddedExercise + dataLabel.Text)
		}
	}()

	// Tab data for the main window.
	dataPage := widget.NewScrollContainer(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), dataLabel))

	// Create tabs with data.
	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon("Activities", theme.HomeIcon(), dataPage),
		widget.NewTabItemWithIcon("Add activity", theme.ContentAddIcon(), ActivityView(window, XMLData, dataLabel, newAddedExercise)),
		widget.NewTabItemWithIcon("Settings", theme.SettingsIcon(), SettingsView(window, app, XMLData, dataLabel)),
		// TODO: Add an about page with logo, name and version number.
	)

	// Set the tabs to be on top of the page.
	tabs.SetTabLocation(widget.TabLocationTop)

	// Adapt the window to a good size and make it resizable again.
	window.SetFixedSize(false)
	window.Resize(fyne.NewSize(800, 500))

	// Set the content to show and do so in a scroll container for the exercieses to show correctly.
	window.SetContent(tabs)
}