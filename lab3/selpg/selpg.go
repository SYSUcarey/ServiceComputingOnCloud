package main

import(
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	flag "github.com/spf13/pflag"
)

/*
 * Global Var and flag Args 
 */
var err error
var progname string
var fileName string = "default"
var startPage int = -1
var endPage int = -1
var pageLines int = 72
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
		os.Exit(2)
	}
	if endPage < 1 || endPage > (int(^uint(0) >> 1)-1) {
		fmt.Println("Error: End Page Invalid!\n-h for help!")
		os.Exit(3)
	}
	if endPage < startPage {
		fmt.Println("It must obey that End Page >= Start Page!\n-h for help!")
		os.Exit(4)
	}
	if pageLines != 72 {
		if pageLines < 1 {
			fmt.Println("Page's Lines Invalid\n-h for help!")
			os.Exit(5)
		}
	}

	/* there is one more arg */
	if flag.NArg() > 0 {
		fileName = flag.Arg(0)
		/* check if file exists */
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Println("input file \"", fileName, "\" does not exist\n-h for help!")
			os.Exit(6)
		}
		/* check if file is readable */
		file, err = os.OpenFile(fileName, os.O_RDONLY, 0666)
		if err != nil {
			if os.IsPermission(err) {
				fmt.Println("input file \"", fileName,"\" exist but cannot be read\n-h for help!")
				os.Exit(7)
			}
		}
		file.Close()
	}
}

func readAndWrite() {
	fin := os.Stdin
	fout := os.Stdout
	var (
		 page_ctr int
		 line_ctr int
		 err error
		 err1 error
		 err2 error
		 line string
		 cmd *exec.Cmd
		 stdin io.WriteCloser
	)
	/* set the input source */
	if fileName != "" {
		fin, err1 = os.Open(fileName)
		if err1 != nil {
			fmt.Println("could not open file" , fileName, " \n-h for help!")
			os.Exit(11)
		}
	}

	if printDst != "" {
		cmd = exec.Command("cat", "-n")
		stdin, err = cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		stdin = nil
	}

/* begin one of two main loops based on page type */
	rd := bufio.NewReader(fin)
	if flagPage == false {
		line_ctr = 0
		page_ctr = 1
		for true {
			line, err2 = rd.ReadString('\n')
			if err2 != nil { /* error or EOF */
				break
			}
			line_ctr++
			if line_ctr > pageLines {
				page_ctr++
				line_ctr = 1
			}
			if page_ctr >= startPage && page_ctr <= endPage {
				fmt.Fprintf(fout, "%s", line)
			}
		}
	} else {
		page_ctr = 1
		for true {
			c, err3 := rd.ReadByte()
			if err3 != nil { /* error or EOF */
				break
			}
			if c == '\f' {
				page_ctr++
			}
			if page_ctr >= startPage && page_ctr <= endPage {
				fmt.Fprintf(fout, "%c", c)
			}
		}
		fmt.Print("\n")
	}

	/* end main loop */
	if page_ctr < startPage {
		fmt.Println(progname, ": start_page (", startPage, ") greater than total pages (", page_ctr, "), no output written\n-h for help")
	} else if page_ctr < endPage {
		fmt.Println(progname, ": end_page (", endPage, ") greater than total pages (", page_ctr, "), less output than expected\n-h for help")
	}
	
	if printDst != "" {
		stdin.Close()
		cmd.Stdout = fout
		cmd.Run()
	}
	fmt.Println("\n---------------\nProcess end\n")
	fin.Close()
	fout.Close()
}



/*
 * func main()
 */

func main() {
	flag.Usage = my_usage
	progname = os.Args[0]
	if len(os.Args) == 2 && os.Args[1] == "-h" {
		flag.Usage()
		os.Exit(1)
	}
	pflag_Parse()
	fileName = flag.Arg(0)
	args_Handler()
	readAndWrite()
	fmt.Println("Print Completed!")
}

