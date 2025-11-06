package functemplates

import "reflect"


func defaultFunc(value, defaultValue interface{}) interface{} {
    v := reflect.ValueOf(value)

    // 1. Lida com interface{}(nil)
    if !v.IsValid() {
        return defaultValue
    }

    // 2. Lida com ponteiros
    for v.Kind() == reflect.Ptr {
        // Se o ponteiro for nil, é "zero"
        if v.IsNil() {
            return defaultValue
        }
        v = v.Elem()

        if !v.IsValid() {
            return defaultValue
        }
    }

    // 3. Verifique se o valor concreto é o "valor zero" do seu tipo
    if v.IsZero() {
        return defaultValue
    }

    // Retorna o valor original.
    return value
}