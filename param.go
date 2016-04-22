package scaffold

import (
	"fmt"
	"strconv"

	"golang.org/x/net/context"
)

// Param is a parameter
type Param string

func paramName(name string) string {
	return fmt.Sprintf("scaffold_param_%s", name)
}

//StoreParam stores a variable in a context
func StoreParam(ctx context.Context, name string, value string) context.Context {
	n := paramName(name)
	return context.WithValue(ctx, n, Param(value))
}

// GetParam retrieves a param from a context
func GetParam(ctx context.Context, name string) Param {
	n := paramName(name)
	if p, ok := ctx.Value(n).(Param); ok {
		return p
	}
	return ""
}

// String returns param as string
func (p Param) String() (string, error) {
	return string(p), nil
}

// Int returns param as int
func (p Param) Int() (int, error) {
	return strconv.Atoi(string(p))
}

// Int32 returns param as int32
func (p Param) Int32() (int32, error) {
	i, err := strconv.ParseInt(string(p), 10, 32)
	return int32(i), err
}

// Int64 returns param as int64
func (p Param) Int64() (int64, error) {
	return strconv.ParseInt(string(p), 10, 64)
}

// UInt returns param as uint
func (p Param) UInt() (uint, error) {
	i, err := strconv.ParseUint(string(p), 10, strconv.IntSize)
	return uint(i), err
}

// UInt32 returns param as uint32
func (p Param) UInt32() (uint32, error) {
	i, err := strconv.ParseUint(string(p), 10, 32)
	return uint32(i), err
}

// UInt64 returns param as uint64
func (p Param) UInt64() (uint64, error) {
	return strconv.ParseUint(string(p), 10, 64)
}

// Float32 returns param as float32
func (p Param) Float32() (float32, error) {
	f, err := strconv.ParseFloat(string(p), 32)
	return float32(f), err
}

// Float64 returns param as float64
func (p Param) Float64() (float64, error) {
	return strconv.ParseFloat(string(p), 64)
}
