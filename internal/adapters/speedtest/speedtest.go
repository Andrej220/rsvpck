package speedtest

//import (
//	"context"
//
//	"github.com/azargarov/rsvpck/pkg/autostr"
//	"github.com/showwin/speedtest-go/speedtest"
//	"errors"
//	"fmt"
//)
//
//type SpeedtestResult struct{
//	Latency 		string  `string:"include"`
//	DownloadMbps	string		   `string:"include"`
//	UploadMbps		string		   `string:"include"`
//}
//
//type SpeedtestChecker struct{}
//
//func (s *SpeedtestResult)String() string{
//	autostrCfg := autostr.Config{Separator: autostr.Ptr("\n"), FieldValueSeparator: autostr.Ptr(" : "), PrettyPrint: true}
//	return autostr.String(s,autostrCfg)
//}
//
//func (s *SpeedtestChecker) Run(ctx context.Context) (*SpeedtestResult, error) {
//	client := speedtest.New()
//	servers, err := client.FetchServers()
//	if err != nil {
//		return nil, err
//	}
//	targets, err := servers.FindServer(nil) 
//	if err != nil {
//		return nil, err
//	}
//	if len(targets) == 0 {
//		return nil, errors.New("no speedtest server found")
//	}
//	server := targets[0]
//	if err := server.PingTest(nil); err != nil {
//		return nil, err
//	}
//	if err := server.DownloadTest(); err != nil {
//		return nil, err
//	}
//	if err := server.UploadTest(); err != nil {
//		return nil, err
//	}
//	return &SpeedtestResult{
//		Latency: fmt.Sprintf("%0.2f ms",float64(server.Latency.Seconds() * 1000)),
//		DownloadMbps: toMbps(server.DLSpeed),
//		UploadMbps:   toMbps(server.ULSpeed),
//	}, nil
//}
//
//func toMbps(bytesPerSec speedtest.ByteRate) string {
//    return fmt.Sprintf("%0.1f",bytesPerSec * 8 / 1e6)
//}