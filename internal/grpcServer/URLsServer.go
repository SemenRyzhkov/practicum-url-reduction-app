package grpcServer

import (
	"context"
	"errors"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/urlservice"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/proto"
)

type URLsServer struct {
	proto.UnimplementedURLsServer
	urlService urlservice.URLService
}

// NewServerImpl конструктор.
func NewServerImpl(urlService urlservice.URLService) *URLsServer {
	return &URLsServer{urlService: urlService}
}

// GetURLByID выполняет получение URL по его ID.
func (u *URLsServer) GetURLByID(ctx context.Context, in *proto.GetURLByIDRequest) (*proto.GetURLByIDResponse, error) {
	urlID := in.UrlId
	var response proto.GetURLByIDResponse
	url, err := u.urlService.GetURLByID(ctx, urlID)
	if err != nil {
		var de *myerrors.DeletedError
		if errors.As(err, &de) {
			return nil, status.Errorf(codes.FailedPrecondition, `URL с ID %s удален`, urlID)
		}
		return nil, status.Errorf(codes.NotFound, `URL с ID %s не найден`, urlID)
	}
	response.Url = url
	return &response, nil
}

// GetAllURL выполняет получение всех сокращенных URL по ID юзера.
func (u *URLsServer) GetAllURL(ctx context.Context, in *proto.GetAllByUserIDRequest) (*proto.GetAllByUserIDResponse, error) {
	userID := in.UserId

	userURLList, notFoundErr := u.urlService.GetAllByUserID(ctx, userID)
	if notFoundErr != nil {
		return nil, status.Errorf(codes.NotFound, ` Не найдено сохраненных URL для пользователя с ID %s `, userID)
	}
	var response proto.GetAllByUserIDResponse
	var protoUserURLList = make([]*proto.FullURL, 0)

	for _, v := range userURLList {
		protoFullUrl := proto.FullURL{OriginalUrl: v.OriginalURL, ShortUrl: v.ShortURL}
		protoUserURLList = append(protoUserURLList, &protoFullUrl)
	}
	response.FullUrls = protoUserURLList
	return &response, nil
}

// ReduceURL выполняет сокращение URL, переданного в текстовом формате.
func (u *URLsServer) ReduceURL(ctx context.Context, in *proto.ReduceAndSaveURLRequest) (*proto.ReduceAndSaveURLResponse, error) {
	userID := in.UserId
	URL := in.Url

	reduceURL, err := u.urlService.ReduceAndSaveURL(ctx, userID, URL)
	if err != nil {
		var ve *myerrors.ViolationError
		if errors.As(err, &ve) {
			return &proto.ReduceAndSaveURLResponse{ShortUrl: URL},
				status.Errorf(codes.AlreadyExists, `URL %s для пользователя с ID %s уже существует`, URL, userID)
		}
		return nil, status.Errorf(codes.Unknown, `Некорректный запрос`)
	}

	return &proto.ReduceAndSaveURLResponse{ShortUrl: reduceURL}, nil
}

// ReduceURLTOJSON выполняет сокращение URL, переданного в JSON-формате.
func (u *URLsServer) ReduceURLTOJSON(ctx context.Context, in *proto.ReduceURLToJSONRequest) (*proto.ReduceURLToJSONResponse, error) {
	urlToReduce := in.UrlRequest.Url
	urlRequest := entity.URLRequest{URL: urlToReduce}
	userID := in.UserId

	urlResponse, err := u.urlService.ReduceURLToJSON(ctx, userID, urlRequest)

	if err != nil {
		var ve *myerrors.ViolationError
		if errors.As(err, &ve) {
			return nil,
				status.Errorf(codes.AlreadyExists, `URL %s для пользователя с ID %s уже существует`, urlToReduce, userID)
		}
		return nil, status.Errorf(codes.Unknown, `Некорректный запрос`)
	}

	return &proto.ReduceURLToJSONResponse{UrlResponse: &proto.URLResponse{Result: urlResponse.Result}}, nil
}

// ReduceSeveralURL выполняет сокращение нескольких URL, переданныф в в JSON-формате.
func (u *URLsServer) ReduceSeveralURL(ctx context.Context, in *proto.ReduceSeveralURLRequest) (*proto.ReduceSeveralURLResponse, error) {
	userID := in.UserId
	var urlWithIDLList = make([]entity.URLWithIDRequest, 0)

	for _, v := range in.UrlWithIdRequests {
		urlWithID := entity.URLWithIDRequest{CorrelationID: v.CorrelationId, OriginalURL: v.OriginalUrl}
		urlWithIDLList = append(urlWithIDLList, urlWithID)
	}
	urlWithIDResponseList, err := u.urlService.ReduceSeveralURL(ctx, userID, urlWithIDLList)

	if err != nil {
		return nil, status.Errorf(codes.Unknown, `Некорректный запрос`)
	}

	var protoUrlWithIDResponseList = make([]*proto.URLWithIDResponse, 0)

	for _, v := range urlWithIDResponseList {
		protoUrlWithIDResp := proto.URLWithIDResponse{CorrelationId: v.CorrelationID, ShortUrl: v.ShortURL}
		protoUrlWithIDResponseList = append(protoUrlWithIDResponseList, &protoUrlWithIDResp)
	}

	return &proto.ReduceSeveralURLResponse{UrlWithIdResponses: protoUrlWithIDResponseList}, nil
}

// RemoveAll выполняет удаление нескольких URL по их ID.
func (u *URLsServer) RemoveAll(ctx context.Context, in *proto.RemoveAllRequest) (*proto.RemoveAllResponse, error) {
	removingErr := u.urlService.RemoveAll(ctx, in.UserId, in.RemovingList)
	if removingErr != nil {
		return nil, status.Errorf(codes.Unknown, `Некорректный запрос`)
	}
	return nil, nil
}

// PingConnection пинг
func (u *URLsServer) PingConnection(ctx context.Context, in *proto.PingConnectionRequest) (*proto.PingConnectionResponse, error) {
	err := u.urlService.PingConnection()

	if err != nil {
		return nil, status.Errorf(codes.Unknown, `Некорректный запрос`)

	}
	return nil, nil
}

// GetStats возвращает количество сокращенных URL и число пользователей сервиса, если IP пользователя
// входит в доверенную подсеть.
func (u *URLsServer) GetStats(ctx context.Context, in *proto.GetStatsRequest) (*proto.GetStatsResponse, error) {
	ipInSIDR, err := ipInSIDR(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, `Некорректный запрос`)
	}

	if len(utils.GetTrustedSubnet()) == 0 || !ipInSIDR {
		return nil, status.Errorf(codes.PermissionDenied, `Ваш айпишник не входит в нашу доверенную подсеть`)

	}

	stats, err := u.urlService.GetStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, `Некорректный запрос`)
	}
	return &proto.GetStatsResponse{Stats: &proto.Stats{Users: int64(stats.Users), Urls: int64(stats.Urls)}}, nil
}

func ipInSIDR(ctx context.Context) (bool, error) {
	var ip string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("X-Real-IP")
		if len(values) > 0 {
			ip = values[0]
		}
	}
	parsedIP := net.ParseIP(ip)

	_, ipv4Net, err := net.ParseCIDR(utils.GetTrustedSubnet())
	if err != nil {
		return false, err
	}

	return ipv4Net.Contains(parsedIP), nil
}
