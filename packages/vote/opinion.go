package vote

import (
	"context"
)

// OpinionGiver gives opinions about the given IDs.
type OpinionGiver interface {
	// Query queries the OpinionGiver for its opinions on the given IDs.
	// The passed in context can be used to signal cancellation of the query.
	Query(ctx context.Context, ids []string) (Opinions, error)
	// ID returns the ID of the opinion giver.
	ID() string
}

// QueriedOpinions represents queried opinions from a given opinion giver.
type QueriedOpinions struct {
	// The ID of the opinion giver.
	OpinionGiverID string `json:"opinion_giver_id"`
	// The map of IDs to opinions.
	Opinions map[string]Opinion `json:"opinions"`
	// The amount of times the opinion giver's opinion has counted.
	// Usually this number is 1 but due to randomization of the queried opinion givers,
	// the same opinion giver's opinions might be taken into account multiple times.
	TimesCounted int `json:"times_counted"`
}

// OpinionGiverFunc is a function which gives a slice of OpinionGivers or an error.
type OpinionGiverFunc func() ([]OpinionGiver, error)

// Opinions is a slice of Opinion.
type Opinions []Opinion

// Opinion is an opinion about a given thing.
type Opinion byte

const (
	Like    Opinion = 1 << 0
	Dislike Opinion = 1 << 1
	Unknown Opinion = 1 << 2
)

// ConvertInt32Opinion converts the given int32 to an Opinion.
func ConvertInt32Opinion(x int32) Opinion {
	switch {
	case x == 1<<0:
		return Like
	case x == 1<<1:
		return Dislike
	}
	return Unknown
}
