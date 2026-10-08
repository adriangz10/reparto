export namespace domain {
	
	export class Billetes {
	    b20000: number;
	    b10000: number;
	    b2000: number;
	    b1000: number;
	    b500: number;
	
	    static createFrom(source: any = {}) {
	        return new Billetes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.b20000 = source["b20000"];
	        this.b10000 = source["b10000"];
	        this.b2000 = source["b2000"];
	        this.b1000 = source["b1000"];
	        this.b500 = source["b500"];
	    }
	}
	export class Cantidades {
	    bidon20: number;
	    bidon6: number;
	    soda: number;
	
	    static createFrom(source: any = {}) {
	        return new Cantidades(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bidon20 = source["bidon20"];
	        this.bidon6 = source["bidon6"];
	        this.soda = source["soda"];
	    }
	}
	export class PeriodoEstadistica {
	    etiqueta: string;
	    ganado: number;
	    gastos: number;
	    parteA: number;
	    parteB: number;
	    cantidad: number;
	
	    static createFrom(source: any = {}) {
	        return new PeriodoEstadistica(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.etiqueta = source["etiqueta"];
	        this.ganado = source["ganado"];
	        this.gastos = source["gastos"];
	        this.parteA = source["parteA"];
	        this.parteB = source["parteB"];
	        this.cantidad = source["cantidad"];
	    }
	}
	export class Estadisticas {
	    dia: PeriodoEstadistica;
	    semana: PeriodoEstadistica;
	    mes: PeriodoEstadistica;
	
	    static createFrom(source: any = {}) {
	        return new Estadisticas(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dia = this.convertValues(source["dia"], PeriodoEstadistica);
	        this.semana = this.convertValues(source["semana"], PeriodoEstadistica);
	        this.mes = this.convertValues(source["mes"], PeriodoEstadistica);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Precios {
	    bidon20: number;
	    bidon6: number;
	    soda: number;
	
	    static createFrom(source: any = {}) {
	        return new Precios(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bidon20 = source["bidon20"];
	        this.bidon6 = source["bidon6"];
	        this.soda = source["soda"];
	    }
	}
	export class Jornada {
	    fecha: string;
	    precios: Precios;
	    carga: Cantidades;
	    devuelto: Cantidades;
	    transferencias: number;
	    fiados: number;
	    gastos: number;
	    billetes: Billetes;
	
	    static createFrom(source: any = {}) {
	        return new Jornada(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fecha = source["fecha"];
	        this.precios = this.convertValues(source["precios"], Precios);
	        this.carga = this.convertValues(source["carga"], Cantidades);
	        this.devuelto = this.convertValues(source["devuelto"], Cantidades);
	        this.transferencias = source["transferencias"];
	        this.fiados = source["fiados"];
	        this.gastos = source["gastos"];
	        this.billetes = this.convertValues(source["billetes"], Billetes);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class Totales {
	    bidon20: number;
	    bidon6: number;
	    soda: number;
	
	    static createFrom(source: any = {}) {
	        return new Totales(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bidon20 = source["bidon20"];
	        this.bidon6 = source["bidon6"];
	        this.soda = source["soda"];
	    }
	}
	export class Resultado {
	    vendidos: Cantidades;
	    totales: Totales;
	    totalVendido: number;
	    efectivoTeorico: number;
	    parteA: number;
	    parteB: number;
	    subBilletes: number[];
	    efectivoReal: number;
	    diferencia: number;
	    estado: string;
	    cajonSoda: number;
	
	    static createFrom(source: any = {}) {
	        return new Resultado(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vendidos = this.convertValues(source["vendidos"], Cantidades);
	        this.totales = this.convertValues(source["totales"], Totales);
	        this.totalVendido = source["totalVendido"];
	        this.efectivoTeorico = source["efectivoTeorico"];
	        this.parteA = source["parteA"];
	        this.parteB = source["parteB"];
	        this.subBilletes = source["subBilletes"];
	        this.efectivoReal = source["efectivoReal"];
	        this.diferencia = source["diferencia"];
	        this.estado = source["estado"];
	        this.cajonSoda = source["cajonSoda"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

