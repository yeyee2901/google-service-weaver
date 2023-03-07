package service

type GetBattleItemRequest struct {
    Name string `form:"name" binding:"required"`
	Tier string `form:"tier"`
}
