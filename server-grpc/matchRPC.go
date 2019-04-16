package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"strings"

	leagueapi "github.com/dannyrsu/league-api"
	pb "github.com/dannyrsu/league-grpc-server/leagueservice"
	"github.com/golang/protobuf/jsonpb"
)

func constructMatchResponse(match interface{}) *pb.GetMatchResponse {
	m := map[string]interface{}{
		"match": match,
	}

	jsonbytes, err := json.Marshal(m)

	if err != nil {
		log.Fatalf("Error marshaling data: %v", err)
		panic(err)
	}

	result := &pb.GetMatchResponse{}

	r := strings.NewReader(string(jsonbytes))
	if err := jsonpb.Unmarshal(r, result); err != nil {
		log.Fatalf("Error unmarshaling to GetMatchResponse: %v", err)
		panic(err)
	}

	return result
}

func (*leagueServer) GetMatch(ctx context.Context, req *pb.GetMatchRequest) (*pb.GetMatchResponse, error) {
	match := leagueapi.GetGameData(req.GetMatchId(), req.GetRegion())
	res := constructMatchResponse(match)

	return res, nil
}

func (*leagueServer) GetMatchBiDirectional(stream pb.LeagueApi_GetMatchBiDirectionalServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		match := leagueapi.GetGameData(req.GetMatchId(), req.GetRegion())
		sendErr := stream.Send(constructMatchResponse(match))

		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", err)
		}
	}
}
