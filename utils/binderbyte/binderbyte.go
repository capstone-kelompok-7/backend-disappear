package binderbyte

import (
	"encoding/json"
	"github.com/capstone-kelompok-7/backend-disappear/config"
	"net/http"
)

func TrackingPackages(courier, awb string) (map[string]interface{}, error) {
	url := "https://api.binderbyte.com/v1/track"
	apiKey := config.InitConfig().ResiKey
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("courier", courier)
	query.Add("awb", awb)
	query.Add("api_key", apiKey)
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var trackingInfo map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&trackingInfo)
	if err != nil {
		return nil, err
	}

	return trackingInfo, nil
}
