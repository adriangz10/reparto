package service

import (
	"math"
	"reparto/internal/domain"
)

// Calculadora encapsula las reglas de negocio y fórmulas de la planilla.
type Calculadora struct{}

func NewCalculadora() *Calculadora {
	return &Calculadora{}
}

// Calcular procesa todas las fórmulas de ventas, rendición y arqueo.
func (c *Calculadora) Calcular(j domain.Jornada) domain.Resultado {
	var r domain.Resultado

	r.Vendidos = domain.Cantidades{
		Bidon20: j.Carga.Bidon20 - j.Devuelto.Bidon20,
		Bidon6:  j.Carga.Bidon6 - j.Devuelto.Bidon6,
		Soda:    j.Carga.Soda - j.Devuelto.Soda,
	}

	r.Totales = domain.Totales{
		Bidon20: float64(r.Vendidos.Bidon20) * j.Precios.Bidon20,
		Bidon6:  float64(r.Vendidos.Bidon6) * j.Precios.Bidon6,
		Soda:    float64(r.Vendidos.Soda) * j.Precios.Soda,
	}

	r.TotalVendido = r.Totales.Bidon20 + r.Totales.Bidon6 + r.Totales.Soda
	r.CajonSoda = j.Precios.Soda * 6

	r.EfectivoTeorico = r.TotalVendido - j.Transferencias - j.Fiados - j.Gastos
	r.ParteA = r.EfectivoTeorico / 2
	r.ParteB = r.EfectivoTeorico / 2

	b := j.Billetes
	r.SubBilletes = []float64{
		float64(b.B20000) * 20000,
		float64(b.B10000) * 10000,
		float64(b.B2000) * 2000,
		float64(b.B1000) * 1000,
		float64(b.B500) * 500,
	}

	for _, s := range r.SubBilletes {
		r.EfectivoReal += s
	}

	r.Diferencia = r.EfectivoReal - r.EfectivoTeorico

	switch {
	case math.Abs(r.Diferencia) < 0.01:
		r.Diferencia = 0
		r.Estado = "exacta"
	case r.Diferencia < 0:
		r.Estado = "falta"
	default:
		r.Estado = "sobra"
	}

	return r
}
