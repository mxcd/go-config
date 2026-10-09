package config

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(1, 1)

	err := LoadConfigWithOptions([]Value{
		String("STRING"),
		String("STRING_DEFAULT").NotEmpty().Default("string_default"),
		String("STRING_NOT_EMPTY").NotEmpty(),

		StringArray("STRING_ARRAY"),
		StringArray("STRING_ARRAY_DEFAULT").Default([]string{"a", "b"}),
		StringArray("STRING_ARRAY_NOT_EMPTY").NotEmpty(),

		Int("INT"),
		Int("INT_DEFAULT").Default(43),

		Bool("BOOL_TRUE"),
		Bool("BOOL_FALSE"),
		Bool("BOOL_DEFAULT_TRUE").Default(true),
		Bool("BOOL_DEFAULT_FALSE").Default(false),
	}, &LoadConfigOptions{
		DotEnvFile: "test.env",
	})

	assert.Nil(err)

	assert.Equal("string", Get().String("STRING"))
	assert.Equal("string_default", Get().String("STRING_DEFAULT"))
	assert.Equal("string_not_empty", Get().String("STRING_NOT_EMPTY"))

	assert.Equal([]string{"foo", "bar"}, Get().StringArray("STRING_ARRAY"))
	assert.Equal([]string{"a", "b"}, Get().StringArray("STRING_ARRAY_DEFAULT"))
	assert.Equal([]string{"fizz", "buzz"}, Get().StringArray("STRING_ARRAY_NOT_EMPTY"))

	assert.Equal(42, Get().Int("INT"))
	assert.Equal(43, Get().Int("INT_DEFAULT"))

	assert.Equal(true, Get().Bool("BOOL_TRUE"))
	assert.Equal(false, Get().Bool("BOOL_FALSE"))
	assert.Equal(true, Get().Bool("BOOL_DEFAULT_TRUE"))
	assert.Equal(false, Get().Bool("BOOL_DEFAULT_FALSE"))
}

func TestSensitiveDefaultIsMasked(t *testing.T) {
	descriptor := String("SECRET").Default("hunter2").Sensitive().Descriptor()
	assert.NotContains(t, GetSanatizedDefaultValue(descriptor), "hunter2")
}

func TestSetStringArray(t *testing.T) {
	t.Setenv("STRING_ARRAY_SET", "a,b")
	assert.NoError(t, LoadConfig([]Value{StringArray("STRING_ARRAY_SET")}))

	assert.NotPanics(t, func() { Set().StringArray("STRING_ARRAY_SET", "c, d") })
	assert.Equal(t, []string{"c", "d"}, Get().StringArray("STRING_ARRAY_SET"))
}

func TestBindPFlagInt(t *testing.T) {
	t.Setenv("PORT", "8080")
	assert.NoError(t, LoadConfig([]Value{Int("PORT")}))

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.Int("port", 0, "")
	assert.NoError(t, flags.Set("port", "9090"))

	assert.NotPanics(t, func() { BindPFlag("PORT", flags.Lookup("port")) })
	assert.Equal(t, 9090, Get().Int("PORT"))
}

func TestLoadConfigReturnsValidationError(t *testing.T) {
	t.Setenv("INVALID_INT", "abc")
	assert.NotPanics(t, func() {
		assert.Error(t, LoadConfig([]Value{Int("INVALID_INT")}))
	})
}
