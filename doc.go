// Copyright 2024 Dražen Golić. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package kombinat implements generic combinatorial functions and generators
// for producing combinations, permutations, multiset permutations and variations
// out of elements in a slice of any type.
//
// The goal of this library is the ability to efficiently permute elements of a slice for
// different testing and probing scenarios. If you need a mathematics library that also has
// many more features, use [Gonum] instead.
//
// The package consists of the following functions and generators:
//   - Combinations via Chase's Twiddle algorithm, based on this [implementation in C]
//   - Permutations without repetition by using non-recursive [Heap's algorithm]
//   - Multiset permutations (permutations with repetition) based on [Algorithm 1]
//     found in "Loopless Generation of Multiset Permutations using a Constant Number of
//     Variables by Prefix Shifts." by Aaron Williams, 2009
//   - Variations (custom)
//
// Generators are generaly recommended as they are not only faster, but also memory efficient,
// and can store results into different slices. If you need to reuse the results many times,
// functions that generate the entire result set are also available.
//
// Consult README.md for more details.
//
// [Gonum]: https://www.gonum.org/
// [implementation in C]: https://web.archive.org/web/20221024045742/http://www.netlib.no/netlib/toms/382
// [Heap's algorithm]: https://en.wikipedia.org/wiki/Heap's_algorithm
// [Algorithm 1]: https://dl.acm.org/doi/10.5555/1496770.1496877
package kombinat
