package fetch

import (
	"errors"
	"testing"
)

func TestRetryWithExponentialBackoff(t *testing.T) {
	// Mock operation that always returns an error
	mockOperation := func() (interface{}, error) {
		return nil, errors.New("mock error")
	}

	// Call the function under test
	result, err := retryWithExponentialBackoff(mockOperation)

	// Check if the error is expected
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	// Check if the result is nil
	if result != nil {
		t.Errorf("Expected result to be nil, but got %v", result)
	}

	// Mock operation that succeeds after 3 retries
	retryCount := 0
	mockOperation = func() (interface{}, error) {
		retryCount++
		if retryCount <= 3 {
			return nil, errors.New("mock error")
		}
		return "success", nil
	}

	// Call the function under test
	result, err = retryWithExponentialBackoff(mockOperation)

	// Check if the error is nil
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}

	// Check if the result is as expected
	expectedResult := "success"
	if result != expectedResult {
		t.Errorf("Expected result to be %v, but got %v", expectedResult, result)
	}
}
