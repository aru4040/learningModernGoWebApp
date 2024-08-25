package services

import (
	"net/http"
	"time"

	"github.com/aru4040/learningModernGoWebApp/runners/models"
)

type ResultsService struct {
	resultsRepository *repositories.ResultsRepository
	runnersRepository *repositories.RunnersRepository
}

func NewResultsService(
	resultsRepository *repository.ResultsRepository,
	runnersRepository *repository.RunnersRepository,
) *ResultsService {
	return &ResultsService{
		resultsRepository: resultsRepository,
		runnersRepository: runnersRepository,
	}
}

func (rs ResultsService) CreateResult(result *models.Result) (*models.Result, *models.ResponseError) {
	if result.RunnerID == "" {
		return nil, &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}

	if result.RaceResult == "" {
		return nil, &models.ResponseError{
			Message: "Invalid race result",
			Status:  http.StatusBadRequest,
		}
	}

	if result.Location == "" {
		return nil, &models.ResponseError{
			Message: "Invalid location",
			Status:  http.StatusBadRequest,
		}
	}

	if result.Position < 0 {
		return nil, &models.ResponseError{
			Message: "Invalid location",
			Status:  http.StatusBadRequest,
		}
	}

	currentYear := time.Now().Year()
	if result.Year < 0 || result.Year > currentYear {
		return nil, &models.ResponseError{
			Message: "Invalid year",
			Status:  http.StatusBadRequest,
		}
	}

	raceResult, err := parseRaceResult(result.RaceResult)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid race result",
			Status:  http.StatusBadRequest,
		}
	}

	// Resultの作成処理
	response, responseErr := rs.resultsRepository.CreateResult(result)
	if responseErr != nil {
		return nil, responseErr
	}

	// 以下はRunnerのベスト記録の更新処理
	runner, responseErr := rs.runnersRepository.GetRunner(result.RunnerID)
	if responseErr != nil {
		return nil, responseErr
	}

	if runner == nil {
		return nil, &models.ResponseError{
			Message: "Runner not found",
			Status:  http.StatusNotFound,
		}
	}

	// 個人のベスト記録を更新した場合
	if runner.PersonalBest == "" {
		runner.PersonalBest = result.RaceResult
	} else {
		personalBest, err := parseRaceResult(runner.PersonalBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Faild to parse" + "parsonal best",
				Status:  http.StatusInternalServerError,
			}
		}
		if raceResult < personalBest {
			runner.PersonalBest = result.RaceResult
		}
	}

	// 個人のシーズンベスト記録を更新した場合
	if result.Year == currentYear {
		if runner.SeasonBest == "" {
			runner.SeasonBest = result.RaceResult
		} else {
			seasonBest, err := parseRaceResult(runner.SeasonBest)
			if err != nil {
				return nil, &models.ResponseError{
					Message: "Failed to parse" + "season best",
					Status:  http.StatusInternalServerError,
				}
			}
			if raceResult < seasonBest {
				runner.SeasonBest = raceResult
			}
		}

	}

	responseErr = rs.runnersRepository.UpdateRunnerResults(runner)
	if responseErr != nil {
		return nil, responseErr
	}
	return response, nil
}

func (rs ResultsService) DeleteResult(resultId string) *models.ResponseError {
	if resultId == "" {
		return &models.ResponseError{
			Message: "Invalid result Id",
			Status:  http.StatusBadRequest,
		}
	}

	result, responseErr := rs.resultsRepository.DeleteResult(resultId)
	if responseErr != nil {
		return responseErr
	}

	runner, responseErr := rs.runnersRepository.GetRunner(result.RunnerID)
	if responseErr != nil {
		return responseErr
	}

	// 削除した記録が個人ベストに登録されている記録だった場合
	if runner.PersonalBest == result.RaceResult {
		personalBest, responseErr := rs.resultsRepository.GetPersonalBestResults(result.RunnerID)
		if responseErr != nil {
			return responseErr
		}
		runner.PersonalBest = personalBest
	}

	// 削除した記録がシーズンベストに登録されている記録だった場合
	currentYear := time.Now().Year()
	if runner.SeasonBest == result.RaceResult && result.Year == currentYear {
		seasonBest, responseErr := rs.resultsRepository.GetSeasonBestResults(result.RunnerID, result.Year)
		if responseErr != nil {
			return responseErr
		}
		runner.SeasonBest = seasonBest
	}

	responseErr = rs.runnersRepository.UpdateRunnerResults(runner)
	if responseErr != nil {
		return responseErr
	}

	return nil
}

func parseRaceResult(timeString string) (time.Duration, error) {
	return time.ParseDuration(
		timeString[0:2] + "h" +
			timeString[3:5] + "m" +
			timeString[6:8] + "s")
}
