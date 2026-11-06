package main

import (
	"context"
	"fmt"

	apiv1 "github.com/bradley-adams/halfback/gen/go/halfback/api/v1"
	matchv1 "github.com/bradley-adams/halfback/gen/go/halfback/match/v1"
	matchv1connect "github.com/bradley-adams/halfback/gen/go/halfback/match/v1/matchv1connect"
	scoreboardv1 "github.com/bradley-adams/halfback/gen/go/halfback/scoreboard/v1"
	scoreboardv1connect "github.com/bradley-adams/halfback/gen/go/halfback/scoreboard/v1/scoreboardv1connect"
)

// MatchPublicServer implements MatchPublicServiceHandler, the public
// contract served as REST. It calls the internal services over Connect
// and maps their canonical responses to the public shape.
type MatchPublicServer struct {
	matchClient      matchv1connect.MatchServiceClient
	scoreboardClient scoreboardv1connect.ScoreboardServiceClient
}

func NewMatchPublicServer(
	matchClient matchv1connect.MatchServiceClient,
	scoreboardClient scoreboardv1connect.ScoreboardServiceClient,
) *MatchPublicServer {
	return &MatchPublicServer{
		matchClient:      matchClient,
		scoreboardClient: scoreboardClient,
	}
}

func (s *MatchPublicServer) ListMatchEvents(
	ctx context.Context,
	req *apiv1.ListMatchEventsRequest,
) (*apiv1.ListMatchEventsResponse, error) {
	internalReq := &matchv1.ListMatchEventsRequest{
		MatchId: req.MatchId,
	}

	internalResp, err := s.matchClient.ListMatchEvents(ctx, internalReq)
	if err != nil {
		return nil, fmt.Errorf("call match service: %w", err)
	}

	resp := &apiv1.ListMatchEventsResponse{}
	for _, e := range internalResp.Events {
		resp.Events = append(resp.Events, &apiv1.MatchEventSummary{
			EventId:          e.EventId,
			EventType:        e.EventType,
			EventTimeSeconds: e.EventTimeSeconds,
		})
	}

	return resp, nil
}

func (s *MatchPublicServer) ListMatches(
	ctx context.Context,
	_ *apiv1.ListMatchesRequest,
) (*apiv1.ListMatchesResponse, error) {
	internalResp, err := s.scoreboardClient.ListMatches(ctx, &scoreboardv1.ListMatchesRequest{})
	if err != nil {
		return nil, fmt.Errorf("call scoreboard service: %w", err)
	}

	resp := &apiv1.ListMatchesResponse{}
	for _, m := range internalResp.Matches {
		resp.Matches = append(resp.Matches, &apiv1.MatchScore{
			MatchId:     m.MatchId,
			EventCount:  m.EventCount,
			TotalPoints: m.TotalPoints,
		})
	}

	return resp, nil
}
