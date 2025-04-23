package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/DarRo9/pvz_service/internal/metrics"
	"github.com/DarRo9/pvz_service/internal/service"
	"github.com/DarRo9/pvz_service/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type HTTPHandler struct {
	service service.ServiceInterface
}

func NewHTTPHandler(service service.ServiceInterface) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}

func validateRole(ctx context.Context, w http.ResponseWriter, allowedRoles []string) bool {
	user, ok := ctx.Value("user").(jwt.MapClaims)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Unauthorized")
	}

	for _, role := range allowedRoles {
		if user["role"] == role {
			return true
		}
	}

	WriteError(w, http.StatusForbidden, "Forbidden")
	return false
}

func writeResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Error{Message: message})
}

func (h *HTTPHandler) PostDummyLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request in PostDummyLogin")

	var request PostDummyLoginJSONBody
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error decoding request body:", err)
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if !h.service.IsValidRole(service.UserRole(request.Role)) {
		log.Println("Invalid role:", request.Role)
		WriteError(w, http.StatusBadRequest, "Invalid role")
		return
	}

	token, err := utils.GenerateJWT("dummy_id", "dummy_email", string(request.Role))
	if err != nil {
		log.Println("Error generating token:", err)
		WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	log.Println("Token generated")
	response := Token(token)
	writeResponse(w, http.StatusOK, response)
}

func (h *HTTPHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request in PostLogin")
	ctx := r.Context()

	var request PostLoginJSONBody
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error decoding request body:", err)
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.GetUserByEmail(ctx, string(request.Email))
	if err != nil {
		log.Println("Error getting user by email:", err)
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.CheckPassword(ctx, user, string(request.Password)); err != nil {
		log.Println("Invalid password:", err)
		WriteError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		log.Println("Error generating token:", err)
		WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	log.Println("Token generated")
	response := Token(token)
	writeResponse(w, http.StatusOK, response)
}

// Добавление товара в текущую приемку (только для сотрудников ПВЗ)
// (POST /products)
func (h *HTTPHandler) PostProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request in PostProducts")
	ctx := r.Context()

	if !validateRole(ctx, w, []string{"employee"}) {
		log.Println("Unauthorized")
		return
	}

	var request Product
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error decoding request body:", err)
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := h.service.CreateProduct(
		ctx,
		request.ReceptionId.String(),
		string(request.Type),
	)
	if err != nil {
		log.Println("Error creating product:", err)
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Println("Product created")
	metrics.ProductsAddedTotal.Inc()
	response := productRepositoryToHTTP(product)
	writeResponse(w, http.StatusCreated, response)
}

// Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией
// (GET /pvz)
func (h *HTTPHandler) GetPvz(w http.ResponseWriter, r *http.Request, params GetPvzParams) {
	log.Println("Got request in GetPvz")
	ctx := r.Context()
	if !validateRole(ctx, w, []string{"employee", "moderator"}) {
		log.Println("Unauthorized")
		return
	}

	page := 1
	limit := 10

	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	if page <= 0 {
		WriteError(w, http.StatusBadRequest, "Page must be greater than 0")
		return
	}

	if limit <= 0 || limit >= 31 {
		WriteError(w, http.StatusBadRequest, "Limit must be between 1 and 30")
		return
	}

	pvzs, err := h.service.ListPVZ(ctx, params.StartDate, params.EndDate, page, limit)
	if err != nil {
		log.Println("Error getting PVZ list:", err)
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := make([]*PVZWithReceptions, len(pvzs))
	for i := range pvzs {
		response[i] = pvzWithReceptionsRepositoryToHTTP(pvzs[i])
	}
	log.Println("PVZ list retrieved")
	writeResponse(w, http.StatusOK, response)
}

// Создание ПВЗ (только для модераторов)
// (POST /pvz)
func (h *HTTPHandler) PostPvz(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request in PostPvz")
	ctx := r.Context()

	if !validateRole(ctx, w, []string{"moderator"}) {
		log.Println("Unauthorized")
		return
	}

	var request PVZ
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error decoding request body:", err)
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	pvz, err := h.service.CreatePVZ(
		ctx,
		string(request.City),
	)
	if err != nil {
		log.Println("Error creating PVZ:", err)
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	metrics.PVZCreatedTotal.Inc()

	log.Println("PVZ created")
	response := pvzRepositoryToHTTP(pvz)
	writeResponse(w, http.StatusCreated, response)
}

// Закрытие последней открытой приемки товаров в рамках ПВЗ
// (POST /pvz/{pvzId}/close_last_reception)
func (h *HTTPHandler) PostPvzPvzIdCloseLastReception(w http.ResponseWriter, r *http.Request, pvzId openapi_types.UUID) {
	log.Println("Got request in PostPvzPvzIdCloseLastReception")
	ctx := r.Context()
	if !validateRole(ctx, w, []string{"employee"}) {
		log.Println("Unauthorized")
		return
	}

	rc, err := h.service.CloseReception(ctx, pvzId.String())
	if err != nil {
		log.Println("Error closing reception:", err)
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Reception closed")
	response := receptionRepositoryToHTTP(rc)
	writeResponse(w, http.StatusOK, response)
}

// Удаление последнего добавленного товара из текущей приемки (LIFO, только для сотрудников ПВЗ)
// (POST /pvz/{pvzId}/delete_last_product)
func (h *HTTPHandler) PostPvzPvzIdDeleteLastProduct(w http.ResponseWriter, r *http.Request, pvzId openapi_types.UUID) {
	log.Println("Got request in PostPvzPvzIdDeleteLastProduct")
	ctx := r.Context()
	if !validateRole(ctx, w, []string{"employee"}) {
		log.Println("Unauthorized")
		return
	}

	_, err := h.service.DeleteProduct(ctx, pvzId.String())
	if err != nil {
		log.Println("Error deleting product:", err)
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Product deleted")
}

// Создание новой приемки товаров (только для сотрудников ПВЗ)
// (POST /receptions)
func (h *HTTPHandler) PostReceptions(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request in PostReceptions")
	ctx := r.Context()

	if !validateRole(ctx, w, []string{"employee"}) {
		log.Println("Unauthorized")
		return
	}

	var request Reception
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error decoding request body:", err)
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	rc, err := h.service.CreateReception(ctx, request.PvzId.String())
	if err != nil {
		log.Println("Error creating reception:", err)
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Println("Reception created")
	metrics.ReceptionsCreatedTotal.Inc()
	response := receptionRepositoryToHTTP(rc)
	writeResponse(w, http.StatusCreated, response)
}

// Регистрация пользователя
// (POST /register)
func (h *HTTPHandler) PostRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request in PostRegister")
	ctx := r.Context()

	var request PostRegisterJSONBody
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error decoding request body:", err)
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.RegisterUser(ctx, string(request.Email), string(request.Password), string(request.Role))
	if err != nil {
		log.Println("Error registering user:", err)
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Println("User registered")
	response := userRepositoryToHTTP(user)
	writeResponse(w, http.StatusCreated, response)
}
