package installer

import (
	"encoding/base64"
	"strconv"
	"unicode/utf16"
)

// script is a PowerShell script and the environment it is given. The
// app's names and paths reach it only through the environment, never
// written into its text, so none of them can be read as code.
type script struct {
	text string
	env  []string
}

// windowsInstall makes the Start menu's shortcut and the installed
// apps' entry, which runs the uninstaller; kib is the size installed.
func windowsInstall(app App, l layout, kib int64) script {
	text := `$ErrorActionPreference = 'Stop'
$s = (New-Object -ComObject WScript.Shell).CreateShortcut($env:SGE_SHORTCUT)
$s.TargetPath = $env:SGE_EXE
$s.WorkingDirectory = $env:SGE_DIR
$s.IconLocation = $env:SGE_ICON
$s.Description = $env:SGE_NAME
$s.Save()
$k = New-Item -Path 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall' -Name $env:SGE_ID -Force
foreach ($v in 'DisplayName','DisplayVersion','Publisher','InstallLocation','DisplayIcon','UninstallString','QuietUninstallString') {
	Set-ItemProperty -LiteralPath $k.PSPath -Name $v -Value ([Environment]::GetEnvironmentVariable('SGE_' + $v))
}
foreach ($v in 'NoModify','NoRepair') { Set-ItemProperty -LiteralPath $k.PSPath -Name $v -Value 1 -Type DWord }
Set-ItemProperty -LiteralPath $k.PSPath -Name EstimatedSize -Value ([int]$env:SGE_SIZE) -Type DWord
`
	uninstall := `"` + l.join(l.Dir, l.Uninstaller) + `" --uninstall`
	return script{text: text, env: []string{
		"SGE_SHORTCUT=" + l.Shortcut, "SGE_EXE=" + l.Exe, "SGE_DIR=" + l.Dir, "SGE_ICON=" + l.IconLoc,
		"SGE_NAME=" + app.Name, "SGE_ID=" + app.ID, "SGE_DisplayName=" + app.Name,
		"SGE_DisplayVersion=" + app.Version, "SGE_Publisher=" + app.Publisher,
		"SGE_InstallLocation=" + l.Dir, "SGE_DisplayIcon=" + l.IconLoc,
		"SGE_UninstallString=" + uninstall, "SGE_QuietUninstallString=" + uninstall + " --quiet",
		"SGE_SIZE=" + strconv.FormatInt(min(kib, 1<<31-1), 10),
	}}
}

// windowsUninstall takes the installed apps' entry away.
func windowsUninstall(app App) script {
	text := `$k = 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\' + $env:SGE_ID
if (Test-Path -LiteralPath $k) { Remove-Item -LiteralPath $k -Recurse -Force }
`
	return script{text: text, env: []string{"SGE_ID=" + app.ID}}
}

// runScript runs sc in the system's own PowerShell, by its full path so
// no program of the same name elsewhere is run instead, the script
// encoded as PowerShell takes it whole, with no quoting to get wrong.
func runScript(s Setup, sc script) error {
	ps := joinOn(Windows, s.Place.SystemRoot, "System32", "WindowsPowerShell", "v1.0", "powershell.exe")
	args := []string{"-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-EncodedCommand", encoded(sc.text)}
	return s.runner()(ps, args, sc.env)
}

// encoded is text as -EncodedCommand takes it: UTF-16LE in base64.
func encoded(text string) string {
	units := utf16.Encode([]rune(text))
	b := make([]byte, 0, 2*len(units))
	for _, u := range units {
		b = append(b, byte(u), byte(u>>8))
	}
	return base64.StdEncoding.EncodeToString(b)
}
