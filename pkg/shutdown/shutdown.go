package shutdown

import (
	"io"
	"os"
	"os/signal"

	"github.com/IvanKyrylov/user-game-api/pkg/logging"
)

func Graceful(signals []os.Signal, closeItems ...io.Closer) {

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)
	sig := <-sigc
	logging.CommonLog.Printf("Caught signal %s. Shutting down...", sig)

	for _, closer := range closeItems {
		if err := closer.Close(); err != nil {
			logging.ErrorLog.Fatalf("failed to close %v: %v", closer, err)
		}
	}
}
