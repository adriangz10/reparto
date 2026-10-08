@echo off
echo ========================================================
echo   Compilando Reparto para Produccion (Wails + SQLite)
echo ========================================================

wails build -clean -s

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================================
    echo   Compilacion exitosa!
    echo   Ejecutable generado en: build\bin\reparto.exe
    echo ========================================================
) else (
    echo.
    echo Error durante la compilacion.
)

if "%~1"=="--nopause" goto end
if "%CI%"=="true" goto end
pause
:end
