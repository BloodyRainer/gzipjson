package gzipjson

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"strconv"
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
	c, err := CompressWitMinSize(buf, p, 1400)

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
	c, err := CompressWitMinSize(bufCompressed, p, 0)

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

func TestGzipJsonCompressed2(t *testing.T) {
	// Arrange
	p := &Person{
		Age:       44,
	}

	const expect = `{"age":44}`
	bufCompressed := new(bytes.Buffer)

	// Act
	c, err := CompressWitMinSize(bufCompressed, p, 0)

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

func TestGzipJsonCompressedMultiple(t *testing.T) {

	for i := 1; i < 1000; i++ {
		// Arrange
		p := &Person{
			Name: "Lenni Linux",
			Age: i,
		}

		expect := `{"name":"Lenni Linux","age":` + strconv.Itoa(i) + `}`
		bufCompressed := new(bytes.Buffer)

		// Act
		c, err := CompressWitMinSize(bufCompressed, p, 0)

		// Assert
		assert.Nil(t, err)

		bufUncompressed := new(bytes.Buffer)
		gr, err := gzip.NewReader(bufCompressed)
		if err != nil {
			log.Fatal("failed to gunzip content", err.Error())
		}
		io.Copy(bufUncompressed, gr)

		assert.True(t, c)
		assert.Equal(t, expect, strings.TrimSpace(bufUncompressed.String()))

	}
}

// run with 'go test -v -bench=. -benchmem -run=^a'
func BenchmarkGzipJsonCompressedMultiple(b *testing.B) {
	b.StopTimer()

	for i := 1; i < b.N; i++ {
		// Arrange
		p := &Person{
			Name: "Lenni Linux",
			Age: i,
			Hobbies: []string{"gaming", "youtube", "eating", "coding"},
		}

		bufCompressed := new(bytes.Buffer)

		// Act
		b.StartTimer()
		c, err := CompressWitMinSize(bufCompressed, p, 0)
		b.StopTimer()

		assert.Nil(b, err)
		assert.True(b, c)
	}
}
