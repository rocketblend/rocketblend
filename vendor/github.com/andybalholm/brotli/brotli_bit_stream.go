package brotli

import (
	"math"
	"sync"
)

const maxHuffmanTreeSize = (2*numCommandSymbols + 1)

/*
The maximum size of Huffman dictionary for distances assuming that

	NPOSTFIX = 0 and NDIRECT = 0.
*/
const maxSimpleDistanceAlphabetSize = 140

/*
Represents the range of values belonging to a prefix code:

	[offset, offset + 2^nbits)
*/
type prefixCodeRange struct {
	offset uint32
	nbits  uint32
}

var kBlockLengthPrefixCode = [numBlockLenSymbols]prefixCodeRange{
	{1, 2},
	{5, 2},
	{9, 2},
	{13, 2},
	{17, 3},
	{25, 3},
	{33, 3},
	{41, 3},
	{49, 4},
	{65, 4},
	{81, 4},
	{97, 4},
	{113, 5},
	{145, 5},
	{177, 5},
	{209, 5},
	{241, 6},
	{305, 6},
	{369, 7},
	{497, 8},
	{753, 9},
	{1265, 10},
	{2289, 11},
	{4337, 12},
	{8433, 13},
	{16625, 24},
}

func blockLengthPrefixCode(len uint32) uint32 {
	var code uint32
	if len >= 177 {
		if len >= 753 {
			code = 20
		} else {
			code = 14
		}
	} else if len >= 41 {
		code = 7
	} else {
		code = 0
	}
	for code < (numBlockLenSymbols-1) && len >= kBlockLengthPrefixCode[code+1].offset {
		code++
	}
	return code
}

func getBlockLengthPrefixCode(len uint32, code *uint, n_extra *uint32, extra *uint32) {
	*code = uint(blockLengthPrefixCode(uint32(len)))
	*n_extra = kBlockLengthPrefixCode[*code].nbits
	*extra = len - kBlockLengthPrefixCode[*code].offset
}

type blockTypeCodeCalculator struct {
	last_type        uint
	second_last_type uint
}

func initBlockTypeCodeCalculator(self *blockTypeCodeCalculator) {
	self.last_type = 1
	self.second_last_type = 0
}

func nextBlockTypeCode(calculator *blockTypeCodeCalculator, type_ byte) uint {
	var type_code uint
	if uint(type_) == calculator.last_type+1 {
		type_code = 1
	} else if uint(type_) == calculator.second_last_type {
		type_code = 0
	} else {
		type_code = uint(type_) + 2
	}
	calculator.second_last_type = calculator.last_type
	calculator.last_type = uint(type_)
	return type_code
}

/*
|nibblesbits| represents the 2 bits to encode MNIBBLES (0-3)

	REQUIRES: length > 0
	REQUIRES: length <= (1 << 24)
*/
func encodeMlen(length uint, bits *uint64, numbits *uint, nibblesbits *uint64) {
	var lg uint
	if length == 1 {
		lg = 1
	} else {
		lg = uint(log2FloorNonZero(uint(uint32(length-1)))) + 1
	}
	var tmp uint
	if lg < 16 {
		tmp = 16
	} else {
		tmp = (lg + 3)
	}
	var mnibbles uint = tmp / 4
	assert(length > 0)
	assert(length <= 1<<24)
	assert(lg <= 24)
	*nibblesbits = uint64(mnibbles) - 4
	*numbits = mnibbles * 4
	*bits = uint64(length) - 1
}

func storeCommandExtra(cmd *command, bw *bitWriter) {
	var copylen_code uint32 = commandCopyLenCode(cmd)
	var inscode uint16 = getInsertLengthCode(uint(cmd.insert_len_))
	var copycode uint16 = getCopyLengthCode(uint(copylen_code))
	var insnumextra uint32 = getInsertExtra(inscode)
	var insextraval uint64 = uint64(cmd.insert_len_) - uint64(getInsertBase(inscode))
	var copyextraval uint64 = uint64(copylen_code) - uint64(getCopyBase(copycode))
	var bits uint64 = copyextraval<<insnumextra | insextraval
	bw.writeBits(uint(insnumextra+getCopyExtra(copycode)), bits)
}

/*
Data structure that stores almost everything that is needed to encode each

	block switch command.
*/
type blockSplitCode struct {
	type_code_calculator blockTypeCodeCalculator
	type_depths          [maxBlockTypeSymbols]byte
	type_bits            [maxBlockTypeSymbols]uint16
	length_depths        [numBlockLenSymbols]byte
	length_bits          [numBlockLenSymbols]uint16
}

/* Stores a number between 0 and 255. */
func storeVarLenUint8(n uint, bw *bitWriter) {
	if n == 0 {
		bw.writeBits(1, 0)
	} else {
		var nbits uint = uint(log2FloorNonZero(n))
		bw.writeBits(1, 1)
		bw.writeBits(3, uint64(nbits))
		bw.writeBits(nbits, uint64(n)-(uint64(uint(1))<<nbits))
	}
}

/*
Stores the compressed meta-block header.

	REQUIRES: length > 0
	REQUIRES: length <= (1 << 24)
*/
func storeCompressedMetaBlockHeader(is_final_block bool, length uint, bw *bitWriter) {
	var lenbits uint64
	var nlenbits uint
	var nibblesbits uint64
	var is_final uint64
	if is_final_block {
		is_final = 1
	} else {
		is_final = 0
	}

	/* Write ISLAST bit. */
	bw.writeBits(1, is_final)

	/* Write ISEMPTY bit. */
	if is_final_block {
		bw.writeBits(1, 0)
	}

	encodeMlen(length, &lenbits, &nlenbits, &nibblesbits)
	bw.writeBits(2, nibblesbits)
	bw.writeBits(nlenbits, lenbits)

	if !is_final_block {
		/* Write ISUNCOMPRESSED bit. */
		bw.writeBits(1, 0)
	}
}

/*
Stores the uncompressed meta-block header.

	REQUIRES: length > 0
	REQUIRES: length <= (1 << 24)
*/
func storeUncompressedMetaBlockHeader(length uint, bw *bitWriter) {
	var lenbits uint64
	var nlenbits uint
	var nibblesbits uint64

	/* Write ISLAST bit.
	   Uncompressed block cannot be the last one, so set to 0. */
	bw.writeBits(1, 0)

	encodeMlen(length, &lenbits, &nlenbits, &nibblesbits)
	bw.writeBits(2, nibblesbits)
	bw.writeBits(nlenbits, lenbits)

	/* Write ISUNCOMPRESSED bit. */
	bw.writeBits(1, 1)
}

var storeHuffmanTreeOfHuffmanTreeToBitMask_kStorageOrder = [codeLengthCodes]byte{1, 2, 3, 4, 0, 5, 17, 6, 16, 7, 8, 9, 10, 11, 12, 13, 14, 15}

var storeHuffmanTreeOfHuffmanTreeToBitMask_kHuffmanBitLengthHuffmanCodeSymbols = [6]byte{0, 7, 3, 2, 1, 15}
var storeHuffmanTreeOfHuffmanTreeToBitMask_kHuffmanBitLengthHuffmanCodeBitLengths = [6]byte{2, 4, 3, 2, 2, 4}

func storeHuffmanTreeOfHuffmanTreeToBitMask(num_codes int, code_length_bitdepth []byte, bw *bitWriter) {
	var skip_some uint = 0
	var codes_to_store uint = codeLengthCodes
	/* The bit lengths of the Huffman code over the code length alphabet
	   are compressed with the following static Huffman code:
	     Symbol   Code
	     ------   ----
	     0          00
	     1        1110
	     2         110
	     3          01
	     4          10
	     5        1111 */

	/* Throw away trailing zeros: */
	if num_codes > 1 {
		for ; codes_to_store > 0; codes_to_store-- {
			if code_length_bitdepth[storeHuffmanTreeOfHuffmanTreeToBitMask_kStorageOrder[codes_to_store-1]] != 0 {
				break
			}
		}
	}

	if code_length_bitdepth[storeHuffmanTreeOfHuffmanTreeToBitMask_kStorageOrder[0]] == 0 && code_length_bitdepth[storeHuffmanTreeOfHuffmanTreeToBitMask_kStorageOrder[1]] == 0 {
		skip_some = 2 /* skips two. */
		if code_length_bitdepth[storeHuffmanTreeOfHuffmanTreeToBitMask_kStorageOrder[2]] == 0 {
			skip_some = 3 /* skips three. */
		}
	}

	bw.writeBits(2, uint64(skip_some))
	{
		var i uint
		for i = skip_some; i < codes_to_store; i++ {
			var l uint = uint(code_length_bitdepth[storeHuffmanTreeOfHuffmanTreeToBitMask_kStorageOrder[i]])
			bw.writeBits(uint(storeHuffmanTreeOfHuffmanTreeToBitMask_kHuffmanBitLengthHuffmanCodeBitLengths[l]), uint64(storeHuffmanTreeOfHuffmanTreeToBitMask_kHuffmanBitLengthHuffmanCodeSymbols[l]))
		}
	}
}

func storeHuffmanTreeToBitMask(huffman_tree_size uint, huffman_tree []byte, huffman_tree_extra_bits []byte, code_length_bitdepth []byte, code_length_bitdepth_symbols []uint16, bw *bitWriter) {
	var i uint
	for i = 0; i < huffman_tree_size; i++ {
		var ix uint = uint(huffman_tree[i])
		bw.writeBits(uint(code_length_bitdepth[ix]), uint64(code_length_bitdepth_symbols[ix]))

		/* Extra bits */
		switch ix {
		case repeatPreviousCodeLength:
			bw.writeBits(2, uint64(huffman_tree_extra_bits[i]))

		case repeatZeroCodeLength:
			bw.writeBits(3, uint64(huffman_tree_extra_bits[i]))
		}
	}
}

func storeSimpleHuffmanTree(depths []byte, symbols []uint, num_symbols uint, max_bits uint, bw *bitWriter) {
	/* value of 1 indicates a simple Huffman code */
	bw.writeBits(2, 1)

	bw.writeBits(2, uint64(num_symbols)-1) /* NSYM - 1 */
	{
		/* Sort */
		var i uint
		for i = 0; i < num_symbols; i++ {
			var j uint
			for j = i + 1; j < num_symbols; j++ {
				if depths[symbols[j]] < depths[symbols[i]] {
					var tmp uint = symbols[j]
					symbols[j] = symbols[i]
					symbols[i] = tmp
				}
			}
		}
	}

	if num_symbols == 2 {
		bw.writeBits(max_bits, uint64(symbols[0]))
		bw.writeBits(max_bits, uint64(symbols[1]))
	} else if num_symbols == 3 {
		bw.writeBits(max_bits, uint64(symbols[0]))
		bw.writeBits(max_bits, uint64(symbols[1]))
		bw.writeBits(max_bits, uint64(symbols[2]))
	} else {
		bw.writeBits(max_bits, uint64(symbols[0]))
		bw.writeBits(max_bits, uint64(symbols[1]))
		bw.writeBits(max_bits, uint64(symbols[2]))
		bw.writeBits(max_bits, uint64(symbols[3]))

		/* tree-select */
		var tmp int
		if depths[symbols[0]] == 1 {
			tmp = 1
		} else {
			tmp = 0
		}
		bw.writeBits(1, uint64(tmp))
	}
}

/*
num = alphabet size

	depths = symbol depths
*/
func storeHuffmanTree(depths []byte, num uint, tree []huffmanTree, bw *bitWriter) {
	var huffman_tree [numCommandSymbols]byte
	var huffman_tree_extra_bits [numCommandSymbols]byte
	var huffman_tree_size uint = 0
	var code_length_bitdepth = [codeLengthCodes]byte{0}
	var code_length_bitdepth_symbols [codeLengthCodes]uint16
	var huffman_tree_histogram = [codeLengthCodes]uint32{0}
	var i uint
	var num_codes int = 0
	/* Write the Huffman tree into the brotli-representation.
	   The command alphabet is the largest, so this allocation will fit all
	   alphabets. */

	var code uint = 0

	assert(num <= numCommandSymbols)

	writeHuffmanTree(depths, num, &huffman_tree_size, huffman_tree[:], huffman_tree_extra_bits[:])

	/* Calculate the statistics of the Huffman tree in brotli-representation. */
	for i = 0; i < huffman_tree_size; i++ {
		huffman_tree_histogram[huffman_tree[i]]++
	}

	for i = 0; i < codeLengthCodes; i++ {
		if huffman_tree_histogram[i] != 0 {
			if num_codes == 0 {
				code = i
				num_codes = 1
			} else if num_codes == 1 {
				num_codes = 2
				break
			}
		}
	}

	/* Calculate another Huffman tree to use for compressing both the
	   earlier Huffman tree with. */
	createHuffmanTree(huffman_tree_histogram[:], codeLengthCodes, 5, tree, code_length_bitdepth[:])

	convertBitDepthsToSymbols(code_length_bitdepth[:], codeLengthCodes, code_length_bitdepth_symbols[:])

	/* Now, we have all the data, let's start storing it */
	storeHuffmanTreeOfHuffmanTreeToBitMask(num_codes, code_length_bitdepth[:], bw)

	if num_codes == 1 {
		code_length_bitdepth[code] = 0
	}

	/* Store the real Huffman tree now. */
	storeHuffmanTreeToBitMask(huffman_tree_size, huffman_tree[:], huffman_tree_extra_bits[:], code_length_bitdepth[:], code_length_bitdepth_symbols[:], bw)
}

/*
Builds a Huffman tree from histogram[0:length] into depth[0:length] and

	bits[0:length] and stores the encoded tree to the bit stream.
*/
func buildAndStoreHuffmanTree(histogram []uint32, histogram_length uint, alphabet_size uint, tree []huffmanTree, depth []byte, bits []uint16, bw *bitWriter) {
	var count uint = 0
	var s4 = [4]uint{0}
	var i uint
	var max_bits uint = 0
	for i = 0; i < histogram_length; i++ {
		if histogram[i] != 0 {
			if count < 4 {
				s4[count] = i
			} else if count > 4 {
				break
			}

			count++
		}
	}
	{
		var max_bits_counter uint = alphabet_size - 1
		for max_bits_counter != 0 {
			max_bits_counter >>= 1
			max_bits++
		}
	}

	if count <= 1 {
		bw.writeBits(4, 1)
		bw.writeBits(max_bits, uint64(s4[0]))
		depth[s4[0]] = 0
		bits[s4[0]] = 0
		return
	}

	for i := 0; i < int(histogram_length); i++ {
		depth[i] = 0
	}
	createHuffmanTree(histogram, histogram_length, 15, tree, depth)
	convertBitDepthsToSymbols(depth, histogram_length, bits)

	if count <= 4 {
		storeSimpleHuffmanTree(depth, s4[:], count, max_bits, bw)
	} else {
		storeHuffmanTree(depth, histogram_length, tree, bw)
	}
}

func sortHuffmanTree1(v0 huffmanTree, v1 huffmanTree) bool {
	return v0.total_count_ < v1.total_count_
}

var huffmanTreePool sync.Pool

func buildAndStoreHuffmanTreeFast(histogram []uint32, histogram_total uint, max_bits uint, depth []byte, bits []uint16, bw *bitWriter) {
	var count uint = 0
	var symbols = [4]uint{0}
	var length uint = 0
	var total uint = histogram_total
	for total != 0 {
		if histogram[length] != 0 {
			if count < 4 {
				symbols[count] = length
			}

			count++
			total -= uint(histogram[length])
		}

		length++
	}

	if count <= 1 {
		bw.writeBits(4, 1)
		bw.writeBits(max_bits, uint64(symbols[0]))
		depth[symbols[0]] = 0
		bits[symbols[0]] = 0
		return
	}

	for i := 0; i < int(length); i++ {
		depth[i] = 0
	}
	{
		var max_tree_size uint = 2*length + 1
		tree, _ := huffmanTreePool.Get().(*[]huffmanTree)
		if tree == nil || cap(*tree) < int(max_tree_size) {
			tmp := make([]huffmanTree, max_tree_size)
			tree = &tmp
		} else {
			*tree = (*tree)[:max_tree_size]
		}
		var count_limit uint32
		for count_limit = 1; ; count_limit *= 2 {
			var node int = 0
			var l uint
			for l = length; l != 0; {
				l--
				if histogram[l] != 0 {
					if histogram[l] >= count_limit {
						initHuffmanTree(&(*tree)[node:][0], histogram[l], -1, int16(l))
					} else {
						initHuffmanTree(&(*tree)[node:][0], count_limit, -1, int16(l))
					}

					node++
				}
			}
			{
				var n int = node
				/* Points to the next leaf node. */ /* Points to the next non-leaf node. */
				var sentinel huffmanTree
				var i int = 0
				var j int = n + 1
				var k int

				sortHuffmanTreeItems(*tree, uint(n), huffmanTreeComparator(sortHuffmanTree1))

				/* The nodes are:
				   [0, n): the sorted leaf nodes that we start with.
				   [n]: we add a sentinel here.
				   [n + 1, 2n): new parent nodes are added here, starting from
				                (n+1). These are naturally in ascending order.
				   [2n]: we add a sentinel at the end as well.
				   There will be (2n+1) elements at the end. */
				initHuffmanTree(&sentinel, math.MaxUint32, -1, -1)

				(*tree)[node] = sentinel
				node++
				(*tree)[node] = sentinel
				node++

				for k = n - 1; k > 0; k-- {
					var left int
					var right int
					if (*tree)[i].total_count_ <= (*tree)[j].total_count_ {
						left = i
						i++
					} else {
						left = j
						j++
					}

					if (*tree)[i].total_count_ <= (*tree)[j].total_count_ {
						right = i
						i++
					} else {
						right = j
						j++
					}

					/* The sentinel node becomes the parent node. */
					(*tree)[node-1].total_count_ = (*tree)[left].total_count_ + (*tree)[right].total_count_

					(*tree)[node-1].index_left_ = int16(left)
					(*tree)[node-1].index_right_or_value_ = int16(right)

					/* Add back the last sentinel node. */
					(*tree)[node] = sentinel
					node++
				}

				if setDepth(2*n-1, *tree, depth, 14) {
					/* We need to pack the Huffman tree in 14 bits. If this was not
					   successful, add fake entities to the lowest values and retry. */
					break
				}
			}
		}

		huffmanTreePool.Put(tree)
	}

	convertBitDepthsToSymbols(depth, length, bits)
	if count <= 4 {
		var i uint

		/* value of 1 indicates a simple Huffman code */
		bw.writeBits(2, 1)

		bw.writeBits(2, uint64(count)-1) /* NSYM - 1 */

		/* Sort */
		for i = 0; i < count; i++ {
			var j uint
			for j = i + 1; j < count; j++ {
				if depth[symbols[j]] < depth[symbols[i]] {
					var tmp uint = symbols[j]
					symbols[j] = symbols[i]
					symbols[i] = tmp
				}
			}
		}

		if count == 2 {
			bw.writeBits(max_bits, uint64(symbols[0]))
			bw.writeBits(max_bits, uint64(symbols[1]))
		} else if count == 3 {
			bw.writeBits(max_bits, uint64(symbols[0]))
			bw.writeBits(max_bits, uint64(symbols[1]))
			bw.writeBits(max_bits, uint64(symbols[2]))
		} else {
			bw.writeBits(max_bits, uint64(symbols[0]))
			bw.writeBits(max_bits, uint64(symbols[1]))
			bw.writeBits(max_bits, uint64(symbols[2]))
			bw.writeBits(max_bits, uint64(symbols[3]))

			/* tree-select */
			bw.writeSingleBit(depth[symbols[0]] == 1)
		}
	} else {
		var previous_value byte = 8
		var i uint

		/* Complex Huffman Tree */
		storeStaticCodeLengthCode(bw)

		/* Actual RLE coding. */
		for i = 0; i < length; {
			var value byte = depth[i]
			var reps uint = 1
			var k uint
			for k = i + 1; k < length && depth[k] == value; k++ {
				reps++
			}

			i += reps
			if value == 0 {
				bw.writeBits(uint(kZeroRepsDepth[reps]), kZeroRepsBits[reps])
			} else {
				if previous_value != value {
					bw.writeBits(uint(kCodeLengthDepth[value]), uint64(kCodeLengthBits[value]))
					reps--
				}

				if reps < 3 {
					for reps != 0 {
						reps--
						bw.writeBits(uint(kCodeLengthDepth[value]), uint64(kCodeLengthBits[value]))
					}
				} else {
					reps -= 3
					bw.writeBits(uint(kNonZeroRepsDepth[reps]), kNonZeroRepsBits[reps])
				}

				previous_value = value
			}
		}
	}
}

func indexOf(v []byte, v_size uint, value byte) uint {
	var i uint = 0
	for ; i < v_size; i++ {
		if v[i] == value {
			return i
		}
	}

	return i
}

func moveToFront(v []byte, index uint) {
	var value byte = v[index]
	var i uint
	for i = index; i != 0; i-- {
		v[i] = v[i-1]
	}

	v[0] = value
}

func moveToFrontTransform(v_in []uint32, v_size uint, v_out []uint32) {
	var i uint
	var mtf [256]byte
	var max_value uint32
	if v_size == 0 {
		return
	}

	max_value = v_in[0]
	for i = 1; i < v_size; i++ {
		if v_in[i] > max_value {
			max_value = v_in[i]
		}
	}

	assert(max_value < 256)
	for i = 0; uint32(i) <= max_value; i++ {
		mtf[i] = byte(i)
	}
	{
		var mtf_size uint = uint(max_value + 1)
		for i = 0; i < v_size; i++ {
			var index uint = indexOf(mtf[:], mtf_size, byte(v_in[i]))
			assert(index < mtf_size)
			v_out[i] = uint32(index)
			moveToFront(mtf[:], index)
		}
	}
}

/*
Finds runs of zeros in v[0..in_size) and replaces them with a prefix code of

	the run length plus extra bits (lower 9 bits is the prefix code and the rest
	are the extra bits). Non-zero values in v[] are shifted by
	*max_length_prefix. Will not create prefix codes bigger than the initial
	value of *max_run_length_prefix. The prefix code of run length L is simply
	Log2Floor(L) and the number of extra bits is the same as the prefix code.
*/
func runLengthCodeZeros(in_size uint, v []uint32, out_size *uint, max_run_length_prefix *uint32) {
	var max_reps uint32 = 0
	var i uint
	var max_prefix uint32
	for i = 0; i < in_size; {
		var reps uint32 = 0
		for ; i < in_size && v[i] != 0; i++ {
		}
		for ; i < in_size && v[i] == 0; i++ {
			reps++
		}

		max_reps = brotli_max_uint32_t(reps, max_reps)
	}

	if max_reps > 0 {
		max_prefix = log2FloorNonZero(uint(max_reps))
	} else {
		max_prefix = 0
	}
	max_prefix = brotli_min_uint32_t(max_prefix, *max_run_length_prefix)
	*max_run_length_prefix = max_prefix
	*out_size = 0
	for i = 0; i < in_size; {
		assert(*out_size <= i)
		if v[i] != 0 {
			v[*out_size] = v[i] + *max_run_length_prefix
			i++
			(*out_size)++
		} else {
			var reps uint32 = 1
			var k uint
			for k = i + 1; k < in_size && v[k] == 0; k++ {
				reps++
			}

			i += uint(reps)
			for reps != 0 {
				if reps < 2<<max_prefix {
					var run_length_prefix uint32 = log2FloorNonZero(uint(reps))
					var extra_bits uint32 = reps - (1 << run_length_prefix)
					v[*out_size] = run_length_prefix + (extra_bits << 9)
					(*out_size)++
					break
				} else {
					var extra_bits uint32 = (1 << max_prefix) - 1
					v[*out_size] = max_prefix + (extra_bits << 9)
					reps -= (2 << max_prefix) - 1
					(*out_size)++
				}
			}
		}
	}
}

const symbolBits = 9

var encodeContextMap_kSymbolMask uint32 = (1 << symbolBits) - 1

func encodeContextMap(context_map []uint32, context_map_size uint, num_clusters uint, tree []huffmanTree, bw *bitWriter) {
	var i uint
	var rle_symbols []uint32
	var max_run_length_prefix uint32 = 6
	var num_rle_symbols uint = 0
	var histogram [maxContextMapSymbols]uint32
	var depths [maxContextMapSymbols]byte
	var bits [maxContextMapSymbols]uint16

	storeVarLenUint8(num_clusters-1, bw)

	if num_clusters == 1 {
		return
	}

	rle_symbols = make([]uint32, context_map_size)
	moveToFrontTransform(context_map, context_map_size, rle_symbols)
	runLengthCodeZeros(context_map_size, rle_symbols, &num_rle_symbols, &max_run_length_prefix)
	histogram = [maxContextMapSymbols]uint32{}
	for i = 0; i < num_rle_symbols; i++ {
		histogram[rle_symbols[i]&encodeContextMap_kSymbolMask]++
	}
	{
		var use_rle bool = (max_run_length_prefix > 0)
		bw.writeSingleBit(use_rle)
		if use_rle {
			bw.writeBits(4, uint64(max_run_length_prefix)-1)
		}
	}

	buildAndStoreHuffmanTree(histogram[:], uint(uint32(num_clusters)+max_run_length_prefix), uint(uint32(num_clusters)+max_run_length_prefix), tree, depths[:], bits[:], bw)
	for i = 0; i < num_rle_symbols; i++ {
		var rle_symbol uint32 = rle_symbols[i] & encodeContextMap_kSymbolMask
		var extra_bits_val uint32 = rle_symbols[i] >> symbolBits
		bw.writeBits(uint(depths[rle_symbol]), uint64(bits[rle_symbol]))
		if rle_symbol > 0 && rle_symbol <= max_run_length_prefix {
			bw.writeBits(uint(rle_symbol), uint64(extra_bits_val))
		}
	}

	bw.writeBits(1, 1) /* use move-to-front */
	rle_symbols = nil
}

/* Stores the block switch command with index block_ix to the bit stream. */
func storeBlockSwitch(code *blockSplitCode, block_len uint32, block_type byte, is_first_block bool, bw *bitWriter) {
	var typecode uint = nextBlockTypeCode(&code.type_code_calculator, block_type)
	var lencode uint
	var len_nextra uint32
	var len_extra uint32
	if !is_first_block {
		bw.writeBits(uint(code.type_depths[typecode]), uint64(code.type_bits[typecode]))
	}

	getBlockLengthPrefixCode(block_len, &lencode, &len_nextra, &len_extra)

	bw.writeBits(uint(code.length_depths[lencode]), uint64(code.length_bits[lencode]))
	bw.writeBits(uint(len_nextra), uint64(len_extra))
}

/*
Builds a BlockSplitCode data structure from the block split given by the

	vector of block types and block lengths and stores it to the bit stream.
*/
func buildAndStoreBlockSplitCode(types []byte, lengths []uint32, num_blocks uint, num_types uint, tree []huffmanTree, code *blockSplitCode, bw *bitWriter) {
	var type_histo [maxBlockTypeSymbols]uint32
	var length_histo [numBlockLenSymbols]uint32
	var i uint
	var type_code_calculator blockTypeCodeCalculator
	for i := 0; i < int(num_types+2); i++ {
		type_histo[i] = 0
	}
	length_histo = [numBlockLenSymbols]uint32{}
	initBlockTypeCodeCalculator(&type_code_calculator)
	for i = 0; i < num_blocks; i++ {
		var type_code uint = nextBlockTypeCode(&type_code_calculator, types[i])
		if i != 0 {
			type_histo[type_code]++
		}
		length_histo[blockLengthPrefixCode(lengths[i])]++
	}

	storeVarLenUint8(num_types-1, bw)
	if num_types > 1 { /* TODO: else? could StoreBlockSwitch occur? */
		buildAndStoreHuffmanTree(type_histo[0:], num_types+2, num_types+2, tree, code.type_depths[0:], code.type_bits[0:], bw)
		buildAndStoreHuffmanTree(length_histo[0:], numBlockLenSymbols, numBlockLenSymbols, tree, code.length_depths[0:], code.length_bits[0:], bw)
		storeBlockSwitch(code, lengths[0], types[0], true, bw)
	}
}

/* Stores a context map where the histogram type is always the block type. */
func storeTrivialContextMap(num_types uint, context_bits uint, tree []huffmanTree, bw *bitWriter) {
	storeVarLenUint8(num_types-1, bw)
	if num_types > 1 {
		var repeat_code uint = context_bits - 1
		var repeat_bits uint = (1 << repeat_code) - 1
		var alphabet_size uint = num_types + repeat_code
		var histogram [maxContextMapSymbols]uint32
		var depths [maxContextMapSymbols]byte
		var bits [maxContextMapSymbols]uint16
		var i uint
		for i := 0; i < int(alphabet_size); i++ {
			histogram[i] = 0
		}

		/* Write RLEMAX. */
		bw.writeBits(1, 1)

		bw.writeBits(4, uint64(repeat_code)-1)
		histogram[repeat_code] = uint32(num_types)
		histogram[0] = 1
		for i = context_bits; i < alphabet_size; i++ {
			histogram[i] = 1
		}

		buildAndStoreHuffmanTree(histogram[:], alphabet_size, alphabet_size, tree, depths[:], bits[:], bw)
		for i = 0; i < num_types; i++ {
			var tmp uint
			if i == 0 {
				tmp = 0
			} else {
				tmp = i + context_bits - 1
			}
			var code uint = tmp
			bw.writeBits(uint(depths[code]), uint64(bits[code]))
			bw.writeBits(uint(depths[repeat_code]), uint64(bits[repeat_code]))
			bw.writeBits(repeat_code, uint64(repeat_bits))
		}

		/* Write IMTF (inverse-move-to-front) bit. */
		bw.writeBits(1, 1)
	}
}

/* Manages the encoding of one block category (literal, command or distance). */
type blockEncoder struct {
	histogram_length_ uint
	num_block_types_  uint
	block_types_      []byte
	block_lengths_    []uint32
	num_blocks_       uint
	block_split_code_ blockSplitCode
	block_ix_         uint
	block_len_        uint
	entropy_ix_       uint
	depths_           []byte
	bits_             []uint16
}

var blockEncoderPool sync.Pool

func getBlockEncoder(histogram_length uint, num_block_types uint, block_types []byte, block_lengths []uint32, num_blocks uint) *blockEncoder {
	self, _ := blockEncoderPool.Get().(*blockEncoder)

	if self != nil {
		self.block_ix_ = 0
		self.entropy_ix_ = 0
		self.depths_ = self.depths_[:0]
		self.bits_ = self.bits_[:0]
	} else {
		self = &blockEncoder{}
	}

	self.histogram_length_ = histogram_length
	self.num_block_types_ = num_block_types
	self.block_types_ = block_types
	self.block_lengths_ = block_lengths
	self.num_blocks_ = num_blocks
	initBlockTypeCodeCalculator(&self.block_split_code_.type_code_calculator)
	if num_blocks == 0 {
		self.block_len_ = 0
	} else {
		self.block_len_ = uint(block_lengths[0])
	}

	return self
}

func cleanupBlockEncoder(self *blockEncoder) {
	blockEncoderPool.Put(self)
}

/*
Creates entropy codes of block lengths and block types and stores them

	to the bit stream.
*/
func buildAndStoreBlockSwitchEntropyCodes(self *blockEncoder, tree []huffmanTree, bw *bitWriter) {
	buildAndStoreBlockSplitCode(self.block_types_, self.block_lengths_, self.num_blocks_, self.num_block_types_, tree, &self.block_split_code_, bw)
}

/*
Stores the next symbol with the entropy code of the current block type.

	Updates the block type and block length at block boundaries.
*/
func storeSymbol(self *blockEncoder, symbol uint, bw *bitWriter) {
	if self.block_len_ == 0 {
		self.block_ix_++
		var block_ix uint = self.block_ix_
		var block_len uint32 = self.block_lengths_[block_ix]
		var block_type byte = self.block_types_[block_ix]
		self.block_len_ = uint(block_len)
		self.entropy_ix_ = uint(block_type) * self.histogram_length_
		storeBlockSwitch(&self.block_split_code_, block_len, block_type, false, bw)
	}

	self.block_len_--
	{
		var ix uint = self.entropy_ix_ + symbol
		bw.writeBits(uint(self.depths_[ix]), uint64(self.bits_[ix]))
	}
}

/*
Stores the next symbol with the entropy code of the current block type and

	context value.
	Updates the block type and block length at block boundaries.
*/
func storeSymbolWithContext(self *blockEncoder, symbol uint, context uint, context_map []uint32, bw *bitWriter, context_bits uint) {
	if self.block_len_ == 0 {
		self.block_ix_++
		var block_ix uint = self.block_ix_
		var block_len uint32 = self.block_lengths_[block_ix]
		var block_type byte = self.block_types_[block_ix]
		self.block_len_ = uint(block_len)
		self.entropy_ix_ = uint(block_type) << context_bits
		storeBlockSwitch(&self.block_split_code_, block_len, block_type, false, bw)
	}

	self.block_len_--
	{
		var histo_ix uint = uint(context_map[self.entropy_ix_+context])
		var ix uint = histo_ix*self.histogram_length_ + symbol
		bw.writeBits(uint(self.depths_[ix]), uint64(self.bits_[ix]))
	}
}

func buildAndStoreEntropyCodesLiteral(self *blockEncoder, histograms []histogramLiteral, histograms_size uint, alphabet_size uint, tree []huffmanTree, bw *bitWriter) {
	var table_size uint = histograms_size * self.histogram_length_
	if cap(self.depths_) < int(table_size) {
		self.depths_ = make([]byte, table_size)
	} else {
		self.depths_ = self.depths_[:table_size]
	}
	if cap(self.bits_) < int(table_size) {
		self.bits_ = make([]uint16, table_size)
	} else {
		self.bits_ = self.bits_[:table_size]
	}
	{
		var i uint
		for i = 0; i < histograms_size; i++ {
			var ix uint = i * self.histogram_length_
			buildAndStoreHuffmanTree(histograms[i].data_[0:], self.histogram_length_, alphabet_size, tree, self.depths_[ix:], self.bits_[ix:], bw)
		}
	}
}

func buildAndStoreEntropyCodesCommand(self *blockEncoder, histograms []histogramCommand, histograms_size uint, alphabet_size uint, tree []huffmanTree, bw *bitWriter) {
	var table_size uint = histograms_size * self.histogram_length_
	if cap(self.depths_) < int(table_size) {
		self.depths_ = make([]byte, table_size)
	} else {
		self.depths_ = self.depths_[:table_size]
	}
	if cap(self.bits_) < int(table_size) {
		self.bits_ = make([]uint16, table_size)
	} else {
		self.bits_ = self.bits_[:table_size]
	}
	{
		var i uint
		for i = 0; i < histograms_size; i++ {
			var ix uint = i * self.histogram_length_
			buildAndStoreHuffmanTree(histograms[i].data_[0:], self.histogram_length_, alphabet_size, tree, self.depths_[ix:], self.bits_[ix:], bw)
		}
	}
}

func buildAndStoreEntropyCodesDistance(self *blockEncoder, histograms []histogramDistance, histograms_size uint, alphabet_size uint, tree []huffmanTree, bw *bitWriter) {
	var table_size uint = histograms_size * self.histogram_length_
	if cap(self.depths_) < int(table_size) {
		self.depths_ = make([]byte, table_size)
	} else {
		self.depths_ = self.depths_[:table_size]
	}
	if cap(self.bits_) < int(table_size) {
		self.bits_ = make([]uint16, table_size)
	} else {
		self.bits_ = self.bits_[:table_size]
	}
	{
		var i uint
		for i = 0; i < histograms_size; i++ {
			var ix uint = i * self.histogram_length_
			buildAndStoreHuffmanTree(histograms[i].data_[0:], self.histogram_length_, alphabet_size, tree, self.depths_[ix:], self.bits_[ix:], bw)
		}
	}
}

func storeMetaBlock(input []byte, start_pos uint, length uint, mask uint, prev_byte byte, prev_byte2 byte, is_last bool, params *encoderParams, literal_context_mode int, commands []command, mb *metaBlockSplit, bw *bitWriter) {
	var pos uint = start_pos
	var i uint
	var num_distance_symbols uint32 = params.dist.alphabet_size
	var num_effective_distance_symbols uint32 = num_distance_symbols
	var tree []huffmanTree
	var literal_context_lut contextLUT = getContextLUT(literal_context_mode)
	var dist *distanceParams = &params.dist
	if params.large_window && num_effective_distance_symbols > numHistogramDistanceSymbols {
		num_effective_distance_symbols = numHistogramDistanceSymbols
	}

	storeCompressedMetaBlockHeader(is_last, length, bw)

	tree = make([]huffmanTree, maxHuffmanTreeSize)
	literal_enc := getBlockEncoder(numLiteralSymbols, mb.literal_split.num_types, mb.literal_split.types, mb.literal_split.lengths, mb.literal_split.num_blocks)
	command_enc := getBlockEncoder(numCommandSymbols, mb.command_split.num_types, mb.command_split.types, mb.command_split.lengths, mb.command_split.num_blocks)
	distance_enc := getBlockEncoder(uint(num_effective_distance_symbols), mb.distance_split.num_types, mb.distance_split.types, mb.distance_split.lengths, mb.distance_split.num_blocks)

	buildAndStoreBlockSwitchEntropyCodes(literal_enc, tree, bw)
	buildAndStoreBlockSwitchEntropyCodes(command_enc, tree, bw)
	buildAndStoreBlockSwitchEntropyCodes(distance_enc, tree, bw)

	bw.writeBits(2, uint64(dist.distance_postfix_bits))
	bw.writeBits(4, uint64(dist.num_direct_distance_codes)>>dist.distance_postfix_bits)
	for i = 0; i < mb.literal_split.num_types; i++ {
		bw.writeBits(2, uint64(literal_context_mode))
	}

	if mb.literal_context_map_size == 0 {
		storeTrivialContextMap(mb.literal_histograms_size, literalContextBits, tree, bw)
	} else {
		encodeContextMap(mb.literal_context_map, mb.literal_context_map_size, mb.literal_histograms_size, tree, bw)
	}

	if mb.distance_context_map_size == 0 {
		storeTrivialContextMap(mb.distance_histograms_size, distanceContextBits, tree, bw)
	} else {
		encodeContextMap(mb.distance_context_map, mb.distance_context_map_size, mb.distance_histograms_size, tree, bw)
	}

	buildAndStoreEntropyCodesLiteral(literal_enc, mb.literal_histograms, mb.literal_histograms_size, numLiteralSymbols, tree, bw)
	buildAndStoreEntropyCodesCommand(command_enc, mb.command_histograms, mb.command_histograms_size, numCommandSymbols, tree, bw)
	buildAndStoreEntropyCodesDistance(distance_enc, mb.distance_histograms, mb.distance_histograms_size, uint(num_distance_symbols), tree, bw)
	tree = nil

	for _, cmd := range commands {
		var cmd_code uint = uint(cmd.cmd_prefix_)
		storeSymbol(command_enc, cmd_code, bw)
		storeCommandExtra(&cmd, bw)
		if mb.literal_context_map_size == 0 {
			var j uint
			for j = uint(cmd.insert_len_); j != 0; j-- {
				storeSymbol(literal_enc, uint(input[pos&mask]), bw)
				pos++
			}
		} else {
			var j uint
			for j = uint(cmd.insert_len_); j != 0; j-- {
				var context uint = uint(getContext(prev_byte, prev_byte2, literal_context_lut))
				var literal byte = input[pos&mask]
				storeSymbolWithContext(literal_enc, uint(literal), context, mb.literal_context_map, bw, literalContextBits)
				prev_byte2 = prev_byte
				prev_byte = literal
				pos++
			}
		}

		pos += uint(commandCopyLen(&cmd))
		if commandCopyLen(&cmd) != 0 {
			prev_byte2 = input[(pos-2)&mask]
			prev_byte = input[(pos-1)&mask]
			if cmd.cmd_prefix_ >= 128 {
				var dist_code uint = uint(cmd.dist_prefix_) & 0x3FF
				var distnumextra uint32 = uint32(cmd.dist_prefix_) >> 10
				var distextra uint64 = uint64(cmd.dist_extra_)
				if mb.distance_context_map_size == 0 {
					storeSymbol(distance_enc, dist_code, bw)
				} else {
					var context uint = uint(commandDistanceContext(&cmd))
					storeSymbolWithContext(distance_enc, dist_code, context, mb.distance_context_map, bw, distanceContextBits)
				}

				bw.writeBits(uint(distnumextra), distextra)
			}
		}
	}

	cleanupBlockEncoder(distance_enc)
	cleanupBlockEncoder(command_enc)
	cleanupBlockEncoder(literal_enc)
	if is_last {
		bw.jumpToByteBoundary()
	}
}

func buildHistograms(input []byte, start_pos uint, mask uint, commands []command, lit_histo *histogramLiteral, cmd_histo *histogramCommand, dist_histo *histogramDistance) {
	var pos uint = start_pos
	for _, cmd := range commands {
		var j uint
		histogramAddCommand(cmd_histo, uint(cmd.cmd_prefix_))
		for j = uint(cmd.insert_len_); j != 0; j-- {
			histogramAddLiteral(lit_histo, uint(input[pos&mask]))
			pos++
		}

		pos += uint(commandCopyLen(&cmd))
		if commandCopyLen(&cmd) != 0 && cmd.cmd_prefix_ >= 128 {
			histogramAddDistance(dist_histo, uint(cmd.dist_prefix_)&0x3FF)
		}
	}
}

func storeDataWithHuffmanCodes(input []byte, start_pos uint, mask uint, commands []command, lit_depth []byte, lit_bits []uint16, cmd_depth []byte, cmd_bits []uint16, dist_depth []byte, dist_bits []uint16, bw *bitWriter) {
	var pos uint = start_pos
	for _, cmd := range commands {
		var cmd_code uint = uint(cmd.cmd_prefix_)
		var j uint
		bw.writeBits(uint(cmd_depth[cmd_code]), uint64(cmd_bits[cmd_code]))
		storeCommandExtra(&cmd, bw)
		for j = uint(cmd.insert_len_); j != 0; j-- {
			var literal byte = input[pos&mask]
			bw.writeBits(uint(lit_depth[literal]), uint64(lit_bits[literal]))
			pos++
		}

		pos += uint(commandCopyLen(&cmd))
		if commandCopyLen(&cmd) != 0 && cmd.cmd_prefix_ >= 128 {
			var dist_code uint = uint(cmd.dist_prefix_) & 0x3FF
			var distnumextra uint32 = uint32(cmd.dist_prefix_) >> 10
			var distextra uint32 = cmd.dist_extra_
			bw.writeBits(uint(dist_depth[dist_code]), uint64(dist_bits[dist_code]))
			bw.writeBits(uint(distnumextra), uint64(distextra))
		}
	}
}

func storeMetaBlockTrivial(input []byte, start_pos uint, length uint, mask uint, is_last bool, params *encoderParams, commands []command, bw *bitWriter) {
	var lit_histo histogramLiteral
	var cmd_histo histogramCommand
	var dist_histo histogramDistance
	var lit_depth [numLiteralSymbols]byte
	var lit_bits [numLiteralSymbols]uint16
	var cmd_depth [numCommandSymbols]byte
	var cmd_bits [numCommandSymbols]uint16
	var dist_depth [maxSimpleDistanceAlphabetSize]byte
	var dist_bits [maxSimpleDistanceAlphabetSize]uint16
	var tree []huffmanTree
	var num_distance_symbols uint32 = params.dist.alphabet_size

	storeCompressedMetaBlockHeader(is_last, length, bw)

	histogramClearLiteral(&lit_histo)
	histogramClearCommand(&cmd_histo)
	histogramClearDistance(&dist_histo)

	buildHistograms(input, start_pos, mask, commands, &lit_histo, &cmd_histo, &dist_histo)

	bw.writeBits(13, 0)

	tree = make([]huffmanTree, maxHuffmanTreeSize)
	buildAndStoreHuffmanTree(lit_histo.data_[:], numLiteralSymbols, numLiteralSymbols, tree, lit_depth[:], lit_bits[:], bw)
	buildAndStoreHuffmanTree(cmd_histo.data_[:], numCommandSymbols, numCommandSymbols, tree, cmd_depth[:], cmd_bits[:], bw)
	buildAndStoreHuffmanTree(dist_histo.data_[:], maxSimpleDistanceAlphabetSize, uint(num_distance_symbols), tree, dist_depth[:], dist_bits[:], bw)
	tree = nil
	storeDataWithHuffmanCodes(input, start_pos, mask, commands, lit_depth[:], lit_bits[:], cmd_depth[:], cmd_bits[:], dist_depth[:], dist_bits[:], bw)
	if is_last {
		bw.jumpToByteBoundary()
	}
}

func storeMetaBlockFast(input []byte, start_pos uint, length uint, mask uint, is_last bool, params *encoderParams, commands []command, bw *bitWriter) {
	var num_distance_symbols uint32 = params.dist.alphabet_size
	var distance_alphabet_bits uint32 = log2FloorNonZero(uint(num_distance_symbols-1)) + 1

	storeCompressedMetaBlockHeader(is_last, length, bw)

	bw.writeBits(13, 0)

	if len(commands) <= 128 {
		var histogram = [numLiteralSymbols]uint32{0}
		var pos uint = start_pos
		var num_literals uint = 0
		var lit_depth [numLiteralSymbols]byte
		var lit_bits [numLiteralSymbols]uint16
		for _, cmd := range commands {
			var j uint
			for j = uint(cmd.insert_len_); j != 0; j-- {
				histogram[input[pos&mask]]++
				pos++
			}

			num_literals += uint(cmd.insert_len_)
			pos += uint(commandCopyLen(&cmd))
		}

		buildAndStoreHuffmanTreeFast(histogram[:], num_literals, /* max_bits = */
			8, lit_depth[:], lit_bits[:], bw)

		storeStaticCommandHuffmanTree(bw)
		storeStaticDistanceHuffmanTree(bw)
		storeDataWithHuffmanCodes(input, start_pos, mask, commands, lit_depth[:], lit_bits[:], kStaticCommandCodeDepth[:], kStaticCommandCodeBits[:], kStaticDistanceCodeDepth[:], kStaticDistanceCodeBits[:], bw)
	} else {
		var lit_histo histogramLiteral
		var cmd_histo histogramCommand
		var dist_histo histogramDistance
		var lit_depth [numLiteralSymbols]byte
		var lit_bits [numLiteralSymbols]uint16
		var cmd_depth [numCommandSymbols]byte
		var cmd_bits [numCommandSymbols]uint16
		var dist_depth [maxSimpleDistanceAlphabetSize]byte
		var dist_bits [maxSimpleDistanceAlphabetSize]uint16
		histogramClearLiteral(&lit_histo)
		histogramClearCommand(&cmd_histo)
		histogramClearDistance(&dist_histo)
		buildHistograms(input, start_pos, mask, commands, &lit_histo, &cmd_histo, &dist_histo)
		buildAndStoreHuffmanTreeFast(lit_histo.data_[:], lit_histo.total_count_, /* max_bits = */
			8, lit_depth[:], lit_bits[:], bw)

		buildAndStoreHuffmanTreeFast(cmd_histo.data_[:], cmd_histo.total_count_, /* max_bits = */
			10, cmd_depth[:], cmd_bits[:], bw)

		buildAndStoreHuffmanTreeFast(dist_histo.data_[:], dist_histo.total_count_, /* max_bits = */
			uint(distance_alphabet_bits), dist_depth[:], dist_bits[:], bw)

		storeDataWithHuffmanCodes(input, start_pos, mask, commands, lit_depth[:], lit_bits[:], cmd_depth[:], cmd_bits[:], dist_depth[:], dist_bits[:], bw)
	}

	if is_last {
		bw.jumpToByteBoundary()
	}
}

/*
This is for storing uncompressed blocks (simple raw storage of

	bytes-as-bytes).
*/
func storeUncompressedMetaBlock(is_final_block bool, input []byte, position uint, mask uint, len uint, bw *bitWriter) {
	var masked_pos uint = position & mask
	storeUncompressedMetaBlockHeader(uint(len), bw)
	bw.jumpToByteBoundary()

	if masked_pos+len > mask+1 {
		var len1 uint = mask + 1 - masked_pos
		bw.writeBytes(input[masked_pos:][:len1])
		len -= len1
		masked_pos = 0
	}

	bw.writeBytes(input[masked_pos:][:len])

	/* Since the uncompressed block itself may not be the final block, add an
	   empty one after this. */
	if is_final_block {
		bw.writeBits(1, 1) /* islast */
		bw.writeBits(1, 1) /* isempty */
		bw.jumpToByteBoundary()
	}
}
