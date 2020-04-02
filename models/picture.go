package models

// PictureFormat specifies the format for a given picture.
type PictureFormat string

const (
	// PictureFormatJPEG specifies pictures with jpeg format.
	PictureFormatJPEG PictureFormat = "image/jpeg"

	// PictureFormatPNG specifies pictures with png format.
	PictureFormatPNG = "image/png"

	// PictureFormatGIF specifies pictures with gif format.
	PictureFormatGIF = "image/gif"
)

// Picture represents a picture that holds the format and data as bytes.
type Picture struct {
	ID     string        `json:"id"`
	Format PictureFormat `json:"format"`
	Data   []byte        `db:"blob" json:"data"`
}
