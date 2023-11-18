package file

import (
	"io"
	"mime/multipart"
	"os"
	"os/exec"
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

	new_name, err := fs.generateRandomName()

	if fileType == "image" && extType != "svg+xml" {
		webp, err := bimg.NewImage(bytes).Convert(bimg.WEBP)

		if err != nil {
			return nil, err
		}

		newFilePath := path.Join(uploadPath, new_name+".webp")

		err = os.WriteFile(newFilePath, webp, 0700)

		if err != nil {
			return nil, err
		}

		absPath := "/static/" + dir + new_name + ".webp"

		return &absPath, nil
	}

	newFilePath := path.Join(uploadPath, new_name+ext)

	err = os.WriteFile(newFilePath, bytes, 0700)

	if err != nil {
		return nil, err
	}

	absPath := "/static/" + dir + "/" + new_name + ext

	return &absPath, nil

}

func (fs *FileService) getFilenameWIthoutExt(name string) string {
	return strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
}

func (fs *FileService) getFileExt(name string) string {
	return strings.TrimPrefix(filepath.Ext(name), filepath.Base(name))
}

func (fs *FileService) generateRandomName() (string, error) {

	newUUID, err := exec.Command("uuidgen").Output()
	return string(newUUID), err
}
