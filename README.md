# GzipJSON

Encodes JSON and compresses the output only if the it is large enough to exceed a given size. Per Default the size is 1400 bytes, which is a good value for webserver responses.

# Usage
JSON object is big enough to be compressed. Outputs a compressed file 'file.gz'.
```
type JsonPerson struct {
	Name       string   `json:"name,omitempty"`
	Age        int      `json:"age,omitempty"`
	Profession string   `json:"profession,omitempty"`
}

func main() {
	f2, err := os.Create("file.gz")
	if err != nil {
		log.Fatal("failed to create file", err.Error())
	}

    // make a big JSON object
	ps := make([]Person, 0)
	for i := 0; i < 100; i++ {
		p := Person{
			Name:       "Lenni Linux",
			Age:        44,
			Profession: "Software Developer",
		}
		ps = append(ps, p)
	}

	c, err := gzipjson.Compress(f2, &ps)
	if err != nil {
		log.Fatal("failed to gzip json", err.Error())
	}

	log.Println("is compressed: ", c) // is compressed: true

}
```

JSON object is _not_ big enough to be compressed. Outputs a file 'file' that is _not_ compressed.

```
func main() {
	f1, err := os.Create("file")
	if err != nil {
		log.Fatal("failed to create file", err.Error())
	}

	p := Person{
		Name:       "Lenni Linux",
		Age:        44,
		Profession: "Software Developer",
	}

	c, err := gzipjson.Compress(f1, &p)
	if err != nil {
		log.Fatal("failed to gzip json", err.Error())
	}

	log.Println("is compressed: ", c) // is compressed: false

}

```

You can configure the size yourself by using 'CompressWitMinSize' instead of 'Compress'.
