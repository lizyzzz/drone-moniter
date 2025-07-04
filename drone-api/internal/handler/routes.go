// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.4

package handler

import (
	"net/http"

	"drone-api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/drone/status",
				Handler: HandleDroneStatusHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/ws",
				Handler: HandleWebSocketHandler(serverCtx),
			},
		},
	)
}
