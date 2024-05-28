package handlers_test

import (
	"testing"

	"github.com/iamnoturkkitty/shortener/internal/app/links"
	"github.com/iamnoturkkitty/shortener/internal/app/links/handlers"
	"github.com/iamnoturkkitty/shortener/internal/config"
	linksDomain "github.com/iamnoturkkitty/shortener/internal/domain/links"
	"github.com/iamnoturkkitty/shortener/internal/echomiddleware"
	linksInfra "github.com/iamnoturkkitty/shortener/internal/infrastructure/links"
	"github.com/iamnoturkkitty/shortener/internal/workers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type LinksSuite struct {
	suite.Suite
	e    *echo.Echo
	repo handlers.Repository
}

func (s *LinksSuite) SetupSuite() {
	s.e = echo.New()
	s.repo = linksInfra.NewInMemoryRepo()

	s.e.Use(echomiddleware.InitJWTMiddleware())
	deleteQueue := make(chan linksDomain.DeleteLinkTask)

	go workers.DeleteLinksWorker(deleteQueue, s.repo)
	// TODO: Убрать хардкод конфига, подумать как получать конфиг для тестов
	links.Setup(s.e, s.repo, &config.Config{Address: "localhost:8080", BaseAddress: "http://localhost:8080"}, deleteQueue)
}

func TestLinksSuite(t *testing.T) {
	suite.Run(t, new(LinksSuite))
}
