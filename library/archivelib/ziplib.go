package archivelib

import (
	"archive/zip"
	"chatgpt-web/library/textcoding"
	"chatgpt-web/library/util"
	"compress/flate"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func Unzip(zipFile string, destDir string, cb func(p float64)) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	if cb != nil {
		cb(0)
	}
	n := len(zipReader.File)
	for i, f := range zipReader.File {
		if cb != nil {
			cb(100 * float64(i) / float64(n))
		}
		fname := string(textcoding.GetUTF8([]byte(f.Name)))
		// logrus.Info("filename=", f.Name)
		// logrus.Info("transformed filename=", fname)

		fpath := filepath.Join(destDir, fname)
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, 0755)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				_ = inFile.Close()
				return err
			}

			logrus.Info("unzip file ", fname, " to ", fpath)
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				_ = inFile.Close()
				_ = outFile.Close()
				return err
			}
			_ = inFile.Close()
			_ = outFile.Close()
		}
	}
	if cb != nil {
		cb(100)
	}

	return nil
}

func ZipArchive(baseDir string, outZipFile string, excludeFileExt []string) error {
	// Get a Buffer to Write To
	zipFile, err := os.Create(outZipFile)
	if err != nil {
		fmt.Println(err)
	}
	defer zipFile.Close()

	// Create a new zip archive.
	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	// Register a custom Deflate compressor.
	zw.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	// Add some files to the archive.
	return addArchiveFiles(zw, baseDir, "", excludeFileExt)
}

func addArchiveFiles(zw *zip.Writer, baseDir string, baseInZip string, excludeFileExt []string) error {
	// Open the Directory
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			ext := path.Ext(file.Name())
			if ok, _ := util.InArray(ext, excludeFileExt); ok {
				continue
			}
			dat, err := ioutil.ReadFile(baseDir + file.Name())
			if err != nil {
				fmt.Println(err)
				return err
			}
			fileHeader, err := zip.FileInfoHeader(file)
			if err != nil {
				fmt.Println(err)
				return err
			}
			fileHeader.Name = baseInZip + file.Name()
			f, err := zw.CreateHeader(fileHeader)
			// Add some files to the archive.
			//f, err := zw.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
				return err
			}
		} else if file.IsDir() {
			// Recurse
			newBase := baseDir + file.Name() + "/"
			//fmt.Println("Recursing and Adding SubDir: " + file.Name())
			//fmt.Println("Recursing and Adding SubDir: " + newBase)
			addArchiveFiles(zw, newBase, baseInZip+file.Name()+"/", excludeFileExt)
		}
	}

	return nil
}
