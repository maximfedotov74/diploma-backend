package file

import (
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/h2non/bimg"
)

type FileService struct {
	uploadDir string
}

func NewFileService(uploadDir string) *FileService {
	return &FileService{uploadDir: uploadDir}
}

func (fs *FileService) Upload(dir string, header *multipart.FileHeader) (*string, error) {
	file, err := header.Open()

	ext := fs.getFileExt(header.Filename)

	contentType := header.Header.Get("Content-type")

	splittedType := strings.Split(contentType, "/")

	fileType := splittedType[0]
	extType := splittedType[1]

	if err != nil {
		return nil, err
	}

	uploadPath := path.Join(fs.uploadDir, dir)

	_, err = os.Stat(uploadPath)

	if os.IsNotExist(err) {
		os.Mkdir(uploadPath, 0700)
	}

	bytes, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	new_name := fs.generateRandomName()

	if fileType == "image" && extType != "svg+xml" {

		options := bimg.Options{
			Quality: 50,
			Type:    bimg.WEBP,
		}

		webp, err := bimg.Resize(bytes, options)

		if err != nil {
			return nil, err
		}

		newFilePath := path.Join(uploadPath, new_name+".webp")

		err = os.WriteFile(newFilePath, webp, 0700)

		if err != nil {
			return nil, err
		}

		absPath := "/" + path.Join("static", dir, new_name) + ".webp"

		return &absPath, nil
	}

	newFilePath := path.Join(uploadPath, new_name+ext)

	err = os.WriteFile(newFilePath, bytes, 0700)

	if err != nil {
		return nil, err
	}

	absPath := "/" + path.Join("static", dir, new_name) + ext

	return &absPath, nil

}

func (fs *FileService) getFilenameWIthoutExt(name string) string {
	return strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
}

func (fs *FileService) getFileExt(name string) string {
	return strings.TrimPrefix(filepath.Ext(name), filepath.Base(name))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func (fs *FileService) generateRandomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func (fs *FileService) generateRandomName() string {

	s := make([]string, 0, 3)

	for i := 0; i < 3; i++ {
		s = append(s, fs.generateRandomString(23))
	}

	return strings.Join(s, "-")
}
