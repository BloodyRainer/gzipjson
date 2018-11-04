package gzipjson

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"strings"
	"testing"
)

type Person struct {
	Name       string   `json:"name,omitempty"`
	Age        int      `json:"age,omitempty"`
	Profession string   `json:"profession,omitempty"`
	Hobbies    []string `json:"hobbies,omitempty"`
}

func TestGzipJsonUncompressed(t *testing.T) {
	// Arrange
	p := &Person{
		Name:       "Lenni Linux",
		Age:        35,
		Profession: "Software Developer",
		Hobbies:    []string{"gaming", "youtube", "eating", "coding"},
	}

	const expect = `{"name":"Lenni Linux","age":35,"profession":"Software Developer","hobbies":["gaming","youtube","eating","coding"]}`
	buf := new(bytes.Buffer)

	// Act
	c, err := CompressWitMinhSize(buf, p, 1400)

	// Assert
	assert.Nil(t, err)
	assert.False(t, c)
	assert.Equal(t, expect , strings.TrimSpace(buf.String()))

}

func TestGzipJsonCompressed(t *testing.T) {
	// Arrange
	p := &Person{
		Name:       "Lenni Linux",
		Hobbies:    []string{"gaming", "youtube", "eating", "coding"},
	}

	const expect = `{"name":"Lenni Linux","hobbies":["gaming","youtube","eating","coding"]}`
	bufCompressed := new(bytes.Buffer)

	// Act
	c, err := CompressWitMinhSize(bufCompressed, p, 0)

	// Assert
	assert.Nil(t, err)

	bufUncompressed := new(bytes.Buffer)
	gr, err := gzip.NewReader(bufCompressed)
	if err != nil {
		log.Fatal("failed to gunzip content", err.Error())
	}

	io.Copy(bufUncompressed, gr)

	assert.True(t, c)
	assert.Equal(t, expect , strings.TrimSpace(bufUncompressed.String()))

}
