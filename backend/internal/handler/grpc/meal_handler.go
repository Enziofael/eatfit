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

type MealHandler struct {
	pb.UnimplementedMealServiceServer
	mealService *service.MealService
}

func NewMealHandler(mealService *service.MealService) *MealHandler {
	return &MealHandler{mealService: mealService}
}

func (h *MealHandler) CreateMeal(ctx context.Context, req *pb.CreateMealRequest) (*pb.CreateMealResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	createReq := &model.CreateMealRequest{
		UserID:      userID,
		Name:        req.GetName(),
		Description: strPtr(req.GetDescription()),
		Recipe:      strPtr(req.GetRecipe()),
		ImageURL:    strPtr(req.GetImageUrl()),
		Calories:    req.GetCalories(),
		Proteins:    req.GetProteins(),
		Fats:        req.GetFats(),
		Carbs:       req.GetCarbs(),
		Water:       req.GetWater(),
	}

	for _, c := range req.GetComponents() {
		compID, err := uuid.Parse(c.GetComponentMealId())
		if err != nil {
			continue
		}
		createReq.Components = append(createReq.Components, model.MealComponentInput{
			ComponentMealID: compID,
			Amount:          c.GetAmount(),
		})
	}

	meal, err := h.mealService.CreateMeal(ctx, createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create meal: %v", err)
	}

	return &pb.CreateMealResponse{
		Success: true,
		Meal:    convertMealToProto(meal),
	}, nil
}

func (h *MealHandler) GetMeal(ctx context.Context, req *pb.GetMealRequest) (*pb.GetMealResponse, error) {
	mealID, err := uuid.Parse(req.GetMealId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid meal ID")
	}

	meal, err := h.mealService.GetMeal(ctx, mealID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "meal not found")
	}

	return &pb.GetMealResponse{Meal: convertMealToProto(meal)}, nil
}

func (h *MealHandler) UpdateMeal(ctx context.Context, req *pb.UpdateMealRequest) (*pb.UpdateMealResponse, error) {
	mealID, err := uuid.Parse(req.GetMealId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid meal ID")
	}

	updateReq := &model.UpdateMealRequest{
		ID:          mealID,
		Name:        strPtr(req.GetName()),
		Description: strPtr(req.GetDescription()),
		Recipe:      strPtr(req.GetRecipe()),
		ImageURL:    strPtr(req.GetImageUrl()),
		Calories:    floatPtr(req.GetCalories()),
		Proteins:    floatPtr(req.GetProteins()),
		Fats:        floatPtr(req.GetFats()),
		Carbs:       floatPtr(req.GetCarbs()),
		Water:       floatPtr(req.GetWater()),
	}

	for _, c := range req.GetComponents() {
		compID, err := uuid.Parse(c.GetComponentMealId())
		if err != nil {
			continue
		}
		updateReq.Components = append(updateReq.Components, model.MealComponentInput{
			ComponentMealID: compID,
			Amount:          c.GetAmount(),
		})
	}

	meal, err := h.mealService.UpdateMeal(ctx, updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update meal: %v", err)
	}

	return &pb.UpdateMealResponse{Success: true, Meal: convertMealToProto(meal)}, nil
}

func (h *MealHandler) DeleteMeal(ctx context.Context, req *pb.DeleteMealRequest) (*pb.DeleteMealResponse, error) {
	mealID, err := uuid.Parse(req.GetMealId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid meal ID")
	}

	if err := h.mealService.DeleteMeal(ctx, mealID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete meal: %v", err)
	}

	return &pb.DeleteMealResponse{Success: true}, nil
}

func (h *MealHandler) ListMeals(ctx context.Context, req *pb.ListMealsRequest) (*pb.ListMealsResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	meals, total, err := h.mealService.ListMeals(ctx, &model.ListMealsRequest{
		UserID:    userID,
		SortBy:    req.GetSortBy(),
		SortOrder: req.GetSortOrder(),
		Page:      int(req.GetPage()),
		PageSize:  int(req.GetPageSize()),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list meals: %v", err)
	}

	var pbMeals []*pb.MealData
	for _, m := range meals {
		pbMeals = append(pbMeals, convertMealToProto(&m))
	}

	return &pb.ListMealsResponse{Meals: pbMeals, Total: int32(total)}, nil
}

func (h *MealHandler) SearchMeals(ctx context.Context, req *pb.SearchMealsRequest) (*pb.SearchMealsResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	meals, total, err := h.mealService.SearchMeals(ctx, &model.SearchMealsRequest{
		UserID:    userID,
		Query:     req.GetQuery(),
		SortBy:    req.GetSortBy(),
		SortOrder: req.GetSortOrder(),
		Page:      int(req.GetPage()),
		PageSize:  int(req.GetPageSize()),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search meals: %v", err)
	}

	var pbMeals []*pb.MealData
	for _, m := range meals {
		pbMeals = append(pbMeals, convertMealToProto(&m))
	}

	return &pb.SearchMealsResponse{Meals: pbMeals, Total: int32(total)}, nil
}

func convertMealToProto(m *model.Meal) *pb.MealData {
	meal := &pb.MealData{
		MealId:    m.ID.String(),
		UserId:    m.UserID.String(),
		Name:      m.Name,
		Calories:  m.Calories,
		Proteins:  m.Proteins,
		Fats:      m.Fats,
		Carbs:     m.Carbs,
		Water:     m.Water,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
	if m.Description != nil {
		meal.Description = *m.Description
	}
	if m.Recipe != nil {
		meal.Recipe = *m.Recipe
	}
	if m.ImageURL != nil {
		meal.ImageUrl = *m.ImageURL
	}

	for _, c := range m.Components {
		meal.Components = append(meal.Components, &pb.MealComponentData{
			ComponentMealId: c.ComponentMealID.String(),
			ComponentName:   c.ComponentName,
			Amount:          c.Amount,
			Calories:        c.Calories,
			Proteins:        c.Proteins,
			Fats:            c.Fats,
			Carbs:           c.Carbs,
			Water:           c.Water,
		})
	}

	return meal
}
