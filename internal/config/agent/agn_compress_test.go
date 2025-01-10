package agent

import "testing"

func TestCompress(t *testing.T) {
	data := []byte("test data")
	compressed, err := compress(data)

	if err != nil {
		t.Fatalf("Ошибка сжатия: %v", err)
	}

	if len(compressed) >= len(data) {
		t.Error("Сжатые данные не меньше исходных")
	}
}
