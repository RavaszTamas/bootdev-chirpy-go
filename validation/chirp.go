package validation

import (
	"strings"
)

const MAX_LENGTH int = 140

var filteredWords = map[string]interface{}{
	"kerfuffle": nil,
	"sharbert":  nil,
	"fornax":    nil,
}

// type ValidationError struct {
// 	message string
// }

// // Error implements [error].
// func (v ValidationError) Error() string {
// 	return fmt.Sprintf("Calidation error: %s", v.message)
// }

// func validateLength(data string, w http.ResponseWriter, r *http.Request) (string, error) {
// 	if len(data) > MAX_LENGTH {
// 		log.Printf("Chirp is too long %d", len(data))

// 		return "", ValidationError{
// 			message: "Chirp is too long",
// 		}
// 	}
// 	return data, nil
// }

func ReplaceBadWords(data string) string {
	words := strings.Split(data, " ")

	for i := 0; i < len(words); i++ {
		lowered := strings.ToLower(words[i])
		if _, ok := filteredWords[lowered]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

// func validateBadWords(data string, w http.ResponseWriter, r *http.Request) (string, error) {

// }
