package main

import (
	"fmt"
	"reflect"
)

type user struct {
	Name string
	Age  int
	Int1 int
	F1   float64
	Int3 int
}

func multiplyNumeric(in interface{}) reflect.Value {
	typeA := reflect.TypeOf(in)
	valueA := reflect.ValueOf(in)
	strField := make([]reflect.StructField, typeA.NumField())
	valArr := make([]reflect.Value, typeA.NumField())
	for i := 0; i < typeA.NumField(); i++ {
		/*fmt.Print(typeA.Field(i).Name, ": ")
		fmt.Println(valueA.Field(i))*/
		strField[i] = typeA.Field(i)
		valArr[i] = valueA.Field(i)
	}
	v := reflect.New(reflect.StructOf(strField)).Elem()
	for i, val := range valArr {
		if typeA.Field(i).Type.Name() == "int" {
			v.Field(i).SetInt(val.Int() * 2)
		} else if typeA.Field(i).Type.Name() == "float64" {
			v.Field(i).SetFloat(val.Float() * 2)
		} else {
			v.Field(i).Set(val)
		}
	}
	return v
}
func main() {
	a := user{"John", 30, 2, 12.3, 14}
	fmt.Println(multiplyNumeric(a))

}
