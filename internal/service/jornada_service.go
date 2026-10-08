package service

import (
	"fmt"
	"time"

	"reparto/internal/domain"
)

var mesesES = []string{
	"", "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
	"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
}

// JornadaService coordina las operaciones y casos de uso del negocio.
type JornadaService struct {
	repo domain.JornadaRepository
	calc *Calculadora
}

func NewJornadaService(repo domain.JornadaRepository, calc *Calculadora) *JornadaService {
	return &JornadaService{
		repo: repo,
		calc: calc,
	}
}

func (s *JornadaService) DefaultPrecios() domain.Precios {
	return domain.Precios{Bidon20: 10000, Bidon6: 5000, Soda: 2000}
}

// Nueva inicializa una jornada con la fecha dada y los últimos precios registrados.
func (s *JornadaService) Nueva(fecha string) domain.Jornada {
	if fecha == "" {
		fecha = time.Now().Format("2006-01-02")
	}

	p := s.DefaultPrecios()
	if s.repo != nil {
		if ultimos, err := s.repo.ObtenerUltimosPrecios(); err == nil && ultimos != nil {
			if ultimos.Bidon20 > 0 || ultimos.Bidon6 > 0 || ultimos.Soda > 0 {
				p = *ultimos
			}
		}
	}

	return domain.Jornada{
		Fecha:   fecha,
		Precios: p,
	}
}

// Calcular delega a la calculadora para obtener los balances.
func (s *JornadaService) Calcular(j domain.Jornada) domain.Resultado {
	return s.calc.Calcular(j)
}

// Guardar valida, calcula y persiste la jornada.
func (s *JornadaService) Guardar(j domain.Jornada) error {
	r := s.calc.Calcular(j)
	return s.repo.Guardar(j, r)
}

// Cargar recupera una jornada por fecha, o devuelve una nueva si no existe.
func (s *JornadaService) Cargar(fecha string) (domain.Jornada, error) {
	if s.repo == nil {
		return s.Nueva(fecha), nil
	}

	guardada, err := s.repo.Cargar(fecha)
	if err != nil {
		return domain.Jornada{}, err
	}
	if guardada != nil {
		return *guardada, nil
	}

	return s.Nueva(fecha), nil
}

// ListarFechas devuelve el historial de fechas guardadas.
func (s *JornadaService) ListarFechas() ([]string, error) {
	if s.repo == nil {
		return []string{}, nil
	}
	return s.repo.ListarFechas()
}

// ObtenerEstadisticas genera el resumen de cierre diario, semanal y mensual.
func (s *JornadaService) ObtenerEstadisticas(fecha string) (domain.Estadisticas, error) {
	var res domain.Estadisticas

	if fecha == "" {
		fecha = time.Now().Format("2006-01-02")
	}
	t, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		t = time.Now()
		fecha = t.Format("2006-01-02")
	}

	if s.repo == nil {
		return res, fmt.Errorf("repositorio no inicializado")
	}

	// 1. Día
	etiquetaDia := t.Format("02/01/2006")
	diaStats, err := s.repo.ConsultarEstadisticas(fecha, fecha, etiquetaDia)
	if err != nil {
		return res, err
	}
	res.Dia = diaStats

	// 2. Semana (Lunes a Domingo)
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	lunes := t.AddDate(0, 0, -(weekday - 1))
	domingo := lunes.AddDate(0, 0, 6)
	desdeSem := lunes.Format("2006-01-02")
	hastaSem := domingo.Format("2006-01-02")
	etiquetaSem := fmt.Sprintf("%s al %s", lunes.Format("02/01"), domingo.Format("02/01/2006"))
	semStats, err := s.repo.ConsultarEstadisticas(desdeSem, hastaSem, etiquetaSem)
	if err != nil {
		return res, err
	}
	res.Semana = semStats

	// 3. Mes (1er día al último día)
	desdeMes := fmt.Sprintf("%04d-%02d-01", t.Year(), t.Month())
	primerDiaSigMes := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	ultimoDiaMes := primerDiaSigMes.AddDate(0, 0, -1)
	hastaMes := ultimoDiaMes.Format("2006-01-02")
	nombreMes := mesesES[t.Month()]
	etiquetaMes := fmt.Sprintf("%s %d", nombreMes, t.Year())
	mesStats, err := s.repo.ConsultarEstadisticas(desdeMes, hastaMes, etiquetaMes)
	if err != nil {
		return res, err
	}
	res.Mes = mesStats

	return res, nil
}
