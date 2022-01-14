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
	//quto_count:=0
	for r.end<r.size{
		switch r.buf[r.end] {
		case '"'://可能有双引号包含
			if r.cellStart==0{
				if r.buf[r.end+1]=='"' && r.end+2<r.size && r.buf[r.end+2]=='"'{//首尾有双引号包含,且该片段也有双引号
					r.cellStart=2
				}else if r.buf[r.end+1]=='"'{//首尾无双引号包含
					r.cellStart=1
				}else{//首尾有双引号包含
					r.cellStart=2
				}
				r.start = r.end +1
				r.cellEnd =0 //列结束标识清空
			}else{
				//if r.buf[r.end-1]=='"' { //内容里面可能也有双引号
				//	if quto_count==0{ //说明此内容开头的N个字带有双引号,可能没有被双引号包含
				//		quto_count=2
				//	}else{
				//		quto_count++
				//	}
				//}else{
				//	quto_count++
				//}
			}
		case ','://可能是列分隔符
			if r.cellStart==2 && r.buf[r.end-1]=='"' { //是双引号包含 则(",)是列结束标识
				//分两种情况
				//1 结尾有双引号 2结尾无双引号
				if (r.buf[r.end-2]=='"' && r.buf[r.end-3]=='"') || (r.buf[r.end-2]!='"' && r.buf[r.end-3]!='"'){
					line=append(line,strings.ReplaceAll(string(r.buf[r.start:r.end-1]),"\"\"","\""))
					r.cellStart=0 //列数据开始标识清空
					r.cellEnd=1 //列数据已经获取完毕标识
					//quto_count=0
					is_has = 1
				}

			}else if r.cellStart==1{//前后无双引号包含
				//结尾内容有双引号或没有双引号
				if r.buf[r.end-1]!='"' || (r.buf[r.end-1]=='"' && r.buf[r.end-2]=='"') {

				}
				line=append(line,strings.ReplaceAll(string(r.buf[r.start:r.end-is_r]),"\"\"","\""))
				r.cellStart=0 //列数据开始标识清空
				r.cellEnd=1 //列数据已经获取完毕标识
				//quto_count=0
				is_has = 1
			}else if r.cellStart==0{
				line=append(line,"")
				//quto_count=0
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
				line=append(line,strings.ReplaceAll(string(r.buf[r.start:r.end-is_r]),"\"\"","\"") )
				r.cellStart=0 //列数据开始标识清空
				if is_has==1{
					line=append(line,"")
				}
				return line,nil
			}else if r.cellStart==2 && r.buf[r.end-1-is_r]=='"'{//被双引号包含了,可能已经过了第二个引号
				//if r.cellEnd!=0 || (r.buf[r.end-1-is_r]=='"' && quto_count<=1) {//单元格结束 或 前1个字符是引号（说明已经结束）
				//	line=append(line,strings.ReplaceAll(string(r.buf[r.start:r.end-is_r-1]),"\"\"","\"") )//-1是减掉结尾的引号
				//	r.cellStart=0 //列数据开始标识清空
				//	return line,nil
				//}
				line=append(line,strings.ReplaceAll(string(r.buf[r.start:r.end-is_r-1]),"\"\"","\"") )//-1是减掉结尾的引号
				r.cellStart=0 //列数据开始标识清空
				return line,nil
			}
		default:
			//如果是双引号 同时是开始行标识
			if r.cellStart==0{
				r.cellStart=1//无双引号包含
				r.cellEnd =0 //列结束标识清空
				r.start = r.end
				is_has = 0
			}else{

			}
		}
		r.end++
	}

	if len(line)>0{

		if r.end-r.start>1{
			line=append(line,string(r.buf[r.start:r.end]))
		}else{
			line=append(line,"")
		}

		return line ,nil
	}
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