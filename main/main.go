package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"

	"github.com/bo-er/screenshot"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

// save *image.RGBA to filePath with PNG format.
func save(img *image.RGBA, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

func main() {
	md := flag.String("md", "", "markdown文件的目录")
	picfolder := flag.String("pic", "", "markdown文件的存储目录")
	flag.Parse()
	fmt.Printf("markdown文件地址是:%s,图片文件夹地址是:%s", *md, *picfolder)
	for 1 == 1 {
		fmt.Println("--- Please press esc ---")
		robotgo.EventHook(hook.KeyDown, []string{"esc"}, func(e hook.Event) {
			fmt.Println("esc")
			robotgo.EventEnd()
		})

		s := robotgo.EventStart()
		<-robotgo.EventProcess(s)

		ok := robotgo.AddEvents("esc")
		if ok {
			fmt.Println("按下了esc!开始截图!")
			// Capture each displays.
			n := screenshot.NumActiveDisplays()
			if n <= 0 {
				panic("Active display not found")
			}

			var all image.Rectangle = image.Rect(0, 0, 0, 0)

			for i := 0; i < n; i++ {
				bounds := screenshot.GetDisplayBounds(i)
				all = bounds.Union(all)

				img, err := screenshot.CaptureRect(bounds)
				if err != nil {
					panic(err)
				}
				fileName := fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
				save(img, fileName)

				fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
			}

			// Capture all desktop region into an image.
			fmt.Printf("%v\n", all)
			img, err := screenshot.Capture(all.Min.X, all.Min.Y, all.Dx(), all.Dy())
			if err != nil {
				panic(err)
			}
			time := rand.Intn(10000000)
			photoPath := fmt.Sprintf("%s/%d.png", *picfolder, time)
			save(img, photoPath)
			insertIntoMarkdown(*md, photoPath)
		}
	}

}

func insertIntoMarkdown(filePath string, photoPath string) {
	f, err := os.OpenFile(filePath, os.O_WRONLY, 0644)
	if err != nil {
		// 打开文件失败处理

	}
	content := fmt.Sprintf("![新建图片](%s)", photoPath)

	// 查找文件末尾的偏移量
	n, _ := f.Seek(0, 2)

	// 从末尾的偏移量开始写入内容
	_, err = f.WriteAt([]byte(content), n)

	defer f.Close()
}
