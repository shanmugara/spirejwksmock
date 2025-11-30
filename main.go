package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"spirejwksmock/jwksmock"
	"syscall"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func main() {
	// set default logger to stdout without timestamp
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	log.Println("Starting spirejwksmock")

	// create a context that is cancelled on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// use a ticker instead of sleep so we can listen for shutdown
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down: context cancelled")
			return
		case <-ticker.C:
			// per-iteration timeout to avoid hanging indefinitely
			itCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			bundles, myjwtsvid, err := jwksmock.FetchMyJWT(itCtx, "omegaaudience")
			cancel()
			if err != nil {
				log.Println("Error fetching JWT:", err)
				// don't exit; wait for next tick (could implement backoff here)
				continue
			}
			log.Println("Fetched JWT Bundle:", bundles)
			for _, b := range bundles.Bundles() {
				bb, err := b.Marshal()
				if err != nil {
					log.Println("Error marshaling bundle:", err)
					continue
				}
				log.Println("Bundle:", string(bb))
			}
			log.Println("Fetched My JWT SVID:", myjwtsvid.Marshal())
			log.Println("Validating My JWT SVID...")
			err = jwksmock.VlidateMyJWT(itCtx, myjwtsvid, bundles)
			if err != nil {
				log.Println("Error validating JWT SVID:", err)
			} else {
				log.Println("JWT SVID validation successful")
			}
			jwtMarshal := myjwtsvid.Marshal()
			// Parse and validate the token locally using the keys in the fetched bundles.
			// The keyFunc looks up the token "kid" header in the bundles to return the
			// corresponding public key for signature verification.
			var claims jwt.MapClaims
			parsedToken, err := jwt.ParseWithClaims(jwtMarshal, &claims, func(token *jwt.Token) (any, error) {
				kid, _ := token.Header["kid"].(string)
				if kid == "" {
					return nil, fmt.Errorf("token missing kid header")
				}
				// search bundles for the key id
				for _, b := range bundles.Bundles() {
					if key, ok := b.FindJWTAuthority(kid); ok {
						return key, nil
					}
				}
				return nil, fmt.Errorf("no key found for kid %q", kid)
			})
			if err != nil {
				log.Println("Error parsing token:", err)
			} else if !parsedToken.Valid {
				log.Println("Token is invalid")
			} else {
				log.Println("Parsed token claims:", claims)
			}

		}
	}
}
