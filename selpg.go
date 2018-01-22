package main

import (
  "flag"
  "fmt"
  "io"
  "os"
  "bufio"
  "os/exec"
)

type selpg_args struct {
  start int
  end int
  page_len int
  isf bool
  output_dest string
  input_file string
}

var pro_name string

func Init(args *selpg_args) {
	flag.IntVar(&args.start, "s", -1, "Start page.")
	flag.IntVar(&args.end, "e", -1, "End page.")
	flag.IntVar(&args.page_len, "l", 72, "Lines per page.")
	flag.BoolVar(&args.isf, "f", false, "Page type")
	flag.StringVar(&args.output_dest, "d", "", "Output destination")
  flag.Usage = Usage
	flag.Parse()
}

func Usage() {
  fmt.Fprintf(os.Stderr,
		`usage: [-s start page(>=1)] [-e end page(>=s)] [-l length of page(default 72)] [-f type of file(default 1)] [-d dest] [filename specify input file]
`)
}

func Process_args(args *selpg_args) {

  if args.start == -1 || args.end == -1 || args.start < 0 || args.end < 0 || args.start > args.end {
    fmt.Fprintf(os.Stderr,"%s:no start or end args,or have bad value",pro_name)
    flag.Usage()
    os.Exit(1)
  }

  if args.page_len != 72 && args.isf {
    fmt.Fprintf(os.Stderr,"%s:-f and -l can not input together",pro_name)
		flag.Usage()
		os.Exit(1)
	}

  if len(flag.Args()) > 1 {
    fmt.Fprintf(os.Stderr,"%s:too much args",pro_name)
    flag.Usage()
    os.Exit(1)
  }

}

func Process_input(args *selpg_args) {

  var stdin io.WriteCloser
	var err error
	var cmd *exec.Cmd

  if args.output_dest != "" {
		cmd = exec.Command("cat", "-n")
		stdin, err = cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		stdin = nil
	}

  if flag.NArg() > 0 {

    args.input_file = flag.Arg(0)
    input,err := os.Open(args.input_file)

    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    bfreader := bufio.NewReader(input)

    if args.isf {
      page_num := 1;
      for page_num <= args.end {

        line_string, err := bfreader.ReadString('\f')

        if err != nil && err != io.EOF {
          fmt.Println(err)
          os.Exit(1)
        }

        if err == io.EOF {
          break
        }

        if page_num >= args.start {
          if args.output_dest != "" {
            stdin.Write([]byte(string(line_string) + "\n"))
	        } else {
		        fmt.Println(string(line_string))
	        }
        }
        page_num ++
      }
    } else {

			line_num := 0
      page_num := 1

			for {
				line_string, _, err := bfreader.ReadLine()
				if err != nil && err != io.EOF {
					fmt.Println(err)
					os.Exit(1)
				}
				if err == io.EOF {
					break
				}
        if page_num >= args.start && page_num <= args.end {
          if args.output_dest != "" {
            stdin.Write([]byte(string(line_string) + "\n"))
          } else {
            fmt.Println(string(line_string))
          }
        }
        line_num ++
        if line_num == args.page_len {
          page_num ++
          line_num = 0
        }
        if page_num > args.end {
          break
        }
			}
		}
	} else {
		bfscanner := bufio.NewScanner(os.Stdin)
		line_num := 0
    page_num := 1
		out_string := ""
		for bfscanner.Scan() {
			line_string := bfscanner.Text()
			line_string += "\n"
      if page_num >= args.start && page_num <= args.end {
        out_string += line_string
      }
      line_num ++
      if line_num == args.page_len {
        page_num ++
        line_num = 0
      }
		}
		if args.output_dest != "" {
      stdin.Write([]byte(string(out_string) + "\n"))
    } else {
      fmt.Println(string(out_string))
    }
	}

	if args.output_dest != "" {
		stdin.Close()
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	pro_name = os.Args[0]
	var args selpg_args
	Init(&args)
	Process_args(&args)
	Process_input(&args)
}
