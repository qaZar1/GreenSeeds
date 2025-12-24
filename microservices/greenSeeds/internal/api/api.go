package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/rs/zerolog/log"
)

type API struct {
	url  string
	http *http.Client
}

func NewAPI(url string) *API {
	return &API{
		url: url,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *API) RequestAI(seed string, buf bytes.Buffer) (models.ResponseAPI, error) {
	if seed == "Amaranth" {
		return models.ResponseAPI{
			PercentOfMatch: 1,
		}, nil
	}
	req, err := http.NewRequest(POST, a.url, &buf)
	if err != nil {
		log.Err(err).Msg("Cannot create request")
	}

	query := req.URL.Query()
	query.Add(SEED, seed)
	req.URL.RawQuery = query.Encode()

	resp, err := a.http.Do(req)
	if err != nil {
		log.Err(err).Msg("Cannot create request")
		return models.ResponseAPI{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Msg("Cannot create request")
		return models.ResponseAPI{}, err
	}

	var response models.ResponseAPI
	if err := json.Unmarshal(data, &response); err != nil {
		log.Err(err).Msg("Cannot create request")
		return models.ResponseAPI{}, err
	}

	response, err = a.CheckAI(seed, buf)
	if err != nil {
		log.Err(err).Msg("Cannot create request")
		return models.ResponseAPI{}, err
	}

	return response, nil
}
