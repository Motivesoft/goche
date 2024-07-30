@echo off
setlocal enabledelayedexpansion

set NAME=goche

for /F "tokens=*" %%g in ('type version') do (
  set VERSION=%%g
)

set PLATFORM_WINDOWS=windows
set PLATFORM_LINUX=linux

rem set PLATFORMS=%PLATFORM_WINDOWS% %PLATFORM_LINUX%
set PLATFORMS=%PLATFORM_WINDOWS%
rem set PLATFORMS=%PLATFORM_LINUX%

for %%p in (%PLATFORMS%) do (
    if %%p == %PLATFORM_WINDOWS% (
        set EXTENSION=.exe
    ) else (
        set EXTENSION=
    )

    set GOOS=%%p
    set GOARCH=amd64
    
    go build -ldflags "-X goche/identification.engineName=%NAME% -X goche/identification.versionName=%VERSION%" -o %NAME%-!GOOS!-!GOARCH!-%VERSION%!EXTENSION! .
)