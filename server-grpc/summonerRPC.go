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

func constructSummonerStatsResponse(summonerProfile leagueapi.SummonerProfile) *pb.GetSummonerStatsResponse {

	m := map[string]interface{}{
		"summonerProfile": summonerProfile,
	}

	jsonbytes, err := json.Marshal(m)

	if err != nil {
		log.Fatalf("Error marshaling data: %v", err)
		panic(err)
	}

	result := &pb.GetSummonerStatsResponse{}

	r := strings.NewReader(string(jsonbytes))
	if err := jsonpb.Unmarshal(r, result); err != nil {
		log.Fatalf("Error unmarshaling to GetSummonerStatsResponse: %v", err)
		panic(err)
	}

	return result
}

func (*leagueServer) GetSummonerStats(ctx context.Context, req *pb.GetSummonerStatsRequest) (*pb.GetSummonerStatsResponse, error) {
	summonerProfile := leagueapi.GetSummonerProfile(req.GetSummonerName(), req.GetRegion())
	res := constructSummonerStatsResponse(summonerProfile)

	return res, nil
}

func (*leagueServer) GetSummonerStatsBiDirectional(stream pb.LeagueApi_GetSummonerStatsBiDirectionalServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		summonerProfile := leagueapi.GetSummonerProfile(req.GetSummonerName(), req.GetRegion())
		sendErr := stream.Send(constructSummonerStatsResponse(summonerProfile))

		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", err)
		}
	}
}
