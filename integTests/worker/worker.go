package worker

import (
	"github.com/gin-gonic/gin"
	"github.com/xcherryio/apis/goapi/xcapi"
	"github.com/xcherryio/sdk-go/xc"
	"log"
	"net/http"
)

type worker struct {
	workerService xc.WorkerService
}

func StartGinWorker(workerService xc.WorkerService) (closeFunc func()) {
	w := worker{
		workerService: workerService,
	}

	router := gin.Default()
	router.POST(xc.ApiPathAsyncStateWaitUntil, w.apiAsyncStateWaitUntil)
	router.POST(xc.ApiPathAsyncStateExecute, w.apiAsyncStateExecute)

	wfServer := &http.Server{
		Addr:    ":" + xc.DefaultWorkerPort,
		Handler: router,
	}
	go func() {
		if err := wfServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return func() { wfServer.Close() }
}

func (w worker) apiAsyncStateWaitUntil(c *gin.Context) {
	var req xcapi.AsyncStateWaitUntilRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := w.workerService.HandleAsyncStateWaitUntil(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (w worker) apiAsyncStateExecute(c *gin.Context) {
	var req xcapi.AsyncStateExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := w.workerService.HandleAsyncStateExecute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
