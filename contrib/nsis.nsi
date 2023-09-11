; SPDX-FileCopyrightText: © Moritz Poldrack
; SPDX-License-Identifier: CC-BY-SA-4.0

; Define your application name
!define APPNAME "Uniview"
!define APPNAMEANDVERSION "Uniview 0.3.0"

!define REGUINSTKEY "sh.moritz.uniview"
!define REGUINST 'HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${REGUINSTKEY}"'

; Main Install settings
Name "${APPNAMEANDVERSION}"
InstallDir "$LocalAppData\Programs\Uniview"
RequestExecutionLevel User
OutFile "${PWD}\uniview-setup.exe"
InstallDirRegKey ${REGUINST} UninstallString
SetCompressor /SOLID lzma

; Modern interface settings
!include "MUI2.nsh"

!define MUI_ABORTWARNING

!define MUI_ICON "${PWD}\icon.ico"

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "${PWD}\AGPL.rtf"
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

; Set languages (first is default language)
!insertmacro MUI_LANGUAGE "English"
!insertmacro MUI_RESERVEFILE_LANGDLL

!define /IfNDef SHPPFW_DIRCREATE 0x01

InstType "Typical" IT_TYPICAL
InstType "Server" IT_SERVER
InstType "Full" IT_FULL

Section "Required Files" Section0
	SectionIn RO
	SectionInstType ${IT_FULL} ${IT_SERVER} ${IT_TYPICAL}
	DetailPrint "Necessary files like an uninstaller and some registry keys"
	System::Call 'SHELL32::SHPathPrepareForWrite(p $hwndParent, p 0, t d, i ${SHPPFW_DIRCREATE})'
	SetOutPath $InstDir
	File "${PWD}\icon.ico"
	WriteRegStr ${REGUINST} "DisplayName" "Uniview"
	WriteRegStr ${REGUINST} "UninstallString" "$INSTDIR\uninstall.exe /S"
	WriteRegStr ${REGUINST} "DisplayIcon" "$INSTDIR\icon.ico"
	WriteRegStr ${REGUINST} "Publisher" "Moritz Poldrack"
	WriteRegStr ${REGUINST} "HelpLink" "https://lists.sr.ht/~mpldr/uniview"
	WriteRegStr ${REGUINST} "Contact" "moritz@poldrack.dev"
	WriteUninstaller "$INSTDIR\uninstall.exe"
SectionEnd

Section "Licenses" Section1
	SectionIn RO
	SectionInstType ${IT_FULL} ${IT_SERVER} ${IT_TYPICAL}
	DetailPrint "The licenses this software is distributed under"
	SetOverwrite on
	SetOutPath "$INSTDIR\"
	File /r "${PWD}\src\LICENSES"
SectionEnd

Section /o "Sourcecode" Section2
	SectionInstType ${IT_FULL}
	DetailPrint "The sourcecode of this program"
	; Set Section properties
	SetOverwrite on

	; Set Section Files and Shortcuts
	SetOutPath "$INSTDIR\src\"
	File /r "${PWD}\src"
SectionEnd

Section "Client" Section3
	SectionInstType ${IT_FULL} ${IT_TYPICAL}
	DetailPrint "The Uniview client installation, including the uniview:// protocol"
	; Set Section properties
	SetOverwrite on

	; Set Section Files and Shortcuts
	SetOutPath "$INSTDIR\"
	File "${PWD}\uniview.exe"
	WriteRegStr HKCU "Software\Classes\uniview" "" "Uniview Client"
	WriteRegStr HKCU "Software\Classes\uniview" "URL Protocol" ""
	WriteRegStr HKCU "Software\Classes\uniview\DefaultIcon" "" "uniview.exe,1"
	WriteRegStr HKCU "Software\Classes\uniview\shell\open\command" "" '"$INSTDIR\uniview.exe" "%1"'
SectionEnd

Section "un.Client"
	DeleteRegKey HKCU "Software\Classes\uniview"
	Delete "$INSTDIR\uniview.exe"
SectionEnd

Section /o "Server" Section4
	SectionInstType ${IT_FULL} ${IT_SERVER}
	DetailPrint "The Uniview server installation. Only get this if you know what you're doing. And if you knew, you'd be using a Linux server."
	; Set Section properties
	SetOverwrite on

	; Set Section Files and Shortcuts
	SetOutPath "$INSTDIR\"
	File "${PWD}\univiewd.exe"
SectionEnd

Section "un.Server"
	Delete "$INSTDIR\univiewd.exe"
SectionEnd

Section /o "mpv (recommended)" SectionMPV
	SectionInstType ${IT_FULL} ${IT_TYPICAL}
	SetOutPath "$INSTDIR\player\mpv\"

	InitPluginsDir

	inetc::get https://github.com/mpvnet-player/mpv.net/releases/download/v6.0.4.0-stable/mpv.net-v6.0.4.0-stable.zip $Temp\mpv.zip
	Pop $0
	StrCmp $0 "OK" dlsucc

	MessageBox MB_ICONEXCLAMATION "Download failed: $R0"
	Quit

dlsucc:
	nsisunz::UnzipToLog "$Temp\mpv.zip" "$INSTDIR\player\mpv\"

	Pop $0
	StrCmp $0 "success" ok
	MessageBox MB_ICONEXCLAMATION "Extraction failed: $0"
ok:
SectionEnd

Section /o "VLC (not recommended)" SectionVLC
	SectionInstType ${IT_FULL}
	InitPluginsDir
	SetOutPath "$INSTDIR\player\"

	inetc::get https://portableapps.com/redir2/?a=VLCPortable&s=s&d=pa&f=VLCPortable_3.0.18.paf.exe $TEMP\vlc.exe
	Pop $0
	StrCmp $0 "OK" dlsucc

	MessageBox MB_ICONEXCLAMATION "Download failed: $R0"
	Quit
dlsucc:
SectionEnd

; Modern install component descriptions
!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
!insertmacro MUI_DESCRIPTION_TEXT ${Section1} ""
!insertmacro MUI_FUNCTION_DESCRIPTION_END

;Uninstall section
Section Uninstall
	;Remove from registry...
	DeleteRegKey ${REGUINST}

	; Delete self
	Delete "$INSTDIR\uninstall.exe"

	; Remove remaining directories
	RMDir /r "$INSTDIR\"
SectionEnd

BrandingText "Uniview – because watching alone just isn't the same"

; eof
