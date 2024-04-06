package yrepository

import (
	"net/url"
	"testing"
)

func TestBuildURLValues(t *testing.T) {
	// Создаем тестовые данные
	testData := []struct {
		input    []string
		expected *url.Values
	}{
		{
			input: []string{"key1", "value1", "key2", "", "key3", "value3"},
			expected: &url.Values{
				"key1": []string{"value1"},
				"key3": []string{"value3"},
			},
		},
		{
			input:    []string{}, // пустые параметры
			expected: &url.Values{},
		},
		{
			input:    []string{"key1", "", "key2", ""}, // только пустые параметры
			expected: &url.Values{},
		},
		{
			input:    []string{"key1", "value1", "key2", "value2"}, // все параметры непустые
			expected: &url.Values{"key1": []string{"value1"}, "key2": []string{"value2"}},
		},
	}

	// Тестирование
	for _, test := range testData {
		result := buildURLValues(test.input...)
		if !urlValuesEqual(t, *result, *test.expected) {
			t.Errorf("Expected %v but got %v", *test.expected, *result)
		}
	}
}

// Вспомогательная функция для сравнения двух url.Values
func urlValuesEqual(t *testing.T, a, b url.Values) bool {
	t.Helper() // Обозначаем эту функцию как вспомогательную для тестирования
	if len(a) != len(b) {
		return false
	}
	for key, valueA := range a {
		valueB, ok := b[key]
		if !ok || len(valueA) != len(valueB) {
			return false
		}
		for i := range valueA {
			if valueA[i] != valueB[i] {
				return false
			}
		}
	}
	return true
}
