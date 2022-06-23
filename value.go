package zetasqlite

import (
	"bytes"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type Value interface {
	Add(Value) (Value, error)
	Sub(Value) (Value, error)
	Mul(Value) (Value, error)
	Div(Value) (Value, error)
	EQ(Value) (bool, error)
	GT(Value) (bool, error)
	GTE(Value) (bool, error)
	LT(Value) (bool, error)
	LTE(Value) (bool, error)
	ToInt64() (int64, error)
	ToString() (string, error)
	ToFloat64() (float64, error)
	ToBool() (bool, error)
	ToArray() (*ArrayValue, error)
	ToStruct() (*StructValue, error)
}

type IntValue int64

func (iv IntValue) Add(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	return ValueOf(int64(iv) + v2)
}

func (iv IntValue) Sub(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	return ValueOf(int64(iv) - v2)
}

func (iv IntValue) Mul(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	return ValueOf(int64(iv) * v2)
}

func (iv IntValue) Div(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	if v2 == 0 {
		return nil, fmt.Errorf("zero divided error ( %d / 0 )", iv)
	}
	return ValueOf(int64(iv) / v2)
}

func (iv IntValue) EQ(v Value) (bool, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to int64", v)
	}
	return int64(iv) == v2, nil
}

func (iv IntValue) GT(v Value) (bool, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to int64", v)
	}
	return int64(iv) > v2, nil
}

func (iv IntValue) GTE(v Value) (bool, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to int64", v)
	}
	return int64(iv) >= v2, nil
}

func (iv IntValue) LT(v Value) (bool, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to int64", v)
	}
	return int64(iv) < v2, nil
}

func (iv IntValue) LTE(v Value) (bool, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to int64", v)
	}
	return int64(iv) <= v2, nil
}

func (iv IntValue) ToInt64() (int64, error) {
	return int64(iv), nil
}

func (iv IntValue) ToString() (string, error) {
	return fmt.Sprint(iv), nil
}

func (iv IntValue) ToFloat64() (float64, error) {
	return float64(iv), nil
}

func (iv IntValue) ToBool() (bool, error) {
	switch iv {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("falied to convert %d to bool type", iv)
	}
}

func (iv IntValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("falied to convert %d to array type", iv)
}

func (iv IntValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("falied to convert %d to struct type", iv)
}

type StringValue string

func (sv StringValue) Add(v Value) (Value, error) {
	v2, err := v.ToString()
	if err != nil {
		return nil, err
	}
	return ValueOf(string(sv) + v2)
}

func (sv StringValue) Sub(v Value) (Value, error) {
	return nil, fmt.Errorf("sub operation is unsupported for string %v", sv)
}

func (sv StringValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for string %v", sv)
}

func (sv StringValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for string %v", sv)
}

func (sv StringValue) EQ(v Value) (bool, error) {
	v2, err := v.ToString()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to string", v)
	}
	return string(sv) == v2, nil
}

func (sv StringValue) GT(v Value) (bool, error) {
	v2, err := v.ToString()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to string", v)
	}
	return string(sv) > v2, nil
}

func (sv StringValue) GTE(v Value) (bool, error) {
	v2, err := v.ToString()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to string", v)
	}
	return string(sv) >= v2, nil
}

func (sv StringValue) LT(v Value) (bool, error) {
	v2, err := v.ToString()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to string", v)
	}
	return string(sv) < v2, nil
}

func (sv StringValue) LTE(v Value) (bool, error) {
	v2, err := v.ToString()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to string", v)
	}
	return string(sv) <= v2, nil
}

func (sv StringValue) ToInt64() (int64, error) {
	if sv == "" {
		return 0, nil
	}
	return strconv.ParseInt(string(sv), 10, 64)
}

func (sv StringValue) ToString() (string, error) {
	return string(sv), nil
}

func (sv StringValue) ToFloat64() (float64, error) {
	if sv == "" {
		return 0, nil
	}
	return strconv.ParseFloat(string(sv), 64)
}

func (sv StringValue) ToBool() (bool, error) {
	if sv == "" {
		return false, nil
	}
	return strconv.ParseBool(string(sv))
}

func (sv StringValue) ToArray() (*ArrayValue, error) {
	if sv == "" {
		return nil, nil
	}
	return nil, fmt.Errorf("failed to convert array from string: %v", sv)
}

func (sv StringValue) ToStruct() (*StructValue, error) {
	if sv == "" {
		return nil, nil
	}
	return nil, fmt.Errorf("failed to convert struct from string: %v", sv)
}

type FloatValue float64

func (fv FloatValue) Add(v Value) (Value, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return nil, err
	}
	return ValueOf(float64(fv) + v2)
}

func (fv FloatValue) Sub(v Value) (Value, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return nil, err
	}
	return ValueOf(float64(fv) - v2)
}

func (fv FloatValue) Mul(v Value) (Value, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return nil, err
	}
	return ValueOf(float64(fv) * v2)
}

func (fv FloatValue) Div(v Value) (Value, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return nil, err
	}
	if v2 == 0 {
		return nil, fmt.Errorf("zero divided error ( %f / 0 )", fv)
	}
	return ValueOf(float64(fv) / v2)
}

func (fv FloatValue) EQ(v Value) (bool, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to float64", v)
	}
	return float64(fv) == v2, nil
}

func (fv FloatValue) GT(v Value) (bool, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to float64", v)
	}
	return float64(fv) > v2, nil
}

func (fv FloatValue) GTE(v Value) (bool, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to float64", v)
	}
	return float64(fv) >= v2, nil
}

func (fv FloatValue) LT(v Value) (bool, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to float64", v)
	}
	return float64(fv) < v2, nil
}

func (fv FloatValue) LTE(v Value) (bool, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to float64", v)
	}
	return float64(fv) <= v2, nil
}

func (fv FloatValue) ToInt64() (int64, error) {
	return int64(fv), nil
}

func (fv FloatValue) ToString() (string, error) {
	return fmt.Sprint(fv), nil
}

func (fv FloatValue) ToFloat64() (float64, error) {
	return float64(fv), nil
}

func (fv FloatValue) ToBool() (bool, error) {
	return false, fmt.Errorf("falied to convert %f to bool type", fv)
}

func (fv FloatValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert array from float64: %v", fv)
}

func (fv FloatValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert struct from float64: %v", fv)
}

type BoolValue bool

func (bv BoolValue) Add(v Value) (Value, error) {
	return nil, fmt.Errorf("add operation is unsupported for bool %v", bv)
}

func (bv BoolValue) Sub(v Value) (Value, error) {
	return nil, fmt.Errorf("sub operation is unsupported for bool %v", bv)
}

func (bv BoolValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for bool %v", bv)
}

func (bv BoolValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for bool %v", bv)
}

func (bv BoolValue) EQ(v Value) (bool, error) {
	v2, err := v.ToBool()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to bool", v)
	}
	return bool(bv) == v2, nil
}

func (bv BoolValue) GT(v Value) (bool, error) {
	return false, fmt.Errorf("gt operation is unsupported for bool %v", bv)
}

func (bv BoolValue) GTE(v Value) (bool, error) {
	return false, fmt.Errorf("gte operation is unsupported for bool %v", bv)
}

func (bv BoolValue) LT(v Value) (bool, error) {
	return false, fmt.Errorf("lt operation is unsupported for bool %v", bv)
}

func (bv BoolValue) LTE(v Value) (bool, error) {
	return false, fmt.Errorf("lte operation is unsupported for bool %v", bv)
}

func (bv BoolValue) ToInt64() (int64, error) {
	if bv {
		return 1, nil
	}
	return 0, nil
}

func (bv BoolValue) ToString() (string, error) {
	return fmt.Sprint(bv), nil
}

func (bv BoolValue) ToFloat64() (float64, error) {
	if bv {
		return 1, nil
	}
	return 0, nil
}

func (bv BoolValue) ToBool() (bool, error) {
	return bool(bv), nil
}

func (bv BoolValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert bool from array: %v", bv)
}

func (bv BoolValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert bool from struct: %v", bv)
}

type ArrayValue struct {
	values []Value
}

func (av *ArrayValue) Has(v Value) (bool, error) {
	for _, val := range av.values {
		cond, err := val.EQ(v)
		if err != nil {
			return false, err
		}
		if cond {
			return true, nil
		}
	}
	return false, nil
}

func (av *ArrayValue) Add(v Value) (Value, error) {
	return nil, fmt.Errorf("add operation is unsupported for array %v", av)
}

func (av *ArrayValue) Sub(v Value) (Value, error) {
	return nil, fmt.Errorf("sub operation is unsupported for array %v", av)
}

func (av *ArrayValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for array %v", av)
}

func (av *ArrayValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for array %v", av)
}

func (av *ArrayValue) EQ(v Value) (bool, error) {
	arr, err := v.ToArray()
	if err != nil {
		return false, err
	}
	if len(arr.values) != len(av.values) {
		return false, nil
	}
	for idx, value := range av.values {
		cond, err := arr.values[idx].EQ(value)
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (av *ArrayValue) GT(v Value) (bool, error) {
	arr, err := v.ToArray()
	if err != nil {
		return false, err
	}
	if len(arr.values) != len(av.values) {
		return false, nil
	}
	for idx, value := range av.values {
		cond, err := arr.values[idx].GT(value)
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (av *ArrayValue) GTE(v Value) (bool, error) {
	arr, err := v.ToArray()
	if err != nil {
		return false, err
	}
	if len(arr.values) != len(av.values) {
		return false, nil
	}
	for idx, value := range av.values {
		cond, err := arr.values[idx].GTE(value)
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (av *ArrayValue) LT(v Value) (bool, error) {
	arr, err := v.ToArray()
	if err != nil {
		return false, err
	}
	if len(arr.values) != len(av.values) {
		return false, nil
	}
	for idx, value := range av.values {
		cond, err := arr.values[idx].LT(value)
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (av *ArrayValue) LTE(v Value) (bool, error) {
	arr, err := v.ToArray()
	if err != nil {
		return false, err
	}
	if len(arr.values) != len(av.values) {
		return false, nil
	}
	for idx, value := range av.values {
		cond, err := arr.values[idx].LTE(value)
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (av *ArrayValue) ToInt64() (int64, error) {
	return 0, fmt.Errorf("failed to convert int64 from array %v", av)
}

func (av *ArrayValue) ToString() (string, error) {
	return "", fmt.Errorf("failed to convert string from array %v", av)
}

func (av *ArrayValue) ToFloat64() (float64, error) {
	return 0, fmt.Errorf("failed to convert float64 from array %v", av)
}

func (av *ArrayValue) ToBool() (bool, error) {
	return false, fmt.Errorf("failed to convert bool from array %v", av)
}

func (av *ArrayValue) ToArray() (*ArrayValue, error) {
	return av, nil
}

func (av *ArrayValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert struct from array %v", av)
}

type StructValue struct {
	keys   []string
	values []Value
	m      map[string]Value
}

func (sv *StructValue) Add(v Value) (Value, error) {
	return nil, fmt.Errorf("add operation is unsupported for struct %v", sv)
}

func (sv *StructValue) Sub(v Value) (Value, error) {
	return nil, fmt.Errorf("sub operation is unsupported for struct %v", sv)
}

func (sv *StructValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for struct %v", sv)
}

func (sv *StructValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for struct %v", sv)
}

func (sv *StructValue) EQ(v Value) (bool, error) {
	st, err := v.ToStruct()
	if err != nil {
		return false, err
	}
	if len(st.m) != len(sv.m) {
		return false, nil
	}
	for key := range sv.m {
		cond, err := st.m[key].EQ(sv.m[key])
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (sv *StructValue) GT(v Value) (bool, error) {
	st, err := v.ToStruct()
	if err != nil {
		return false, err
	}
	if len(st.m) != len(sv.m) {
		return false, nil
	}
	for key := range sv.m {
		cond, err := st.m[key].GT(sv.m[key])
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (sv *StructValue) GTE(v Value) (bool, error) {
	st, err := v.ToStruct()
	if err != nil {
		return false, err
	}
	if len(st.m) != len(sv.m) {
		return false, nil
	}
	for key := range sv.m {
		cond, err := st.m[key].GTE(sv.m[key])
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (sv *StructValue) LT(v Value) (bool, error) {
	st, err := v.ToStruct()
	if err != nil {
		return false, err
	}
	if len(st.m) != len(sv.m) {
		return false, nil
	}
	for key := range sv.m {
		cond, err := st.m[key].LT(sv.m[key])
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (sv *StructValue) LTE(v Value) (bool, error) {
	st, err := v.ToStruct()
	if err != nil {
		return false, err
	}
	if len(st.m) != len(sv.m) {
		return false, nil
	}
	for key := range sv.m {
		cond, err := st.m[key].LTE(sv.m[key])
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (sv *StructValue) ToInt64() (int64, error) {
	return 0, fmt.Errorf("failed to convert int64 from struct %v", sv)
}

func (sv *StructValue) ToString() (string, error) {
	return "", fmt.Errorf("failed to convert string from struct %v", sv)
}

func (sv *StructValue) ToFloat64() (float64, error) {
	return 0, fmt.Errorf("failed to convert float64 from struct %v", sv)
}

func (sv *StructValue) ToBool() (bool, error) {
	return false, fmt.Errorf("failed to convert bool from struct %v", sv)
}

func (sv *StructValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert array from struct %v", sv)
}

func (sv *StructValue) ToStruct() (*StructValue, error) {
	return sv, nil
}

const (
	ArrayValueHeader  = "zetasqlitearray:"
	StructValueHeader = "zetasqlitestruct:"
)

func ValueOf(v interface{}) (Value, error) {
	switch vv := v.(type) {
	case int:
		return IntValue(int64(vv)), nil
	case int8:
		return IntValue(int64(vv)), nil
	case int16:
		return IntValue(int64(vv)), nil
	case int32:
		return IntValue(int64(vv)), nil
	case int64:
		return IntValue(vv), nil
	case uint:
		return IntValue(int64(vv)), nil
	case uint8:
		return IntValue(int64(vv)), nil
	case uint16:
		return IntValue(int64(vv)), nil
	case uint32:
		return IntValue(int64(vv)), nil
	case uint64:
		return IntValue(int64(vv)), nil
	case string:
		switch {
		case isArrayValue(vv):
			return ArrayValueOf(vv)
		case isStructValue(vv):
			return StructValueOf(vv)
		}
		return StringValue(vv), nil
	case []byte:
		return StringValue(string(vv)), nil
	case float32:
		return FloatValue(float64(vv)), nil
	case float64:
		return FloatValue(vv), nil
	case bool:
		return BoolValue(vv), nil
	}
	return nil, fmt.Errorf("failed to convert value from %T", v)
}

func isArrayValue(v string) bool {
	if len(v) < len(ArrayValueHeader) {
		return false
	}
	if v[0] == '"' {
		return strings.HasPrefix(v[1:], ArrayValueHeader)
	}
	return strings.HasPrefix(v, ArrayValueHeader)
}

func isStructValue(v string) bool {
	if len(v) < len(StructValueHeader) {
		return false
	}
	if v[0] == '"' {
		return strings.HasPrefix(v[1:], StructValueHeader)
	}
	return strings.HasPrefix(v, StructValueHeader)
}

func ArrayValueOf(v string) (Value, error) {
	arr, err := arrayValueFromEncodedString(v)
	if err != nil {
		return nil, fmt.Errorf("failed to get array value from encoded string: %w", err)
	}
	values := make([]Value, 0, len(arr))
	for _, a := range arr {
		val, err := ValueOf(a)
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return &ArrayValue{values: values}, nil
}

func StructValueOf(v string) (Value, error) {
	if len(v) == 0 {
		return nil, nil
	}
	if v[0] == '"' {
		unquoted, err := strconv.Unquote(v)
		if err != nil {
			return nil, fmt.Errorf("failed to unquote value %q: %w", v, err)
		}
		v = unquoted
	}
	content := v[len(StructValueHeader):]
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode for struct value %q: %w", content, err)
	}
	dec := json.NewDecoder(bytes.NewBuffer(decoded))
	t, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to decode struct value %q: %w", decoded, err)
	}
	if t != json.Delim('{') {
		return nil, fmt.Errorf("invalid delimiter of struct value %q", decoded)
	}
	var (
		keys   []string
		values []Value
		valMap = map[string]Value{}
	)
	for {
		k, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to decode struct key %q: %w", decoded, err)
		}
		if k == json.Delim('}') {
			break
		}
		key := k.(string)
		var value interface{}
		if err := dec.Decode(&value); err != nil {
			return nil, fmt.Errorf("failed to decode struct value %q: %w", decoded, err)
		}
		keys = append(keys, key)
		val, err := ValueOf(value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert value from %v: %w", value, err)
		}
		values = append(values, val)
		valMap[key] = val
	}
	return &StructValue{keys: keys, values: values, m: valMap}, nil
}

func SQLiteValue(v interface{}) (interface{}, error) {
	rv := reflect.TypeOf(v)
	switch rv.Kind() {
	case reflect.Int:
		return int64(v.(int)), nil
	case reflect.Int8:
		return int64(v.(int8)), nil
	case reflect.Int16:
		return int64(v.(int16)), nil
	case reflect.Int32:
		return int64(v.(int32)), nil
	case reflect.Uint:
		return int64(v.(uint)), nil
	case reflect.Uint8:
		return int64(v.(uint8)), nil
	case reflect.Uint16:
		return int64(v.(uint16)), nil
	case reflect.Uint32:
		return int64(v.(uint32)), nil
	case reflect.Uint64:
		return int64(v.(uint64)), nil
	case reflect.Float32:
		return float64(v.(float32)), nil
	case reflect.Slice:
		if rv.Elem().Kind() == reflect.Uint8 {
			return string(v.([]byte)), nil
		}
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to encode value %v: %w", v, err)
		}
		return toArrayValueFromJSONString(string(b)), nil
	case reflect.Array:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to encode value %v: %w", v, err)
		}
		return toArrayValueFromJSONString(string(b)), nil
	case reflect.Struct:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to encode value %v: %w", v, err)
		}
		return toStructValueFromJSONString(string(b)), nil
	}
	return v, nil
}

func toArrayValueFromJSONString(json string) string {
	return strconv.Quote(
		fmt.Sprintf(
			"%s%s",
			ArrayValueHeader,
			base64.StdEncoding.EncodeToString([]byte(json)),
		),
	)
}

func arrayValueFromEncodedString(v string) ([]interface{}, error) {
	if len(v) == 0 {
		return nil, nil
	}
	if v[0] == '"' {
		unquoted, err := strconv.Unquote(v)
		if err != nil {
			return nil, fmt.Errorf("failed to unquote value %q: %w", v, err)
		}
		v = unquoted
	}
	content := v[len(ArrayValueHeader):]
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode for array value %q: %w", content, err)
	}
	var arr []interface{}
	if err := json.Unmarshal(decoded, &arr); err != nil {
		return nil, fmt.Errorf("failed to decode array: %w", err)
	}
	return arr, nil
}

func jsonArrayFromEncodedString(v string) ([]byte, error) {
	if len(v) == 0 {
		return nil, nil
	}
	if v[0] == '"' {
		unquoted, err := strconv.Unquote(v)
		if err != nil {
			return nil, fmt.Errorf("failed to unquote value %q: %w", v, err)
		}
		v = unquoted
	}
	content := v[len(ArrayValueHeader):]
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode for array value %q: %w", content, err)
	}
	return decoded, nil
}

func toStructValueFromJSONString(json string) string {
	return strconv.Quote(
		fmt.Sprintf(
			"%s%s",
			StructValueHeader,
			base64.StdEncoding.EncodeToString([]byte(json)),
		),
	)
}

func isNULLValue(v interface{}) bool {
	vv, ok := v.([]byte)
	if !ok {
		return false
	}
	return len(vv) == 0
}

func convertNamedValues(v []driver.NamedValue) ([]driver.NamedValue, error) {
	ret := make([]driver.NamedValue, 0, len(v))
	for _, vv := range v {
		converted, err := convertNamedValue(vv)
		if err != nil {
			return nil, fmt.Errorf("failed to convert value from %+v: %w", vv, err)
		}
		ret = append(ret, converted)
	}
	return ret, nil
}

func convertNamedValue(v driver.NamedValue) (driver.NamedValue, error) {
	value, err := SQLiteValue(v.Value)
	if err != nil {
		return driver.NamedValue{}, err
	}
	return driver.NamedValue{
		Name:    strings.ToLower(v.Name),
		Ordinal: v.Ordinal,
		Value:   value,
	}, nil
}

func convertValues(v []driver.Value) ([]driver.Value, error) {
	ret := make([]driver.Value, 0, len(v))
	for _, vv := range v {
		value, err := SQLiteValue(vv)
		if err != nil {
			return nil, err
		}
		ret = append(ret, value)
	}
	return ret, nil
}