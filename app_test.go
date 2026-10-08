package main

import (
	"os"
	"path/filepath"
	"testing"

	"reparto/internal/domain"
	"reparto/internal/repository"
	"reparto/internal/service"
)

func TestAppArchitectureAndFlow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "reparto-arch-test-*")
	if err != nil {
		t.Fatalf("Error creando tmpDir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 1. Instanciar capas
	repo, err := repository.NewSQLiteJornadaRepository(tmpDir)
	if err != nil {
		t.Fatalf("Error creando repositorio: %v", err)
	}
	defer repo.Close()

	calc := service.NewCalculadora()
	svc := service.NewJornadaService(repo, calc)

	app := NewApp()
	app.repo = repo
	app.service = svc

	// 2. Probar Nueva()
	nueva := app.Nueva("2026-10-07")
	if nueva.Fecha != "2026-10-07" {
		t.Errorf("Fecha esperada 2026-10-07, obtenida %s", nueva.Fecha)
	}
	if nueva.Precios.Bidon20 != 10000 {
		t.Errorf("Precio por defecto esperado 10000, obtenido %v", nueva.Precios.Bidon20)
	}

	// 3. Probar Guardar()
	j1 := domain.Jornada{
		Fecha: "2026-10-07",
		Precios: domain.Precios{
			Bidon20: 10000,
			Bidon6:  5000,
			Soda:    2000,
		},
		Carga: domain.Cantidades{
			Bidon20: 20,
			Bidon6:  10,
			Soda:    30,
		},
		Devuelto: domain.Cantidades{
			Bidon20: 5,
			Bidon6:  2,
			Soda:    6,
		},
		Transferencias: 10000,
		Fiados:         5000,
		Gastos:         20000,
		Billetes: domain.Billetes{
			B20000: 10,
			B10000: 5,
		},
	}

	if err := app.Guardar(j1); err != nil {
		t.Fatalf("Error guardando j1: %v", err)
	}

	// 4. Probar Cargar()
	cargada, err := app.Cargar("2026-10-07")
	if err != nil {
		t.Fatalf("Error cargando j1: %v", err)
	}
	if cargada.Fecha != "2026-10-07" || cargada.Carga.Bidon20 != 20 || cargada.Gastos != 20000 {
		t.Fatalf("Datos cargados incorrectos: %+v", cargada)
	}

	// 5. Probar FechasGuardadas()
	fechas := app.FechasGuardadas()
	if len(fechas) != 1 || fechas[0] != "2026-10-07" {
		t.Fatalf("Fechas guardadas esperadas ['2026-10-07'], obtenido: %v", fechas)
	}

	// 6. Probar Estadísticas
	stats, err := app.ObtenerEstadisticas("2026-10-07")
	if err != nil {
		t.Fatalf("Error obteniendo estadísticas: %v", err)
	}

	// Vendidos:
	// Bidon20: 15 * 10000 = 150000
	// Bidon6: 8 * 5000 = 40000
	// Soda: 24 * 2000 = 48000
	// Total Vendido = 238000
	// Gastos = 20000
	if stats.Dia.Ganado != 238000 {
		t.Errorf("stats.Dia.Ganado esperado 238000, obtenido %v", stats.Dia.Ganado)
	}
	if stats.Dia.Gastos != 20000 {
		t.Errorf("stats.Dia.Gastos esperado 20000, obtenido %v", stats.Dia.Gastos)
	}
}

func TestCalculadoraPureDomain(t *testing.T) {
	calc := service.NewCalculadora()

	j := domain.Jornada{
		Fecha: "2026-10-07",
		Precios: domain.Precios{
			Bidon20: 10000,
			Bidon6:  5000,
			Soda:    2000,
		},
		Carga: domain.Cantidades{
			Bidon20: 10,
			Bidon6:  5,
			Soda:    12,
		},
		Devuelto: domain.Cantidades{
			Bidon20: 2,
			Bidon6:  1,
			Soda:    2,
		},
		Transferencias: 10000,
		Fiados:         5000,
		Gastos:         15000,
		Billetes: domain.Billetes{
			B20000: 4,
			B10000: 1,
		},
	}

	res := calc.Calcular(j)

	// Vendidos: 8 * 10000 + 4 * 5000 + 10 * 2000 = 80000 + 20000 + 20000 = 120000
	if res.TotalVendido != 120000 {
		t.Errorf("TotalVendido esperado 120000, obtenido %v", res.TotalVendido)
	}
	// Teorico: 120000 - 10000 - 5000 - 15000 = 90000
	if res.EfectivoTeorico != 90000 {
		t.Errorf("EfectivoTeorico esperado 90000, obtenido %v", res.EfectivoTeorico)
	}
	// Partes: 45000 cada una
	if res.ParteA != 45000 || res.ParteB != 45000 {
		t.Errorf("Partes A y B esperadas 45000, obtenido A=%v, B=%v", res.ParteA, res.ParteB)
	}
	// Efectivo real: 4*20000 + 1*10000 = 90000
	if res.EfectivoReal != 90000 {
		t.Errorf("EfectivoReal esperado 90000, obtenido %v", res.EfectivoReal)
	}
	if res.Estado != "exacta" {
		t.Errorf("Estado esperado exacta, obtenido %s", res.Estado)
	}
}

func TestMigrationInRepository(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "reparto-repo-migration-*")
	if err != nil {
		t.Fatalf("Error creando tmpDir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	jornadasDir := filepath.Join(tmpDir, "jornadas")
	_ = os.MkdirAll(jornadasDir, 0o755)

	jsonSample := `{
		"fecha": "2026-09-15",
		"precios": {"bidon20": 10000, "bidon6": 5000, "soda": 2000},
		"carga": {"bidon20": 10, "bidon6": 0, "soda": 0},
		"devuelto": {"bidon20": 0, "bidon6": 0, "soda": 0},
		"transferencias": 0,
		"fiados": 0,
		"gastos": 5000,
		"billetes": {"b20000": 5, "b10000": 0, "b2000": 0, "b1000": 0, "b500": 0}
	}`
	err = os.WriteFile(filepath.Join(jornadasDir, "2026-09-15.json"), []byte(jsonSample), 0o644)
	if err != nil {
		t.Fatalf("Error escribiendo JSON: %v", err)
	}

	repo, err := repository.NewSQLiteJornadaRepository(tmpDir)
	if err != nil {
		t.Fatalf("Error creando repo: %v", err)
	}
	defer repo.Close()

	cargada, err := repo.Cargar("2026-09-15")
	if err != nil {
		t.Fatalf("Error cargando jornada: %v", err)
	}
	if cargada == nil || cargada.Fecha != "2026-09-15" || cargada.Carga.Bidon20 != 10 {
		t.Fatalf("Datos migrados incorrectos: %+v", cargada)
	}
}
