package main

import (
	"fmt"
	TDADiccionario "hash/diccionario"
)

func main() {
	dic := TDADiccionario.CrearHash[string, string]()
	claves := make([]string, 10)
	valores := make([]string, 10)
	cadena := "%d~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~" +
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
	for i := 0; i < 10; i++ {
		claves[i] = fmt.Sprintf(cadena, i)
		valores[i] = fmt.Sprintf("valor_%d", i)
	}
	for i := 0; i < 10; i++ {
		dic.Guardar(claves[i], valores[i])
	}
	fmt.Println("Verificando claves y valores...")
	ok := true
	for i := 0; i < 10 && ok; i++ {
		valorObtenido := dic.Obtener(claves[i])
		fmt.Printf("Clave: %s, Valor esperado: %s, Valor obtenido: %s\n", claves[i], valores[i], valorObtenido)
		ok = valorObtenido == valores[i]
	}

	if ok {
		fmt.Println("Todas las claves largas se obtuvieron correctamente.")
	} else {
		fmt.Println("Error: No se pudieron obtener todas las claves largas correctamente.")
	}
}
