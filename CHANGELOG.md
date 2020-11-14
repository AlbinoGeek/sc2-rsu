# Changelog

## Unreleased

**Added**

- Checks replay status until the replay has been processed

**Fixed**

- Bug where empty replay files were sent before SC2 wrote them out
- Bug where accounts could not be found right after locating replays root

## v0.2

**Added**

- More concise non-verbose logging
- More detailed verbose logging
- More descriptive error messages
- `update` command to check for and optionally download updates
- Periodic check for updates (configurable)
- Automatically download updates when found (configurable)

**Changed**

- Faster replays directory searching on Windows and Linux
- Faster browser login process (less QuerySelector)
- SC2ReplayStats functions moved to own package

**Fixed**

- Bug where the root directory was offered as Reply directory in some cases
- Configuration not being detected in home directory

**Removed**

- Preliminary support for multiple API keys removed

## v0.1

- Initial Tagged Build
