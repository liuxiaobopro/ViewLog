package res

type LogIndexShowFoldsRes struct {
	Title    string                  `json:"title,omitempty"`
	Id       string                  `json:"id,omitempty"`
	Path     string                  `json:"path,omitempty"`
	PathId   string                  `json:"pathId,omitempty"`
	Checked  bool                    `json:"checked,omitempty"`
	Spread   bool                    `json:"spread,omitempty"`
	Name     string                  `json:"name,omitempty"`
	Children []*LogIndexShowFoldsRes `json:"children,omitempty"`
}
