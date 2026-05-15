package grpc

import (
	"context"

	pb "github.com/Enziofael/eatfit/backend/api/gen/eatfit/v1"
	"github.com/Enziofael/eatfit/backend/internal/model"
	"github.com/Enziofael/eatfit/backend/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProfileHandler struct {
	pb.UnimplementedProfileServiceServer
	profileService *service.ProfileService
}

func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	p, err := h.profileService.GetProfile(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profile: %v", err)
	}

	return &pb.GetProfileResponse{Profile: convertToProto(p)}, nil
}

func (h *ProfileHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	return &pb.UpdateProfileResponse{Success: true}, nil
}

func (h *ProfileHandler) SetWeight(ctx context.Context, req *pb.SetWeightRequest) (*pb.SetWeightResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	var note *string
	if n := req.GetNote(); n != "" {
		note = &n
	}

	weight, err := h.profileService.UpdateWeight(ctx, userID, req.GetWeight(), note)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update weight: %v", err)
	}

	return &pb.SetWeightResponse{Success: true, CurrentWeight: weight}, nil
}

func (h *ProfileHandler) SetNorms(ctx context.Context, req *pb.SetNormsRequest) (*pb.SetNormsResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	normsReq := &model.SetNormsRequest{
		Calories: req.GetCalories(),
		Proteins: req.GetProteins(),
		Fats:     req.GetFats(),
		Carbs:    req.GetCarbs(),
		Water:    req.GetWater(),
	}

	norms, err := h.profileService.SetNorms(ctx, userID, normsReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update norms: %v", err)
	}

	return &pb.SetNormsResponse{
		Success: true,
		Norms: &pb.NormsData{
			Calories:  norms.Calories,
			Proteins:  norms.Proteins,
			Fats:      norms.Fats,
			Carbs:     norms.Carbs,
			Water:     norms.Water,
			UpdatedAt: timestamppb.New(norms.UpdatedAt),
		},
	}, nil
}

func (h *ProfileHandler) GetWeightHistory(ctx context.Context, req *pb.GetWeightHistoryRequest) (*pb.GetWeightHistoryResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	entries, err := h.profileService.GetWeightHistory(ctx, userID, int(req.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get weight history: %v", err)
	}

	var pbEntries []*pb.WeightEntry
	for _, e := range entries {
		pe := &pb.WeightEntry{
			Weight:     e.Weight,
			RecordedAt: timestamppb.New(e.RecordedAt),
		}
		if e.Note != nil {
			pe.Note = *e.Note
		}
		pbEntries = append(pbEntries, pe)
	}

	return &pb.GetWeightHistoryResponse{Entries: pbEntries}, nil
}

func convertToProto(p *model.Profile) *pb.ProfileData {
	pd := &pb.ProfileData{
		UserId:    p.UserID.String(),
		Login:     p.Login,
		Email:     p.Email,
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}
	if p.FirstName != nil {
		pd.FirstName = *p.FirstName
	}
	if p.LastName != nil {
		pd.LastName = *p.LastName
	}
	if p.AvatarURL != nil {
		pd.AvatarUrl = *p.AvatarURL
	}
	if p.Height != nil {
		pd.Height = *p.Height
	}
	if p.BirthDate != nil {
		pd.BirthDate = p.BirthDate.Format("2006-01-02")
	}
	if p.Age != nil {
		pd.Age = int32(*p.Age)
	}
	if p.Gender != nil {
		pd.Gender = *p.Gender
	}
	if p.Bio != nil {
		pd.Bio = *p.Bio
	}
	if p.CurrentWeight != nil {
		pd.CurrentWeight = *p.CurrentWeight
	}
	if p.Norms != nil {
		pd.Norms = &pb.NormsData{
			Calories:  p.Norms.Calories,
			Proteins:  p.Norms.Proteins,
			Fats:      p.Norms.Fats,
			Carbs:     p.Norms.Carbs,
			Water:     p.Norms.Water,
			UpdatedAt: timestamppb.New(p.Norms.UpdatedAt),
		}
	}
	return pd
}
