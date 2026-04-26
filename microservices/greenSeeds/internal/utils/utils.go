package utils

import (
	"encoding/json"
	"net/http"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/infrastructure"
	"github.com/rs/zerolog/log"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func WriteString(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func GetUuid(r *http.Request, infra *infrastructure.Infrastructure) string {
	token := r.Header.Get("Authorization")[7:]
	claims, err := infra.GetTokenClaims(token)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		return ""
	}
	return claims.Subject
}

func WriteImage(w http.ResponseWriter, status int, image []byte) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(status)
	w.Write(image)
}

func WriteStream(w http.ResponseWriter, status int, video []byte) {
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.WriteHeader(status)
	w.Write(video)
}
