package fsutil

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/urionz/goutil/strutil"
)

type ZipRawFile struct {
	FileName string `json:"file_name"`
	Raw      []byte `json:"raw"`
}

func ZipAppendRawsToZip(src string, raws []ZipRawFile, includeSrcPathArg ...bool) error {
	var err error

	tmpDir := strings.Replace(os.TempDir(), "\\", "/", -1) + "/" + strutil.Md5(src)

	if err = os.RemoveAll(tmpDir); err != nil {
		return err
	}

	if err = ZipDeCompress(src, tmpDir); err != nil {
		return err
	}

	for _, raw := range raws {
		fileName := strings.TrimPrefix(raw.FileName, ".")
		fileName = strings.TrimPrefix(raw.FileName, "/")
		fileName = strings.TrimPrefix(fileName, "\\")
		if err = os.MkdirAll(filepath.Dir(tmpDir+"/"+fileName), os.ModePerm); err != nil {
			return err
		}
		var appendFile *os.File
		if err = os.RemoveAll(tmpDir + "/" + fileName); err != nil {
			return err
		}
		appendFile, err = os.Create(tmpDir + "/" + fileName)
		if err != nil {
			return err
		}
		defer appendFile.Close()
		if _, err = appendFile.Write(raw.Raw); err != nil {
			return err
		}
	}

	return ZipCompress(tmpDir, src, includeSrcPathArg...)
}

func ZipCompress(src, dst string, includeSrcPathArg ...bool) error {
	var includeSrcPath bool
	if len(includeSrcPathArg) > 0 {
		includeSrcPath = includeSrcPathArg[0]
	}
	zFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer zFile.Close()

	archive := zip.NewWriter(zFile)
	defer archive.Close()

	return filepath.Walk(src, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if p == src {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		normalizeSrc := strings.Replace(src, "\\", "/", -1)
		normalizeSrc = strings.Replace(normalizeSrc, ".", "", -1)

		p = strings.Replace(p, "\\", "/", -1)

		if !includeSrcPath {
			header.Name = strings.TrimPrefix(p, strings.TrimPrefix(normalizeSrc, "/")+"/")
		} else {
			header.Name = strings.TrimPrefix(p, filepath.Dir(strings.TrimPrefix(normalizeSrc, "/"))+"/")
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(p)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})
}

func ZipDeCompress(src, dst string) error {
	zReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zReader.Close()

	for _, f := range zReader.File {
		fPath := filepath.Join(dst, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fPath, os.ModePerm); err != nil {
				return err
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
