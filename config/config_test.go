package config

import (
	"testing"

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
