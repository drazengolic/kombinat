# Kombinat

Package kombinat implements generic combinatorics functions and generators for producing combinations, permutations, multiset permutations and variations from elements in a slice of any type in Go language.

The goal of this library is the ability to efficiently permute elements of a slice for different testing and probing scenarios. If you need a mathematics library that also has many more features, use [Gonum](https://www.gonum.org/) instead.

The package consists of the following functions and generators:
  - **Combinations** via Chase's Twiddle algorithm, based on this [implementation in C](https://web.archive.org/web/20221024045742/http://www.netlib.no/netlib/toms/382)
  - **Permutations without repetition** by using non-recursive [Heap's algorithm](https://en.wikipedia.org/wiki/Heap's_algorithm)
  - **Multiset permutations** (permutations with repetition) based on [Algorithm 1](https://dl.acm.org/doi/10.5555/1496770.1496877) found in "Loopless Generation of Multiset Permutations using a Constant Number of Variables by Prefix Shifts." by Aaron Williams, 2009
  - **Variations** (custom)

Generators are generaly recommended as they are not only faster, but also memory efficient, and can store results into different slices. If you need to reuse the results many times, functions that generate the entire result set are also available.

# Usage

## Generators

### Initialization

To create & initialize, for example, a `CombinationGenerator`, you can use the following code:

```Go
items := []string{"A", "B", "C"}

gen, err := NewCombinationGenerator(2, items)
// or
gen2 := new(CombinationGenerator[string])
err2 := gen2.Init(2, items)
```

Generator instance can be reused by invoking `Init` with different input arguments of the same type. Previous data will be overwritten.

### Iteration

To have the generator produce another result, call the `Next` method. The method returns `true` until it reaches the end of the sequence, and therefore can be used in a `for` loop.

To access the new result, call `Current` or `CurrentCopy` after calling `Next` (it has to be called at least once!). `Current` will return the internaly used slice, whereas `CurrentCopy` will return a copy of the internal slice. Use the former if you are not going to modify the slice in order to avoid allocation. But if modification of the slice is expected, use `CurrentCopy` in order not to interfere with the generator and cause unexpected behavior.

For example:

```Go
for gen.Next() {
  fmt.Printf("%v\n", gen.Current())
  list = append(list, gen.CurrentCopy())
}
```

It is possible to set a slice that will be updated whenever `Next` is called by using `SetDest` method. This also enables partial updates of a larger slice, for example:

```Go
items := []string{"A", "B", "C"}
dest := []string{"i", "j", "k", "l"}

gen, err := NewCombinationGenerator(2, items)
gen.SetDest(dest[1:3])

gen.Next()
fmt.Printf("%v", dest) // [i, B, C, l]
```

After a new destination slice is set, subsequent calls to `Current` will always return the slice that was set, making the call redundant.

To make the generator start from the beginning, call `Reset`. Reset preserves the destination slice.

## Functions

Function call returns a slice of slices, containing the entire result set. For example: 

```Go
ps, err := Permutations([]int{1, 2, 3, 4, 5}) //[][]int
```

**Warning:** be careful when using these functions, as they could use up a lot of memory for larger values.

# Benchmarks

Some benchmarks run on Apple M1 Pro can be found below:

```
goos: darwin
goarch: arm64
pkg: github.com/drazengolic/kombinat
```

## Combinations

**go test -bench BenchmarkCombinations -benchmem -benchtime=10s**

```
BenchmarkCombinations/c(2,6)=15-8           26261770         442.7 ns/op       704 B/op       18 allocs/op
BenchmarkCombinations/c(3,6)=20-8           20279744         591.3 ns/op      1048 B/op       23 allocs/op
BenchmarkCombinations/c(4,6)=15-8           25988402         459.6 ns/op       960 B/op       18 allocs/op
BenchmarkCombinations/c(5,6)=6-8            53533160         222.0 ns/op       544 B/op        9 allocs/op
BenchmarkCombinations/c(6,6)=1-8            123885088         96.82 ns/op      184 B/op        4 allocs/op
```

**go test -bench BenchmarkCombinationGenerator -benchmem -benchtime=10s**

```
BenchmarkCombinationGenerator/c(2,6)=15-8           81718262         126.7 ns/op        80 B/op        2 allocs/op
BenchmarkCombinationGenerator/c(3,6)=20-8           75011426         157.6 ns/op        88 B/op        2 allocs/op
BenchmarkCombinationGenerator/c(4,6)=15-8           93862041         126.8 ns/op        96 B/op        2 allocs/op
BenchmarkCombinationGenerator/c(5,6)=6-8            151733102         78.96 ns/op      112 B/op        2 allocs/op
BenchmarkCombinationGenerator/c(6,6)=1-8            224967640         53.35 ns/op      112 B/op        2 allocs/op
```

## Permutations

**go test -bench BenchmarkPermutations -benchmem -benchtime=10s**

```
BenchmarkPermutations/p(2)=2-8          133187949         90.06 ns/op      112 B/op        5 allocs/op
BenchmarkPermutations/p(3)=6-8          56977621         208.8 ns/op       336 B/op        9 allocs/op
BenchmarkPermutations/p(4)=24-8         18749431         637.7 ns/op      1408 B/op       27 allocs/op
BenchmarkPermutations/p(5)=120-8         4090677        2933 ns/op        8928 B/op      123 allocs/op
BenchmarkPermutations/p(6)=720-8          668253       17674 ns/op       53088 B/op      723 allocs/op
```

**go test -bench BenchmarkPermutationGenerator -benchmem -benchtime=10s**

```
BenchmarkPermutationGenerator/p(2)=2-8          310151914         38.68 ns/op       32 B/op        2 allocs/op
BenchmarkPermutationGenerator/p(3)=6-8          218699462         54.79 ns/op       48 B/op        2 allocs/op
BenchmarkPermutationGenerator/p(4)=24-8         92384619         128.1 ns/op        64 B/op        2 allocs/op
BenchmarkPermutationGenerator/p(5)=120-8        24015909         499.2 ns/op        96 B/op        2 allocs/op
BenchmarkPermutationGenerator/p(6)=720-8         4256427        2819 ns/op          96 B/op        2 allocs/op
```

## MultiPermutations

**go test -bench BenchmarkMultiPermutations -benchmem -benchtime=10s**

```
BenchmarkMultiPermutations/P(AABCCCDD)-8            123399       97058 ns/op      342145 B/op     1700 allocs/op
BenchmarkMultiPermutations/P(AABB)-8              28815830         414.8 ns/op       832 B/op       14 allocs/op
BenchmarkMultiPermutations/P(ABCD)-8               9078172        1314 ns/op        3136 B/op       34 allocs/op
```

**go test -bench BenchmarkMultiPermutationGenerator -benchmem -benchtime=10s**

```
BenchmarkMultiPermutationGenerator/P(AABCCCDD)-8            447722       24886 ns/op         192 B/op        9 allocs/op
BenchmarkMultiPermutationGenerator/P(AABB)-8              76625335         154.6 ns/op        96 B/op        5 allocs/op
BenchmarkMultiPermutationGenerator/P(ABCD)-8              34821980         343.3 ns/op        96 B/op        5 allocs/op
```

## Variations

**go test -bench BenchmarkVariations -benchmem -benchtime=10s**

```
BenchmarkVariations/v(2,4)=16-8           38591599         296.7 ns/op       656 B/op       18 allocs/op
BenchmarkVariations/v(3,4)=64-8            9806036        1217 ns/op        3096 B/op       66 allocs/op
BenchmarkVariations/v(4,4)=256-8           2516913        4775 ns/op       14368 B/op      258 allocs/op
BenchmarkVariations/v(5,4)=1024-8           577398       20656 ns/op       73776 B/op     1026 allocs/op
BenchmarkVariations/v(6,4)=4096-8           141417       85309 ns/op      294960 B/op     4098 allocs/op
```

**go test -bench BenchmarkVariationGenerator -benchmem -benchtime=10s**

```
BenchmarkVariationGenerator/v(2,4)=16-8           227415219         52.79 ns/op        0 B/op        0 allocs/op
BenchmarkVariationGenerator/v(3,4)=64-8           36819321         325.0 ns/op         0 B/op        0 allocs/op
BenchmarkVariationGenerator/v(4,4)=256-8           8336451        1438 ns/op           0 B/op        0 allocs/op
BenchmarkVariationGenerator/v(5,4)=1024-8          1733893        6920 ns/op           0 B/op        0 allocs/op
BenchmarkVariationGenerator/v(6,4)=4096-8           367791       32557 ns/op           0 B/op        0 allocs/op
```
