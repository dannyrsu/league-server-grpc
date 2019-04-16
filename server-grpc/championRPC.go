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

func constructChampionResponse(champion interface{}) *pb.GetChampionByKeyResponse {
	m := map[string]interface{}{
		"champion": champion,
	}

	jsonbytes, err := json.Marshal(m)

	if err != nil {
		log.Fatalf("Error marshaling data: %v", err)
		panic(err)
	}

	result := &pb.GetChampionByKeyResponse{}

	r := strings.NewReader(string(jsonbytes))
	if err := jsonpb.Unmarshal(r, result); err != nil {
		log.Fatalf("Error unmarshaling to GetChampionByKeyResult: %v", err)
		panic(err)
	}

	return result
}

func (*leagueServer) GetChampionByKey(ctx context.Context, req *pb.GetChampionByKeyRequest) (*pb.GetChampionByKeyResponse, error) {
	res := constructChampionResponse(leagueapi.GetChampionByKey(req.GetChampionKey()))

	return res, nil
}

func (*leagueServer) GetChampionByKeyBiDirectional(stream pb.LeagueApi_GetChampionByKeyBiDirectionalServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		sendErr := stream.Send(constructChampionResponse(leagueapi.GetChampionByKey(req.GetChampionKey())))

		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", err)
		}
	}
}
