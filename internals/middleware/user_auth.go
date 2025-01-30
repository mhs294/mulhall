package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhs294/mulhall/internals/repos"
	"github.com/mhs294/mulhall/internals/types"
)

// UserAuthMiddleware is responsible for handling user authentication in view/API requests.
type UserAuthMiddleware struct {
	logger   *log.Logger
	sessRepo *repos.SessionRepository
}

// NewUserAuthMiddleware creates a new UserAuthMiddleware instance and returns a pointer to it.
//
// l is the pointer to the [log.Logger] that will be used at runtime by the UserAuthMiddleware.
//
// r is the SessionRepository used to look up Sessions for User authentication.
func NewUserAuthMiddleware(l *log.Logger, r *repos.SessionRepository) *UserAuthMiddleware {
	return &UserAuthMiddleware{logger: l, sessRepo: r}
}

// ViewAuth handles the validation of user authentication for View (webpage) requests,
// which will involve redirection to a Login page if unsuccessful.
func (m *UserAuthMiddleware) ViewAuth(ctx *gin.Context) {
	sess, err := m.userAuth(ctx)
	if err != nil {
		// User is unauthorized, redirect to Login view
		ctx.Header("Location", "/login")
		ctx.AbortWithStatus(http.StatusFound)
		return
	}

	// Add the User's Session to the request context, continue with subsequent middleware/request handling
	ctx.Set("session", sess)
	ctx.Next()
}

// APIAuth handles the validation of user authentication for API requests,
// which should defer to the calling process to determine how any authentication
// failures should be handled.
func (m *UserAuthMiddleware) APIAuth(ctx *gin.Context) {
	sess, err := m.userAuth(ctx)
	if err != nil {
		// User is unauthorized, return status to caller
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Add the User's Session to the request context, continue with subsequent middleware/request handling
	ctx.Set("session", sess)
	ctx.Next()
}

func (m *UserAuthMiddleware) userAuth(ctx *gin.Context) (*types.Session, error) {
	// Read the Session ID cookie
	sessCookie, err := ctx.Cookie("mulhall.sessionID")

	if err != nil {
		// Couldn't read session cookie, assume user is unauthorized
		ctx.AbortWithStatus(http.StatusUnauthorized)
		m.logger.Printf("failed to read session cookie: %v", err)
		return nil, err
	}

	// Verify the Session ID is present
	sessID := types.SessionID(sessCookie)
	if len(sessID) == 0 {
		return nil, &types.MissingSessionIDError{}
	}

	// Look up the Session corresponding to the provided ID
	sess, err := m.sessRepo.GetByID(sessID)
	if err != nil {
		m.logger.Printf("failed to load session from database: %v", err)
		return nil, err
	}

	return sess, nil
}
