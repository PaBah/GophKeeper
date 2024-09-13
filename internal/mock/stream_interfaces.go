package mock

import (
	pb "github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
	"google.golang.org/grpc"
)

type MockServerStream interface {
	grpc.ServerStream
}

type GophKeeperService_UploadFileClient interface {
	Send(request *pb.UploadFileRequest) error
	Recv() (*pb.UploadFileResponse, error)
	grpc.ClientStream
}

type GophKeeperService_SubscribeToChangesClient interface {
	Recv() (*pb.SubscribeToChangesResponse, error)
	grpc.ClientStream
}

type GophKeeperService_DownloadFileClient interface {
	Recv() (*pb.DownloadFileResponse, error)
	grpc.ClientStream
}

type GophKeeperService_SubscribeToChangesServer interface {
	Send(*pb.SubscribeToChangesResponse) error
	grpc.ServerStream
}

type GophKeeperService_UploadFileServer interface {
	Recv() (*pb.UploadFileRequest, error)
	Send(*pb.UploadFileResponse) error
	grpc.ServerStream
}

type GophKeeperService_DownloadFileServer interface {
	Send(*pb.DownloadFileResponse) error
	grpc.ServerStream
}
