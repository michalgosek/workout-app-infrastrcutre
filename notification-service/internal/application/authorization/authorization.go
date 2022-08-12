package authorization

import (
	"context"
	"errors"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

type Claim string

func (c Claim) String() string {
	return string(c)
}

const (
	NotifyParticipants             Claim = "trainer:create:notifications"
	ViewParticipantNotifications   Claim = "participant:read:notifications"
	DeleteParticipantNotifications Claim = "participant:delete:notifications"
)

type CustomClaims struct {
	Permissions []string `json:"permissions"`
}

func (c CustomClaims) HasScope(expectedClaims ...Claim) bool {
	if len(expectedClaims) == 0 {
		return false
	}

	lookup := make(map[Claim]bool)
	for _, c := range expectedClaims {
		lookup[c] = true
	}

	sufficient := len(expectedClaims)
	n := 0
	for _, p := range c.Permissions {
		if _, ok := lookup[Claim(p)]; ok {
			n++
		}
	}
	return sufficient == n
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func ValidateJWT() func(next http.Handler) http.Handler {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env auth0 file: %s", err)
	}

	issuerURL, err := url.Parse(os.Getenv("AUTH0_DOMAIN"))
	if err != nil {
		log.Fatalf("failed to parse the issuer url: %s", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(JWTMiddlewareErrorHandler),
	)
	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}

func JWTMiddlewareErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Encountered error while validating JWT: %v", err)

	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, "failed to validate JWT.", http.StatusUnauthorized)
}

func CreateUserClaimsFromToken(ctx context.Context) (*CustomClaims, error) {
	token := ctx.Value(jwtmiddleware.ContextKey{})
	if token == nil {
		return nil, errors.New("missing jwt context key")
	}
	validatedClaims, ok := token.(*validator.ValidatedClaims)
	if !ok {
		return nil, errors.New("invalid ctx token key value")
	}
	claims, ok := validatedClaims.CustomClaims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid validated claims value type")
	}
	return claims, nil
}
