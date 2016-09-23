package barfile

import (
	"os"
	"gopkg.in/cheggaaa/pb.v1"
	"time"
)

func CreateBar(file *os.File) (bar *pb.ProgressBar, err error) {
	var fi os.FileInfo
	if fi, err = file.Stat(); err != nil {
		return
	}
	bar = pb.New64(fi.Size()).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	return
}