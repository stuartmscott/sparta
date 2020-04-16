package gui

import (
	"sparta/crypto"
	"sparta/gui/widgets"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// SettingsView contains the gui information for the settings screen.
func (u *user) settingsView(w fyne.Window, a fyne.App) fyne.CanvasObject {

	// TODO: Add setting for changing language.

	// Make it possible for the user to switch themes.
	themeSwitcher := widget.NewSelect([]string{"Dark", "Light"}, func(selected string) {
		switch selected {
		case "Dark":
			a.Settings().SetTheme(theme.DarkTheme())
		case "Light":
			a.Settings().SetTheme(theme.LightTheme())
		}

		// Set the theme to the selected one and save it using the preferences api in fyne.
		a.Preferences().SetString("Theme", selected)
	})

	// Default theme is light and thus we set the placeholder to that and then refresh it (without a refresh, it doesn't show until hovering on to widget).
	themeSwitcher.SetSelected(a.Preferences().StringWithFallback("Theme", "Light"))

	// Add the theme switcher next to a label.
	themeChanger := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Application Theme"), themeSwitcher)

	// An entry for typing the new username.
	usernameEntry := widgets.NewAdvancedEntry("New Username", false)

	// Create the button used for changing the username.
	usernameButton := widget.NewButtonWithIcon("Change Username", theme.ConfirmIcon(), func() {
		// Check that the username is valid.
		if usernameEntry.Text == u.password || usernameEntry.Text == "" {
			dialog.ShowInformation("Please enter a valid username", "Usernames need to not be empty and not the same as the password.", w)
		} else {
			// Ask the user to confirm what we are about to do.
			dialog.ShowConfirm("Are you sure that you want to continue?", "The action will permanently change your username.", func(change bool) {
				if change {
					// Replace the password hash in a new storage location.
					a.Preferences().RemoveValue("Username:" + u.username)
					a.Preferences().SetString("Username:"+usernameEntry.Text, u.passwordHash)

					// Set the username  to the updated username.
					u.username = usernameEntry.Text

					// Clear out the text inside the entry.
					usernameEntry.SetText("")

					// Write the data encrypted using the new key and do so concurrently.
					go u.data.Write(&u.encryptionKey, u.username)
				}
			}, w)
		}

	})

	// Create the entry for updating the password.
	passwordEntry := widgets.NewAdvancedEntry("New Password", true)

	// Create the button used for changing the password.
	passwordButton := widget.NewButtonWithIcon("Change Password", theme.ConfirmIcon(), func() {
		// Check that the password is valid.
		if len(passwordEntry.Text) < 8 || passwordEntry.Text == u.username {
			dialog.ShowInformation("Please enter a valid password", "Passwords need to be at least eight characters long.", w)
		} else {
			// Ask the user to confirm what we are about to do.
			dialog.ShowConfirm("Are you sure that you want to continue?", "The action will permanently change your password.", func(change bool) {
				if change {
					// Calculate and store the new hashes.
					u.encryptionKey = crypto.GenerateEncryptionKey(passwordEntry.Text)
					u.passwordHash = crypto.GeneratePasswordHash(passwordEntry.Text)

					// Update the password hash in storage.
					a.Preferences().SetString("Username:"+u.username, u.passwordHash)

					// Clear out the text inside the entry.
					passwordEntry.SetText("")

					// Write the data encrypted using the new key and do so concurrently.
					go u.data.Write(&u.encryptionKey, u.username)
				}
			}, w)
		}
	})

	// Extend our extended buttons with array entry switching and enter to change.
	usernameEntry.InitExtend(*usernameButton, widgets.MoveAction{Down: true, DownEntry: passwordEntry, Window: w})
	passwordEntry.InitExtend(*passwordButton, widgets.MoveAction{Up: true, UpEntry: usernameEntry, Window: w})

	// revertToDefaultSettings reverts all settings to their default values.
	revertToDefaultSettings := widget.NewButtonWithIcon("Reset settings to default values", theme.ViewRefreshIcon(), func() {
		// Update theme and saved settings for theme change.
		if a.Preferences().String("Theme") != "Light" {
			themeSwitcher.SetSelected("Light")

			// Set the visible theme to the light theme.
			a.Settings().SetTheme(theme.LightTheme())

			// Set the saved theme to Light.
			a.Preferences().SetString("Theme", "Light")
		}
	})
	// Create a button for clearing the data of a given profile.
	deleteButton := widget.NewButtonWithIcon("Delete all saved activities", theme.DeleteIcon(), func() {

		// Ask the user to confirm what we are about to do.
		dialog.ShowConfirm("Are you sure that you want to continue?", "Deleting your data will remove all of your exercises and activities.", func(remove bool) {
			if remove {
				// Run the delete function and do it concurrently to avoid stalling the thread with file io.
				go u.data.Delete(u.username)

				// Notify the label that we have removed the data.
				u.emptyExercises <- true
			}
		}, w)
	})

	// userInterfaceSettings is a group holding widgets related to user interface settings such as theme.
	userInterfaceSettings := widget.NewGroup("User Interface Settings", themeChanger)

	// credentialSettings groups together all settings related to usernames and passwords.
	credentialSettings := widget.NewGroup("Login Credential Settings", fyne.NewContainerWithLayout(layout.NewGridLayout(2), usernameEntry, usernameButton, passwordEntry, passwordButton))

	// advancedSettings is a group holding widgets related to advanced settings.
	advancedSettings := widget.NewGroup("Advanced Settings", revertToDefaultSettings, widget.NewLabel(""), deleteButton)

	// settingsContentView holds all widget groups and content for the settings page.
	settingsContentView := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), userInterfaceSettings, layout.NewSpacer(), credentialSettings, layout.NewSpacer(), advancedSettings)

	return widget.NewScrollContainer(settingsContentView)
}