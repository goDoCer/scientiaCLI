package scientia

import (
	"errors"
	"net/http"
)

func checkResponse(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Error: Response Status is not 200 but is " + resp.Status)
	}

	return nil
}
