package ws

import (
	"errors"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type WSRouter struct {
	routes     map[models.WSMessageType]HandlerFunc
	middleware []MiddlewareFunc
}

type (
	HandlerFunc    func(s *Server, client *Client, req models.WSRequest)
	MiddlewareFunc func(HandlerFunc) HandlerFunc
)

func NewWsRouter() *WSRouter {
	router := &WSRouter{
		routes: make(map[models.WSMessageType]HandlerFunc),
	}

	router.Handle(models.TypeAuth, handleAuth)

	authRouter := &WSRouter{
		routes: make(map[models.WSMessageType]HandlerFunc),
	}

	authRouter.Use(
		WsAuthMiddleware(),
	)

	authRouter.Handle("BOOT", handleBoot)
	authRouter.Handle("STATUS", handleStatus)
	authRouter.Handle("SET STATUS READY", handleSetStatusReady)
	authRouter.Handle("START", handlePlanting)

	for msgType, handler := range authRouter.routes {
		router.routes[msgType] = handler
	}

	return router
}

func (r *WSRouter) Use(mw ...MiddlewareFunc) {
	r.middleware = append(r.middleware, mw...)
}

func (r *WSRouter) Handle(msgType models.WSMessageType, handler HandlerFunc) {
	for _, mw := range r.middleware {
		handler = mw(handler)
	}
	r.routes[msgType] = handler
}

func (r *WSRouter) WsRouter(req models.WSRequest) (HandlerFunc, error) {
	handler, exists := r.routes[req.Type]
	if !exists {
		return nil, errors.New("Handler not found for message type: " + string(req.Type))
	}

	return handler, nil
}
