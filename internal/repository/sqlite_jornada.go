package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"

	"reparto/internal/domain"
)

type SQLiteJornadaRepository struct {
	db      *sql.DB
	baseDir string
}

func NewSQLiteJornadaRepository(baseDir string) (*SQLiteJornadaRepository, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("error creando directorio base: %w", err)
	}
	_ = os.MkdirAll(filepath.Join(baseDir, "jornadas"), 0o755)

	dbPath := filepath.Join(baseDir, "reparto.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error abriendo base de datos SQLite: %w", err)
	}

	// SQLite funciona mejor con 1 conexión concurrente para escrituras
	db.SetMaxOpenConns(1)

	repo := &SQLiteJornadaRepository{
		db:      db,
		baseDir: baseDir,
	}

	if err := repo.initDB(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("error inicializando esquema de base de datos: %w", err)
	}

	repo.migrarJSON()

	return repo, nil
}

func (r *SQLiteJornadaRepository) initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS jornadas (
		fecha TEXT PRIMARY KEY,
		precio_bidon20 REAL NOT NULL DEFAULT 0,
		precio_bidon6 REAL NOT NULL DEFAULT 0,
		precio_soda REAL NOT NULL DEFAULT 0,
		carga_bidon20 INTEGER NOT NULL DEFAULT 0,
		carga_bidon6 INTEGER NOT NULL DEFAULT 0,
		carga_soda INTEGER NOT NULL DEFAULT 0,
		devuelto_bidon20 INTEGER NOT NULL DEFAULT 0,
		devuelto_bidon6 INTEGER NOT NULL DEFAULT 0,
		devuelto_soda INTEGER NOT NULL DEFAULT 0,
		transferencias REAL NOT NULL DEFAULT 0,
		fiados REAL NOT NULL DEFAULT 0,
		gastos REAL NOT NULL DEFAULT 0,
		b20000 INTEGER NOT NULL DEFAULT 0,
		b10000 INTEGER NOT NULL DEFAULT 0,
		b2000 INTEGER NOT NULL DEFAULT 0,
		b1000 INTEGER NOT NULL DEFAULT 0,
		b500 INTEGER NOT NULL DEFAULT 0,
		total_vendido REAL NOT NULL DEFAULT 0,
		efectivo_teorico REAL NOT NULL DEFAULT 0,
		efectivo_real REAL NOT NULL DEFAULT 0,
		parte_a REAL NOT NULL DEFAULT 0,
		parte_b REAL NOT NULL DEFAULT 0,
		diferencia REAL NOT NULL DEFAULT 0,
		estado TEXT NOT NULL DEFAULT '',
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_jornadas_fecha ON jornadas(fecha);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteJornadaRepository) migrarJSON() {
	files, _ := filepath.Glob(filepath.Join(r.baseDir, "jornadas", "*.json"))
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		var j domain.Jornada
		if err := json.Unmarshal(data, &j); err != nil {
			continue
		}
		var count int
		_ = r.db.QueryRow("SELECT COUNT(*) FROM jornadas WHERE fecha = ?", j.Fecha).Scan(&count)
		if count == 0 {
			// Calcular mínimos si faltan
			totales := float64(j.Carga.Bidon20-j.Devuelto.Bidon20)*j.Precios.Bidon20 +
				float64(j.Carga.Bidon6-j.Devuelto.Bidon6)*j.Precios.Bidon6 +
				float64(j.Carga.Soda-j.Devuelto.Soda)*j.Precios.Soda
			teorico := totales - j.Transferencias - j.Fiados - j.Gastos
			res := domain.Resultado{
				TotalVendido:    totales,
				EfectivoTeorico: teorico,
				ParteA:          teorico / 2,
				ParteB:          teorico / 2,
			}
			_ = r.Guardar(j, res)
		}
	}
}

func (r *SQLiteJornadaRepository) Guardar(j domain.Jornada, res domain.Resultado) error {
	query := `
	INSERT INTO jornadas (
		fecha,
		precio_bidon20, precio_bidon6, precio_soda,
		carga_bidon20, carga_bidon6, carga_soda,
		devuelto_bidon20, devuelto_bidon6, devuelto_soda,
		transferencias, fiados, gastos,
		b20000, b10000, b2000, b1000, b500,
		total_vendido, efectivo_teorico, efectivo_real,
		parte_a, parte_b, diferencia, estado,
		updated_at
	) VALUES (
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP
	) ON CONFLICT(fecha) DO UPDATE SET
		precio_bidon20 = excluded.precio_bidon20,
		precio_bidon6 = excluded.precio_bidon6,
		precio_soda = excluded.precio_soda,
		carga_bidon20 = excluded.carga_bidon20,
		carga_bidon6 = excluded.carga_bidon6,
		carga_soda = excluded.carga_soda,
		devuelto_bidon20 = excluded.devuelto_bidon20,
		devuelto_bidon6 = excluded.devuelto_bidon6,
		devuelto_soda = excluded.devuelto_soda,
		transferencias = excluded.transferencias,
		fiados = excluded.fiados,
		gastos = excluded.gastos,
		b20000 = excluded.b20000,
		b10000 = excluded.b10000,
		b2000 = excluded.b2000,
		b1000 = excluded.b1000,
		b500 = excluded.b500,
		total_vendido = excluded.total_vendido,
		efectivo_teorico = excluded.efectivo_teorico,
		efectivo_real = excluded.efectivo_real,
		parte_a = excluded.parte_a,
		parte_b = excluded.parte_b,
		diferencia = excluded.diferencia,
		estado = excluded.estado,
		updated_at = CURRENT_TIMESTAMP;
	`
	_, err := r.db.Exec(query,
		j.Fecha,
		j.Precios.Bidon20, j.Precios.Bidon6, j.Precios.Soda,
		j.Carga.Bidon20, j.Carga.Bidon6, j.Carga.Soda,
		j.Devuelto.Bidon20, j.Devuelto.Bidon6, j.Devuelto.Soda,
		j.Transferencias, j.Fiados, j.Gastos,
		j.Billetes.B20000, j.Billetes.B10000, j.Billetes.B2000, j.Billetes.B1000, j.Billetes.B500,
		res.TotalVendido, res.EfectivoTeorico, res.EfectivoReal,
		res.ParteA, res.ParteB, res.Diferencia, res.Estado,
	)
	if err != nil {
		return err
	}

	// Backup en archivo JSON para redundancia
	if len(j.Fecha) == 10 && !strings.ContainsAny(j.Fecha, `/\.`) {
		ruta := filepath.Join(r.baseDir, "jornadas", j.Fecha+".json")
		if data, err := json.MarshalIndent(j, "", "  "); err == nil {
			_ = os.WriteFile(ruta, data, 0o644)
		}
	}
	pd, _ := json.Marshal(j.Precios)
	_ = os.WriteFile(filepath.Join(r.baseDir, "precios.json"), pd, 0o644)

	return nil
}

func (r *SQLiteJornadaRepository) Cargar(fecha string) (*domain.Jornada, error) {
	var j domain.Jornada
	query := `
	SELECT fecha,
	       precio_bidon20, precio_bidon6, precio_soda,
	       carga_bidon20, carga_bidon6, carga_soda,
	       devuelto_bidon20, devuelto_bidon6, devuelto_soda,
	       transferencias, fiados, gastos,
	       b20000, b10000, b2000, b1000, b500
	FROM jornadas WHERE fecha = ?`

	err := r.db.QueryRow(query, fecha).Scan(
		&j.Fecha,
		&j.Precios.Bidon20, &j.Precios.Bidon6, &j.Precios.Soda,
		&j.Carga.Bidon20, &j.Carga.Bidon6, &j.Carga.Soda,
		&j.Devuelto.Bidon20, &j.Devuelto.Bidon6, &j.Devuelto.Soda,
		&j.Transferencias, &j.Fiados, &j.Gastos,
		&j.Billetes.B20000, &j.Billetes.B10000, &j.Billetes.B2000, &j.Billetes.B1000, &j.Billetes.B500,
	)
	if err == nil {
		return &j, nil
	}
	if err == sql.ErrNoRows {
		// Fallback por si existe en archivo JSON previo
		ruta := filepath.Join(r.baseDir, "jornadas", fecha+".json")
		if data, err := os.ReadFile(ruta); err == nil {
			var fj domain.Jornada
			if err := json.Unmarshal(data, &fj); err == nil {
				return &fj, nil
			}
		}
		return nil, nil
	}
	return nil, err
}

func (r *SQLiteJornadaRepository) ListarFechas() ([]string, error) {
	rows, err := r.db.Query("SELECT fecha FROM jornadas ORDER BY fecha DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fechas []string
	for rows.Next() {
		var f string
		if err := rows.Scan(&f); err == nil {
			fechas = append(fechas, f)
		}
	}

	if len(fechas) > 0 {
		return fechas, nil
	}

	// Fallback por si hay archivos
	files, _ := filepath.Glob(filepath.Join(r.baseDir, "jornadas", "*.json"))
	for _, f := range files {
		fechas = append(fechas, strings.TrimSuffix(filepath.Base(f), ".json"))
	}
	sort.Sort(sort.Reverse(sort.StringSlice(fechas)))
	return fechas, nil
}

func (r *SQLiteJornadaRepository) ObtenerUltimosPrecios() (*domain.Precios, error) {
	var p domain.Precios
	err := r.db.QueryRow("SELECT precio_bidon20, precio_bidon6, precio_soda FROM jornadas ORDER BY fecha DESC LIMIT 1").Scan(
		&p.Bidon20, &p.Bidon6, &p.Soda,
	)
	if err == nil {
		return &p, nil
	}

	// Fallback precios.json
	if data, err := os.ReadFile(filepath.Join(r.baseDir, "precios.json")); err == nil {
		if err := json.Unmarshal(data, &p); err == nil {
			return &p, nil
		}
	}
	return nil, nil
}

func (r *SQLiteJornadaRepository) ConsultarEstadisticas(desde, hasta, etiqueta string) (domain.PeriodoEstadistica, error) {
	var p domain.PeriodoEstadistica
	p.Etiqueta = etiqueta

	query := `
	SELECT COALESCE(SUM(total_vendido), 0),
	       COALESCE(SUM(gastos), 0),
	       COALESCE(SUM(parte_a), 0),
	       COALESCE(SUM(parte_b), 0),
	       COUNT(*)
	FROM jornadas
	WHERE fecha >= ? AND fecha <= ?`

	err := r.db.QueryRow(query, desde, hasta).Scan(
		&p.Ganado,
		&p.Gastos,
		&p.ParteA,
		&p.ParteB,
		&p.Cantidad,
	)
	return p, err
}

func (r *SQLiteJornadaRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
