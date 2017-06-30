package vkapi

const (
	// Full description at https://vk.com/dev/permissions
	ScopeFriends        = 2
	ScopePhotos         = 4
	ScopeAudio          = 8
	ScopeVideo          = 16
	ScopePages          = 128
	ScopeAddLink        = 256
	ScopeStatus         = 1024
	ScopeNotes          = 2048
	ScopeMessages       = 4096
	ScopeWall           = 8192
	ScopeAds            = 32768
	ScopeOffline        = 65536
	ScopeDocs           = 131072
	ScopeGroupsOrManage = 262144
	ScopeNotifications  = 524288
	ScopeStats          = 1048576
	ScopeEmail          = 4194304
	ScopeMarket         = 134217728

	ScopeAll = ScopeFriends | ScopePhotos | ScopeAudio | ScopeVideo | ScopePages | ScopeAddLink | ScopeStatus | ScopeNotes | ScopeMessages | ScopeWall | ScopeAds | ScopeOffline | ScopeDocs | ScopeGroupsOrManage | ScopeNotifications | ScopeStats | ScopeEmail | ScopeMarket
)
