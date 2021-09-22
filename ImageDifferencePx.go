package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
	"sort"
	"strconv"
	"time"
)
type outline struct {
	x int
	y int
}
type Outlinesort []outline
type Outlinesortx []outline

func (o Outlinesort) Len() int {
	//返回传入数据的总数
	return len(o)
}
func (o Outlinesort) Swap(i, j int) {
	//两个对象满足Less()则位置对换
	//表示执行交换数组中下标为i的数据和下标为j的数据
	o[i], o[j] = o[j], o[i]
}
func (o Outlinesort) Less(i, j int) bool {
	//按字段比较大小,此处是降序排序
	//返回数组中下标为i的数据是否小于下标为j的数据
	return o[i].y < o[j].y
}
func (o Outlinesortx) Len() int {
	//返回传入数据的总数
	return len(o)
}
func (o Outlinesortx) Swap(i, j int) {
	//两个对象满足Less()则位置对换
	//表示执行交换数组中下标为i的数据和下标为j的数据
	o[i], o[j] = o[j], o[i]
}
func (o Outlinesortx) Less(i, j int) bool {
	//按字段比较大小,此处是降序排序
	//返回数组中下标为i的数据是否小于下标为j的数据
	return o[i].x < o[j].x
}
func main() {
	t2 := time.Now().Nanosecond()
	file1, _ := os.Open(".\\img\\a0.jpg")
	file2, _ := os.Open(".\\img\\a8.jpg")
	defer file1.Close()
	defer file2.Close()

	img1, _ := jpeg.Decode(file1)
	img2, _ := jpeg.Decode(file2)

	width, high := img1.Bounds().Dx(),img1.Bounds().Dy()
	var status,same, z,h,w int = 0,1,0,1,1   //status same划线状态，gap允许色差 z多少个差别  h单个色块高 w单个色块宽
	var gap float64 = 1
	var outlines []outline = make([]outline,0,(width*high)/32)

	b := img1.Bounds()
	//根据b画布的大小新建一个新图像
	m := image.NewRGBA(b)
	draw.Draw(m, b, img1, b.Min, draw.Over)

	for i:= 0;i < width ; i+=w {
		for j:=0 ; j < high ; j+=h  {
			subimg1px := rgb2gray1px(img1.At(i,j))
			subimg2px := rgb2gray1px(img2.At(i,j))
			//AverageHash  DifferenceHash  PerceptionHash 三种常用算法,适合比对相似图片。我们是求差别
			distance := math.Abs(subimg1px - subimg2px)

			if distance > gap {
				z++
				outlines = append(outlines, outline{
					x:i,
					y:j,
				})

				if status == 0 && same == 1 {
					//outlines = append(outlines, outline{
					//	x:i,
					//	y:j,
					//})
					drawline(i,j,4,2,w,m)	//竖向画框
					status = 1
				}
			}
			if  status == 1 &&  distance <= gap {
				outlines = append(outlines, outline{
					x:i,
					y:j,
				})
				drawline(i,j,4,3,w,m)	//竖向画框
				status,same = 0, 1
			}
		} //w
	}//h

	name1 := strconv.Itoa(int(time.Now().Unix()))
	imgw, err := os.Create(name1 + "shuidaan.jpg")
	if err != nil {
		fmt.Println(err)
	}

	sort.Sort(Outlinesort(outlines))
	sortline(outlines)
	for k,v := range outlines{
		if k == 0 {
			status, same= 0,0
		}
		if k+1 == len(outlines) {
			drawline(outlines[k].x,outlines[k].y,4,1,w,m)
		}
		if status == 0 && same == 0 {
			drawline(v.x,v.y,4,0,w,m)
			same, status = v.x,1
			continue
		}
		if v.x - same == w {
			same, status= v.x,1
		}
		if (v.x - same > w || v.y != outlines[k-1].y ) && status == 1{
			drawline(outlines[k-1].x,outlines[k-1].y,4,1,w,m)
			same,status = 0, 0
		}
	}
	for _,v := range outlines{
		m.Set(v.x,v.y,img2.At(v.x,v.y))
	}

	jpeg.Encode(imgw, m, &jpeg.Options{100})
	defer imgw.Close()
	t3:= time.Now().Nanosecond() -t2
	fmt.Printf("This picture width is %d,height is %d pixels. The program runs for %d milliseconds. There are %d pixels that are different \n",width,high,t3/1e6,len(outlines))
	fmt.Printf("图片宽 %d,高 %d 像素. 程序运行耗时 %d 毫秒. 相片有 %d 像素不同 \n",width,high,t3/1e6,len(outlines))
}

//由点划线
func drawline(x, y, size, dire, zone int, m *image.RGBA) error {
	//x,y划线起点坐标  size线粗  dire线方向 zone对比像素大小
	size+=4
	switch dire {
	case 0:
		for dot := 4;dot < size;dot++ {
			for z:= 0;z < zone ;z++  {
				m.Set(x-dot,y+z,color.RGBA{255, 0, 0, 255})
			}
		}
	case 1:
		for dot := 4;dot < size;dot++ {
			for z:= 0;z < zone ;z++  {
				m.Set(x+dot+zone,y+z,color.RGBA{0, 255, 0, 255})
			}
		}
	case 2:
		for dot := 4;dot < size;dot++ {
			for z:= 0;z < zone ;z++  {
				m.Set(x+z,y-dot,color.RGBA{255, 0, 0, 255})
			}
		}
	case 3:
		for dot := 4;dot < size;dot++ {
			for z:= 0;z < zone ;z++  {
				m.Set(x+z,y+dot,color.RGBA{0, 255, 0, 255})
			}
		}
	default:
		return errors.New("Parameter error")
	}
	return nil
}

// 排序，用于框出差异，优化减少重复设置像素  切片指针传递
func sortline(outlines Outlinesort) {
	oy,startkey := -1,0
	if len(outlines) > 0 {
		oy = outlines[0].y
	}
	var sortx  Outlinesort
	for key,value := range outlines {
		if value.y != oy {
			sortx = outlines[startkey:key]
			sort.Sort(Outlinesortx(sortx))
			startkey,oy = key,value.y
		}
		if key == outlines.Len() {
			sortx = outlines[startkey:key]
			sort.Sort(Outlinesortx(sortx))
		}
	}
}
//将rgb像素转化为gray，用于对比色差
func rgb2gray1px(colorImg color.Color) float64 {
	r, g, b, _ := colorImg.RGBA()
	lum := 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/256)
	return lum
}
