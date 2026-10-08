# 💧 Reparto — Sistema de Control de Carga, Ventas y Arqueo de Caja

> **Solución de escritorio moderna, robusta y eficiente para la administración de reparto de agua y soda.**  
> Desarrollado a medida por **Apolo Studio** con **Go**, **Wails v2** y motor de base de datos **SQLite**.

---

## 📌 Tabla de Contenidos
- [Características del Sistema](#-características-del-sistema)
- [Guía de Uso para el Usuario](#-guía-de-uso-para-el-usuario)
- [Instalación y Puesta en Marcha](#-instalación-y-puesta-en-marcha)
- [Seguridad y Respaldos de Datos](#-seguridad-y-respaldos-de-datos)
- [Arquitectura del Backend (Para Desarrolladores)](#-arquitectura-del-backend)
- [Compilación para Producción](#-compilación-para-producción)
- [Soporte y Créditos](#-soporte-y-créditos)

---

## 🚀 Características del Sistema

### 1. Control Integral de Carga y Ventas
- Registro de carga inicial y unidades devueltas por producto (**Bidón 20L**, **Bidón 6L** y **Soda**).
- Cálculo automático de unidades vendidas y recaudación por producto y global.
- Equivalencia configurable de cajones de soda (6 unidades).

### 2. Rendición Financiera y Arqueo de Caja
- Deducción inmediata de **transferencias bancarias**, **fiados / créditos** y **gastos del día**.
- Cálculo en vivo del **efectivo teórico** y división equitativa de utilidades:
  - **Parte A (50%)**
  - **Parte B (50%)**
- Desglose y arqueo detallado por denominación de billetes ($20.000, $10.000, $2.000, $1.000, $500).
- Detección visual instantánea del estado de caja:
  - ✔ **Caja exacta**
  - ✖ **Falta dinero**
  - ▲ **Sobra dinero**

### 3. Base de Datos SQLite y Persistencia Local
- Almacenamiento seguro, rápido y libre de dependencias externas en base de datos local SQLite (`modernc.org/sqlite`).
- Optimizado en modo WAL (*Write-Ahead Logging*) para máxima integridad de datos.
- Migración y respaldo redundante automático en archivos JSON.

### 4. Calendario y Navegación Inteligente
- **Calendario interactivo desplegable**: Días con cierres guardados resaltados en **color verde distintivo** para un control visual rápido del mes.
- **Historial dinámico**: Menú desplegable con límite visual de 7 registros y barra de desplazamiento suave (*scroll*).
- Carga instantánea de planillas históricas con un solo clic.

### 5. Estadísticas de Cierre (Día, Semana y Mes)
- Panel de estadísticas claras y directas para la fecha de trabajo:
  - 📅 **Cierre Diario**: Cuánto se ganó, cuánto se gastó y cuánto correspondió a cada parte.
  - 📆 **Cierre Semanal**: Resumen acumulado de lunes a domingo con cantidad de cierres registrados.
  - 🗓️ **Cierre Mensual**: Rendimiento total del mes en curso.

### 6. Control de Calidad y Bloqueo de Impresión
- **Validación estricta**: Garantiza que no queden campos en blanco ni valores inconsistentes (por ejemplo, unidades devueltas superiores a la carga).
- **Protección de impresión**: No permite imprimir o exportar el resumen en PDF sin haber guardado previamente el cierre en la base de datos.
- **Entrada ágil**: Eliminación de flechas molestas en números y borrado automático del cero al seleccionar cualquier campo.
- **Modo Oscuro / Claro**: Selector de tema con memoria para comodidad visual.

---

## 📖 Guía de Uso para el Usuario

### Paso a Paso para la Jornada Diaria:

1. **Seleccionar Fecha:**
   - La aplicación inicia por defecto con la fecha de hoy. Podés abrir el calendario (`📅`) para elegir otra fecha o revisar días anteriores marcados en verde.
2. **Revisar Precios:**
   - Verificá los precios unitarios de cada producto (se recuerdan automáticamente del último cierre).
3. **Ingresar Carga y Devolución:**
   - Completá la cantidad de bidones y sodas que salieron en el reparto y los que volvieron en el camión.
4. **Completar Rendición:**
   - Ingresá transferencias recibidas, cuentas fiadas y gastos del día (combustible, almuerzo, etc.). Si no hubo gastos o fiados, pueden quedar en `0`.
5. **Realizar Arqueo de Billetes:**
   - Contá y cargá la cantidad de cada billete en la tabla de arqueo. El sistema te indicará si la caja está exacta, si falta o si sobra dinero.
6. **Guardar el Cierre:**
   - Hacé clic en **Guardar**. El sistema validará que todo esté completo y te mostrará el mensaje `✔ Cierre guardado con éxito`. El día quedará marcado en el calendario y se actualizarán las estadísticas semanales y mensuales.
7. **Imprimir / Guardar en PDF:**
   - Con el cierre guardado, hacé clic en **Imprimir / PDF** para obtener el comprobante físico o guardarlo en formato digital.

---

## 💾 Instalación y Puesta en Marcha

El sistema está empaquetado como un **ejecutable único portable**:

- **Ubicación:** `build\bin\reparto.exe`
- **Requisitos:** Windows 10 o Windows 11 (64 bits).
- **Instalación:** No requiere instalador ni librerías externas. Podés colocar `reparto.exe` en el Escritorio, en una carpeta dedicada o ejecutarlo desde un pendrive.

---

## 🔒 Seguridad y Respaldos de Datos

Todos los registros se almacenan localmente en la computadora del usuario dentro del directorio seguro de Windows:

```
%APPDATA%\RepartoAguaSoda\
├── reparto.db         <-- Base de datos principal SQLite
├── precios.json       <-- Última configuración de precios
└── jornadas\          <-- Copias de respaldo en formato JSON por fecha
    ├── 2026-10-07.json
    └── ...
```

> **Nota:** Para realizar una copia de seguridad periódica, simplemente copiá la carpeta `%APPDATA%\RepartoAguaSoda` a un dispositivo externo o servicio en la nube (Google Drive, OneDrive, etc.).

---

## 🏗 Arquitectura del Backend

El código fuente en Go sigue una **Arquitectura en Capas (Layered / Clean Architecture)** para garantizar mantenibilidad y modularidad:

```
reparto/
├── main.go                     # Punto de entrada Wails v2
├── app.go                      # Adaptador de aplicación / Fachada Wails
├── build.bat                   # Script de compilación de producción en 1 clic
├── internal/
│   ├── domain/                 # Entidades de negocio y contratos (Repository)
│   │   ├── models.go
│   │   └── repository.go
│   ├── service/                # Lógica de cálculo y casos de uso
│   │   ├── calculator.go       # Fórmulas matemáticas de caja y arqueo
│   │   └── jornada_service.go  # Orquestación de guardado, carga y estadísticas
│   └── repository/             # Acceso a datos e infraestructura
│       └── sqlite_jornada.go   # Manejador SQLite (driver pure-Go modernc)
└── frontend/dist/              # Interfaz gráfica moderna (HTML5, CSS3, JS Vanilla)
```

---

## 🛠 Compilación para Producción

Si disponés del código fuente y necesitás generar un nuevo binario:

### Requisitos previos:
- [Go](https://golang.org/) 1.21 o superior instalado.
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`).

### Compilación rápida:
Hacé doble clic en el archivo **`build.bat`** o ejecutá en la terminal:

```bash
wails build -clean -s
```

El ejecutable optimizado sin ventana de consola se generará automáticamente en `build\bin\reparto.exe`.

---

## 🤝 Soporte y Créditos

Este software fue diseñado y desarrollado con altos estándares de calidad por **Apolo Studio**.

Para asistencia técnica, solicitudes de soporte o personalizaciones, comunicate con el equipo de soporte.
