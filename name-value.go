// Package namevalue is a package to handle name-value pairs
//
//	Author: Elizalde G. Baguinon
//	Created: October 17, 2021
package namevalue

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	ssd "github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

type (
	// NameValue is a struct to contain a name-value pair
	NameValue[T any] struct {
		Name  string `json:"name,omitempty"`
		Value T      `json:"value,omitempty"`
	}
	// NameValues is a struct to manage value structs
	NameValues struct {
		Pair     map[string]any
		prepared bool
	}
)

const (
	INTERPOLATE_PATTERN string = `\$\{(\w*)\}` // search for ${*}
)

func (nvp *NameValues) prepare() {
	np := make(map[string]any)
	for n := range nvp.Pair {
		ln := strings.ToLower(n)
		np[ln] = nvp.Pair[n]
		delete(nvp.Pair, n)
	}
	for n := range np {
		nvp.Pair[n] = np[n]
	}
	nvp.prepared = true
}

// Exists checks if the key or name exists. It returns the index of the element if found, -1 if not found.
func (nvp *NameValues) Exists(name string) bool {
	if !nvp.prepared {
		nvp.prepare()
	}
	name = strings.ToLower(name)
	_, exists := nvp.Pair[name]
	return exists
}

// NameValueGet gets the value from the collection of NameValues by name
//
// This function requires version 1.18+
func NameValueGet[T constraints.Ordered | bool](nvs NameValues, name string) T {
	if !nvs.prepared {
		nvs.prepare()
	}

	name = strings.ToLower(name)
	tmp := nvs.Pair[name]

	tpt := any(*new(T))
	value := *new(T)

	// If the value is a string and the inferred type is otherwise
	// try to convert, else just convery via inferred type
	switch t := tmp.(type) {
	case string:
		switch tpt.(type) {
		case int:
			val, _ := strconv.ParseInt(t, 10, 32)
			value = any(int(val)).(T)
		case int64:
			val, _ := strconv.ParseInt(t, 10, 64)
			value = any(val).(T)
		case bool:
			val, _ := strconv.ParseBool(t)
			value = any(val).(T)
		case float32:
			val, _ := strconv.ParseFloat(t, 32)
			value = any(val).(T)
		case float64:
			val, _ := strconv.ParseFloat(t, 64)
			value = any(val).(T)
		default:
			if tmp != nil {
				value = tmp.(T)
			} else {
				value = getZero[T]()
			}
		}
	default:
		if t != nil {
			value = t.(T)
		} else {
			value = getZero[T]()
		}
	}

	return value
}

// NameValueGetPtr gets the value from the collection of NameValues by name as pointer
//
// This function requires version 1.18+
func NameValueGetPtr[T constraints.Ordered | bool](nvs NameValues, name string) *T {
	value := NameValueGet[T](nvs, name)
	return &value
}

// String returns the name value as string. The second result returns the existence.
func (nvp *NameValues) String(name string) (string, bool) {
	if !nvp.prepared {
		nvp.prepare()
	}
	var (
		tmp          any
		exists, conv bool
		str, val     string
	)
	name = strings.ToLower(name)
	tmp, exists = nvp.Pair[name]
	if !exists {
		return str, exists
	}
	val, conv = tmp.(string)
	if !conv {
		return val, exists
	}
	return val, exists
}

// Strings returns the values as a string array.
// If the name does not exist, this function will return an empty string array
// If the value is comma-separated, all elements delimited by the comma will be returned as an array
func (nvp *NameValues) Strings(name string) []string {
	value, exists := nvp.String(name)
	if !exists {
		return []string{}
	}
	if strings.Contains(value, ",") {
		return strings.Split(value, ",")
	}
	var result [1]string
	result[0] = value
	return result[:]
}

// Int returns the name value as int. The second result returns the existence.
func (nvp *NameValues) Int(name string) (int, bool) {
	if !nvp.prepared {
		nvp.prepare()
	}
	var (
		conv, exists bool
		tmp          any
		val          int
		str          string
		err          error
	)
	name = strings.ToLower(name)
	// Check if the key exists in the map
	tmp, exists = nvp.Pair[name]
	if !exists {
		return val, exists
	}
	// Attempt to convert the interface to int
	val, conv = tmp.(int)
	if conv {
		return val, exists
	}
	// If it did not succeed, try converting to string first
	// and manually convert to int
	str, conv = tmp.(string)
	if !conv {
		return val, exists
	}
	val, err = strconv.Atoi(str)
	if err != nil {
		return val, exists
	}
	return val, exists
}

// Ints returns the values as an int array
// If the name does not exist, this function will return an empty int array
func (nvp *NameValues) Ints(name string) []int {
	value, exists := nvp.Int(name)
	if !exists {
		return []int{}
	}
	var result [1]int
	result[0] = value
	return result[:]
}

// Int64 returns the name value as int64. The second result returns the existence.
func (nvp *NameValues) Int64(name string) (int64, bool) {
	if !nvp.prepared {
		nvp.prepare()
	}
	var (
		conv, exists bool
		tmp          any
		val          int64
		str          string
		err          error
	)
	name = strings.ToLower(name)
	// Check if the key exists in the map
	tmp, exists = nvp.Pair[name]
	if !exists {
		return val, exists
	}
	// Attempt to convert the interface to int64
	val, conv = tmp.(int64)
	if conv {
		return val, exists
	}
	// If it did not succeed, try converting to string first
	// and manually convert to int64
	str, conv = tmp.(string)
	if !conv {
		return val, exists
	}
	val, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		return val, exists
	}
	return val, exists
}

// Int64s returns the values as an int64 array
// If the name does not exist, this function will return an empty int64 array
func (nvp *NameValues) Int64s(name string) []int64 {
	value, exists := nvp.Int64(name)
	if !exists {
		return []int64{}
	}
	var result [1]int64
	result[0] = value
	return result[:]
}

// Plain returns the name value as interface{}. The second result returns the existence.
func (nvp *NameValues) Plain(name string) (interface{}, bool) {
	if !nvp.prepared {
		nvp.prepare()
	}
	name = strings.ToLower(name)
	tmp, exists := nvp.Pair[name]
	return tmp, exists
}

// Bool returns the name value as boolean. It automatically convers 'true', 'yes', '1', '-1' and 'on' to boolean.
// The second result returns the existence.
func (nvp *NameValues) Bool(name string) (bool, bool) {
	value, exists := nvp.String(name)
	if !exists {
		return exists, exists
	}
	return (value == "true" || value == "yes" || value == "1" || value == "-1" || value == "on"), true
}

// Bools returns the values as a boolean array
// If the name does not exist, this function will return an empty boolean array
func (nvp *NameValues) Bools(name string) []bool {
	value, exists := nvp.Bool(name)
	if !exists {
		return []bool{}
	}
	var result [1]bool
	result[0] = value
	return result[:]
}

// Float64 returns the name value as float64. The second result returns the existence.
func (nvp *NameValues) Float64(name string) (float64, bool) {
	if !nvp.prepared {
		nvp.prepare()
	}
	var (
		conv, exists bool
		tmp          any
		val          float64
		str          string
		err          error
	)
	name = strings.ToLower(name)
	// Check if the key exists in the map
	tmp, exists = nvp.Pair[name]
	if !exists {
		return val, exists
	}
	// Attempt to convert the interface to float64
	val, conv = tmp.(float64)
	if conv {
		return val, exists
	}
	// If it did not succeed, try converting to string first
	// and manually convert to int64
	str, conv = tmp.(string)
	if !conv {
		return val, exists
	}
	val, err = strconv.ParseFloat(str, 64)
	if err != nil {
		return val, exists
	}
	return val, exists
}

// Float64s returns the values as a float64 array
// If the name does not exist, this function will return an empty float64 array
func (nvp *NameValues) Float64s(name string) []float64 {
	value, exists := nvp.Float64(name)
	if !exists {
		return []float64{}
	}
	var result [1]float64
	result[0] = value
	return result[:]
}

// Decimal returns the name value as shopspring.Decimal. The second result returns the existence.
func (nvp *NameValues) Decimal(name string) (ssd.Decimal, bool) {
	if !nvp.prepared {
		nvp.prepare()
	}
	var (
		exists bool
		tmp    any
		val    ssd.Decimal
		err    error
	)

	name = strings.ToLower(name)
	tmp, exists = nvp.Pair[name]
	if !exists {
		return val, exists
	}
	// Try getting if this is really a decimal
	// Then try as a string again
	switch t := tmp.(type) {
	case string:
		t = strings.ReplaceAll(t, ",", "")
		t = strings.ReplaceAll(t, " ", "")
		val, err = ssd.NewFromString(t)
		if err != nil {
			return val, exists
		}
	case int:
		val = ssd.NewFromInt(int64(t))
		return val, exists
	case int64:
		val = ssd.NewFromInt(t)
		return val, exists
	case float32:
		val = ssd.NewFromFloat(float64(t))
		return val, exists
	case float64:
		val = ssd.NewFromFloat(t)
		return val, exists
	case ssd.Decimal:
		val = t
		return val, exists
	}

	return val, exists
}

// Decimals returns the values as a decimal array
// If the name does not exist, this function will return an empty decimal array
func (nvp *NameValues) Decimals(name string) []ssd.Decimal {
	value, exists := nvp.Decimal(name)
	if !exists {
		return []ssd.Decimal{}
	}
	var result [1]ssd.Decimal
	result[0] = value
	return result[:]
}

// **************************************************************
//   Pointer outputs
// **************************************************************

// PtrString returns the name value as pointer to string. The second result returns the existence.
func (nvp *NameValues) PtrString(name string) (*string, bool) {
	value, exists := nvp.String(name)
	if !exists {
		return nil, exists
	}
	return &value, exists
}

// PtrInt returns the name value as pointer to int. The second result returns the existence.
func (nvp *NameValues) PtrInt(name string) (*int, bool) {
	value, exists := nvp.Int(name)
	if !exists {
		return nil, exists
	}
	return &value, exists
}

// PtrInt64 returns the name value as pointer to int64. The second result returns the existence.
func (nvp *NameValues) PtrInt64(name string) (*int64, bool) {
	value, exists := nvp.Int64(name)
	if !exists {
		return nil, exists
	}
	return &value, exists
}

// PtrPlain returns the name value as pointer to interface{}. The second result returns the existence.
func (nvp *NameValues) PtrPlain(name string) (*interface{}, bool) {
	value, exists := nvp.Plain(name)
	if !exists {
		return nil, exists
	}
	return &value, exists
}

// PtrBool returns the name value as pointer to bool. The second result returns the existence.
func (nvp *NameValues) PtrBool(name string) (*bool, bool) {
	value, exists := nvp.Bool(name)
	if !exists {
		return nil, exists
	}
	return &value, exists
}

// PtrFloat64 returns the name value as pointer to int64. The second result returns the existence.
func (nvp *NameValues) PtrFloat64(name string) (*float64, bool) {
	value, exists := nvp.Float64(name)
	if !exists {
		return nil, exists
	}
	return &value, exists
}

// PtrFloat64 returns the name value as pointer to int64. The second result returns the existence.
func (nvp *NameValues) PtrDecimal(name string) (*ssd.Decimal, bool) {
	value, exists := nvp.Decimal(name)
	if !exists {
		return nil, exists
	}
	return &value, exists
}

// ToInterfaceArray converts name values to interface array
func (nvp *NameValues) ToInterfaceArray() []interface{} {
	return NameValuesToInterfaceArray(*nvp)
}

// Interpolate - interpolate string with values from with base string
func (nvp *NameValues) Interpolate(base string) (string, []interface{}) {
	return Interpolate(base, *nvp)
}

// SortByKey sort name values by key order array
func (nvp *NameValues) SortByKey(keyOrder *[]string) NameValues {
	return SortByKey(nvp, keyOrder)
}

// **************************************************************
//   Miscellaneous functions
// **************************************************************

// Interpolate interpolates string with the name value pairs
func Interpolate(base string, nv NameValues) (string, []any) {
	var (
		val  any
		sval string
	)

	nstr := base
	re := regexp.MustCompile(INTERPOLATE_PATTERN)
	matches := re.FindAllString(base, -1)
	vals := make([]any, len(matches))
	for i, match := range matches {
		val = "0"
		sval = "0"
		for n, v := range nv.Pair {
			if strings.EqualFold(match, `${`+n+`}`) {
				sval = anyToStr(v)
				val = v
				break
			}
		}
		nstr = strings.Replace(nstr, match, sval, -1)
		vals[i] = val //a string 0 would cater to both string and number columns
	}
	return nstr, vals
}

// NameValuesToInterfaceArray converts name values to interface array
func NameValuesToInterfaceArray(values NameValues) []interface{} {
	args := make([]interface{}, len(values.Pair))
	i := 0
	for _, v := range values.Pair {
		args[i] = v
		i++
	}
	return args
}

// SortByKey reorders keys and values based on a keyOrder array sequence
func SortByKey(values *NameValues, keyOrder *[]string) NameValues {
	if keyOrder == nil {
		return *values
	}
	ko := *keyOrder
	if len(ko) == 0 {
		return *values
	}
	ret := NameValues{
		Pair: make(map[string]any),
	}
	for i := 0; i < len(ko); i++ {
		for k, v := range values.Pair {
			if strings.EqualFold(ko[i], k) {
				ret.Pair[k] = v
				break
			}
		}
	}
	return ret
}

func getZero[T constraints.Ordered | bool]() T {
	var r T
	return r
}

// anyToStr converts any variable to string
func anyToStr(value interface{}) string {
	var b string
	if value == nil {
		return ""
	}
	switch t := value.(type) {
	case string:
		b = t
	case int:
		b = strconv.FormatInt(int64(t), 10)
	case int8:
		b = strconv.FormatInt(int64(t), 10)
	case int16:
		b = strconv.FormatInt(int64(t), 10)
	case int32:
		b = strconv.FormatInt(int64(t), 10)
	case int64:
		b = strconv.FormatInt(t, 10)
	case uint:
		b = strconv.FormatUint(uint64(t), 10)
	case uint8:
		b = strconv.FormatUint(uint64(t), 10)
	case uint16:
		b = strconv.FormatUint(uint64(t), 10)
	case uint32:
		b = strconv.FormatUint(uint64(t), 10)
	case uint64:
		b = strconv.FormatUint(uint64(t), 10)
	case float32:
		b = fmt.Sprintf("%f", t)
	case float64:
		b = fmt.Sprintf("%f", t)
	case bool:
		if t {
			return "true"
		} else {
			return "false"
		}
	case time.Time:
		b = "'" + t.Format(time.RFC3339) + "'"
	case *string:
		if t == nil {
			return ""
		}
		b = *t
	case *int:
		if t == nil {
			return "0"
		}
		b = strconv.FormatInt(int64(*t), 10)
	case *int8:
		if t == nil {
			return "0"
		}
		b = strconv.FormatInt(int64(*t), 10)
	case *int16:
		if t == nil {
			return "0"
		}
		b = strconv.FormatInt(int64(*t), 10)
	case *int32:
		if t == nil {
			return "0"
		}
		b = strconv.FormatInt(int64(*t), 10)
	case *int64:
		if t == nil {
			return "0"
		}
		b = strconv.FormatInt(*t, 10)
	case *uint:
		if t == nil {
			return "0"
		}
		b = strconv.FormatUint(uint64(*t), 10)
	case *uint8:
		if t == nil {
			return "0"
		}
		b = strconv.FormatUint(uint64(*t), 10)
	case *uint16:
		if t == nil {
			return "0"
		}
		b = strconv.FormatUint(uint64(*t), 10)
	case *uint32:
		if t == nil {
			return "0"
		}
		b = strconv.FormatUint(uint64(*t), 10)
	case *uint64:
		if t == nil {
			return "0"
		}
		b = strconv.FormatUint(uint64(*t), 10)
	case *float32:
		if t == nil {
			return "0"
		}
		b = fmt.Sprintf("%f", *t)
	case *float64:
		if t == nil {
			return "0"
		}
		b = fmt.Sprintf("%f", *t)
	case *bool:
		if t == nil || !*t {
			return "false"
		}
		return "true"
	case *time.Time:
		if t == nil {
			return "'" + time.Time{}.Format(time.RFC3339) + "'"
		}
		tm := *t
		b = "'" + tm.Format(time.RFC3339) + "'"
	}

	return b
}
