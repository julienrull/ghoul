package ghoul
import (
	"net/http"
	"os"
	"os/signal"
    "syscall"

)
func New(config ...Config) *Router {
    newConfig := defaultConfiguration
    if len(config) > 0 {
        newConfig = config[0]
    }
    mux := http.NewServeMux()
    done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    router := &Router{
        Handle: mux,
        Mux: mux,
        BaseUrl: "",
        signalOut: done,
        isRoot: true,
        Server: &http.Server{
            Addr:                         newConfig.Addr, 
            Handler:                      newConfig.Handler, 
            DisableGeneralOptionsHandler: newConfig.DisableGeneralOptionsHandler, 
            TLSConfig:                    newConfig.TLSConfig, 
            ReadTimeout:                  newConfig.ReadTimeout, 
            ReadHeaderTimeout:            newConfig.ReadHeaderTimeout, 
            WriteTimeout:                 newConfig.WriteTimeout, 
            IdleTimeout:                  newConfig.IdleTimeout, 
            MaxHeaderBytes:               newConfig.MaxHeaderBytes, 
            TLSNextProto:                 newConfig.TLSNextProto, 
            ConnState:                    newConfig.ConnState, 
            ErrorLog:                     newConfig.ErrorLog, 
            BaseContext:                  newConfig.BaseContext, 
            ConnContext:                  newConfig.ConnContext,  
        },
    }
    router.Root = router
    return router
}
