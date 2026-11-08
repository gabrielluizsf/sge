// Package installer puts a game on the player's system, for the player
// alone, with no administrator: its files in a folder of their own and
// an entry where the system lists its programs.
//
//   - Linux: the files in ~/.local/share/<ID>, a desktop entry in
//     ~/.local/share/applications and a link in ~/.local/bin.
//   - Windows: the files in %LOCALAPPDATA%\Programs\<ID>, a shortcut in the
//     Start menu and an entry in the installed apps, which uninstalls it.
//   - macOS: an application bundle, ~/Applications/<Name>.app.
//
// An installer is the game's own small program carrying the game's files
// after it: Pack writes one, Main runs it. It shows a window with a bar
// while it copies the files and then runs the game's own Steps, such as
// baking its pictures, so it is the game's first loading screen too; and
// then it starts the game. Run with --uninstall, it takes it all away
// again. What it copies is checked as untrusted: no name may leave the
// game's folder, no link is made, and sizes are bounded.
package installer
