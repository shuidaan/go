package main

import (
	"errors"
	"fmt"
	"github.com/corona10/goimagehash"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
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
	file1, _ := os.Open("E:\\wq\\img\\a0.jpg")
	file2, _ := os.Open("E:\\wq\\img\\a8.jpg")
	defer file1.Close()
	defer file2.Close()

	img1, _ := jpeg.Decode(file1)
	img2, _ := jpeg.Decode(file2)

	width, high := img1.Bounds().Dx(),img1.Bounds().Dy()
	var status,same, gap, z,h,w int = 0,1,1,0,1,1    //status same划线状态，gap允许色差 z多少个差别  h单个色块高 w单个色块宽
	var outlines []outline = make([]outline,0,(width+high)/64)
	b := img1.Bounds()
	//根据b画布的大小新建一个新图像
	m := image.NewRGBA(b)
	fmt.Println(width,high)
	draw.Draw(m, b, img1, b.Min, draw.Over)

	////测试被裁剪的小图是否全部加入对比
	//sb1 := img1.Bounds()
	////根据b画布的大小新建一个新图像
	//sm1 := image.NewRGBA(sb1)
	//
	//sb2 := img1.Bounds()
	////根据b画布的大小新建一个新图像
	//sm2 := image.NewRGBA(sb2)

	for i:= 0;i < width ; i+=w {
		//if (width - i < 16)  i = width- i
		//x = i
		//if (width - i < 16) {
		//	x = width- i
		//}
		for j:=0 ; j < high ; j+=h  {
			//y=j
			//if (high - j < 16) {
			//	y = high- j
			//
			//}
			subimg1,err := clip(img1,i,j,w,h)
			if err != nil {
				fmt.Println(err)
			}
			subimg2,err := clip(img2,i,j,w,h)
			if err != nil {
				fmt.Println(err)
			}
			//soffet1 := image.Pt(i,j)
			//ssr1 := subimg2.Bounds()
			//draw.Draw(sm1,sb1,subimg1,ssr1.Min.Sub(soffet1),draw.Over)
			//soffet2 := image.Pt(i,j)
			//ssr2 := subimg2.Bounds()
			//draw.Draw(sm2,sb2,subimg2,ssr2.Min.Sub(soffet2),draw.Over)


			hash1, _ := goimagehash.DifferenceHash(subimg1)  //AverageHash  DifferenceHash  PerceptionHash 三种常用算法
			hash2, _ := goimagehash.DifferenceHash(subimg2)
			distance, err := hash1.Distance(hash2)
			if err != nil {
				fmt.Println(err)
			}



			if distance > gap {
				//z++
				//fmt.Printf("Distance between images:: %v . x: %d ,y: %d ,z: %d \n", distance,i,j,z)

				offet := image.Pt(i,j)
				sr := subimg2.Bounds()

				outlines = append(outlines, outline{
					x:i,
					y:j,
				})
				draw.Draw(m,b,subimg2,sr.Min.Sub(offet),draw.Over)
				if status == 0 && same == 1 {

					outlines = append(outlines, outline{
						x:i,
						y:j,
					})
					drawline(i,j,4,2,w,m)
					status = 1
				}
				offet,sr = image.ZP ,image.ZR

			}
			if  status == 1 &&  distance <= gap {
				//fmt.Printf("Distance between images: %v . x: %d ,y: %d ,z: %d ______\n", distance,i,j,z)
				outlines = append(outlines, outline{
					x:i,
					y:j,
				})
				drawline(i,j,4,3,w,m)
				status,same = 0, 1
			}
			z++
		} //w
		//fmt.Println(x ,y)
	}//h

	name1 := strconv.Itoa(int(time.Now().Unix()))
	fmt.Println(name1)
	imgw, err := os.Create(name1 + "shuidaan.jpg")
	if err != nil {
		fmt.Println(err)
	}
	//测试被裁剪的小图是否全部加入对比
	//simgw1, err := os.Create(name1 + "new1.jpg")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//simgw2, err := os.Create(name1 + "new2.jpg")
	//if err != nil {
	//	fmt.Println(err)
	//}


	//sort.Sort(Outlinesortx(outlines))
	sort.Sort(Outlinesort(outlines))
	fmt.Println(outlines)
	sortlx := sortline(outlines)
	fmt.Println(outlines)
	fmt.Println(sortlx)
	for k,v := range sortlx{
		if k == 0 {
			status, same= 0,0
		}
		if k+1 == sortlx.Len() {
			drawline(outlines[k].x,outlines[k].y,4,1,w,m)
		}
		if status == 0 && same == 0 {
			drawline(v.x,v.y,4,0,w,m)
			same, status = v.x,1
			//fmt.Println(same,v.y,status)
			//fmt.Println(same,v.x,v.y,oy,"___")
			continue
		}
		if v.x - same == w {
			same, status= v.x,1
		}
		//fmt.Println(same,v.x,sortlx[k].x,v.y,  sortlx[k].y ,"__1_")
		if (v.x - same > w || v.y != sortlx[k-1].y ) && status == 1{
			//fmt.Println(k)
			//fmt.Println(v.x,v.y)
			// 取K值，无法应用
			drawline(outlines[k-1].x,outlines[k-1].y,4,1,w,m)
			same,status = 0, 0
			//fmt.Println(same,v.x,v.y,oy)
		}
	}

	//fmt.Println(outlines)

	jpeg.Encode(imgw, m, &jpeg.Options{100})
	defer imgw.Close()
	//	//	测试被裁剪的小图是否全部加入对比
	//jpeg.Encode(simgw1, sm1, &jpeg.Options{100})
	//defer simgw1.Close()
	//
	//jpeg.Encode(simgw2, sm2, &jpeg.Options{100})
	//defer simgw2.Close()


	fmt.Println(time.Now().Unix() , cap(outlines),z,w)
}

//func clip(in io.Reader, out io.Writer, x0, y0, x1, y1, quality int) error {

func clip(src image.Image, x, y, w, h int) (image.Image, error) {

	var subImg image.Image

	if rgbImg, ok := src.(*image.YCbCr); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.RGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.RGBA) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.NRGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA) //图片裁剪x0 y0 x1 y1
	} else {

		return subImg, errors.New("图片解码失败")
	}

	return subImg, nil
}
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

func sortline(outlines Outlinesort)  (Outlinesort) {
	oy,startkey := -1,0
	var sortx,sortxx  Outlinesort
	for key,value := range outlines {
		//if oy == -1 {
		//	oy = value.y
		//}
		if value.y != oy {
			sortx = outlines[startkey:key]
			//fmt.Println(startkey,key)
			sort.Sort(Outlinesortx(sortx))
			sortxx = append(sortxx,sortx...)
			startkey,oy = key,value.y
		}
		if key == outlines.Len() {
			sortx = outlines[startkey:key]
			//fmt.Println(startkey,key)
			sort.Sort(Outlinesortx(sortx))
			sortxx = append(sortxx,sortx...)
		}
	}
	return sortxx
}
