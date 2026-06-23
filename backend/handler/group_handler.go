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
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MemberIDs   []string `json:"member_ids"`
}

type addMemberRequest struct {
	Email string `json:"email"`
}

type memberResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
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

func (h *GroupHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")

	users, err := h.groupService.GetMembers(r.Context(), groupID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar membros")
		return
	}

	members := make([]memberResponse, 0, len(users))
	for _, u := range users {
		members = append(members, memberResponse{
			UserID: u.ID,
			Name:   u.Name,
			Email:  u.Email,
		})
	}

	writeJSON(w, http.StatusOK, members)
}

func (h *GroupHandler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := chi.URLParam(r, "id")

	err := h.groupService.JoinGroup(r.Context(), groupID, userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "você entrou no grupo"})
}

func (h *GroupHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := chi.URLParam(r, "id")

	var req addMemberRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" {
		writeError(w, http.StatusBadRequest, "email é obrigatório")
		return
	}

	err = h.groupService.AddMemberByEmail(r.Context(), groupID, req.Email, userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "membro adicionado com sucesso"})
}
