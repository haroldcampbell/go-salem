package salem

import (
	"fmt"
	"reflect"
)

type RunType uint

const (
	Invalid RunType = iota

	MinRun
	MaxRun
	ExactRun
)

type PlanRun struct {
	RunType RunType
	Count   int
}

type fieldSetter struct {
	fptr GenType
}
type Plan struct {
	fixedFields map[string]fieldSetter // fields set via ensure
	generators  map[reflect.Kind]func() interface{}

	run        PlanRun
	parentName string
}

func NewPlan() *Plan {
	p := &Plan{}

	p.fixedFields = make(map[string]fieldSetter)
	p.generators = make(map[reflect.Kind]func() interface{})

	p.initDefaultGenerators()

	return p
}

func (p *Plan) RequireFieldValue(fieldName string, sharedValue interface{}) {
	p.fixedFields[fieldName] = fieldSetter{
		fptr: func() interface{} {
			return sharedValue
		},
	}
}

func (p *Plan) SetRunCount(runType RunType, n int) {
	p.run = PlanRun{
		RunType: runType,
		Count:   n,
	}
}

func (p *Plan) generateRandomMock(f *factory) interface{} {
	v := reflect.ValueOf(f.rootType)
	newMockPtr := reflect.New(v.Type())
	newElm := newMockPtr.Elem()
	typeOfT := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		k := field.Kind()
		ff := newElm.Field(i)
		fieldName := typeOfT.Field(i).Name

		qualifiedName := distinctFileName(p.parentName, fieldName)

		var generator GenType
		if p.fixedFields[qualifiedName].fptr != nil {
			generator = p.fixedFields[qualifiedName].fptr
		} else {
			generator = p.generators[k]
		}

		p.updateFieldValue(k, generator, ff, qualifiedName)
	}

	return newMockPtr.Elem().Interface()
}

func (p *Plan) updateFieldValue(k reflect.Kind, generator GenType, ff reflect.Value, qualifiedName string) {
	if !ff.CanSet() {
		return
	}

	if isPrimativeKind(k) {
		val := generator()
		ff.Set(reflect.ValueOf(val))

		return
	}

	fieldName := ff.Type().Name()

	// Complex Types
	switch k {
	case reflect.Array:
		fmt.Printf("%v (array) \n", fieldName)
	case reflect.Slice:
		fmt.Printf("%v (slice) \n", fieldName)

	case reflect.Struct:
		if generator != nil {
			val := generator()
			ff.Set(reflect.ValueOf(val))

		} else {
			ff.Interface()
			mock := Mock(ff.Interface())

			mock.plan.parentName = qualifiedName
			mock.plan.CopyParentRequiredFields(p)
			// WithExactItems(5)

			results := mock.Execute()

			ff.Set(reflect.ValueOf(results[0]))
		}

	default:
		fmt.Printf("%v (Found type) \n", fieldName)
	}
}

func (p *Plan) CopyParentRequiredFields(pp *Plan) {
	for k, v := range pp.fixedFields {
		p.fixedFields[k] = v
	}

}
func distinctFileName(parentFieldName string, fieldName string) string {
	if parentFieldName == "" {
		return fieldName
	}

	return fmt.Sprintf("%s.%s", parentFieldName, fieldName)
}

func (p *Plan) Run(f *factory) []interface{} {
	items := make([]interface{}, 0)

	for i := 0; i < p.run.Count; i++ {
		items = append(items, p.generateRandomMock(f))
	}

	return items
}

func isPrimativeKind(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:

		return true
	}

	return false
}
