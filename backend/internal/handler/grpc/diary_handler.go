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

type DiaryHandler struct {
	pb.UnimplementedDiaryServiceServer
	diaryService *service.DiaryService
}

func NewDiaryHandler(diaryService *service.DiaryService) *DiaryHandler {
	return &DiaryHandler{diaryService: diaryService}
}

func (h *DiaryHandler) AddConsumption(ctx context.Context, req *pb.AddConsumptionRequest) (*pb.AddConsumptionResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	record, err := h.diaryService.AddConsumption(ctx, &model.AddConsumptionRequest{
		UserID:   userID,
		MealID:   req.GetMealId(),
		MealName: req.GetMealName(),
		Amount:   req.GetAmount(),
		Calories: req.GetCalories(),
		Proteins: req.GetProteins(),
		Fats:     req.GetFats(),
		Carbs:    req.GetCarbs(),
		Water:    req.GetWater(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add: %v", err)
	}

	return &pb.AddConsumptionResponse{
		Success: true,
		Record:  convertRecordToProto(record),
	}, nil
}

func (h *DiaryHandler) ListConsumptions(ctx context.Context, req *pb.ListConsumptionsRequest) (*pb.ListConsumptionsResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	groups, err := h.diaryService.ListConsumptions(ctx, userID, int(req.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list: %v", err)
	}

	var pbGroups []*pb.ConsumptionGroup
	for _, g := range groups {
		var pbRecords []*pb.ConsumptionRecord
		for _, r := range g.Records {
			pbRecords = append(pbRecords, convertRecordToProto(&r))
		}
		pbGroups = append(pbGroups, &pb.ConsumptionGroup{
			Date:    g.Date,
			Records: pbRecords,
		})
	}

	return &pb.ListConsumptionsResponse{Groups: pbGroups}, nil
}

func (h *DiaryHandler) DeleteConsumption(ctx context.Context, req *pb.DeleteConsumptionRequest) (*pb.DeleteConsumptionResponse, error) {
	recordID, err := uuid.Parse(req.GetRecordId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid record ID")
	}

	if err := h.diaryService.DeleteConsumption(ctx, recordID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete: %v", err)
	}

	return &pb.DeleteConsumptionResponse{Success: true}, nil
}

func convertRecordToProto(r *model.ConsumptionRecord) *pb.ConsumptionRecord {
	rec := &pb.ConsumptionRecord{
		RecordId:   r.ID.String(),
		MealName:   r.MealName,
		Amount:     r.Amount,
		Calories:   r.Calories,
		Proteins:   r.Proteins,
		Fats:       r.Fats,
		Carbs:      r.Carbs,
		Water:      r.Water,
		ConsumedAt: timestamppb.New(r.ConsumedAt),
	}
	if r.MealID != nil {
		rec.MealId = r.MealID.String()
	}
	return rec
}
