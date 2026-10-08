# Reparto — Carga, Ventas y Arqueo de Caja (Go + Wails v2)

Requisitos: Go 1.21+, Wails v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`), `wails doctor`.

    cd reparto
    go mod tidy
    wails dev        # desarrollo
    wails build      # genera el .exe en build/bin

Datos: JSON por día en `%AppData%/RepartoAguaSoda/jornadas/` (Windows).
