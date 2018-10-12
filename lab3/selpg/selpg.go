package main

import(
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	flag "github.com/spf13/pflag"
)

/*
 * Global Var and flag Args 
 */
var fileName string = "fileName"
var startPage int = -1
var endPage int = -1
var pageLines int = -1
var flagPage bool = false
var printDst string = "default"

/*
 * func programParse() 
 * flag.Parse()
 */
func pflag_Parse() {
	flag.IntVarP(&startPage,"start", "s", -1, "首页")
	flag.IntVarP(&endPage,"end","e", -1, "尾页")
	flag.IntVarP(&pageLines,"linenum", "l", 72, "打印的每页行数")
	flag.BoolVarP(&flagPage,"printdes","f", false, "是否用换页符换页")
	flag.StringVarP(&printDst, "othertype","d", "default", "打印目的地")
	flag.Parse()
}

/*
 * func my_usage()
 * init my personal usage in pflag
 */
func my_usage() {
	fmt.Println("Usage:\tselpg -s [Number] -e [Number] [options] [filename]\n")
	fmt.Println("\t-s=[Number]\t开始页数(1<=开始<=结束)")
	fmt.Println("\t-e=[Number]\t结束页数(1<=开始<=结束)")
	fmt.Println("\t-l=[Number]\t每页行数(可选)，默认72")
	fmt.Println("\t-f=[true,false]\t是否用换页符来换页(可选)")
	fmt.Println("\t[filename]\t从文件读，省略为标准输入\n")
}

/* 
 * func args_Handler()
 * args_error handing
 */
func args_Handler() {
	if startPage == -1 || endPage == -1 {
		fmt.Println("Error: No Enough Arguments!\n-h for help!")
		os.Exit(1)
	}
	if startPage < 1 || startPage > (int(^uint(0) >> 1)-1) {
		fmt.Println("Error: Start Page Invalid!\n-h for help!")
		os.Exit(1)
	}
	if endPage < 1 || endPage > (int(^uint(0) >> 1)-1) {
		fmt.Println("Error: End Page Invalid!\n-h for help!")
		os.Exit(1)
	}
	if endPage < startPage {
		fmt.Println("It must obey that End Page >= Start Page!\n-h for help!")
		os.Exit(1)
	}
	if pageLines != 72 {
		if pageLines < 1 {
			fmt.Println("pageLines Invalid\n-h for help!")
			os.Exit(1)
		}
	}
}

func readAndWrite() (string) {
	result := ""
	pageCount := 1
	lineCount := 0
	reader := bufio.NewReader(os.Stdin)

	// set the input source
	if flag.NArg() > 0 {
		// accept one file each time
		fileName = flag.Arg(0)
		file, err := os.Open(fileName)
		if err != nil {
			return "ERROR1"
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	}

	// process the input
	if flagPage {
		pageCount = 1
		for {
			str, err := reader.ReadString('\f')
			pageCount++
			if err == io.EOF {
				return "ERROR2"
			}		
			if pageCount >= startPage && pageCount <= endPage {
				result = strings.Join([]string{result, str}, "")
			}
		}
	} else {
		pageCount = 1
		lineCount = 0

		for {
			str, err := reader.ReadString('\n')
			lineCount++
			if err == io.EOF {
				return "ERROR3"
			}		
			if lineCount > pageLines {
				pageCount++
				lineCount = 1
			}
			if pageCount >= startPage && pageCount <= endPage {
				result = strings.Join([]string{result, str}, "")
			}
		}
	}

	// handle invalid input option
	/*
	if pageCount < startPage {
		msg := fmt.Sprintf("start page: (%d) greater than total pages: (%d)",
			*startPage, pageCount)
		return "", errors.New(msg)
	} else if pageCount < *endPage {
		msg := fmt.Sprintf("end page: (%d) greater than total pages: (%d)",
			*endPage, pageCount)
		return "", errors.New(msg)
	}
	*/
	// set the output source
	if printDst != "default" {
		cmd := exec.Command("lp", "-d"+printDst)
		cmd.Stdin = strings.NewReader(result)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		//err = cmd.Run()
		//if err != nil {
		//	return "", errors.New(fmt.Sprint(err) + " : " + stderr.String())
		//}
	}


	//return result, nil
	return result
}



/*
 * func main()
 */

func main() {
	flag.Usage = my_usage
	if len(os.Args) == 2 && os.Args[1] == "-h" {
		flag.Usage()
		os.Exit(1)
	}
	fmt.Println(fileName, startPage, endPage, pageLines, flagPage, printDst)
	pflag_Parse()
	fileName = flag.Arg(0)
	fmt.Println(fileName, startPage, endPage, pageLines, flagPage, printDst)
	args_Handler()
	message := readAndWrite()
	fmt.Println(message)

	fmt.Println("Print Completed!")

	
}

