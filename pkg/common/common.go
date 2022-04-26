package common

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func JwtValid(token string) error {

	claims := jwt.StandardClaims{}
	parser := &jwt.Parser{}
	t, _, err := parser.ParseUnverified(token, &claims)

	// there was an issue parsing the token
	if err != nil {
		return fmt.Errorf("unable to parse the jwt token: %w", err)
	}

	return t.Claims.Valid()
}

const (
	RepoUrlAnnotationKey  = "arlon.io/repo-url"
	RepoPathAnnotationKey = "arlon.io/repo-path"

	ProfileAnnotationKey     = "arlon.io/profile"
	ClusterSpecAnnotationKey = "arlon.io/clusterspec"
)
