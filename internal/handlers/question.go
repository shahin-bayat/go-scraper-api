package handlers

import (
	"errors"
	"fmt"
	"github.com/shahin-bayat/scraper-api/internal/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shahin-bayat/scraper-api/internal/middlewares"
	"github.com/shahin-bayat/scraper-api/internal/utils"
)

var (
	ErrorMissingCategoryId       = errors.New("category id is required")
	ErrorUnsupportedLanguage     = errors.New("language not supported")
	ErrorMissingQuestionId       = errors.New("question id is required")
	ErrorMissingFilename         = errors.New("filename is required")
	ErrorFileNotFound            = errors.New("file not found")
	ErrorUnsupportedQuestionType = errors.New("unsupported question type")
)

var (
	freeQuestionIds    = [3]uint{14, 50, 55}
	supportedLanguages = map[string]string{"en": "English"}
)

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) error {
	categories, err := h.store.QuestionRepository().GetCategories()
	if err != nil {
		return err
	}
	utils.WriteJSON(w, http.StatusOK, categories, nil)
	return nil
}

func (h *Handler) GetSupportedLanguages(w http.ResponseWriter, r *http.Request) error {
	utils.WriteJSON(w, http.StatusOK, supportedLanguages, nil)
	return nil
}

func (h *Handler) GetCategoryDetail(w http.ResponseWriter, r *http.Request) error {
	categoryId := chi.URLParam(r, "categoryId")
	questionType := r.URL.Query().Get("type")
	if questionType != "" && questionType != "image" {
		return utils.NewAPIError(
			http.StatusUnprocessableEntity, ErrorUnsupportedQuestionType,
		)
	}

	if categoryId == "" {
		return utils.NewAPIError(http.StatusUnprocessableEntity, ErrorMissingCategoryId)
	}
	intCategoryId, err := strconv.Atoi(categoryId)
	if err != nil {
		return utils.NewAPIError(http.StatusUnprocessableEntity, err)
	}

	userId, err := middlewares.GetUserIdFromContext(r.Context())
	if err != nil {
		category, err := h.store.QuestionRepository().GetFreeCategoryDetail(uint(intCategoryId), freeQuestionIds)
		if err != nil {
			return err
		}
		utils.WriteJSON(w, http.StatusOK, category, nil)
		return nil
	}
	if questionType == "image" && userId == 0 {
		return utils.NewAPIError(http.StatusUnauthorized, h.services.AuthService.ErrorUnauthorized())
		//TODO: check subscription status for free questions: && subscription was valid
	} else if questionType == "image" && userId != 0 {
		category, err := h.store.QuestionRepository().GetCategoryDetail(uint(intCategoryId), questionType)
		if err != nil {
			return err
		}

		utils.WriteJSON(w, http.StatusOK, category, nil)
		return nil
	}

	category, err := h.store.QuestionRepository().GetCategoryDetail(uint(intCategoryId), "")
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, category, nil)
	return nil
}

func (h *Handler) GetQuestionDetail(w http.ResponseWriter, r *http.Request) error {
	questionId := chi.URLParam(r, "questionId")
	lang := r.URL.Query().Get("lang")

	if lang != "" && !utils.KeyInMap(supportedLanguages, lang) {
		return utils.NewAPIError(
			http.StatusUnprocessableEntity, ErrorUnsupportedLanguage,
		)
	}

	if questionId == "" {
		return utils.NewAPIError(http.StatusUnprocessableEntity, ErrorMissingQuestionId)
	}
	intQuestionId, err := strconv.Atoi(questionId)
	if err != nil {
		return utils.NewAPIError(http.StatusUnprocessableEntity, err)
	}

	userId, err := middlewares.GetUserIdFromContext(r.Context())
	if err != nil {
		if !utils.UintInSlice(freeQuestionIds[:], uint(intQuestionId)) {
			return utils.NewAPIError(http.StatusUnauthorized, h.services.AuthService.ErrorUnauthorized())
		}
	}
	//TODO: check subscription status for bookmarks
	question, err := h.store.QuestionRepository().GetQuestionDetail(
		uint(intQuestionId), userId, utils.TrimSpaceLower(lang), h.appConfig.APIBaseURL,
	)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, question, nil)
	return nil
}

func (h *Handler) ToggleBookmark(w http.ResponseWriter, r *http.Request) error {
	userId, err := middlewares.GetUserIdFromContext(r.Context())
	if err != nil {
		return utils.NewAPIError(http.StatusUnauthorized, h.services.AuthService.ErrorUnauthorized())
	}
	// TODO: check subscription status for bookmarks

	var req models.BookmarkRequest
	if err := utils.DecodeRequestBody(r, &req); err != nil {
		return utils.InvalidJSON()
	}
	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		return utils.InvalidRequestData(validationErrors)
	}

	bookmarkId, err := h.store.QuestionRepository().BookmarkQuestion(
		req.QuestionId, userId,
	)
	if err != nil {
		return err
	}
	if bookmarkId == 0 {
		utils.WriteJSON(w, http.StatusNoContent, nil, nil)
		return nil
	} else {
		utils.WriteJSON(w, http.StatusCreated, nil, nil)
		return nil
	}
}

func (h *Handler) GetBookmarks(w http.ResponseWriter, r *http.Request) error {
	userId, err := middlewares.GetUserIdFromContext(r.Context())
	if err != nil {
		return utils.NewAPIError(http.StatusUnauthorized, h.services.AuthService.ErrorUnauthorized())
	}

	// TODO: check subscription status for bookmarks

	bookmarks, err := h.store.QuestionRepository().GetBookmarks(userId)
	if err != nil {
		return err
	}
	utils.WriteJSON(w, http.StatusOK, bookmarks, nil)
	return nil
}

func (h *Handler) GetImage(w http.ResponseWriter, r *http.Request) error {
	filename := chi.URLParam(r, "filename")
	if filename == "" {
		return utils.NewAPIError(http.StatusUnprocessableEntity, ErrorMissingFilename)
	}
	filenameSanitized := filepath.Clean(filename)
	filePath := fmt.Sprintf("assets/images/%s", filenameSanitized)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return utils.NewAPIError(http.StatusNotFound, ErrorFileNotFound)
	} else if err != nil {
		return err
	}

	http.ServeFile(w, r, fmt.Sprintf("assets/images/%s", filename))
	return nil
}

func (h *Handler) GetFailedQuestions(w http.ResponseWriter, r *http.Request) error {
	userId, err := middlewares.GetUserIdFromContext(r.Context())
	if err != nil {
		return utils.NewAPIError(http.StatusUnauthorized, h.services.AuthService.ErrorUnauthorized())
	}

	failedQuestions, err := h.store.QuestionRepository().GetFailedQuestions(userId)
	if err != nil {
		return err
	}
	utils.WriteJSON(w, http.StatusOK, failedQuestions, nil)
	return nil
}

func (h *Handler) AddOrRemoveFailedQuestion(w http.ResponseWriter, r *http.Request) error {
	questionId := chi.URLParam(r, "questionId")
	if questionId == "" {
		return utils.NewAPIError(http.StatusUnprocessableEntity, ErrorMissingQuestionId)
	}
	intQuestionId, err := strconv.Atoi(questionId)
	if err != nil {
		return utils.NewAPIError(http.StatusUnprocessableEntity, err)
	}
	userId, err := middlewares.GetUserIdFromContext(r.Context())
	if err != nil {
		return utils.NewAPIError(http.StatusUnauthorized, h.services.AuthService.ErrorUnauthorized())
	}
	if err := h.store.QuestionRepository().AddOrRemoveFailedQuestion(userId, uint(intQuestionId)); err != nil {
		return err
	}
	utils.WriteJSON(w, http.StatusNoContent, nil, nil)
	return nil
}
