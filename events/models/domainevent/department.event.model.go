package domainevent

type DepartmentInsertEventModel struct {
	ID       string `json:"id"`
	ParentID string `json:"parentId,omitempty"`
	Type     string `json:"type"`
	IsActive bool   `json:"isActive"`
	Version  int    `json:"version"`
}

type DepartmentUpdateEventModel struct {
	ID       string `json:"id"`
	ParentID string `json:"parentId,omitempty"`
	Type     string `json:"type"`
	IsActive bool   `json:"isActive"`
	Version  int    `json:"version"`
}

type DepartmentDeleteEventModel struct {
	ID string `json:"id"`
}
