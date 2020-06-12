package changeset

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Fs02/changeset/params"
	"github.com/stretchr/testify/assert"
)

func ExampleCast() {
	type User struct {
		ID   int
		Name string
	}

	user := User{}
	input := params.Map{
		"id":   1,
		"name": "name",
	}

	ch := Cast(user, input, []string{"name"})
	fmt.Println(ch.Changes())
	// Output: map[name:name]
}

func ExampleCast_invalidType() {
	type User struct {
		ID   int
		Name string
	}

	user := User{}
	input := params.Map{
		"id":   1,
		"name": true,
	}

	ch := Cast(user, input, []string{"name"})
	fmt.Println(ch.Error())
	// Output: name is invalid
}

func ExampleCast_invalidTypeWithCustomError() {
	type User struct {
		ID   int
		Name string
	}

	user := User{}
	input := params.Map{
		"id":   1,
		"name": true,
	}

	ch := Cast(user, input, []string{"name"}, Message("{field} tidak valid"))
	fmt.Println(ch.Error())
	// Output: name tidak valid
}

func TestCast(t *testing.T) {
	var data struct {
		Field1 int `db:"field1"`
		Field2 string
		Field3 bool
		Field4 bool `db:"-"`
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": true,
		"field4": true,
		"field5": "ignore please",
	}

	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": true,
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf(false),
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
		"field3": false,
	}

	ch := Cast(data, input, []string{"field1", "field2", "field3", "field4"})
	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)

	ch = Cast(&data, input, []string{"field1", "field2", "field3", "field4"})
	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)

	assert.NotNil(t, ch.Changes())
	assert.NotNil(t, ch.Values())
}

func TestCast_existingChangeset(t *testing.T) {
	var data struct {
		Field1 int `db:"field1"`
		Field2 string
		Field3 bool
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": true,
		"field4": "ignore please",
	}

	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": true,
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf(false),
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
		"field3": false,
	}

	ch := Cast(data, input, []string{})
	assert.Nil(t, ch.Errors())
	assert.Equal(t, 0, len(ch.Changes()))

	ch = Cast(ch, input, []string{"field1", "field2"})
	assert.Nil(t, ch.Errors())
	assert.Equal(t, 2, len(ch.Changes()))

	ch = Cast(*ch, input, []string{"field1", "field3"})
	assert.Nil(t, ch.Errors())
	assert.Equal(t, 3, len(ch.Changes()))
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
}

func TestCast_unchanged(t *testing.T) {
	data := struct {
		Field1 int `db:"field1"`
		Field2 string
		Field3 bool
		Field4 *bool
	}{
		Field1: 1,
		Field2: "1",
		Field3: true,
		Field4: nil,
	}

	input := params.Map{
		"field1": 1,
		"field2": "1",
		"field3": true,
		"field4": nil,
	}

	expectedChanges := map[string]interface{}{}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf(false),
		"field4": reflect.TypeOf(false),
	}

	expectedValues := map[string]interface{}{
		"field1": 1,
		"field2": "1",
		"field3": true,
	}

	ch := Cast(data, input, []string{"field1", "field2", "field3", "field4"})
	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
}

func TestCast_error(t *testing.T) {
	var data struct {
		Field1 int
	}

	input := params.Map{
		"field1": "1",
	}

	ch := Cast(data, input, []string{"field1"})
	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field1 is invalid", ch.Error().Error())
}

func TestCast_panic(t *testing.T) {
	input := params.Map{
		"field1": "1",
	}

	assert.Panics(t, func() {
		Cast("data", input, []string{"field1"})
	})
}

var input = params.Map{
	"field1":  true,
	"field2":  2,
	"field3":  3,
	"field4":  4,
	"field5":  5,
	"field6":  6,
	"field7":  7,
	"field8":  8,
	"field9":  9,
	"field10": 10,
	"field11": 11,
	"field12": 12,
	"field13": 13,
	"field14": 14,
	"field15": "15",
}

var expectedChanges = map[string]interface{}{
	"field1":  true,
	"field2":  2,
	"field3":  int8(3),
	"field4":  int16(4),
	"field5":  int32(5),
	"field6":  int64(6),
	"field7":  uint(7),
	"field8":  uint8(8),
	"field9":  uint16(9),
	"field10": uint32(10),
	"field11": uint64(11),
	"field12": uintptr(12),
	"field13": float32(13),
	"field14": float64(14),
	"field15": "15",
}

var expectedValues = map[string]interface{}{
	"field1":  false,
	"field2":  0,
	"field3":  int8(0),
	"field4":  int16(0),
	"field5":  int32(0),
	"field6":  int64(0),
	"field7":  uint(0),
	"field8":  uint8(0),
	"field9":  uint16(0),
	"field10": uint32(0),
	"field11": uint64(0),
	"field12": uintptr(0),
	"field13": float32(0),
	"field14": float64(0),
	"field15": "",
}

var expectedTypes = map[string]reflect.Type{
	"field1":  reflect.TypeOf(false),
	"field2":  reflect.TypeOf(0),
	"field3":  reflect.TypeOf(int8(0)),
	"field4":  reflect.TypeOf(int16(0)),
	"field5":  reflect.TypeOf(int32(0)),
	"field6":  reflect.TypeOf(int64(0)),
	"field7":  reflect.TypeOf(uint(0)),
	"field8":  reflect.TypeOf(uint8(0)),
	"field9":  reflect.TypeOf(uint16(0)),
	"field10": reflect.TypeOf(uint32(0)),
	"field11": reflect.TypeOf(uint64(0)),
	"field12": reflect.TypeOf(uintptr(0)),
	"field13": reflect.TypeOf(float32(0)),
	"field14": reflect.TypeOf(float64(0)),
	"field15": reflect.TypeOf(""),
}

func TestCast_basic(t *testing.T) {
	var data struct {
		Field1  bool
		Field2  int
		Field3  int8
		Field4  int16
		Field5  int32
		Field6  int64
		Field7  uint
		Field8  uint8
		Field9  uint16
		Field10 uint32
		Field11 uint64
		Field12 uintptr
		Field13 float32
		Field14 float64
		Field15 string
	}

	ch := Cast(data, input, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedTypes, ch.types)
}

func TestCast_basicWithValue(t *testing.T) {
	data := struct {
		Field1  bool
		Field2  int
		Field3  int8
		Field4  int16
		Field5  int32
		Field6  int64
		Field7  uint
		Field8  uint8
		Field9  uint16
		Field10 uint32
		Field11 uint64
		Field12 uintptr
		Field13 float32
		Field14 float64
		Field15 string
	}{
		false,
		1,
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		uintptr(1),
		float32(1),
		float64(1),
		"1",
	}

	expectedValues := map[string]interface{}{
		"field1":  false,
		"field2":  1,
		"field3":  int8(1),
		"field4":  int16(1),
		"field5":  int32(1),
		"field6":  int64(1),
		"field7":  uint(1),
		"field8":  uint8(1),
		"field9":  uint16(1),
		"field10": uint32(1),
		"field11": uint64(1),
		"field12": uintptr(1),
		"field13": float32(1),
		"field14": float64(1),
		"field15": "1",
	}

	ch := Cast(data, input, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedTypes, ch.types)
}

func TestCast_basicZeroValue(t *testing.T) {
	data := struct {
		Field1  bool
		Field2  int
		Field3  int8
		Field4  int16
		Field5  int32
		Field6  int64
		Field7  uint
		Field8  uint8
		Field9  uint16
		Field10 uint32
		Field11 uint64
		Field12 uintptr
		Field13 float32
		Field14 float64
		Field15 string
	}{}

	var input = params.Map{
		"field1":  false,
		"field2":  0,
		"field3":  0,
		"field4":  0,
		"field5":  0,
		"field6":  0,
		"field7":  0,
		"field8":  0,
		"field9":  0,
		"field10": 0,
		"field11": 0,
		"field12": 0,
		"field13": 0,
		"field14": 0,
		"field15": "",
	}

	expectedValues := map[string]interface{}{
		"field1":  false,
		"field2":  0,
		"field3":  int8(0),
		"field4":  int16(0),
		"field5":  int32(0),
		"field6":  int64(0),
		"field7":  uint(0),
		"field8":  uint8(0),
		"field9":  uint16(0),
		"field10": uint32(0),
		"field11": uint64(0),
		"field12": uintptr(0),
		"field13": float32(0),
		"field14": float64(0),
		"field15": "",
	}

	// Empty string is ignored because it's default of empty values
	expectedChanges := map[string]interface{}{
		"field1":  false,
		"field2":  0,
		"field3":  int8(0),
		"field4":  int16(0),
		"field5":  int32(0),
		"field6":  int64(0),
		"field7":  uint(0),
		"field8":  uint8(0),
		"field9":  uint16(0),
		"field10": uint32(0),
		"field11": uint64(0),
		"field12": uintptr(0),
		"field13": float32(0),
		"field14": float64(0),
	}

	ch := Cast(data, input, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedTypes, ch.types)
}

func TestCast_ptr(t *testing.T) {
	var data struct {
		Field1  *bool
		Field2  *int
		Field3  *int8
		Field4  *int16
		Field5  *int32
		Field6  *int64
		Field7  *uint
		Field8  *uint8
		Field9  *uint16
		Field10 *uint32
		Field11 *uint64
		Field12 *uintptr
		Field13 *float32
		Field14 *float64
		Field15 *string
	}

	ch := Cast(data, input, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, map[string]interface{}{}, ch.values)
	assert.Equal(t, expectedTypes, ch.types)
}

func TestCast_ptrWithValue(t *testing.T) {
	var (
		vbool    = false
		vint     = int(1)
		vint8    = int8(1)
		vint16   = int16(1)
		vint32   = int32(1)
		vint64   = int64(1)
		vuint    = uint(1)
		vuint8   = uint8(1)
		vuint16  = uint16(1)
		vuint32  = uint32(1)
		vuint64  = uint64(1)
		vuintptr = uintptr(1)
		vfloat32 = float32(1)
		vfloat64 = float64(1)
		vstring  = "1"
	)
	data := struct {
		Field1  *bool
		Field2  *int
		Field3  *int8
		Field4  *int16
		Field5  *int32
		Field6  *int64
		Field7  *uint
		Field8  *uint8
		Field9  *uint16
		Field10 *uint32
		Field11 *uint64
		Field12 *uintptr
		Field13 *float32
		Field14 *float64
		Field15 *string
	}{
		&vbool,
		&vint,
		&vint8,
		&vint16,
		&vint32,
		&vint64,
		&vuint,
		&vuint8,
		&vuint16,
		&vuint32,
		&vuint64,
		&vuintptr,
		&vfloat32,
		&vfloat64,
		&vstring,
	}

	expectedValues := map[string]interface{}{
		"field1":  false,
		"field2":  1,
		"field3":  int8(1),
		"field4":  int16(1),
		"field5":  int32(1),
		"field6":  int64(1),
		"field7":  uint(1),
		"field8":  uint8(1),
		"field9":  uint16(1),
		"field10": uint32(1),
		"field11": uint64(1),
		"field12": uintptr(1),
		"field13": float32(1),
		"field14": float64(1),
		"field15": "1",
	}

	ch := Cast(data, input, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedTypes, ch.types)
}

func TestCast_ptrWithNilValue(t *testing.T) {
	var (
		vbool    = true
		vint     = int(1)
		vint8    = int8(1)
		vint16   = int16(1)
		vint32   = int32(1)
		vint64   = int64(1)
		vuint    = uint(1)
		vuint8   = uint8(1)
		vuint16  = uint16(1)
		vuint32  = uint32(1)
		vuint64  = uint64(1)
		vuintptr = uintptr(1)
		vfloat32 = float32(1)
		vfloat64 = float64(1)
		vstring  = "1"
	)
	data := struct {
		Field1  *bool
		Field2  *int
		Field3  *int8
		Field4  *int16
		Field5  *int32
		Field6  *int64
		Field7  *uint
		Field8  *uint8
		Field9  *uint16
		Field10 *uint32
		Field11 *uint64
		Field12 *uintptr
		Field13 *float32
		Field14 *float64
		Field15 *string
	}{
		&vbool,
		&vint,
		&vint8,
		&vint16,
		&vint32,
		&vint64,
		&vuint,
		&vuint8,
		&vuint16,
		&vuint32,
		&vuint64,
		&vuintptr,
		&vfloat32,
		&vfloat64,
		&vstring,
	}

	input := params.Map{
		"field1":  nil,
		"field2":  nil,
		"field3":  nil,
		"field4":  nil,
		"field5":  nil,
		"field6":  nil,
		"field7":  nil,
		"field8":  nil,
		"field9":  nil,
		"field10": nil,
		"field11": nil,
		"field12": nil,
		"field13": nil,
		"field14": nil,
		"field15": nil,
	}

	expectedChanges := map[string]interface{}{
		"field1":  nil,
		"field2":  nil,
		"field3":  nil,
		"field4":  nil,
		"field5":  nil,
		"field6":  nil,
		"field7":  nil,
		"field8":  nil,
		"field9":  nil,
		"field10": nil,
		"field11": nil,
		"field12": nil,
		"field13": nil,
		"field14": nil,
		"field15": nil,
	}

	expectedValues := map[string]interface{}{
		"field1":  true,
		"field2":  1,
		"field3":  int8(1),
		"field4":  int16(1),
		"field5":  int32(1),
		"field6":  int64(1),
		"field7":  uint(1),
		"field8":  uint8(1),
		"field9":  uint16(1),
		"field10": uint32(1),
		"field11": uint64(1),
		"field12": uintptr(1),
		"field13": float32(1),
		"field14": float64(1),
		"field15": "1",
	}

	ch := Cast(data, input, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedChanges, ch.Changes())
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedTypes, ch.types)
}

var sliceParams = params.Map{
	"field1":  []bool{true},
	"field2":  []int{2},
	"field3":  []int8{3},
	"field4":  []int16{4},
	"field5":  []int32{5},
	"field6":  []int64{6},
	"field7":  []uint{7},
	"field8":  []uint8{8},
	"field9":  []uint16{9},
	"field10": []uint32{10},
	"field11": []uint64{11},
	"field12": []uintptr{12},
	"field13": []float32{13},
	"field14": []float64{14},
	"field15": []string{"15"},
}

var sliceExpectedChanges = map[string]interface{}{
	"field1":  []bool{true},
	"field2":  []int{2},
	"field3":  []int8{3},
	"field4":  []int16{4},
	"field5":  []int32{5},
	"field6":  []int64{6},
	"field7":  []uint{7},
	"field8":  []uint8{8},
	"field9":  []uint16{9},
	"field10": []uint32{10},
	"field11": []uint64{11},
	"field12": []uintptr{12},
	"field13": []float32{13},
	"field14": []float64{14},
	"field15": []string{"15"},
}

var sliceExpectedTypes = map[string]reflect.Type{
	"field1":  reflect.TypeOf([]bool{}),
	"field2":  reflect.TypeOf([]int{}),
	"field3":  reflect.TypeOf([]int8{}),
	"field4":  reflect.TypeOf([]int16{}),
	"field5":  reflect.TypeOf([]int32{}),
	"field6":  reflect.TypeOf([]int64{}),
	"field7":  reflect.TypeOf([]uint{}),
	"field8":  reflect.TypeOf([]uint8{}),
	"field9":  reflect.TypeOf([]uint16{}),
	"field10": reflect.TypeOf([]uint32{}),
	"field11": reflect.TypeOf([]uint64{}),
	"field12": reflect.TypeOf([]uintptr{}),
	"field13": reflect.TypeOf([]float32{}),
	"field14": reflect.TypeOf([]float64{}),
	"field15": reflect.TypeOf([]string{}),
}

func TestCast_slice(t *testing.T) {
	var data struct {
		Field1  []bool
		Field2  []int
		Field3  []int8
		Field4  []int16
		Field5  []int32
		Field6  []int64
		Field7  []uint
		Field8  []uint8
		Field9  []uint16
		Field10 []uint32
		Field11 []uint64
		Field12 []uintptr
		Field13 []float32
		Field14 []float64
		Field15 []string
	}

	expectedValues := map[string]interface{}{
		"field1":  []bool(nil),
		"field2":  []int(nil),
		"field3":  []int8(nil),
		"field4":  []int16(nil),
		"field5":  []int32(nil),
		"field6":  []int64(nil),
		"field7":  []uint(nil),
		"field8":  []uint8(nil),
		"field9":  []uint16(nil),
		"field10": []uint32(nil),
		"field11": []uint64(nil),
		"field12": []uintptr(nil),
		"field13": []float32(nil),
		"field14": []float64(nil),
		"field15": []string(nil),
	}

	ch := Cast(data, sliceParams, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, sliceExpectedChanges, ch.Changes())
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, sliceExpectedTypes, ch.types)
}

func TestCast_sliceWithValue(t *testing.T) {
	data := struct {
		Field1  []bool
		Field2  []int
		Field3  []int8
		Field4  []int16
		Field5  []int32
		Field6  []int64
		Field7  []uint
		Field8  []uint8
		Field9  []uint16
		Field10 []uint32
		Field11 []uint64
		Field12 []uintptr
		Field13 []float32
		Field14 []float64
		Field15 []string
	}{
		[]bool{true},
		[]int{1},
		[]int8{1},
		[]int16{1},
		[]int32{1},
		[]int64{1},
		[]uint{1},
		[]uint8{1},
		[]uint16{1},
		[]uint32{1},
		[]uint64{1},
		[]uintptr{1},
		[]float32{1},
		[]float64{1},
		[]string{"1"},
	}

	expectedValues := map[string]interface{}{
		"field1":  []bool{true},
		"field2":  []int{1},
		"field3":  []int8{1},
		"field4":  []int16{1},
		"field5":  []int32{1},
		"field6":  []int64{1},
		"field7":  []uint{1},
		"field8":  []uint8{1},
		"field9":  []uint16{1},
		"field10": []uint32{1},
		"field11": []uint64{1},
		"field12": []uintptr{1},
		"field13": []float32{1},
		"field14": []float64{1},
		"field15": []string{"1"},
	}

	ch := Cast(data, sliceParams, []string{
		"field1",
		"field2",
		"field3",
		"field4",
		"field5",
		"field6",
		"field7",
		"field8",
		"field9",
		"field10",
		"field11",
		"field12",
		"field13",
		"field14",
		"field15",
	})

	assert.Nil(t, ch.Errors())
	assert.Equal(t, sliceExpectedChanges, ch.Changes())
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, sliceExpectedTypes, ch.types)
}

func TestCast_isDeepZero(t *testing.T) {
	var emtptyIntSlice []int

	tests := []struct {
		Value  reflect.Value
		Result bool
		Depth  int
	}{
		{
			Value:  reflect.ValueOf(1),
			Result: false,
		},
		{
			Value:  reflect.ValueOf(0),
			Result: true,
		},
		{
			Value:  reflect.ValueOf(uint(1)),
			Result: false,
		},
		{
			Value:  reflect.ValueOf(uint(0)),
			Result: true,
		},
		{
			Value:  reflect.ValueOf(float64(1)),
			Result: false,
		},
		{
			Value:  reflect.ValueOf(float64(0)),
			Result: true,
		},
		{
			Value:  reflect.ValueOf(complex(1, 1)),
			Result: false,
		},
		{
			Value:  reflect.ValueOf(complex(0, 0)),
			Result: true,
		},
		{
			Value:  reflect.ValueOf(true),
			Result: false,
		},
		{
			Value:  reflect.ValueOf(false),
			Result: true,
		},
		{
			Value:  reflect.ValueOf("Hello World!"),
			Result: false,
		},
		{
			Value:  reflect.ValueOf(""),
			Result: true,
		},
		{
			Value:  reflect.ValueOf([]int{1}),
			Result: false,
			Depth:  1,
		},
		{
			Value:  reflect.ValueOf(emtptyIntSlice),
			Result: true,
			Depth:  1,
		},
		{
			Value:  reflect.ValueOf([1]string{"hello world"}),
			Result: false,
			Depth:  1,
		},
		{
			Value:  reflect.ValueOf([1]string{}),
			Result: true,
			Depth:  1,
		},
		{
			Value:  reflect.ValueOf(struct{ Text string }{Text: "Hello World"}),
			Result: false,
			Depth:  1,
		},
		{
			Value:  reflect.ValueOf(struct{ Text string }{Text: ""}),
			Result: true,
			Depth:  1,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.Result, isDeepZero(tt.Value, tt.Depth))
	}
}
