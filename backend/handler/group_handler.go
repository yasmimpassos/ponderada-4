package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"racha-historico/service"
)

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

type createGroupRequest struct {
	Name      string   `json:"name"`
	MemberIDs []string `json:"member_ids"`
}

type addMemberRequest struct {
	UserID string `json:"user_id"`
}

func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

	groups, err := h.groupService.GetUserGroups(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar grupos")
		return
	}

	writeJSON(w, http.StatusOK, groups)
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

	var req createGroupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" {
		writeError(w, http.StatusBadRequest, "nome do grupo é obrigatório")
		return
	}

	group, err := h.groupService.CreateGroup(r.Context(), req.Name, userID, req.MemberIDs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao criar grupo")
		return
	}

	writeJSON(w, http.StatusCreated, group)
}

func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := chi.URLParam(r, "id")

	group, err := h.groupService.GetGroupDetails(r.Context(), groupID, userID)
	if err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, group)
}

func (h *GroupHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := chi.URLParam(r, "id")

	var req addMemberRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.UserID == "" {
		writeError(w, http.StatusBadRequest, "user_id é obrigatório")
		return
	}

	err = h.groupService.AddMember(r.Context(), groupID, req.UserID, userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "membro adicionado com sucesso"})
}
