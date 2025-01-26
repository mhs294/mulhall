package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mhs294/mulhall/internals/repos"
	"github.com/mhs294/mulhall/internals/types"
)

// UserAuthMiddleware is responsible for handling user authentication in view/API requests.
type UserAuthMiddleware struct {
	sessRepo *repos.SessionRepository
}

// NewUserAuthMiddleware creates a new UserAuthMiddleware instance and returns a pointer to it.
//
// r is the SessionRepository used to look up Sessions for User authentication.
func NewUserAuthMiddleware(r *repos.SessionRepository) *UserAuthMiddleware {
	return &UserAuthMiddleware{sessRepo: r}
}

// ViewAuth handles the validation of user authentication for View (webpage) requests,
// which will involve redirection to a Login page if unsuccessful.
func (m *UserAuthMiddleware) ViewAuth(ctx *gin.Context) {
	sessID := types.SessionID(ctx.GetHeader("Session-ID"))
	if err := m.userAuth(sessID); err != nil {
		// User is unauthorized, redirect to Login view
		setLoginRedirect(ctx)
		return
	}

	// Continue with subsequent middleware/request handling
	ctx.Next()
}

func (m *UserAuthMiddleware) APIAuth(ctx *gin.Context) {
	sessID := types.SessionID(ctx.GetHeader("Session-ID"))
	if err := m.userAuth(sessID); err != nil {
		// User is unauthorized, return status to caller
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Continue with subsequent middleware/request handling
	ctx.Next()
}

func (m *UserAuthMiddleware) userAuth(sessID types.SessionID) error {
	// Verify the Session ID is present
	if len(sessID) == 0 {
		return &types.MissingSessionIDError{}
	}

	// Look up the Session corresponding to the provided ID
	sess, err := m.sessRepo.GetSession(sessID)
	if err != nil {
		return err
	}
	if sess == nil {
		return &types.SessionNotFoundError{}
	}

	// Verify the Session is active (i.e. - has not expired)
	if sess.Expiration.Before(time.Now().UTC()) {
		return &types.SessionExpiredError{}
	}

	return nil
}

func setLoginRedirect(ctx *gin.Context) {
	ctx.Header("Location", "/login")
	ctx.AbortWithStatus(http.StatusFound)
}
