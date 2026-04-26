package shared

import (
	"encoding/json"

	content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"
)

func MarshalRichContent(value content.RichContent) ([]byte, error) {
	if value.Empty() {
		return nil, nil
	}
	return json.Marshal(value)
}

func UnmarshalRichContent(raw []byte) (content.RichContent, error) {
	if len(raw) == 0 {
		return content.RichContent{}, nil
	}

	var value content.RichContent
	if err := json.Unmarshal(raw, &value); err != nil {
		return content.RichContent{}, err
	}
	return value, nil
}
