package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"log"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

type LoadConfigOptions struct {
	DotEnvFile  string
	DotEnvFiles []string
}

type Value interface {
	Descriptor() *Descriptor
}

type ValueType string

const (
	StringType      ValueType = "string"
	StringArrayType ValueType = "string_array"
	BoolType        ValueType = "bool"
	IntType         ValueType = "int"
)

type TypeInfo struct {
	Type ValueType
}

type Descriptor struct {
	EnvionmentVariable string
	Default            any
	Provided           bool
	Value              any
	NotEmpty           bool
	Sensitive          bool
	TypeInfo           TypeInfo
}

type applicationConfigType struct {
	StringValues      map[string]*Descriptor
	StringArrayValues map[string]*Descriptor
	BoolValues        map[string]*Descriptor
	IntValues         map[string]*Descriptor
}

type valueSetter struct{}

var applicationConfig applicationConfigType

func Get() *applicationConfigType {
	return &applicationConfig
}

func Set() *valueSetter {
	return &valueSetter{}
}

func (s *valueSetter) String(key string, value string) {
	descriptor, ok := applicationConfig.StringValues[key]
	if !ok {
		log.Panicf("Config value %s not found", key)
	}
	descriptor.Value = value
}

func (s *valueSetter) StringArray(key string, value string) {
	descriptor, ok := applicationConfig.StringValues[key]
	if !ok {
		log.Panicf("Config value %s not found", key)
	}
	descriptor.Value = value
}

func (s *valueSetter) Bool(key string, value bool) {
	descriptor, ok := applicationConfig.BoolValues[key]
	if !ok {
		log.Panicf("Config value %s not found", key)
	}
	descriptor.Value = value
}

func (s *valueSetter) Int(key string, value int) {
	descriptor, ok := applicationConfig.IntValues[key]
	if !ok {
		log.Panicf("Config value %s not found", key)
	}
	descriptor.Value = value
}

func (c *applicationConfigType) String(key string) string {
	descriptor, ok := c.StringValues[key]
	if !ok {
		log.Panicf("Config value %s not found", key)
	}
	return descriptor.Value.(string)
}

func (c *applicationConfigType) StringArray(key string) []string {
	descriptor, ok := c.StringArrayValues[key]
	if !ok {
		log.Panicf("Config value %s not found", key)
	}
	return descriptor.Value.([]string)
}

func (c *applicationConfigType) Bool(key string) bool {
	descriptor, ok := c.BoolValues[key]
	if !ok {
		log.Panicf("Config value %s not found", key)
	}
	return descriptor.Value.(bool)
}

func (c *applicationConfigType) Int(key string) int {
	descriptor, ok := c.IntValues[key]
	if !ok {
		log.Panicf("Config value %s not found", key)
	}
	return descriptor.Value.(int)
}

type stringBuilder struct {
	desc *Descriptor
}

func (b *stringBuilder) Descriptor() *Descriptor {
	return b.desc
}
func String(environmentVariable string) *stringBuilder {
	return &stringBuilder{
		desc: &Descriptor{
			EnvionmentVariable: environmentVariable,
			NotEmpty:           false,
			Sensitive:          false,
			TypeInfo:           TypeInfo{Type: StringType},
		},
	}
}
func (b *stringBuilder) Default(defaultValue string) *stringBuilder {
	b.desc.Default = defaultValue
	return b
}
func (b *stringBuilder) NotEmpty() *stringBuilder {
	b.desc.NotEmpty = true
	return b
}
func (b *stringBuilder) Sensitive() *stringBuilder {
	b.desc.Sensitive = true
	return b
}

type stringArrayBuilder struct {
	desc *Descriptor
}

func (b *stringArrayBuilder) Descriptor() *Descriptor {
	return b.desc
}
func StringArray(environmentVariable string) *stringArrayBuilder {
	return &stringArrayBuilder{
		desc: &Descriptor{
			EnvionmentVariable: environmentVariable,
			NotEmpty:           false,
			Sensitive:          false,
			TypeInfo:           TypeInfo{Type: StringArrayType},
		},
	}
}
func (b *stringArrayBuilder) Default(defaultValue []string) *stringArrayBuilder {
	b.desc.Default = defaultValue
	return b
}
func (b *stringArrayBuilder) NotEmpty() *stringArrayBuilder {
	b.desc.NotEmpty = true
	return b
}
func (b *stringArrayBuilder) Sensitive() *stringArrayBuilder {
	b.desc.Sensitive = true
	return b
}

type boolBuilder struct {
	desc *Descriptor
}

func (b *boolBuilder) Descriptor() *Descriptor {
	return b.desc
}
func Bool(environmentVariable string) *boolBuilder {
	return &boolBuilder{
		desc: &Descriptor{
			EnvionmentVariable: environmentVariable,
			NotEmpty:           false,
			Sensitive:          false,
			TypeInfo:           TypeInfo{Type: BoolType},
		},
	}
}
func (b *boolBuilder) Default(defaultValue bool) *boolBuilder {
	b.desc.Default = defaultValue
	return b
}
func (b *boolBuilder) Sensitive() *boolBuilder {
	b.desc.Sensitive = true
	return b
}

type intBuilder struct {
	desc *Descriptor
}

func (b *intBuilder) Descriptor() *Descriptor {
	return b.desc
}
func Int(environmentVariable string) *intBuilder {
	return &intBuilder{
		desc: &Descriptor{
			EnvionmentVariable: environmentVariable,
			NotEmpty:           false,
			Sensitive:          false,
			TypeInfo:           TypeInfo{Type: IntType},
		},
	}
}
func (b *intBuilder) Default(defaultValue int) *intBuilder {
	b.desc.Default = defaultValue
	return b
}
func (b *intBuilder) Sensitive() *intBuilder {
	b.desc.Sensitive = true
	return b
}

func LoadConfig(values []Value) error {
	return LoadConfigWithOptions(values, &LoadConfigOptions{})
}

func LoadConfigWithOptions(values []Value, options *LoadConfigOptions) error {

	dotEnvFiles := []string{}
	if options.DotEnvFile != "" {
		dotEnvFiles = append(dotEnvFiles, options.DotEnvFile)
	}

	if len(options.DotEnvFiles) > 0 {
		dotEnvFiles = append(dotEnvFiles, options.DotEnvFiles...)
	}

	if err := loadDotEnvFile(dotEnvFiles); err != nil {
		log.Printf("No .env file found. Relying on environment variables.")
	} else {
		log.Printf("Loaded .env file.")
	}

	applicationConfig = applicationConfigType{
		StringValues:      make(map[string]*Descriptor),
		StringArrayValues: make(map[string]*Descriptor),
		BoolValues:        make(map[string]*Descriptor),
		IntValues:         make(map[string]*Descriptor),
	}

	for _, value := range values {
		if err := loadConfigValue(value.Descriptor()); err != nil {
			log.Panic(err)
		}
	}

	return nil
}

func loadDotEnvFile(dotEnvFiles []string) error {
	return godotenv.Load(dotEnvFiles...)
}

func loadConfigValue(valueDescriptor *Descriptor) error {
	switch valueDescriptor.TypeInfo.Type {
	case StringType:
		return loadStringValue(valueDescriptor)
	case StringArrayType:
		return loadStringArrayValue(valueDescriptor)
	case BoolType:
		return loadBoolValue(valueDescriptor)
	case IntType:
		return loadIntValue(valueDescriptor)
	default:
		return fmt.Errorf("unknown type %s", valueDescriptor.TypeInfo.Type)
	}
}

func loadStringValue(valueDescriptor *Descriptor) error {
	value, ok := os.LookupEnv(valueDescriptor.EnvionmentVariable)

	if !ok {
		if valueDescriptor.Default != nil {
			valueDescriptor.Value = valueDescriptor.Default.(string)
			applicationConfig.StringValues[valueDescriptor.EnvionmentVariable] = valueDescriptor
			return nil
		} else {
			return fmt.Errorf("missing environment variable %s", valueDescriptor.EnvionmentVariable)
		}
	}

	if valueDescriptor.NotEmpty && value == "" {
		return fmt.Errorf("environment variable %s must not be empty", valueDescriptor.EnvionmentVariable)
	}

	valueDescriptor.Provided = true
	valueDescriptor.Value = value
	applicationConfig.StringValues[valueDescriptor.EnvionmentVariable] = valueDescriptor

	return nil
}

func loadStringArrayValue(valueDescriptor *Descriptor) error {
	value, ok := os.LookupEnv(valueDescriptor.EnvionmentVariable)

	if !ok {
		if valueDescriptor.Default != nil {
			valueDescriptor.Value = valueDescriptor.Default.([]string)
			applicationConfig.StringArrayValues[valueDescriptor.EnvionmentVariable] = valueDescriptor
			return nil
		} else {
			return fmt.Errorf("missing environment variable %s", valueDescriptor.EnvionmentVariable)
		}
	}

	stringItems := strings.Split(value, ",")
	stringArray := []string{}
	for _, item := range stringItems {
		if item != "" {
			stringArray = append(stringArray, strings.TrimSpace(item))
		}
	}

	if valueDescriptor.NotEmpty && len(stringArray) == 0 {
		return fmt.Errorf("environment variable %s must not be empty", valueDescriptor.EnvionmentVariable)
	}

	valueDescriptor.Provided = true
	valueDescriptor.Value = stringArray
	applicationConfig.StringArrayValues[valueDescriptor.EnvionmentVariable] = valueDescriptor

	return nil
}

func loadBoolValue(valueDescriptor *Descriptor) error {
	value, ok := os.LookupEnv(valueDescriptor.EnvionmentVariable)

	if !ok {
		if valueDescriptor.Default != nil {
			valueDescriptor.Value = valueDescriptor.Default.(bool)
			applicationConfig.BoolValues[valueDescriptor.EnvionmentVariable] = valueDescriptor
			return nil
		} else {
			return fmt.Errorf("missing environment variable %s", valueDescriptor.EnvionmentVariable)
		}
	}

	if strings.ToLower(value) != "true" && strings.ToLower(value) != "false" {
		return fmt.Errorf("environment variable %s must be either 'true' or 'false' | received '%s'", valueDescriptor.EnvionmentVariable, value)
	}

	valueDescriptor.Provided = true
	valueDescriptor.Value = strings.ToLower(value) == "true"
	applicationConfig.BoolValues[valueDescriptor.EnvionmentVariable] = valueDescriptor

	return nil
}

func loadIntValue(valueDescriptor *Descriptor) error {
	value, ok := os.LookupEnv(valueDescriptor.EnvionmentVariable)

	if !ok {
		if valueDescriptor.Default != nil {
			valueDescriptor.Value = valueDescriptor.Default.(int)
			applicationConfig.IntValues[valueDescriptor.EnvionmentVariable] = valueDescriptor
			return nil
		} else {
			return fmt.Errorf("missing environment variable %s", valueDescriptor.EnvionmentVariable)
		}
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("environment variable %s must be a valid integer | received '%s'", valueDescriptor.EnvionmentVariable, value)
	}

	valueDescriptor.Provided = true
	valueDescriptor.Value = number
	applicationConfig.IntValues[valueDescriptor.EnvionmentVariable] = valueDescriptor

	return nil
}

func Print() {
	t := table.NewWriter()
	t.SetTitle("Application Configuration")
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ENV VAR", "Type", "Default Value", "Value Provided", "Current Value"})

	for _, value := range applicationConfig.StringValues {
		t.AppendRow(table.Row{value.EnvionmentVariable, "string", GetSanatizedDefaultValue(value), value.Provided, MaskSensitiveString(value, GetTruncated(value, value.Value))})
	}
	for _, value := range applicationConfig.BoolValues {
		t.AppendRow(table.Row{value.EnvionmentVariable, "bool", GetSanatizedDefaultValue(value), value.Provided, MaskSensitiveString(value, GetTruncated(value, value.Value))})
	}
	for _, value := range applicationConfig.IntValues {
		t.AppendRow(table.Row{value.EnvionmentVariable, "int", GetSanatizedDefaultValue(value), value.Provided, MaskSensitiveString(value, GetTruncated(value, value.Value))})
	}
	for _, value := range applicationConfig.StringArrayValues {
		t.AppendRow(table.Row{value.EnvionmentVariable, "string", GetSanatizedDefaultValue(value), value.Provided, MaskSensitiveString(value, GetTruncated(value, value.Value))})
	}
	t.SortBy([]table.SortBy{
		{Name: "ENV VAR", Mode: table.Asc},
	})
	t.SetStyle(table.StyleColoredYellowWhiteOnBlack)
	t.Render()
}

func MaskSensitiveString(valueDescriptor *Descriptor, value string) string {
	if valueDescriptor.Sensitive {
		switch valueDescriptor.TypeInfo.Type {
		case StringType:
			return strings.Repeat("*", len(value))
		case BoolType:
			return "****"
		case IntType:
			return strings.Repeat("*", len(value))
		default:
			return "****"
		}
	} else {
		return value
	}
}

func GetSanatizedDefaultValue(valueDescriptor *Descriptor) string {
	if valueDescriptor.Default == nil {
		return "-"
	} else {
		return GetTruncated(valueDescriptor, valueDescriptor.Default)
	}
}

func GetTruncated(valueDescriptor *Descriptor, value any) string {
	truncate := func(s string) string {
		if len(s) > 50 {
			return s[:47] + "..."
		}
		return s
	}
	switch valueDescriptor.TypeInfo.Type {
	case StringType:
		return truncate(value.(string))
	case BoolType:
		return fmt.Sprintf("%v", value)
	case IntType:
		valueString := fmt.Sprintf("%v", value)
		return truncate(valueString)
	case StringArrayType:
		valueString := strings.Join(value.([]string), ", ")
		return truncate(valueString)
	default:
		return ""
	}
}

func BindPFlag(key string, flag *pflag.Flag) {
	if flag == nil {
		return
	}

	if valueDescriptor, ok := applicationConfig.StringValues[key]; ok {
		valueDescriptor.Value = flag.Value.String()
		applicationConfig.StringValues[key] = valueDescriptor
		return
	}

	if valueDescriptor, ok := applicationConfig.BoolValues[key]; ok {
		valueDescriptor.Value = strings.ToLower(flag.Value.String()) == "true"
		applicationConfig.BoolValues[key] = valueDescriptor
		return
	}

	if valueDescriptor, ok := applicationConfig.IntValues[key]; ok {
		number, err := strconv.Atoi(flag.Value.String())
		if err != nil {
			log.Panicf("provided flag %s must be a valid integer | received '%s'", flag.Name, flag.Value.String())
		}
		valueDescriptor.Value = number
		applicationConfig.IntValues[key] = valueDescriptor
	}

	log.Panicf("Config value %s not found", key)
}
