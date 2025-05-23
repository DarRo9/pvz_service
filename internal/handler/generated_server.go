// Package http_handler provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Получение тестового токена
	// (POST /dummyLogin)
	PostDummyLogin(w http.ResponseWriter, r *http.Request)
	// Авторизация пользователя
	// (POST /login)
	PostLogin(w http.ResponseWriter, r *http.Request)
	// Добавление товара в текущую приемку (только для сотрудников ПВЗ)
	// (POST /products)
	PostProducts(w http.ResponseWriter, r *http.Request)
	// Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией
	// (GET /pvz)
	GetPvz(w http.ResponseWriter, r *http.Request, params GetPvzParams)
	// Создание ПВЗ (только для модераторов)
	// (POST /pvz)
	PostPvz(w http.ResponseWriter, r *http.Request)
	// Закрытие последней открытой приемки товаров в рамках ПВЗ
	// (POST /pvz/{pvzId}/close_last_reception)
	PostPvzPvzIdCloseLastReception(w http.ResponseWriter, r *http.Request, pvzId openapi_types.UUID)
	// Удаление последнего добавленного товара из текущей приемки (LIFO, только для сотрудников ПВЗ)
	// (POST /pvz/{pvzId}/delete_last_product)
	PostPvzPvzIdDeleteLastProduct(w http.ResponseWriter, r *http.Request, pvzId openapi_types.UUID)
	// Создание новой приемки товаров (только для сотрудников ПВЗ)
	// (POST /receptions)
	PostReceptions(w http.ResponseWriter, r *http.Request)
	// Регистрация пользователя
	// (POST /register)
	PostRegister(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Получение тестового токена
// (POST /dummyLogin)
func (_ Unimplemented) PostDummyLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Авторизация пользователя
// (POST /login)
func (_ Unimplemented) PostLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Добавление товара в текущую приемку (только для сотрудников ПВЗ)
// (POST /products)
func (_ Unimplemented) PostProducts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией
// (GET /pvz)
func (_ Unimplemented) GetPvz(w http.ResponseWriter, r *http.Request, params GetPvzParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Создание ПВЗ (только для модераторов)
// (POST /pvz)
func (_ Unimplemented) PostPvz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Закрытие последней открытой приемки товаров в рамках ПВЗ
// (POST /pvz/{pvzId}/close_last_reception)
func (_ Unimplemented) PostPvzPvzIdCloseLastReception(w http.ResponseWriter, r *http.Request, pvzId openapi_types.UUID) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Удаление последнего добавленного товара из текущей приемки (LIFO, только для сотрудников ПВЗ)
// (POST /pvz/{pvzId}/delete_last_product)
func (_ Unimplemented) PostPvzPvzIdDeleteLastProduct(w http.ResponseWriter, r *http.Request, pvzId openapi_types.UUID) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Создание новой приемки товаров (только для сотрудников ПВЗ)
// (POST /receptions)
func (_ Unimplemented) PostReceptions(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Регистрация пользователя
// (POST /register)
func (_ Unimplemented) PostRegister(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// PostDummyLogin operation middleware
func (siw *ServerInterfaceWrapper) PostDummyLogin(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostDummyLogin(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostLogin operation middleware
func (siw *ServerInterfaceWrapper) PostLogin(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostLogin(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostProducts operation middleware
func (siw *ServerInterfaceWrapper) PostProducts(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostProducts(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetPvz operation middleware
func (siw *ServerInterfaceWrapper) GetPvz(w http.ResponseWriter, r *http.Request) {

	var err error

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPvzParams

	// ------------- Optional query parameter "startDate" -------------

	err = runtime.BindQueryParameter("form", true, false, "startDate", r.URL.Query(), &params.StartDate)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "startDate", Err: err})
		return
	}

	// ------------- Optional query parameter "endDate" -------------

	err = runtime.BindQueryParameter("form", true, false, "endDate", r.URL.Query(), &params.EndDate)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "endDate", Err: err})
		return
	}

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", r.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page", Err: err})
		return
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPvz(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostPvz operation middleware
func (siw *ServerInterfaceWrapper) PostPvz(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostPvz(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostPvzPvzIdCloseLastReception operation middleware
func (siw *ServerInterfaceWrapper) PostPvzPvzIdCloseLastReception(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "pvzId" -------------
	var pvzId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "pvzId", chi.URLParam(r, "pvzId"), &pvzId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "pvzId", Err: err})
		return
	}

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostPvzPvzIdCloseLastReception(w, r, pvzId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostPvzPvzIdDeleteLastProduct operation middleware
func (siw *ServerInterfaceWrapper) PostPvzPvzIdDeleteLastProduct(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "pvzId" -------------
	var pvzId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "pvzId", chi.URLParam(r, "pvzId"), &pvzId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "pvzId", Err: err})
		return
	}

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostPvzPvzIdDeleteLastProduct(w, r, pvzId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostReceptions operation middleware
func (siw *ServerInterfaceWrapper) PostReceptions(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostReceptions(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostRegister operation middleware
func (siw *ServerInterfaceWrapper) PostRegister(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostRegister(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/dummyLogin", wrapper.PostDummyLogin)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/login", wrapper.PostLogin)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/products", wrapper.PostProducts)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/pvz", wrapper.GetPvz)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/pvz", wrapper.PostPvz)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/pvz/{pvzId}/close_last_reception", wrapper.PostPvzPvzIdCloseLastReception)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/pvz/{pvzId}/delete_last_product", wrapper.PostPvzPvzIdDeleteLastProduct)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/receptions", wrapper.PostReceptions)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/register", wrapper.PostRegister)
	})

	return r
}
