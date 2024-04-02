package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shahin-bayat/scraper-api/internal/utils"
)

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.store.QARepository().GetCategories()
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, categories)

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

	category, err := h.store.QARepository().GetCategoryDetail(uintCategoryId)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, category)
}

func (h *Handler) GetQuestionDetail(w http.ResponseWriter, r *http.Request) {
	categoryId := chi.URLParam(r, "categoryId")
	questionId := chi.URLParam(r, "questionId")
	if categoryId == "" || questionId == "" {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("category id and question id are required"))
		return
	}
	uintCategoryId, err := strconv.Atoi(categoryId)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	uintQuestionId, err := strconv.Atoi(questionId)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	question, err := h.store.QARepository().GetQuestionDetail(uintCategoryId, uintQuestionId)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, question)
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
