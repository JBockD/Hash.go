package diccionario

import (
	"fmt"
	"hash/fnv"
)

type estado int

const (
	VACIO estado = iota
	OCUPADA
	BORRADA
)
const (
	TRESCUARTOS         = 0.75
	UNCUARTO            = 0.25
	POR_CUANTO_AUMENTAR = 2
)

// PRE: La clave debe ser de un tipo comparable
// POST: Devuelve un slice de bytes representando la clave
func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

// PRE: La clave debe ser de un tipo comparable
// POST: Devuelve un hash de la clave utilizando el algoritmo FNV-1a
func calcularHash[K comparable](clave K) uint64 {
	h := fnv.New64a()
	h.Write(convertirABytes(clave))
	return h.Sum64()
}

// PRE: El hash debe estar inicializado
// POST: Devuelve la posición inicial para una clave dada
func (h *hashCerrado[K, V]) posicionInicial(clave K) int {
	return int(calcularHash(clave) % uint64(h.tam))
}

type celdaHash[K comparable, V any] struct {
	clave  K
	dato   V
	estado estado
}
type hashCerrado[K comparable, V any] struct {
	tabla    []celdaHash[K, V]
	cantidad int
	tam      int
	borrados int
}
type iterador[K comparable, V any] struct {
	hash   *hashCerrado[K, V]
	actual int
}

// PRE: El hash debe estar inicializado
// POST: Devuelve un iterador para el hash cerrado11
func (h *hashCerrado[K, V]) Iterador() IterDiccionario[K, V] {
	it := &iterador[K, V]{
		hash:   h,
		actual: 0,
	}
	for it.actual < len(it.hash.tabla) && h.tabla[it.actual].estado != OCUPADA {
		it.actual++
	}
	return it
}

// PRE: El iterador debe estar inicializado
// POST: Devuelve true si el iterador se encuentra en una posición OCUPADA
func (it *iterador[K, V]) HayAlgoMas() bool {
	if it.actual >= len(it.hash.tabla) {
		return false
	}
	return it.hash.tabla[it.actual].estado == OCUPADA
}

// PRE: El iterador debe estar inicializado
// POST: Avanza el iterador al siguiente elemento OCUPADO
func (it *iterador[K, V]) Avanzar() {
	if !it.HayAlgoMas() {
		panic("El iterador termino de iterar")
	}
	it.actual++
	for it.actual < len(it.hash.tabla) && it.hash.tabla[it.actual].estado != OCUPADA {
		it.actual++
	}
}

// PRE: El iterador debe estar inicializado
// POST: Devuelve la clave y el dato del elemento donde esta el iterador
func (it iterador[K, V]) VerActual() (K, V) {
	if !it.HayAlgoMas() {
		panic("El iterador termino de iterar")
	}
	return it.hash.tabla[it.actual].clave, it.hash.tabla[it.actual].dato

}

// PRE: -
// POST: Crea una tabla
func crearTabla[K comparable, V any](tamano int) []celdaHash[K, V] {
	return make([]celdaHash[K, V], tamano)
}

// PRE: -
// POST: Crea un hash cerrado con un tamaño inicial de 5
func CrearHash[K comparable, V any]() Diccionario[K, V] {
	return &hashCerrado[K, V]{
		tabla:    crearTabla[K, V](5),
		cantidad: 0,
		tam:      5,
		borrados: 0,
	}
}

// PRE: Visitar debe devolver valores validos
// POST: Itera el hash si visitar es true, se detiene si visitar devuelve false
func (hash *hashCerrado[K, V]) Iterar(visitar func(K, V) bool) {
	for i := 0; i < hash.tam; i++ {
		actual := hash.tabla[i]
		if actual.estado == OCUPADA {
			seguir := visitar(actual.clave, actual.dato)
			if !seguir {
				return

			}
		}
	}
}

// PRE: El hash debe estar inicializado
// POST: Devuelve la posición de la clave en la tabla y un booleano indicando si se encontró o no
func (hash *hashCerrado[K, V]) buscar(clave K) (int, bool) {
	posicion := hash.posicionInicial(clave)

	for i := 0; i < hash.tam; i++ {
		posicionActual := (posicion + i) % hash.tam
		celda := hash.tabla[posicionActual]
		if celda.estado == VACIO {
			return posicionActual, false
		} else if celda.estado == OCUPADA && celda.clave == clave {
			return posicionActual, true
		}
	}
	return -1, false
}

// PRE: El hash debe estar inicializado
// POST: Devuelve true si la clave esta en la tabla hash
func (hash *hashCerrado[K, V]) Pertenece(clave K) bool {
	_, encontrada := hash.buscar(clave)
	return encontrada
}

// PRE: El hash debe estar inicializado
// POST: Devuelve el valor asociado a la clave
func (hash *hashCerrado[K, V]) Obtener(clave K) V {
	posicion, encontrada := hash.buscar(clave)

	if encontrada {
		return hash.tabla[posicion].dato
	}
	panic("La clave no pertenece al diccionario")
}

// PRE: La tabla hash debe estar ocupada un 75%
// POST: Copia todos los elementos a una nueva tabla de con el tamañano duplicado y actualiza los valores de BORRADOS y TAMAÑO
func redimensionar[K comparable, V any](hash *hashCerrado[K, V], nuevo_tamano int) {
	nueva_tabla := crearTabla[K, V](nuevo_tamano)

	for i := 0; i < hash.tam; i++ {
		celda := hash.tabla[i]
		if celda.estado == OCUPADA {
			posicion := int(calcularHash(celda.clave) % uint64(nuevo_tamano))
			for nueva_tabla[posicion].estado == OCUPADA {
				posicion = (posicion + 1) % nuevo_tamano
			}
			nueva_tabla[posicion] = celda
		}
	}
	hash.tabla = nueva_tabla
	hash.borrados = 0
	hash.tam = nuevo_tamano

}

// PRE: El hash debe estar inicializado
// POST: Devuelve la carga del hash
func calcularCarga[K comparable, V any](hash *hashCerrado[K, V]) float64 {
	return float64(hash.cantidad+hash.borrados) / float64(hash.tam)
}

// PRE: El hash debe estar inicializado
// POST: Redimensiona la tabla si la carga es mayor a 0.75 y devuelve la nueva posición de la clave
func (hash *hashCerrado[K, V]) verificarYRedimensionar(clave K) (int, bool) {
	if calcularCarga(hash) >= TRESCUARTOS {
		redimensionar(hash, hash.tam*POR_CUANTO_AUMENTAR)
		return hash.buscar(clave)
	}
	return hash.buscar(clave)
}

// PRE: El hash debe tener espacios habilitados
// POST: Guardo un elemento en la tabla Hash si el espacio no esta OCUPADO o si la clave es la misma
func (hash *hashCerrado[K, V]) Guardar(clave K, dato V) {
	posicion, encontrada := hash.verificarYRedimensionar(clave)

	if encontrada {
		hash.tabla[posicion].dato = dato
	} else {
		hash.tabla[posicion] = celdaHash[K, V]{clave, dato, OCUPADA}
		hash.cantidad++
	}
}

// PRE: El hash debe estar inicializado
// POST: Devuelve la cantidad de elementos dentro del diccionario
func (hash *hashCerrado[K, V]) Cantidad() int {

	return hash.cantidad
}

// PRE: El hash debe estar inicializado y la clave debe pertenecer al diccionario
// POST: Devuelve el dato asociado a la clave y cambia el estado de la celda a BORRADA
func (hash *hashCerrado[K, V]) Borrar(clave K) V {
	posicion, encontrada := hash.buscar(clave)
	if !encontrada {
		panic("La clave no pertenece al diccionario")
	}

	posicion, encontrada = hash.verificarYRedimensionar(clave)

	hash.tabla[posicion].estado = BORRADA
	hash.cantidad--
	hash.borrados++

	return hash.tabla[posicion].dato
}
