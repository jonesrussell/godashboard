package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jonesrussell/dashboard/internal/logger"
)

// Handler manages graceful shutdown of the application
type Handler struct {
	logger     logger.Logger
	shutdownWg sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	sigChan    chan os.Signal
}

// New creates a new shutdown handler
func New(log logger.Logger) *Handler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Handler{
		logger: log,
		ctx:    ctx,
		cancel: cancel,
		// Buffer size of 1 is sufficient for signal notification
		sigChan: make(chan os.Signal, 1),
	}
}

// HandleSignals starts listening for OS signals
func (h *Handler) HandleSignals() {
	// Register for common interrupt signals plus SIGQUIT for stack dumps
	signal.Notify(h.sigChan,
		syscall.SIGINT,  // Terminal interrupt (Ctrl+C)
		syscall.SIGTERM, // Termination request
		syscall.SIGHUP,  // Terminal disconnect
		syscall.SIGQUIT, // Terminal quit (Ctrl+\)
	)

	go func() {
		sig := <-h.sigChan
		h.logger.Info("received signal", logger.NewField("signal", sig.String()))

		// Special handling for SIGQUIT to log stack trace
		if sig == syscall.SIGQUIT {
			h.logger.Info("received quit signal, dumping stacks")
			// Let default Go handler print stack trace
			signal.Reset(syscall.SIGQUIT)
			// Resend signal to get stack trace
			p, _ := os.FindProcess(os.Getpid())
			p.Signal(syscall.SIGQUIT)
		}

		h.cancel()
	}()
}

// Context returns the shutdown context
func (h *Handler) Context() context.Context {
	return h.ctx
}

// AddToWaitGroup adds a component to the shutdown wait group
func (h *Handler) AddToWaitGroup() {
	h.shutdownWg.Add(1)
}

// Done marks a component as done in the shutdown wait group
func (h *Handler) Done() {
	h.shutdownWg.Done()
}

// Wait waits for all components to finish cleanup
func (h *Handler) Wait() {
	h.shutdownWg.Wait()
}

// Close stops signal handling and waits for cleanup
func (h *Handler) Close() {
	// Stop receiving signals
	signal.Stop(h.sigChan)
	// Reset all signal handlers to default
	signal.Reset()
	// Cancel context if not already done
	h.cancel()
	// Wait for cleanup
	h.Wait()
}
