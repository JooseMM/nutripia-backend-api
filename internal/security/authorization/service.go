package authorization

import (
	"context"
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/internal/clients"
	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist"
	"github.com/JooseMM/nutripia-backend-api/internal/security/session"
	sessionTypes "github.com/JooseMM/nutripia-backend-api/internal/security/session/types"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/response"
)

const AuthenticationHeader = "x-session-token"
const UserIdKey = "userId"
const RoleIdKey = "roleId"

type IAuthorizationMiddlewares interface {
	NutritionistOnly(next http.HandlerFunc) http.Handler
	ClientOnly(next http.HandlerFunc) http.Handler
}

type AuthorizationMiddlewares struct {
	NutritionistRepo nutritionist.RepositoryManager
	ClientRepo       clients.IClientRepository
	SessionService   session.ISessionService
}

func NewAuthorizationMiddlewares(
	nutritionistRepo nutritionist.RepositoryManager,
	clientRepo clients.IClientRepository,
	sessionService session.ISessionService,
) IAuthorizationMiddlewares {
	return &AuthorizationMiddlewares{
		NutritionistRepo: nutritionistRepo,
		ClientRepo:       clientRepo,
		SessionService:   sessionService,
	}
}

func (s *AuthorizationMiddlewares) NutritionistOnly(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, sessionErr := s.validateSession(r)
		if sessionErr != nil {
			response.WriteJSON(w, sessionErr.StatusCode, sessionErr)
			return
		}

		ctx := r.Context()
		user, userErr := s.NutritionistRepo.GetById(&ctx, &session.UserId)
		if userErr != nil {
			if userErr.ErrorCode == string(core.UNEXPECTED_ERROR) {
				response.WriteJSON(w, userErr.StatusCode, userErr)
				return
			}

			Unauthenticated()
			return
		}

		if !user.IsEmailConfirmed {
			responseErr := NotEnoughPermissions()
			response.WriteJSON(w, responseErr.StatusCode, responseErr)
			return
		}

		newContext := context.WithValue(ctx, UserIdKey, user.ID)
		newContext = context.WithValue(newContext, RoleIdKey, session.Role)
		next.ServeHTTP(w, r.WithContext(newContext))
	})
}

func (s *AuthorizationMiddlewares) ClientOnly(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, sessionErr := s.validateSession(r)
		if sessionErr != nil {
			response.WriteJSON(w, sessionErr.StatusCode, sessionErr)
			return
		}

		user, userErr := s.ClientRepo.GetById(r.Context(), &session.UserId)
		if userErr != nil {
			if userErr.ErrorCode == string(core.UNEXPECTED_ERROR) {
				response.WriteJSON(w, userErr.StatusCode, userErr)
				return
			}

			Unauthenticated()
			return
		}

		// if !user.IsEmailConfirmed {
		// 	responseErr := NotEnoughPermissions()
		// 	response.WriteJSON(w, responseErr.StatusCode, responseErr)
		// 	return
		// }

		ctx := context.WithValue(r.Context(), UserIdKey, user.ID)
		ctx = context.WithValue(ctx, RoleIdKey, session.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *AuthorizationMiddlewares) validateSession(
	r *http.Request,
) (*sessionTypes.Session, *core.BaseError) {
	token := r.Header.Get(AuthenticationHeader)
	if token == "" {
		return nil, Unauthenticated()
	}

	ctx := r.Context()
	session, e := s.SessionService.VerifySession(&token, &ctx)
	if e != nil {
		if e.ErrorCode == string(core.UNEXPECTED_ERROR) {
			return nil, e
		}

		return nil, Unauthenticated()
	}

	return session, nil
}
