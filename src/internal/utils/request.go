package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func ReadID(ctx *gin.Context, param string) (id int, err error) {
	
	valueStr := ctx.Param(param)
	id, err = strconv.Atoi(valueStr)
	return
}

func GetArrayInt(ctx *gin.Context, param string) ([]int, error) {
	value := ctx.QueryArray(param)
	values := make([]int, len(value))
	for i, v := range value {
		vint, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		values[i] = vint
	}
	
	if len(values) == 0 {
		return nil, errors.New("array length is zero")
	}
	
	return values, nil
}

func GetArrayInt32(ctx *gin.Context, param string) ([]int32, error) {
	value := ctx.QueryArray(param)
	values := make([]int32, len(value))
	for i, v := range value {
		vint, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		values[i] = int32(vint)
	}
	
	if len(values) == 0 {
		return nil, errors.New("array length is zero")
	}
	
	return values, nil
}

func ReadBool(ctx *gin.Context, param string) (bool, error) {
	value := ctx.Query(param)
	num, err := strconv.Atoi(value)
	if err != nil {
		return false, err
	}
	if num == 1 {
		return true, nil
	}
	return false, nil
}

func GetArrayInt16(ctx *gin.Context, param string) ([]int16, error) {
	value := ctx.QueryArray(param)
	values := make([]int16, len(value))
	for i, v := range value {
		vint, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		values[i] = int16(vint)
	}
	if len(values) == 0 {
		return nil, errors.New("array length is zero")
	}
	
	return values, nil
}

func GetCSVFloat(value string) ([]float64, error) {
	sep := strings.Split(value, ",")
	values := make([]float64, len(sep))
	for i := 0; i < len(sep); i++ {
		v, err := strconv.ParseFloat(sep[i], 64)
		if err != nil {
			return nil, err
		}
		values[i] = v
	}
	
	return values, nil
}

