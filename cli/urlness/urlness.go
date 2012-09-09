package main

import (
  "flag"
  "fmt"
  "github.com/drio/drio.go/common/files"
  "github.com/drio/drio.go/urlness"
  "io"
  "log"
  "os"
  "runtime/pprof"
  "strings"
)

var usage = []string{
  //"urlness, a tool to find unrelated individuals.",
  "urlness (v0.1)",
  "Usage: urlness -ks <kinship_file> [optional params]",
  "",
  "Madatory params: ",
  "  -ks   : csv file containing phi coefficients per each relation.",
  "",
  "Optional params: ",
  "  -sex     : csv file containing samples' sex",
  "  -phe     : csv file containing the gender of the samples",
  "  -only    : only use this sample when filtering by phi score",
  "  -phi     : maximum phi score allowed between relations",
  "  -optimal : find the biggest subset of unrelated samples ",
  "",
  "Output:",
  "  stdin : csv matrix of phi coefficients for ALL the samples.",
  "  stdout: If -phi used: csv list of samples that pass the phi score",
  "",
  "Examples:",
  "  $ urlness -ks data/kinship.csv > matrix.csv",
  "  $ urlness -ks data/kinship.csv -phi 0.5 > matrix.csv 2> l.csv",
  "  $ urlness -ks data/kinship.csv -phi 0.5 -sex data/441-Gender.csv -phe data/441-Pheno.csv > m.csv 2> l.csv",
  "  $ urlness -ks data/kinship.csv -phi 0.5 -optimal > list.optimal.csv",
}

// For processing the input parameters from the user.
type options struct {
  ksFname, sexFname, pheFname, onlyFname *string
  phiFilter                              *float64
  optimal                                *bool
  cpuProfile, memProfile                 *string
}

func parseArgs() *options {
  o := new(options)
  o.ksFname = flag.String("ks", "", "Kinship csv file.")
  o.sexFname = flag.String("sex", "", "Gender csv file.")
  o.pheFname = flag.String("phe", "", "Gender csv file.")
  o.onlyFname = flag.String("only", "", "Gender csv file.")
  o.phiFilter = flag.Float64("phi", 0, "Maximum phi coefficient allowed.")
  o.optimal = flag.Bool("optimal", false, "Enable optimal.")
  o.cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")
  o.memProfile = flag.String("memprofile", "", "write mem profile to file")

  flag.Parse()

  // TODO: check that the file exists!!!!!!!!!!!!!!!!!!!!
  err := false
  if len(flag.Args()) != 0 {
    fmt.Fprintln(os.Stderr, "Invalid parameter.")
    err = true
  }

  if *o.ksFname == "" {
    fmt.Fprintln(os.Stderr, "ERROR: Need kinship file")
    err = true
  }

  if *o.onlyFname != "" && *o.phiFilter == 0 {
    fmt.Fprintln(os.Stderr, "ERROR: -only requires -phi param")
    err = true
  }

  if *o.optimal && *o.phiFilter == 0 {
    fmt.Fprintln(os.Stderr, "ERROR: -optimal requires -phi param")
    err = true
  }

  if err {
    fmt.Println(strings.Join(usage, "\n"))
    os.Exit(0)
  }

  return o
}

func main() {
  o := parseArgs()
  inputData := new(urlness.Options)

  // link between file paths and their locations in the datastructure
  // that we will pass to the urlness package (inputData)
  fNamesToFiles := map[*string]*io.Reader{
    o.ksFname:   &inputData.KS,
    o.sexFname:  &inputData.Sex,
    o.pheFname:  &inputData.Phe,
    o.onlyFname: &inputData.Only,
  }

  // Open the files and set the readers for them in inputData
  for path, reader := range fNamesToFiles {
    if *path != "" {
      file, r := files.Xopen(*path)
      *reader = r
      defer file.Close()
    }
  }

  // We know the Phi param has to be there for sure
  inputData.PhiFilter = *o.phiFilter

  // CPU Profiling
  if *o.cpuProfile != "" {
    f, err := os.Create(*o.cpuProfile)
    if err != nil {
      log.Fatal(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
  }

  // Mem profile
  if *o.memProfile != "" {
    f, err := os.Create(*o.memProfile)
    if err != nil {
      log.Fatal(err)
    }
    pprof.WriteHeapProfile(f)
    f.Close()
    return
  }

  // Run the appropiate routine
  if *o.optimal {
    fmt.Println(urlness.ComputeOptimal(*inputData))
  } else {
    m, l := urlness.Compute(*inputData) // matrix and list of unrelated samples
    fmt.Println(m)
    if *o.phiFilter != 0 {
      fmt.Fprintln(os.Stderr, l)
    }
  }
}
