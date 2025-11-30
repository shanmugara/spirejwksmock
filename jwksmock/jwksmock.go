package jwksmock

import (
	"context"

	"github.com/shanmugara/spireauthlib"
	"github.com/sirupsen/logrus"
	"github.com/spiffe/go-spiffe/v2/bundle/jwtbundle"
	"github.com/spiffe/go-spiffe/v2/svid/jwtsvid"
)

func FetchMyJWT(ctx context.Context, audience string) (*jwtbundle.Set, *jwtsvid.SVID, error) {
	Cauth := spireauthlib.ClientAuth{Logger: logrus.New()}
	bundle, myjwtsvid, err := Cauth.GetJWT(ctx, audience)
	if err != nil {
		return nil, nil, err
	}
	return bundle, myjwtsvid, nil
}

func VlidateMyJWT(ctx context.Context, myjwtsvid *jwtsvid.SVID, bundle *jwtbundle.Set) error {
	cauth := spireauthlib.ClientAuth{Logger: logrus.New()}
	return cauth.ValidateJWT(bundle, myjwtsvid)
}
