package hocon

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type HoconValue struct {
	values []HoconElement
}

func NewHoconValue() *HoconValue {
	return &HoconValue{}
}

func (p *HoconValue) IsEmpty() bool {
	if len(p.values) == 0 {
		return true
	}

	if first, ok := p.values[0].(*HoconObject); ok {
		if len(first.items) == 0 {
			return true
		}
	}
	return false
}

func (p *HoconValue) AtKey(key string) *HoconRoot {
	obj := NewHoconObject()
	obj.GetOrCreateKey(key)
	obj.items[key] = p
	r := NewHoconValue()
	r.values = append(r.values, obj)
	return NewHoconRoot(r)
}

func (p *HoconValue) IsString() bool {

	strCount := 0
	for _, v := range p.values {
		if v.IsString() {
			strCount += 1
		}
	}

	if strCount > 0 && strCount == len(p.values) {
		return true
	}

	return false
}

func (p *HoconValue) concatString() string {
	concat := ""
	for _, v := range p.values {
		concat += v.GetString()
	}

	if concat == "null" {
		concat = ""
	}

	return strings.TrimSpace(concat)
}

func (p *HoconValue) GetByteSize() int64 {
	res := p.GetString()
	if len(res) > 0 {
		if res[len(res)-1] == 'b' {
			v := res[0 : len(res)-1]
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				panic(err)
			}
			return i
		}
	}
	i, err := strconv.ParseInt(res, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (p *HoconValue) String() string {
	return p.ToString(0)
}

func (p *HoconValue) ToString(indent int) string {
	if p.IsString() {
		return p.quoteIfNeeded(p.GetString())
	}

	if p.IsObject() {
		tmp := strings.Repeat(" ", indent*2)
		return fmt.Sprintf("{\r\n%s%s}", p.GetObject().ToString(indent+1), tmp)
	}

	if p.IsArray() {
		var strs []string
		for _, item := range p.GetArray() {
			strs = append(strs, item.ToString(indent+1))
		}
		return "[" + strings.Join(strs, ",") + "]"
	}

	return "<<unknown value>>"
}

func (p *HoconValue) GetObject() *HoconObject {

	if len(p.values) == 0 {
		return nil
	}

	var raw interface{}
	raw = p.values[0]

	if o, ok := raw.(*HoconObject); ok {
		return o
	}

	if sub, ok := raw.(MightBeAHoconObject); ok {
		if sub != nil && sub.IsObject() {
			return sub.GetObject()
		}
	}

	return nil
}

func (p *HoconValue) IsObject() bool {
	return p.GetObject() != nil
}

func (p *HoconValue) AppendValue(value HoconElement) {
	p.values = append(p.values, value)
}

func (p *HoconValue) Clear() {
	p.values = []HoconElement{}
}

func (p *HoconValue) NewValue(value HoconElement) {
	p.values = []HoconElement{}
	p.values = append(p.values, value)
}

func (p *HoconValue) GetBoolean() bool {
	v := p.GetString()
	switch v {
	case "on":
		return true
	case "off":
		return false
	case "true":
		return true
	case "false":
		return false
	default:
		panic("Unknown boolean format: " + v)
	}
}

func (p *HoconValue) GetString() string {
	if p.IsString() {
		return p.concatString()
	}
	return ""
}

func (p *HoconValue) GetFloat64() float64 {
	val, err := strconv.ParseFloat(p.GetString(), 64)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *HoconValue) GetFloat32() float32 {
	val, err := strconv.ParseFloat(p.GetString(), 32)
	if err != nil {
		panic(err)
	}
	return float32(val)
}

func (p *HoconValue) GetInt64() int64 {
	val, err := strconv.ParseInt(p.GetString(), 10, 64)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *HoconValue) GetInt32() int32 {
	val, err := strconv.ParseInt(p.GetString(), 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(val)
}

func (p *HoconValue) GetByte() byte {
	val, err := strconv.ParseUint(p.GetString(), 10, 8)
	if err != nil {
		panic(err)
	}
	return byte(val)
}

func (p *HoconValue) GetByteList() []byte {
	arrs := p.GetArray()
	var items []byte
	for _, v := range arrs {
		items = append(items, v.GetByte())
	}
	return items
}

func (p *HoconValue) GetInt32List() []int32 {
	arrs := p.GetArray()
	var items []int32
	for _, v := range arrs {
		items = append(items, v.GetInt32())
	}
	return items
}

func (p *HoconValue) GetInt64List() []int64 {
	arrs := p.GetArray()
	var items []int64
	for _, v := range arrs {
		items = append(items, v.GetInt64())
	}
	return items
}

func (p *HoconValue) GetBooleanList() []bool {
	arrs := p.GetArray()
	var items []bool
	for _, v := range arrs {
		items = append(items, v.GetBoolean())
	}
	return items
}

func (p *HoconValue) GetFloat32List() []float32 {
	arrs := p.GetArray()
	var items []float32
	for _, v := range arrs {
		items = append(items, v.GetFloat32())
	}
	return items
}

func (p *HoconValue) GetFloat64List() []float64 {
	arrs := p.GetArray()
	var items []float64
	for _, v := range arrs {
		items = append(items, v.GetFloat64())
	}
	return items
}

func (p *HoconValue) GetStringList() []string {
	arrs := p.GetArray()
	var items []string
	for _, v := range arrs {
		items = append(items, v.GetString())
	}
	return items
}

func (p *HoconValue) GetArray() []*HoconValue {
	var arrs []*HoconValue

	if len(p.values) == 0 {
		return arrs
	}
	arrs = []*HoconValue{}
	for _, v := range p.values {
		if v.IsArray() {
			arrs = append(arrs, v.GetArray()...)
		}
	}

	return arrs
}

func (p *HoconValue) GetChildObject(key string) *HoconValue {
	return p.GetObject().GetKey(key)
}

func (p *HoconValue) IsArray() bool {
	return p.GetArray() != nil
}

func (p *HoconValue) GetTimeDuration(allowInfinite bool) time.Duration {
	res := p.GetString()
	groups, matched := findStringSubmatchMap(res, `^(?P<value>([0-9]+(\.[0-9]+)?))\s*(?P<unit>(nanoseconds|nanosecond|nanos|nano|ns|microseconds|microsecond|micros|micro|us|milliseconds|millisecond|millis|milli|ms|seconds|second|s|minutes|minute|m|hours|hour|h|days|day|d))$`)

	if matched {
		u := groups["unit"]
		strV := groups["value"]
		v := parsePositiveValue(strV)

		switch u {
		case "nanoseconds", "nanosecond", "nanos", "nano", "ns":
			return time.Duration(float64(time.Nanosecond) * v)
		case "microseconds", "microsecond", "micros", "micro":
			return time.Duration(float64(time.Microsecond) * v)
		case "milliseconds", "millisecond", "millis", "milli", "ms":
			return time.Duration(float64(time.Millisecond) * v)
		case "seconds", "second", "s":
			return time.Duration(float64(time.Second) * v)
		case "minutes", "minute", "m":
			return time.Duration(float64(time.Minute) * v)
		case "hours", "hour", "h":
			return time.Duration(float64(time.Hour) * v)
		case "days", "day", "d":
			return time.Duration(float64(time.Hour*24) * v)
		}
	}

	if strings.ToLower(res) == "infinite" {
		if allowInfinite {
			return time.Duration(-1)
		}
		panic("infinite time duration not allowed")
	}

	return time.Duration(float64(time.Millisecond) * parsePositiveValue(res))
}

func (p *HoconValue) quoteIfNeeded(text string) string {
	if len(text) == 0 {
		return "\"\""
	}

	if strings.IndexByte(text, ' ') >= 0 ||
		strings.IndexByte(text, '\t') >= 0 {
		return "\"" + text + "\""
	}

	return text
}

func findStringSubmatchMap(s, exp string) (map[string]string, bool) {
	reg := regexp.MustCompile(exp)
	captures := make(map[string]string)

	match := reg.FindStringSubmatch(s)
	if match == nil {
		return captures, false
	}

	for i, name := range reg.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		captures[name] = match[i]
	}
	return captures, true
}

func parsePositiveValue(v string) float64 {
	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic(err)
	}
	if value < 0 {
		panic("Expected a positive value instead of " + v)
	}
	return value
}