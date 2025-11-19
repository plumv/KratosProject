package service

import (
	"api/base"
	"boss/internal/biz"
	"context"
	"strings"

	pb "api/boss"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type UserService struct {
	pb.UnimplementedUserServer
	uc  *biz.UserUsecase
	log *log.Helper
}

func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	id, err := s.uc.CreateUser(ctx, &biz.User{
		Username: req.Name,
		Password: req.Password,
		Age:      req.Age,
	})
	if err != nil {
		return &pb.CreateUserReply{
			R: base.ERROR.FillMsg("无法保存"),
		}, nil
	}
	user, err := s.uc.FindUser(ctx, id)
	if err != nil {
		return &pb.CreateUserReply{
			R: base.ERROR.FillMsg("保存失败"),
		}, nil
	}
	return &pb.CreateUserReply{
		R: base.SUCCESS,
		Data: &pb.UserReply{
			Id:    user.ID.String(),
			Name:  user.Username,
			Age:   user.Age,
			Email: "",
		},
	}, nil
}
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	id := uuid.MustParse(req.Id)
	err := s.uc.UpdateUser(ctx, id, &biz.User{
		Username: req.Name,
		Age:      req.Age,
	})
	if err != nil {
		return &pb.UpdateUserReply{
			R: base.ERROR.FillMsg("无法更新"),
		}, nil
	}
	user, err := s.uc.FindUser(ctx, id)
	if err != nil {
		return &pb.UpdateUserReply{
			R: base.ERROR.FillMsg("更新失败"),
		}, nil
	}
	return &pb.UpdateUserReply{
		R: base.SUCCESS,
		Data: &pb.UserReply{
			Id:    user.ID.String(),
			Name:  user.Username,
			Age:   user.Age,
			Email: "",
		},
	}, nil
}
func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	err := s.uc.DeleteUser(ctx, uuid.MustParse(req.Id))
	if err != nil {
		return &pb.DeleteUserReply{
			R: base.ERROR.FillMsg("删除失败"),
		}, nil
	}
	return &pb.DeleteUserReply{
		R: base.SUCCESS,
	}, nil
}
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	user, err := s.uc.FindUser(ctx, uuid.MustParse(req.Id))
	if err != nil {
		return &pb.GetUserReply{
			R: base.ERROR.FillMsg("查询失败"),
		}, nil
	}
	return &pb.GetUserReply{
		R: base.SUCCESS,
		Data: &pb.UserReply{
			Id:    user.ID.String(),
			Name:  user.Username,
			Age:   user.Age,
			Email: "",
		},
	}, nil
}
func (s *UserService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserReply, error) {
	filter := biz.UserFilter{}
	if query := req.Query; query != nil {
		if q := query.Name; q != "" {
			filter.Name = &q
		}
		if q := query.Age; q != 0 {
			filter.AgeEQ = &q
		}
	}
	var order []*biz.Order
	if o := req.Sort; o != nil {
		order = make([]*biz.Order, 0)
		if o.Order != "" {
			for _, of := range strings.Split(o.Order, ",") {
				order = append(order, &biz.Order{
					Field: of,
					Desc:  o.Sort,
				})
			}
		}
	}

	users, err := s.uc.ListUser(ctx, &filter, &order)
	if err != nil {
		return &pb.ListUserReply{
			R: base.ERROR.FillMsg("查询失败"),
		}, nil
	}
	var userReply = make([]*pb.UserReply, len(users))
	for i, u := range users {
		userReply[i] = &pb.UserReply{
			Id:    u.ID.String(),
			Name:  u.Username,
			Age:   u.Age,
			Email: "",
		}
	}
	return &pb.ListUserReply{
		R:    base.SUCCESS,
		Data: userReply,
	}, nil
}
func (s *UserService) PageUser(ctx context.Context, req *pb.PageUserRequest) (*pb.PageUserReply, error) {
	filter := biz.UserFilter{}
	if query := req.Query; query != nil {
		if q := query.Name; q != "" {
			filter.Name = &q
		}
		if q := query.Age; q != 0 {
			filter.AgeEQ = &q
		}
	}
	var order []*biz.Order
	if o := req.Sort; o != nil {
		order = make([]*biz.Order, 0)
		if o.Order != "" {
			for _, of := range strings.Split(o.Order, ",") {
				order = append(order, &biz.Order{
					Field: of,
					Desc:  o.Sort,
				})
			}
		}
	}
	page := biz.Page{Page: 0, Limit: 10, Orders: &order}
	if p := req.Page; p != nil {
		page.Page = p.Page
		page.Limit = p.PageSize
	}

	users, total, err := s.uc.PageUser(ctx, &filter, &page)
	if err != nil {
		return &pb.PageUserReply{
			R: base.ERROR.FillMsg("查询失败"),
		}, nil
	}
	var userReply = make([]*pb.UserReply, len(users))
	for i, u := range users {
		userReply[i] = &pb.UserReply{
			Id:    u.ID.String(),
			Name:  u.Username,
			Age:   u.Age,
			Email: "",
		}
	}
	return &pb.PageUserReply{
		R:    base.SUCCESS,
		Page: &base.PageResp{Total: int64(total)},
		Data: userReply,
	}, nil
}
