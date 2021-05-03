package ccfpq

import (
    "encoding/csv"
    "fmt"
    "io"
    "log"
    "os"
    ds "rdf-ccfpq/go/data_structures"
    "time"
    // "runtime/debug"
    "strconv"
)

const (
    TEST_REPS = 10

    PRINT_NONE  = iota
    PRINT_ERROR = iota
    PRINT_ALL   = iota

    NOT_APPLICABLE = -1
)

var (
    printmode uint = PRINT_ALL
    countok   uint
    counterr  uint
)

func QueryAll(G *Grammar, D ds.Graph) []pair {
    Q := []pair{}
    for n := range D.AllNodes().Iterate() {
        Q = append(Q, *newPair(n, G.StartSymbol))
    }
    return Q
}

func QueryGroups(numberOfGroups int, G *Grammar, D ds.Graph) []pair {
    Q := make([]pair, numberOfGroups)
    for i := range Q {
        Q[i].node = f.NewSuperVertex(f.NewVertexSet())
        Q[i].symbol = G.StartSymbol
    }
    i := 0
    for n := range D.AllNodes().Iterate() {
        Q[i].node.(ds.SuperVertex).Vertices.Add(n)
        i = (i + 1) % len(Q)
    }
    return Q
}

func RunExperiment(G *Grammar, D ds.Graph, Q []pair, expected int) {
    resCount := 0

    R, t, m := Run(D, G, Q)

    for _, p := range Q {
        var node ds.Vertex
        if super, isSuperVertex := p.node.(ds.SuperVertex); isSuperVertex {
            node = super.Vertex
        } else {
            node = p.node
        }
        resCount += R.get(node, p.symbol).Objects().Size()
    }

    test_result := (resCount == expected) || (expected == NOT_APPLICABLE)

    if test_result == true {
        countok++
    } else {
        counterr++
    }

    mustprint := (test_result == false && printmode != PRINT_NONE) ||
        (printmode == PRINT_ALL)

    if mustprint {
        PrintTestResults(test_result, G.Name, D.Name(), D.AllNodes().Size(),
            int(D.Size()), resCount, expected, t, m)
    }
}

func RunExperiments(expfile string, factoryType string) {
    csvfile, err := os.Open(expfile)
    if err != nil {
        log.Fatalln("Couldn't open csv file. Error:", err)
    }
    r := csv.NewReader(csvfile)
    r.Comment = '#'
    for {
        record, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }
        G, D := QuickLoad(record[0], record[1], factoryType)
        checksum, err := strconv.Atoi(record[2])
        if err != nil {
            fmt.Println(err)
            fmt.Println("Skipping this test.")
            continue
        }
        Q := QueryAll(G,D)
        RunExperiment(G, D, Q, checksum)
        f.Reset()
    }
}

func bToMb(b uint64) float64 {
    return float64(b) / 1024 / 1024
}

func PrintTestResults(test_result bool, grammarname, graphname string,
    vertices, edges, resultcount, expected int, cputime time.Duration,
    memory uint64) {
    fmt.Printf("%s\t%s\t%d\t%d\t%d\t%d%8.1f\t%s\n",
        grammarname, graphname,
        vertices, edges, resultcount, int64(cputime/time.Millisecond),
        bToMb(memory), f.Type())
    if test_result == false {
        fmt.Println("\texpected:", expected)
    }
}

func truncateAlign(str string, length int) string {
    if length <= 0 {
        return ""
    }
    truncated := ""
    for i := 0; i < length && i < len(str); i++ {
        truncated += string(str[i])
    }
    for i := len(truncated); i < length; i++ {
        truncated += " "
    }
    return truncated
}

func PrintTestHeaders() {
    fmt.Println("test_result, grammarname, graphname, resultcount, cputime, memory, factory")
}
