package worker

import (
	"github.com/gin-gonic/gin"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-sdk/xdb"
	"log"
	"net/http"
)

type worker struct {
	workerService xdb.WorkerService
}

func StartGinWorker(workerService xdb.WorkerService) (closeFunc func()) {
	w := worker{
		workerService: workerService,
	}

	router := gin.Default()
	router.POST(xdb.ApiPathAsyncStateWaitUntil, w.apiAsyncStateWaitUntil)
	router.POST(xdb.ApiPathAsyncStateExecute, w.apiAsyncStateExecute)

	wfServer := &http.Server{
		Addr:    ":" + xdb.DefaultWorkerPort,
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
	var req xdbapi.AsyncStateWaitUntilRequest
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
	return
}
func (w worker) apiAsyncStateExecute(c *gin.Context) {
	var req xdbapi.AsyncStateExecuteRequest
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
	return
}
