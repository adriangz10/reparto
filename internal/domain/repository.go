package domain

// JornadaRepository define el contrato de persistencia para las jornadas.
type JornadaRepository interface {
	// Guardar persiste o actualiza una jornada y su resultado.
	Guardar(j Jornada, r Resultado) error

	// Cargar busca una jornada guardada por fecha. Retorna nil, nil si no existe.
	Cargar(fecha string) (*Jornada, error)

	// ListarFechas devuelve todas las fechas con cierres guardados (más recientes primero).
	ListarFechas() ([]string, error)

	// ObtenerUltimosPrecios recupera los precios de la última jornada guardada.
	ObtenerUltimosPrecios() (*Precios, error)

	// ConsultarEstadisticas calcula el resumen financiero entre dos fechas inclusivas.
	ConsultarEstadisticas(desde, hasta, etiqueta string) (PeriodoEstadistica, error)

	// Close libera recursos de la base de datos si corresponde.
	Close() error
}
