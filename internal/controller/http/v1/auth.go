package v1

import (
	"github.com/Microsoft/go-winio/pkg/guid"
	_ "github.com/cutlery47/auth-service/docs"
	"github.com/cutlery47/auth-service/internal/service"
	"github.com/labstack/echo/v4"
	rip "github.com/vikram1565/request-ip"
)

type authRoutes struct {
	srv service.Service
	e   *errMapper
}

func newAuthRoutes(g *echo.Group, srv service.Service, e *errMapper) {
	r := &authRoutes{
		srv: srv,
		e:   e,
	}

	g.GET("/auth", r.createTokens)
	g.GET("/refresh", r.refreshTokens)
}

type response struct {
	Access  string
	Refresh string
}

// @Summary	Create Tokens
// @Tags		Auth
// @Param		id	query		string	true	"user guid"
// @Success	200	{object}	response
// @Failure	400	{object}	echo.HTTPError
// @Failure	404	{object}	echo.HTTPError
// @Failure	500	{object}	echo.HTTPError
// @Router		/api/v1/auth [get]
func (r *authRoutes) createTokens(c echo.Context) error {
	response := response{}

	ctx := c.Request().Context()
	id := c.QueryParam("id")

	if id == "" {
		r.e.errLog.Error("id was not provided")
		return echo.NewHTTPError(400, "id was not provided")
	}

	guid, err := guid.FromString(id)
	if err != nil {
		r.e.errLog.Error(err.Error())
		return echo.NewHTTPError(400, "provided id is not guid")
	}

	ip := rip.GetClientIP(c.Request())

	access, refresh, err := r.srv.Create(ctx, guid, ip)
	if err != nil {
		return r.e.Map(err)
	}

	response.Access = access
	response.Refresh = refresh

	return c.JSON(200, response)
}

// @Summary	Refresh Tokens
// @Tags		Auth
// @Param		id		query		string	true	"user guid"
// @Param		refresh	query		string	true	"refresh token"
// @Success	200		{object}	response
// @Failure	400		{object}	echo.HTTPError
// @Failure	500		{object}	echo.HTTPError
// @Router		/api/v1/refresh [get]
func (r *authRoutes) refreshTokens(c echo.Context) error {
	response := response{}

	ctx := c.Request().Context()
	id := c.QueryParam("id")
	// по-хорошему, нужно доставать токены из локальных хранилищ
	// но для данного задания, надеюсь, это не принципиально
	refresh := c.QueryParam("refresh")

	if id == "" {
		r.e.errLog.Error("id was not provided")
		return echo.NewHTTPError(400, "id was not provided")
	}

	if refresh == "" {
		r.e.errLog.Error("refresh token not provided")
		return echo.NewHTTPError(400, "refresh token was not provided")
	}

	guid, err := guid.FromString(id)
	if err != nil {
		r.e.errLog.Error(err.Error())
		return echo.NewHTTPError(400, "provided id is not guid")
	}

	ip := rip.GetClientIP(c.Request())

	access, refresh, err := r.srv.Refresh(ctx, guid, ip, refresh)
	if err != nil {
		return r.e.Map(err)
	}

	response.Access = access
	response.Refresh = refresh

	return c.JSON(200, response)
}
