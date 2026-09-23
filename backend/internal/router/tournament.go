package router

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/handler"
	"github.com/esportsbar/backend/internal/middleware"
)

// RegisterTournament 注册赛事路由。
func RegisterTournament(rg *gin.RouterGroup, h *handler.TournamentHandler, jwtSecret string) {
	tournaments := rg.Group("/tournaments", middleware.Auth(jwtSecret))
	{
		tournaments.GET("", h.List)
		tournaments.GET("/:id", h.Get)
		tournaments.POST("", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.Create)
		tournaments.PUT("/:id", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.Update)
		tournaments.DELETE("/:id", middleware.RBAC(constants.RoleAdmin), h.Delete)
		tournaments.POST("/:id/register", h.Register)
		tournaments.GET("/:id/registrations", h.ListRegistrations)
		tournaments.POST("/:id/draw", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.DrawGroups)
		tournaments.GET("/:id/matches", h.ListMatches)

		teams := rg.Group("/teams", middleware.Auth(jwtSecret))
		{
			teams.GET("/mine", h.ListMyTeams)
			teams.POST("", h.CreateTeam)
		}

		matches := rg.Group("/matches", middleware.Auth(jwtSecret))
		{
			matches.POST("/:id/result", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.SubmitMatchResult)
		}
	}
}
