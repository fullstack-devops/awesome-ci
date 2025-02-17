package uploadasset

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type UploadAsset struct {
	File  os.File
	Name  string
	Infos fs.FileInfo
	Hash  []byte
}

// GetAsset takes a string as argument which must be in the format `<assetformat>=<filename>`.
//
// The function will return an `UploadAsset` object or an error.
//
// The supported formats are `file`, `zip` and `tgz`.
//
// For `file` assets, the function will open the file in read only mode and
// create an `UploadAsset` object with the file and the file info.
//
// For `zip` and `tgz` assets, the function will use the `createZipFile` and
// `createTgzFile` functions to create the asset.
//
// An error will be returned if the asset format is not supported or if an
// error occurs while creating the asset.
func GetAsset(assetLocation string) (asset *UploadAsset, err error) {
	assetName := strings.Split(assetLocation, "=")

	switch assetName[0] {

	case "file":
		file, err := os.OpenFile(assetName[1], os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
		if err := asset.New(assetName[1], file); err != nil {
			return nil, err
		}
		return asset, nil
	case "zip":
		return createZipFile(assetName[1])
	case "tgz":
		return createTgzFile(assetName[1])
	default:
		return nil, fmt.Errorf("not an valid asset format")
	}
}

func createZipFile(name string) (asset *UploadAsset, err error) {
	tmpFile, err := os.CreateTemp("awesome-ci", "upload-asset-")
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()
	zipWriter := zip.NewWriter(tmpFile)
	err = filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(filepath.Dir(name), path)
		if err != nil {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		f, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	if err := asset.New(fmt.Sprintf("%s.zip", filepath.Base(name)), tmpFile); err != nil {
		return nil, err
	}

	return
}

func createTgzFile(name string) (asset *UploadAsset, err error) {
	tmpFile, err := os.CreateTemp("awesome-ci", "upload-asset-")
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()
	gzipWriter := gzip.NewWriter(tmpFile)
	tarWriter := tar.NewWriter(gzipWriter)
	err = filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(filepath.Dir(name), path)
		if err != nil {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		header, err := tar.FileInfoHeader(info, relPath)
		if err != nil {
			return err
		}
		err = tarWriter.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = tarWriter.Close()
	if err != nil {
		return nil, err
	}
	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	if err := asset.New(fmt.Sprintf("%s.tgz", filepath.Base(name)), tmpFile); err != nil {
		return nil, err
	}

	return
}

func (ua *UploadAsset) New(name string, file *os.File) (err error) {
	// get file info
	info, err := file.Stat()
	if err != nil {
		return err
	}
	// get hash
	h := sha256.New()
	if _, err = io.Copy(h, file); err != nil {
		return
	}

	// set values
	ua.Name = name
	ua.File = *file
	ua.Infos = info
	ua.Hash = h.Sum(nil)
	return nil
}
