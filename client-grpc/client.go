package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/dannyrsu/league-api/leagueservice"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	c := pb.NewLeagueApiClient(cc)

	getMultpleStats(c)
}

func getSingleStat(c pb.LeagueApiClient) {
	req := &pb.GetSummonerStatsRequest{
		SummonerName: "WhirlingDeth",
		Region:       "na1",
	}

	res, err := c.GetSummonerStats(context.Background(), req)

	if err != nil {
		log.Fatalf("Error calling for stats: %v", err)
	}

	log.Printf("Stats: %v", res.SummonerProfile)
}

func getMultpleStats(c pb.LeagueApiClient) {
	stream, err := c.GetSummonerStatsBiDirectional(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream and calling Stats: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		summonerNames := map[string]string{
			"whirlingdeth": "na1",
			"admiremeaxes": "na1",
		}
		for key, value := range summonerNames {
			sendErr := stream.Send(&pb.GetSummonerStatsRequest{
				SummonerName: key,
				Region:       value,
			})
			time.Sleep(1000 * time.Millisecond)

			if sendErr != nil {
				log.Fatalf("Error sending req to stream: %v", sendErr)
			}
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Problem while reading server stream: %v", err)
			}

			summonerStat := res.GetSummonerProfile()
			fmt.Printf("Received summoner stat: %v\n", summonerStat)
		}

		close(waitc)
	}()

	<-waitc
}
