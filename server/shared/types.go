package shared

type Player struct {
	ID int     `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}

type GameState struct {
	Players map[int]*Player `json:"players"`
	NextID  int             `json:"next_id"`
}

type UpdatePositionRequest struct {
	PlayerID    int     `json:"player_id"`
	SequenceNum int     `json:"sequence_num"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
}

type UpdatePositionResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	SequenceNum int    `json:"sequence_num"`
}

type GetPositionsRequest struct {
	PlayerID    int `json:"player_id"`
	SequenceNum int `json:"sequence_num"`
}

type GetPositionsResponse struct {
	Success     bool      `json:"success"`
	Message     string    `json:"message"`
	SequenceNum int       `json:"sequence_num"`
	Players     []*Player `json:"players"`
}
