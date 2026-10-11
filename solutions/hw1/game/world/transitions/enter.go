package transitions

type ItemID string

type Enter struct {
	Locked         bool
	RequiredItemID ItemID
}

func NewEnter(locked bool, itemID ItemID) *Enter {
	return &Enter{Locked: locked, RequiredItemID: itemID}
}
