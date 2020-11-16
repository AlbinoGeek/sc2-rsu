// +build !windows

package utils

// AttachConsole gives us access to the parent console on Windows, when built
// as a GUI application, where the console would otherwise not be available.
func AttachConsole() error {
	return nil
}

// FixRedirection restores the parent console's ability to redirect this
// program's output through a pipe to a file or device, where without
// doing so, the output would be forced to the console even if redirected.
func FixRedirection() error {
	return nil
}

// FreeConsole would relinquish our control of the parent console on Windows
// ? but is it actually necessary ? requires more testing
// TODO: NOT YET IMPLEMENTED
func FreeConsole() error {
	return nil
}
