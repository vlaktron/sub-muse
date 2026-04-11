package ui

// TabType represents the type of tab in the navigation
type TabType int

const (
	TabSongs TabType = iota
	TabArtists
	TabAlbums
	TabAlbumArtists
)

// Navigation holds the navigation state
type Navigation struct {
	currentTab  TabType
	searchQuery string
	focused     bool
}

// NewNavigation creates a new navigation instance
func NewNavigation() Navigation {
	return Navigation{
		currentTab:  TabSongs,
		searchQuery: "",
		focused:     false,
	}
}

// Init initializes the navigation
func (n *Navigation) Init() {
}

// Update handles navigation updates
func (n *Navigation) Update(msg interface{}) {
}

// View renders the navigation UI
func (n *Navigation) View() string {
	return ""
}

// SetTab sets the current tab
func (n *Navigation) SetTab(tab TabType) {
	n.currentTab = tab
}

// GetTab returns the current tab
func (n *Navigation) GetTab() TabType {
	return n.currentTab
}

// UpdateSearch updates the search query
func (n *Navigation) UpdateSearch(query string) {
	n.searchQuery = query
}

// GetSearchQuery returns the current search query
func (n *Navigation) GetSearchQuery() string {
	return n.searchQuery
}

// ToggleFocus toggles the focus state
func (n *Navigation) ToggleFocus() {
	n.focused = !n.focused
}
