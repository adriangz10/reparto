package domain

// Precios unitarios de cada producto.
type Precios struct {
	Bidon20 float64 `json:"bidon20"`
	Bidon6  float64 `json:"bidon6"`
	Soda    float64 `json:"soda"`
}

// Cantidades de unidades para carga, devuelto o vendidos.
type Cantidades struct {
	Bidon20 int `json:"bidon20"`
	Bidon6  int `json:"bidon6"`
	Soda    int `json:"soda"`
}

// Totales monetarios por producto.
type Totales struct {
	Bidon20 float64 `json:"bidon20"`
	Bidon6  float64 `json:"bidon6"`
	Soda    float64 `json:"soda"`
}

// Billetes cuenta el arqueo de billetes físicos.
type Billetes struct {
	B20000 int `json:"b20000"`
	B10000 int `json:"b10000"`
	B2000  int `json:"b2000"`
	B1000  int `json:"b1000"`
	B500   int `json:"b500"`
}

// Jornada representa la planilla diaria completa cargada por el usuario.
type Jornada struct {
	Fecha          string     `json:"fecha"` // Formato YYYY-MM-DD
	Precios        Precios    `json:"precios"`
	Carga          Cantidades `json:"carga"`
	Devuelto       Cantidades `json:"devuelto"`
	Transferencias float64    `json:"transferencias"`
	Fiados         float64    `json:"fiados"`
	Gastos         float64    `json:"gastos"`
	Billetes       Billetes   `json:"billetes"`
}

// Resultado representa el arqueo y balance calculado de una jornada.
type Resultado struct {
	Vendidos        Cantidades `json:"vendidos"`
	Totales         Totales    `json:"totales"`
	TotalVendido    float64    `json:"totalVendido"`
	EfectivoTeorico float64    `json:"efectivoTeorico"`
	ParteA          float64    `json:"parteA"`
	ParteB          float64    `json:"parteB"`
	SubBilletes     []float64  `json:"subBilletes"` // Subtotales de 20000, 10000, 2000, 1000, 500
	EfectivoReal    float64    `json:"efectivoReal"`
	Diferencia      float64    `json:"diferencia"`
	Estado          string     `json:"estado"` // "exacta" | "falta" | "sobra"
	CajonSoda       float64    `json:"cajonSoda"`
}

// PeriodoEstadistica resume los datos agregados de un periodo (día, semana o mes).
type PeriodoEstadistica struct {
	Etiqueta string  `json:"etiqueta"`
	Ganado   float64 `json:"ganado"`
	Gastos   float64 `json:"gastos"`
	ParteA   float64 `json:"parteA"`
	ParteB   float64 `json:"parteB"`
	Cantidad int     `json:"cantidad"`
}

// Estadisticas contiene el resumen de día, semana y mes.
type Estadisticas struct {
	Dia    PeriodoEstadistica `json:"dia"`
	Semana PeriodoEstadistica `json:"semana"`
	Mes    PeriodoEstadistica `json:"mes"`
}
