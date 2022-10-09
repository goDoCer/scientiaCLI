package scientia

type Course struct {
	Title        string `json:"title"`
	Code         string `json:"code"`
	CanManage    bool   `json:"can_manage"`
	HasMaterials bool   `json:"has_materials"`
}

func (c Course) FullName() string {
	return c.Code + "-" + c.Title
}

type Resource struct {
	ID           int           `json:"id"`
	Downloads    int           `json:"downloads"`
	Index        int           `json:"index"`
	Tags         []interface{} `json:"tags"`
	Path         string        `json:"path"`
	VisibleAfter string        `json:"visible_after"`
	Course       string        `json:"course"`
	Year         string        `json:"year"`
	Category     string        `json:"category"`
	Title        string        `json:"title"`
	Type         string        `json:"type"`
}
