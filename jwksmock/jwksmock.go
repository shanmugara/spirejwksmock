package jwksmock

import (
	"context"

	"github.com/shanmugara/spireauthlib"
	"github.com/sirupsen/logrus"
	"github.com/spiffe/go-spiffe/v2/bundle/jwtbundle"
	"github.com/spiffe/go-spiffe/v2/svid/jwtsvid"
)

func FetchMyJWT(ctx context.Context) (*jwtbundle.Set, *jwtsvid.SVID, error) {
	cauath := spireauthlib.ClientAuth{Logger: logrus.New()}
	bundle, myjwtsvid, err := cauath.GetJWT(ctx)
	if err != nil {
		return nil, nil, err
	}
	return bundle, myjwtsvid, nil
}
