package scientia

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

func checkResponse(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("response status is not 200 but is " + resp.Status)
	}

	return nil
}

//This probably isn't accurate but we cba to find the spec
// Returns something like 1920 for 2019-2020
func getCurrentAcademicYear() string {
	currentYear := time.Now().Year() % 100 //Only works till 2099 - doesn't matter since DoC will probably replace Scientia by 2023
	currentMonth := time.Now().Month()

	if currentMonth >= time.October {
		return fmt.Sprintf("%d%d", currentYear, currentYear+1)
	}

	return fmt.Sprintf("%d%d", currentYear-1, currentYear)
}
