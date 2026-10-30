package main

import (
	"context"
	"fmt"

	apiv1 "github.com/bradley-adams/halfback/gen/go/halfback/api/v1"
	matchv1 "github.com/bradley-adams/halfback/gen/go/halfback/match/v1"
	matchv1connect "github.com/bradley-adams/halfback/gen/go/halfback/match/v1/matchv1connect"
)

// MatchPublicServer implements MatchPublicServiceHandler, the public
// contract served as REST. It calls MatchService internally over
// Connect and maps the canonical response to the public shape.
type MatchPublicServer struct {
	matchClient matchv1connect.MatchServiceClient
}

func NewMatchPublicServer(matchClient matchv1connect.MatchServiceClient) *MatchPublicServer {
	return &MatchPublicServer{matchClient: matchClient}
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
