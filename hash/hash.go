package diccionario
import (
   "hash/fnv"
   "fmt"
)

type estado int 
const (
	VACIO estado = iota
	OCUPADA
	BORRADA
)
func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}
func hash[K comparable](clave K) uint64 {
	h := fnv.New64a()
	h.Write(convertirABytes(clave))
	return h.Sum64()
}
func (h *hashCerrado[K, V]) posicionInicial(clave K) int {
	return int(hash(clave)) % len(h.tabla)
}
type celdaHash[K comparable, V any] struct {
   clave  K
   dato   V
   estado
}

type hashCerrado[K comparable, V any] struct {
   tabla    []celdaHash[K,V]
   cantidad int
   tam      int
   borrados int
}

func (hash *hashCerrado[K, V]) Obtener(clave K) V{
   pos := hash.posicionInicial(clave)
   for i := 0; i < len(hash.tabla) ; i ++{
      posactual := (pos + i) % len(hash.tabla) - 1
      celda := hash.tabla[posactual]
      if celda.estado == VACIO{
         panic("no se encontro esa clave")
      } else if celda.estado == OCUPADA && celda.clave == clave{
         return celda.dato
      
      } 
   }
   panic("no se encontró")

}
func (hash *hashCerrado[K, V]) Pertenece(clave K) bool{
   pos := hash.posicionInicial(clave)
   for i := 0; i < len(hash.tabla) ; i ++{
      posactual := (pos + i) % len(hash.tabla)
      celda := hash.tabla[posactual]
      if celda.estado == VACIO{ 
         return false
      } else if celda.estado == OCUPADA && celda.clave == clave {
         return true
         
      }
   }
   return false
}