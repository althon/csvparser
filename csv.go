package csv

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Reader struct {
	buf []byte
	size int
	cellStart int
	cellEnd int
	start int
	end int
}

func (r *Reader) Read() ([]string,error){
	line:=[]string{}
	is_r :=0
	is_has:=0
	for r.end<r.size{
		switch r.buf[r.end] {
		case '"':
			if r.cellStart==0{
				r.cellStart=2 //代表有双引号包含
				r.cellEnd =0 //列结束标识清空
				r.start = r.end+1
			}
		case ','://可能是列分隔符
			if r.cellStart==2 && r.buf[r.end-1]=='"'{ //是双引号包含 则【",】是列结束标识
				line=append(line,strings.ReplaceAll(string(r.buf[r.start:r.end-1]),"\"\"","\""))
				r.cellStart=0 //列数据开始标识清空
				r.cellEnd=1 //列数据已经获取完毕标识
				is_has = 1
			}else if r.cellStart==1{
				line=append(line,string(r.buf[r.start:r.end-is_r]))
				r.cellStart=0 //列数据开始标识清空
				r.cellEnd=1 //列数据已经获取完毕标识
				is_has = 1
			}else if r.cellStart==0{
				line=append(line,"")
				is_has = 1
			}
		case '\r':
			is_r=1
		case '\n'://可能是一行数据结束标识
			if r.cellEnd==1{//列数据已经获取完毕，遇上了换行符，说明此行已结束
				r.cellEnd=0
				if is_has==1{
					line=append(line,"")
				}
				return line,nil
			}else if r.cellStart==1{ //如果这个列没有被双引号包含，则肯定此行已经结束
				line=append(line,string(r.buf[r.start:r.end-is_r]))
				r.cellStart=0 //列数据开始标识清空
				if is_has==1{
					line=append(line,"")
				}
				return line,nil
			}
		default:
			//如果是双引号 同时是开始行标识
			if r.cellStart==0{
				r.cellStart=1//无双引号包含
				r.cellEnd =0 //列结束标识清空
				r.start = r.end
				is_has = 0
			}
		}
		r.end++
	}

	if len(line)>0{
		line=append(line,string(r.buf[r.start:r.end]))
		return line ,nil
	}
	//if r.cellEnd>0{
	//	r.cellEnd=0
	//	if is_has==1{
	//		line=append(line,"")
	//	}
	//	return line,nil
	//}

	return nil,io.EOF
}

func NewReader(filename string)*Reader{
	buf,err:=os.ReadFile(filename)
	if err!=nil{
		fmt.Println("cannot open csv file:",filename)
		return nil
	}
	return &Reader{buf:buf,size: len(buf)}
}

//func parseCsv(file string) []string{
//
//	buf,err:=os.ReadFile(file)
//	if err!=nil{
//		fmt.Println("cannot open file:",file)
//		return nil
//	}
//
//	line:=[]string{}
//	cell_start:=0
//	cell_end:=0
//	is_r :=0
//	size:= len(buf)
//	start,end:=0,0
//	for end<size{
//		switch buf[end] {
//		case '"':
//			if cell_start==0{
//				cell_start=2 //代表有双引号包含
//				cell_end =0 //列结束标识清空
//				start = end+1
//			}
//		case ','://可能是列分隔符
//			if cell_start==2 && buf[end-1]=='"'{ //是双引号包含 则【",】是列结束标识
//				line=append(line,strings.ReplaceAll(l,"\"\"","\""))
//				cell_start=0 //列数据开始标识清空
//				cell_end=1 //列数据已经获取完毕标识
//			}else if cell_start==1{
//				line=append(line,string(buf[start:end]))
//				cell_start=0 //列数据开始标识清空
//				cell_end=1 //列数据已经获取完毕标识
//			}
//		case '\r':
//			is_r=1
//		case '\n'://可能是一行数据结束标识
//			if cell_end==1{//列数据已经获取完毕，遇上了换行符，说明此行已结束
//				cell_end=0
//			}else if cell_start==1{ //如果这个列没有被双引号包含，则肯定此行已经结束
//				line=append(line,string(buf[start:end-is_r]))
//				cell_start=0 //列数据开始标识清空
//			}
//		default:
//			//如果是双引号 同时是开始行标识
//			if cell_start==0{
//				cell_start=1//无双引号包含
//				cell_end =0 //列结束标识清空
//				start = end
//			}
//		}
//		end++
//	}
//
//	return line
//}
