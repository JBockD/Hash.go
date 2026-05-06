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
	POR_CUANTO_AUMENTAR = 2
)

func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}
func calcular_hash[K comparable](clave K) uint64 {
	h := fnv.New64a()
	h.Write(convertirABytes(clave))
	return h.Sum64()
}
func (h *hashCerrado[K, V]) posicionInicial(clave K) int {
	return int(calcular_hash(clave)) % len(h.tabla)
}

type celdaHash[K comparable, V any] struct {
	clave K
	dato  V
	estado
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
func (it *iterador[K, V]) HayAlgoMas() bool {
	if it.actual >= len(it.hash.tabla) {
		return false
	}
	return it.hash.tabla[it.actual].estado == OCUPADA
}
func (it iterador[K, V]) Avanzar() {
	it.actual++
	if !it.HayAlgoMas() {
		panic("El iterador termino de iterar")
	}

}
func (it iterador[K, V]) VerActual() (K, V) {
	if !it.HayAlgoMas() {
		panic("El iterador termino de iterar")
	}
	return it.hash.tabla[it.actual].clave, it.hash.tabla[it.actual].dato

}

// PRE: -
// POST: Crea un hash cerrado con un tamaño inicial de 5
func CrearHash[K comparable, V any]() *hashCerrado[K, V] {
	return &hashCerrado[K, V]{tabla: make([]celdaHash[K, V], 5), cantidad: 0, tam: 5, borrados: 0}
}
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
	posicion, _ := hash.buscar(clave)

	if hash.Pertenece(clave) {
		return hash.tabla[posicion].dato
	}
	panic("La clave no pertenece al diccionario")
}

// PRE: La tabla hash debe estar ocupada un 75%
// POST: Copia todos los elementos a una nueva tabla de con el tamañano duplicado y actualiza los valores de BORRADOS y TAMAÑO
func redimensionar[K comparable, V any](hash *hashCerrado[K, V], nuevo_tamano int) {
	nueva_tabla := make([]celdaHash[K, V], nuevo_tamano)

	for _, celda := range hash.tabla {
		if celda.estado == OCUPADA {
			posicion := int(calcular_hash(celda.clave)) % nuevo_tamano
			for i := 0; i < nuevo_tamano; i++ {
				posicion_actual := (posicion + i) % nuevo_tamano
				if nueva_tabla[posicion_actual].estado != OCUPADA {
					nueva_tabla[posicion_actual] = celdaHash[K, V]{celda.clave, celda.dato, OCUPADA}
				}
			}
		}
	}

	hash.tabla = nueva_tabla
	hash.borrados = int(VACIO)
	hash.tam = nuevo_tamano
}

// PRE: El hash debe tener espacios habilitados
// POST: Guardo un elemento en la tabla Hash si el espacio no esta OCUPADO o si la clave es la misma
func (hash *hashCerrado[K, V]) Guardar(clave K, dato V) {
	posicion, encontrada := hash.buscar(clave)
	carga := (hash.cantidad + hash.borrados) / hash.tam

	if carga >= int(float64(hash.tam)*TRESCUARTOS) {
		redimensionar(hash, hash.tam*POR_CUANTO_AUMENTAR)
	}

	if !encontrada {
		hash.tabla[posicion] = celdaHash[K, V]{clave, dato, OCUPADA}
		hash.cantidad++
	} else {
		hash.tabla[posicion].dato = dato
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

	hash.tabla[posicion].estado = BORRADA
	hash.cantidad--
	hash.borrados++
	return hash.tabla[posicion].dato
}
