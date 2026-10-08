package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"reparto/internal/domain"
	"reparto/internal/repository"
	"reparto/internal/service"
)

// App actúa como el adaptador/controlador principal conectado al frontend con Wails.
type App struct {
	ctx     context.Context
	service *service.JornadaService
	repo    domain.JornadaRepository
}

// NewApp instancia el adaptador de la aplicación.
func NewApp() *App {
	return &App{}
}

// startup inicializa el repositorio SQLite y los servicios de dominio al arrancar la app.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	base, err := os.UserConfigDir()
	if err != nil {
		base = "."
	}
	appDir := filepath.Join(base, "RepartoAguaSoda")

	repo, err := repository.NewSQLiteJornadaRepository(appDir)
	if err != nil {
		fmt.Printf("Error inicializando repositorio SQLite: %v\n", err)
	}
	a.repo = repo

	calc := service.NewCalculadora()
	a.service = service.NewJornadaService(repo, calc)
}

// shutdown libera los recursos de la base de datos al cerrar la aplicación.
func (a *App) shutdown(ctx context.Context) {
	if a.repo != nil {
		_ = a.repo.Close()
	}
}

// Calcular replica las fórmulas de la planilla en tiempo real para la interfaz.
func (a *App) Calcular(j domain.Jornada) domain.Resultado {
	if a.service == nil {
		return domain.Resultado{}
	}
	return a.service.Calcular(j)
}

// Nueva inicializa una jornada vacía con la fecha especificada y los últimos precios usados.
func (a *App) Nueva(fecha string) domain.Jornada {
	if a.service == nil {
		return domain.Jornada{Fecha: fecha}
	}
	return a.service.Nueva(fecha)
}

// Guardar persiste la jornada en la base de datos SQLite y genera el respaldo.
func (a *App) Guardar(j domain.Jornada) error {
	if a.service == nil {
		return fmt.Errorf("servicio no inicializado")
	}
	return a.service.Guardar(j)
}

// Cargar devuelve la jornada guardada para una fecha o una nueva si no existe.
func (a *App) Cargar(fecha string) (domain.Jornada, error) {
	if a.service == nil {
		return domain.Jornada{}, fmt.Errorf("servicio no inicializado")
	}
	return a.service.Cargar(fecha)
}

// Historial devuelve las fechas guardadas para el selector de historial.
func (a *App) Historial() []string {
	if a.service == nil {
		return []string{}
	}
	fechas, _ := a.service.ListarFechas()
	return fechas
}

// FechasGuardadas devuelve las fechas con cierres para pintar el calendario.
func (a *App) FechasGuardadas() []string {
	return a.Historial()
}

// ObtenerEstadisticas calcula el resumen financiero diario, semanal y mensual.
func (a *App) ObtenerEstadisticas(fecha string) (domain.Estadisticas, error) {
	if a.service == nil {
		return domain.Estadisticas{}, fmt.Errorf("servicio no inicializado")
	}
	return a.service.ObtenerEstadisticas(fecha)
}
