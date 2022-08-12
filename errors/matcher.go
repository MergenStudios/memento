package errors

type FailedExtractionError struct{}
func (_ *FailedExtractionError) Error() string {
	return "Failed to extract string from pattern"
}