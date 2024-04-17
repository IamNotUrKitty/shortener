package handlers_test

import (
	"testing"

	"github.com/iamnoturkkitty/shortener/internal/app/links"
	"github.com/iamnoturkkitty/shortener/internal/app/links/handlers"
	linksInfra "github.com/iamnoturkkitty/shortener/internal/infrastructure/links"
	"github.com/labstack/echo"
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
	links.Setup(s.e, s.repo)
}

func TestLinksSuite(t *testing.T) {
	suite.Run(t, new(LinksSuite))
}
