package service

import (
	"context"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	"github.com/gin-gonic/gin"
	unitedbpapi "github.com/yeyee2901/service-weaver/service/unitedbapi"
)

type service struct {
	ginEngine  *gin.Engine
	weaverRoot weaver.Instance
	unitedbapi unitedbpapi.UniteDBService
}

func NewService(ginMode string) (*service, error) {
	gin.SetMode(ginMode)

	s := &service{
		ginEngine:  gin.New(),
		weaverRoot: weaver.Init(context.Background()),
	}

	// attach unitedb client to service root
	udb, err := weaver.Get[unitedbpapi.UniteDBService](s.weaverRoot)
	if err != nil {
		return nil, err
	}
	s.unitedbapi = udb

	return s, nil
}

func (s *service) RegisterRouting() {
	s.ginEngine.Handle(http.MethodGet, "/unitedb/battle-item", s.GetBattleItem)
}

// Run akan menjalankan server di goroutine dan mengembalikan error channel nya
// ke caller. Error channel dapat di polling untuk mengecek apakah server masih running
func (s *service) Run(addr string) <-chan error {
	errChan := make(chan error)

	listener, err := s.weaverRoot.Listener("unitedb", weaver.ListenerOptions{
		LocalAddress: addr,
	})

	if err != nil {
		errChan <- err
		return errChan
	}

	server := http.Server{
		Handler: weaver.InstrumentHandler("unitedb", s.ginEngine),
	}

	go func() {
		if err = server.Serve(listener); err != nil {
			errChan <- err
			return
		}

		close(errChan)
	}()

	return errChan
}

// GetBattleItem gin handler
func (s *service) GetBattleItem(c *gin.Context) {
	var req GetBattleItemRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := s.unitedbapi.GetBattleItem(c.Request.Context(), req.Name, req.Tier)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
