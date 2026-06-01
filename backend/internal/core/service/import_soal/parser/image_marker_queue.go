package parser

import (
	"fmt"

	content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"
)

type imageNameQueue struct {
	names []string
	index int
	set   map[string]struct{}
}

func newImageNameQueue(data []byte) (*imageNameQueue, error) {
	names, err := extractImageNameQueue(data)
	if err != nil {
		return nil, err
	}

	set := make(map[string]struct{}, len(names))
	for _, name := range names {
		set[name] = struct{}{}
	}
	return &imageNameQueue{names: names, set: set}, nil
}

func (q *imageNameQueue) next() (string, bool) {
	if q.index >= len(q.names) {
		return "", false
	}
	name := q.names[q.index]
	q.index++
	return name, true
}

func (q *imageNameQueue) exists(name string) bool {
	_, ok := q.set[name]
	return ok
}

func (q *imageNameQueue) consumeIfNext(name string) {
	if q.index < len(q.names) && q.names[q.index] == name {
		q.index++
	}
}

func imageContentParts(
	rawBody string,
	bodyContent content.RichContent,
	nextImage func() (string, bool),
	imageExists func(string) bool,
	consumeImageIfNext func(string),
) (string, content.RichContent, error) {
	if filename, ok := leadingImageFilename(rawBody); ok {
		if !imageExists(filename) {
			return "", content.RichContent{}, fmt.Errorf("gambar %q tidak ditemukan di dalam file docx", filename)
		}
		consumeImageIfNext(filename)
		return filename, trimContentPrefix(bodyContent, len(filename)), nil
	}

	name, ok := nextImage()
	if !ok {
		return "", content.RichContent{}, fmt.Errorf("[IMG] tetapi tidak ada gambar tersisa di dokumen")
	}
	return name, bodyContent, nil
}
