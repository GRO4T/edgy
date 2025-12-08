package videos

type videoService interface {
	ListVideos() ([]video, error)
	GetVideoByID(ID string) (video, error)
	CreateVideo(path string) (string, error)
	DeleteVideo(ID string) error
	UpdateVideo(ID string, newPath string, isUpsert bool) error
}
