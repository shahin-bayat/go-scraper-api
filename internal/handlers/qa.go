package handlers

import (
	"fmt"
	"net/http"
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
