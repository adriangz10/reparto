const $ = (id) => document.getElementById(id);
const PRODS = [["bidon20", "Bidón 20L"], ["bidon6", "Bidón 6L"], ["soda", "Soda (unidad)"]];
const DENS = [20000, 10000, 2000, 1000, 500];
const MESES = [
  "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
  "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"
];

const peso = (n) => "$" + (n || 0).toLocaleString("es-AR", { maximumFractionDigits: 2 });
const num = (id) => parseFloat($(id).value) || 0;
const int = (id) => parseInt($(id).value, 10) || 0;
const api = () => window.go.main.App;

// Estado de la aplicación
let fechasGuardadas = [];
let currentCalDate = new Date();
let cierreGuardado = false;

function construirUI() {
  $("precios").innerHTML = PRODS.map(([k, n]) =>
    `<label>${n}<input type="number" min="0" step="any" id="precio_${k}" placeholder="0.00"></label>`).join("");

  $("tablaCarga").querySelector("tbody").innerHTML = PRODS.map(([k, n]) => `
    <tr><td>${n}</td>
    <td><input type="number" min="0" id="carga_${k}" placeholder="0"></td>
    <td><input type="number" min="0" id="dev_${k}" placeholder="0"></td>
    <td id="v_${k}">0</td><td id="t_${k}">$0</td></tr>`).join("");

  $("tablaBilletes").querySelector("tbody").innerHTML = DENS.map((d) => `
    <tr><td>${peso(d)}</td><td><input type="number" min="0" id="b_${d}" placeholder="0"></td><td id="bs_${d}">$0</td></tr>`).join("");

  document.querySelectorAll("main input").forEach((el) => {
    el.addEventListener("input", () => {
      el.classList.remove("input-error");
      cierreGuardado = false;
      actualizarEstadoGuardado();
      refrescar();
    });

    // Al seleccionar / hacer foco: si el valor es 0, se borra automáticamente
    el.addEventListener("focus", function () {
      if (this.value === "0" || this.value === "0.00" || this.value === "0,00") {
        this.value = "";
      } else {
        this.select();
      }
    });

    el.addEventListener("mouseup", function (e) {
      if (this.value === "0" || this.value === "0.00" || this.value === "0,00") {
        this.value = "";
        e.preventDefault();
      }
    });

    // Al salir del campo: si quedó vacío y no es un precio, restaurar 0
    el.addEventListener("blur", function () {
      if (this.value.trim() === "") {
        if (!this.id.startsWith("precio_")) {
          this.value = "0";
        }
        refrescar();
      }
    });
  });
}

function leer() {
  const por = (pref, fn) => Object.fromEntries(PRODS.map(([k]) => [k, fn(`${pref}_${k}`)]));
  return {
    fecha: $("fecha").value,
    precios: por("precio", num),
    carga: por("carga", int),
    devuelto: por("dev", int),
    transferencias: num("transferencias"),
    fiados: num("fiados"),
    gastos: num("gastos"),
    billetes: Object.fromEntries(DENS.map((d) => ["b" + d, int("b_" + d)])),
  };
}

function formatearFechaVisual(iso) {
  if (!iso || !iso.includes("-")) return "--/--/----";
  const [y, m, d] = iso.split("-");
  return `${d}/${m}/${y}`;
}

function volcar(j, esCierreGuardado = false) {
  $("fecha").value = j.fecha;
  $("fechaTexto").textContent = formatearFechaVisual(j.fecha);

  PRODS.forEach(([k]) => {
    $("precio_" + k).value = j.precios[k] !== undefined && j.precios[k] !== null ? j.precios[k] : "";
    $("carga_" + k).value = j.carga && j.carga[k] !== undefined && j.carga[k] !== 0 ? j.carga[k] : (j.carga && j.carga[k] === 0 ? "0" : "");
    $("dev_" + k).value = j.devuelto && j.devuelto[k] !== undefined && j.devuelto[k] !== 0 ? j.devuelto[k] : (j.devuelto && j.devuelto[k] === 0 ? "0" : "");
  });

  $("transferencias").value = j.transferencias !== undefined && j.transferencias !== null && j.transferencias !== 0 ? j.transferencias : (j.transferencias === 0 ? "0" : "");
  $("fiados").value = j.fiados !== undefined && j.fiados !== null && j.fiados !== 0 ? j.fiados : (j.fiados === 0 ? "0" : "");
  $("gastos").value = j.gastos !== undefined && j.gastos !== null && j.gastos !== 0 ? j.gastos : (j.gastos === 0 ? "0" : "");

  DENS.forEach((d) => {
    const val = j.billetes && j.billetes["b" + d];
    $("b_" + d).value = val !== undefined && val !== null && val !== 0 ? val : (val === 0 ? "0" : "");
  });

  // Limpiar marcas de error
  document.querySelectorAll(".input-error").forEach((el) => el.classList.remove("input-error"));

  cierreGuardado = esCierreGuardado;
  actualizarEstadoGuardado();
  refrescar();
  cargarEstadisticas(j.fecha);
}

function actualizarEstadoGuardado() {
  const btn = $("btnGuardar");
  if (cierreGuardado) {
    btn.classList.add("btn-guardado");
    btn.textContent = "✔ Guardado";
  } else {
    btn.classList.remove("btn-guardado");
    btn.textContent = "Guardar";
  }
}

async function refrescar() {
  const j = leer();
  const r = await api().Calcular(j);
  PRODS.forEach(([k]) => {
    $("v_" + k).textContent = r.vendidos[k];
    $("t_" + k).textContent = peso(r.totales[k]);
  });
  $("cajon").textContent = peso(r.cajonSoda);
  $("totalVendido").textContent = peso(r.totalVendido);
  $("teorico").textContent = peso(r.efectivoTeorico);
  $("parteA").textContent = peso(r.parteA);
  $("parteB").textContent = peso(r.parteB);
  DENS.forEach((d, i) => ($("bs_" + d).textContent = peso(r.subBilletes[i])));
  $("real").textContent = peso(r.efectivoReal);

  const el = $("diferencia");
  el.className = "diff " + r.estado;
  el.textContent =
    r.estado === "exacta" ? "✔ Caja exacta" :
    r.estado === "falta" ? `✖ Faltan ${peso(Math.abs(r.diferencia))}` :
    `▲ Sobran ${peso(r.diferencia)}`;

  armarResumen(j, r);
}

function armarResumen(j, r) {
  const filas = PRODS.map(([k, n]) =>
    `<tr><td>${n}</td><td>${j.carga[k] || 0}</td><td>${j.devuelto[k] || 0}</td><td>${r.vendidos[k] || 0}</td><td>${peso(r.totales[k])}</td></tr>`).join("");
  const bill = DENS.map((d, i) =>
    `<tr><td>${peso(d)}</td><td>${j.billetes["b" + d] || 0}</td><td>${peso(r.subBilletes[i])}</td></tr>`).join("");

  $("resumen").innerHTML = `
    <h2>Cierre de reparto — ${formatearFechaVisual(j.fecha)}</h2>
    <table><thead><tr><th>Producto</th><th>Carga</th><th>Devuelto</th><th>Vendido</th><th>Total</th></tr></thead>
      <tbody>${filas}</tbody><tfoot><tr><td colspan="4">Total vendido</td><td>${peso(r.totalVendido)}</td></tr></tfoot></table>
    <table><tbody>
      <tr><td>Transferencias</td><td>- ${peso(j.transferencias)}</td></tr>
      <tr><td>Fiados</td><td>- ${peso(j.fiados)}</td></tr>
      <tr><td>Gastos</td><td>- ${peso(j.gastos)}</td></tr>
      <tr><td><b>Efectivo teórico</b></td><td><b>${peso(r.efectivoTeorico)}</b></td></tr>
      <tr><td>Parte A</td><td>${peso(r.parteA)}</td></tr>
      <tr><td>Parte B</td><td>${peso(r.parteB)}</td></tr></tbody></table>
    <table><thead><tr><th>Billete</th><th>Cantidad</th><th>Subtotal</th></tr></thead><tbody>${bill}</tbody>
      <tfoot><tr><td colspan="2">Efectivo real contado</td><td>${peso(r.efectivoReal)}</td></tr></tfoot></table>
    <p><b>${r.estado === "exacta" ? "Caja exacta" : r.estado === "falta" ? "Faltan " + peso(Math.abs(r.diferencia)) : "Sobran " + peso(r.diferencia)}</b></p>`;
}

function aviso(t, tipo = "ok") {
  const msgEl = $("msg");
  msgEl.textContent = t;
  msgEl.className = tipo === "error" ? "msg-error" : "msg-ok";
  clearTimeout(window._avisoTimer);
  window._avisoTimer = setTimeout(() => {
    msgEl.textContent = "";
    msgEl.className = "";
  }, 4500);
}

// ---------- Validación ----------

function validarTodosLosDatos() {
  let valido = true;
  let primerInvalido = null;
  const errores = [];

  document.querySelectorAll(".input-error").forEach((el) => el.classList.remove("input-error"));

  // Normalizar campos opcionales numéricos vacíos a "0" antes de validar
  document.querySelectorAll("main input[type='number']").forEach((el) => {
    if (!el.id.startsWith("precio_") && !el.id.startsWith("carga_") && el.value.trim() === "") {
      el.value = "0";
    }
  });

  const marcarError = (el, texto) => {
    if (el) {
      el.classList.add("input-error");
      if (!primerInvalido) primerInvalido = el;
    }
    if (texto && !errores.includes(texto)) errores.push(texto);
    valido = false;
  };

  // 1. Fecha requerida
  if (!$("fecha").value) {
    marcarError($("btnCalToggle"), "Debe seleccionar una fecha para el cierre");
  }

  // 2. Precios requeridos y > 0
  PRODS.forEach(([k, n]) => {
    const el = $("precio_" + k);
    const val = parseFloat(el.value);
    if (isNaN(val) || val <= 0) {
      marcarError(el, `El precio de ${n} es requerido (mayor a 0)`);
    }
  });

  // 3. Carga y Devuelto requeridos
  PRODS.forEach(([k, n]) => {
    const elCarga = $("carga_" + k);
    const elDev = $("dev_" + k);
    const valCarga = elCarga.value.trim() === "" ? NaN : parseInt(elCarga.value, 10);
    const valDev = elDev.value.trim() === "" ? NaN : parseInt(elDev.value, 10);

    if (isNaN(valCarga) || valCarga < 0) {
      marcarError(elCarga, `La carga de ${n} es requerida (ingrese 0 si no hubo)`);
    }
    if (isNaN(valDev) || valDev < 0) {
      marcarError(elDev, `El devuelto de ${n} es requerido (ingrese 0 si no hubo)`);
    }
    if (!isNaN(valCarga) && !isNaN(valDev) && valDev > valCarga) {
      marcarError(elDev, `El devuelto de ${n} no puede ser mayor que la carga`);
    }
  });

  // 4. Rendición (no pueden estar vacíos)
  ["transferencias", "fiados", "gastos"].forEach((id) => {
    const el = $(id);
    const val = el.value.trim() === "" ? NaN : parseFloat(el.value);
    if (isNaN(val) || val < 0) {
      marcarError(el, `El campo ${id} es requerido (ingrese 0 si no hubo)`);
    }
  });

  // 5. Billetes (no pueden estar vacíos)
  DENS.forEach((d) => {
    const el = $("b_" + d);
    const val = el.value.trim() === "" ? NaN : parseInt(el.value, 10);
    if (isNaN(val) || val < 0) {
      marcarError(el, `La cantidad de billetes de $${d} es requerida (ingrese 0 si no hubo)`);
    }
  });

  if (!valido && primerInvalido) {
    primerInvalido.focus();
  }
  return { valido, errores };
}

// ---------- Calendario Personalizado ----------

function setupCalendario() {
  const popover = $("calPopover");
  const btnToggle = $("btnCalToggle");

  btnToggle.onclick = (e) => {
    e.stopPropagation();
    if ($("historialDropdown")) $("historialDropdown").style.display = "none";
    const visible = popover.style.display !== "none";
    if (visible) {
      popover.style.display = "none";
    } else {
      // Sincronizar fecha del calendario con la fecha seleccionada
      if ($("fecha").value) {
        const [y, m, d] = $("fecha").value.split("-").map(Number);
        currentCalDate = new Date(y, m - 1, d);
      }
      renderCalendario();
      popover.style.display = "block";
    }
  };

  $("calPrev").onclick = (e) => {
    e.stopPropagation();
    currentCalDate.setMonth(currentCalDate.getMonth() - 1);
    renderCalendario();
  };

  $("calNext").onclick = (e) => {
    e.stopPropagation();
    currentCalDate.setMonth(currentCalDate.getMonth() + 1);
    renderCalendario();
  };

  $("calBtnHoy").onclick = async (e) => {
    e.stopPropagation();
    const hoy = new Date();
    currentCalDate = new Date(hoy);
    const hoyStr = formatISODate(hoy.getFullYear(), hoy.getMonth() + 1, hoy.getDate());
    await cargarPorFecha(hoyStr);
    popover.style.display = "none";
  };

  document.addEventListener("click", (e) => {
    if (!popover.contains(e.target) && e.target !== btnToggle && !btnToggle.contains(e.target)) {
      popover.style.display = "none";
    }
  });
}

function formatISODate(y, m, d) {
  const mm = String(m).padStart(2, "0");
  const dd = String(d).padStart(2, "0");
  return `${y}-${mm}-${dd}`;
}

function renderCalendario() {
  const anio = currentCalDate.getFullYear();
  const mesIndex = currentCalDate.getMonth();
  $("calTituloMes").textContent = `${MESES[mesIndex]} ${anio}`;

  const grid = $("calGrid");
  grid.innerHTML = "";

  // Primer día del mes
  const primerDia = new Date(anio, mesIndex, 1);
  // getDay() devuelve 0 para domingo, 1 para lunes... Ajustar para que Lunes sea 0
  let startDay = primerDia.getDay() - 1;
  if (startDay === -1) startDay = 6;

  const diasEnMes = new Date(anio, mesIndex + 1, 0).getDate();
  const diasEnMesAnterior = new Date(anio, mesIndex, 0).getDate();

  const fechaSeleccionada = $("fecha").value;
  const hoyObj = new Date();
  const hoyStr = formatISODate(hoyObj.getFullYear(), hoyObj.getMonth() + 1, hoyObj.getDate());

  // Relleno días del mes anterior
  for (let i = startDay - 1; i >= 0; i--) {
    const d = diasEnMesAnterior - i;
    const cell = document.createElement("button");
    cell.type = "button";
    cell.className = "cal-day other-month";
    cell.textContent = d;
    cell.disabled = true;
    grid.appendChild(cell);
  }

  // Días del mes actual
  for (let d = 1; d <= diasEnMes; d++) {
    const dateStr = formatISODate(anio, mesIndex + 1, d);
    const cell = document.createElement("button");
    cell.type = "button";
    cell.className = "cal-day";
    cell.textContent = d;

    const tieneCierre = fechasGuardadas.includes(dateStr);
    if (tieneCierre) {
      cell.classList.add("has-cierre");
      cell.title = "Cierre guardado en base de datos";
    }

    if (dateStr === fechaSeleccionada) {
      cell.classList.add("is-selected");
    }

    if (dateStr === hoyStr) {
      cell.classList.add("is-today");
    }

    cell.onclick = async (e) => {
      e.stopPropagation();
      $("calPopover").style.display = "none";
      await cargarPorFecha(dateStr);
    };

    grid.appendChild(cell);
  }
}

async function cargarPorFecha(fechaStr) {
  try {
    const tieneCierre = fechasGuardadas.includes(fechaStr);
    const j = await api().Cargar(fechaStr);
    volcar(j, tieneCierre);
    if ($("historialLabel")) {
      $("historialLabel").textContent = tieneCierre ? formatearFechaVisual(fechaStr) : "Historial…";
    }
    renderHistorialDropdown();
  } catch (e) {
    aviso("Error al cargar fecha: " + e, "error");
  }
}

// ---------- Dropdown de Historial Personalizado ----------

function setupHistorialDropdown() {
  const btn = $("btnHistorialToggle");
  const menu = $("historialDropdown");
  if (!btn || !menu) return;

  btn.onclick = (e) => {
    e.stopPropagation();
    if ($("calPopover")) $("calPopover").style.display = "none";
    const visible = menu.style.display !== "none";
    menu.style.display = visible ? "none" : "block";
  };

  document.addEventListener("click", (e) => {
    if (!menu.contains(e.target) && e.target !== btn && !btn.contains(e.target)) {
      menu.style.display = "none";
    }
  });
}

function renderHistorialDropdown() {
  const container = $("historialItems");
  if (!container) return;
  container.innerHTML = "";

  if (!fechasGuardadas || fechasGuardadas.length === 0) {
    const empty = document.createElement("div");
    empty.className = "dropdown-item empty";
    empty.textContent = "Sin cierres guardados";
    container.appendChild(empty);
    return;
  }

  const fechaActual = $("fecha").value;

  fechasGuardadas.forEach((fechaStr) => {
    const item = document.createElement("div");
    item.className = "dropdown-item";
    if (fechaStr === fechaActual) {
      item.classList.add("is-active");
    }
    item.textContent = `${formatearFechaVisual(fechaStr)} (${fechaStr})`;

    item.onclick = async (e) => {
      e.stopPropagation();
      $("historialDropdown").style.display = "none";
      await cargarPorFecha(fechaStr);
    };

    container.appendChild(item);
  });
}

// ---------- Estadísticas ----------

async function cargarEstadisticas(fechaStr) {
  try {
    const stats = await api().ObtenerEstadisticas(fechaStr || $("fecha").value);
    if (!stats) return;

    // Subtítulo
    $("statsSubtitulo").textContent = `Fecha de referencia: ${formatearFechaVisual(fechaStr)}`;

    // 1. Día
    $("statDiaPeriodo").textContent = stats.dia.etiqueta + (stats.dia.cantidad > 0 ? " (Cierre registrado)" : " (Sin cierre registrado)");
    $("statDiaGanado").textContent = peso(stats.dia.ganado);
    $("statDiaGastos").textContent = peso(stats.dia.gastos);
    $("statDiaParteA").textContent = peso(stats.dia.parteA);
    $("statDiaParteB").textContent = peso(stats.dia.parteB);

    // 2. Semana
    $("statSemanaPeriodo").textContent = stats.semana.etiqueta + ` (${stats.semana.cantidad} ${stats.semana.cantidad === 1 ? "cierre" : "cierres"})`;
    $("statSemanaGanado").textContent = peso(stats.semana.ganado);
    $("statSemanaGastos").textContent = peso(stats.semana.gastos);
    $("statSemanaParteA").textContent = peso(stats.semana.parteA);
    $("statSemanaParteB").textContent = peso(stats.semana.parteB);

    // 3. Mes
    $("statMesPeriodo").textContent = stats.mes.etiqueta + ` (${stats.mes.cantidad} ${stats.mes.cantidad === 1 ? "cierre" : "cierres"})`;
    $("statMesGanado").textContent = peso(stats.mes.ganado);
    $("statMesGastos").textContent = peso(stats.mes.gastos);
    $("statMesParteA").textContent = peso(stats.mes.parteA);
    $("statMesParteB").textContent = peso(stats.mes.parteB);
  } catch (e) {
    console.error("Error cargando estadísticas:", e);
  }
}

async function cargarHistorial() {
  try {
    fechasGuardadas = (await api().FechasGuardadas()) || [];
    renderHistorialDropdown();
    renderCalendario();
  } catch (e) {
    console.error("Error cargando historial:", e);
  }
}

// ---------- Inicialización ----------

async function init() {
  construirUI();
  setupCalendario();
  setupHistorialDropdown();
  await cargarHistorial();

  const jInicial = await api().Nueva("");
  const tieneCierre = fechasGuardadas.includes(jInicial.fecha);
  volcar(jInicial, tieneCierre);
  if ($("historialLabel")) {
    $("historialLabel").textContent = tieneCierre ? formatearFechaVisual(jInicial.fecha) : "Historial…";
  }

  // Botón Guardar
  $("btnGuardar").onclick = async () => {
    const { valido, errores } = validarTodosLosDatos();
    if (!valido) {
      aviso("⚠️ " + errores[0], "error");
      return;
    }

    try {
      const j = leer();
      await api().Guardar(j);
      cierreGuardado = true;
      actualizarEstadoGuardado();
      aviso("✔ Cierre guardado con éxito", "ok");
      await cargarHistorial();
      await cargarEstadisticas(j.fecha);
    } catch (e) {
      aviso("Error al guardar: " + e, "error");
    }
  };

  // Botón Nueva Jornada
  $("btnNueva").onclick = async () => {
    const hoyObj = new Date();
    const hoyStr = formatISODate(hoyObj.getFullYear(), hoyObj.getMonth() + 1, hoyObj.getDate());
    const j = await api().Nueva(hoyStr);
    const tieneCierre = fechasGuardadas.includes(hoyStr);
    volcar(j, tieneCierre);
    if ($("historialLabel")) {
      $("historialLabel").textContent = "Historial…";
    }
    renderHistorialDropdown();
    aviso("Nueva jornada iniciada", "ok");
  };

  // Botón Imprimir / PDF
  $("btnImprimir").onclick = () => {
    // 1. Validar que todos los campos requeridos estén completos
    const { valido, errores } = validarTodosLosDatos();
    if (!valido) {
      aviso("⚠️ Debe completar todos los datos requeridos antes de imprimir.", "error");
      return;
    }

    // 2. Validar que el cierre esté guardado
    if (!cierreGuardado) {
      aviso("⚠️ Debe guardar el cierre antes de poder imprimir el resumen.", "error");
      $("btnGuardar").classList.add("pulse-highlight");
      setTimeout(() => $("btnGuardar").classList.remove("pulse-highlight"), 1200);
      return;
    }

    // Si todo está correcto y guardado:
    window.print();
  };
}

// Toggle Dark Mode
function setupTheme() {
  const btn = $("btnTheme");
  if (!btn) return;

  const setDark = (isDark) => {
    document.documentElement.setAttribute("data-theme", isDark ? "dark" : "light");
    btn.textContent = isDark ? "☀️" : "🌙";
    localStorage.setItem("theme", isDark ? "dark" : "light");
  };

  const saved = localStorage.getItem("theme");
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  setDark(saved ? saved === "dark" : prefersDark);

  btn.onclick = () => {
    const isDark = document.documentElement.getAttribute("data-theme") === "dark";
    setDark(!isDark);
  };
}

window.addEventListener("DOMContentLoaded", () => {
  setupTheme();
  init();
});
