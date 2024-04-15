package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shahin-bayat/scraper-api/internal/middlewares"
	"github.com/shahin-bayat/scraper-api/internal/utils"
)

var freeQuestionIds = [3]uint{14, 50, 55}

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.store.QuestionRepository().GetCategories()
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, categories, nil)
}

func (h *Handler) GetCategoryDetail(w http.ResponseWriter, r *http.Request) {
	categoryId := chi.URLParam(r, "categoryId")
	if categoryId == "" {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("category id is required"))
		return
	}
	uintCategoryId, err := strconv.Atoi(categoryId)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	_, err = middlewares.GetUserIdFromContext(r.Context())
	if err != nil {
		category, err := h.store.QuestionRepository().GetFreeCategoryDetail(uintCategoryId, freeQuestionIds)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusNotFound, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, category, nil)
		return
	}

	category, err := h.store.QuestionRepository().GetCategoryDetail(uintCategoryId)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, category, nil)
}

func (h *Handler) GetQuestionDetail(w http.ResponseWriter, r *http.Request) {
	var SupportedLanguages = []string{"en"}

	questionId := chi.URLParam(r, "questionId")
	lang := r.URL.Query().Get("lang")

	if lang != "" && !utils.StringInSlice(SupportedLanguages, lang) {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("language not supported"))
		return
	}

	if questionId == "" {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("category id and question id are required"))
		return
	}
	uintQuestionId, err := strconv.Atoi(questionId)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	_, err = middlewares.GetUserIdFromContext(r.Context())
	if err != nil {
		if !utils.UintInSlice(freeQuestionIds[:], uint(uintQuestionId)) {
			utils.WriteErrorJSON(w, http.StatusUnauthorized, fmt.Errorf("user is not authorized to view this question"))
			return
		}
	}

	question, err := h.store.QuestionRepository().GetQuestionDetail(uintQuestionId, utils.TrimSpaceLower(lang), h.appConfig.APIBaseURL)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusNotFound, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, question, nil)
}

func (h *Handler) GetImage(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	if filename == "" {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("filename is required"))
		return
	}
	filenameSanitized := filepath.Clean(filename)
	filepath := fmt.Sprintf("assets/images/%s", filenameSanitized)
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		utils.WriteErrorJSON(w, http.StatusNotFound, fmt.Errorf("file not found"))
		return
	} else if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	http.ServeFile(w, r, fmt.Sprintf("assets/images/%s", filename))
}
