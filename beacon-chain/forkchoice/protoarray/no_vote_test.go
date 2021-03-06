package protoarray

import (
	"context"
	"strings"
	"testing"

	"github.com/prysmaticlabs/prysm/shared/params"
)

func TestNoVote_CanFindHead(t *testing.T) {
	balances := make([]uint64, 16)
	f := setup(1, 1)

	// The head should always start at the finalized block.
	r, err := f.Head(context.Background(), 1, params.BeaconConfig().ZeroHash, balances, 1)
	if err != nil {
		t.Fatal(err)
	}
	if r != params.BeaconConfig().ZeroHash {
		t.Errorf("Incorrect head with genesis")
	}

	// Insert block 2 into the tree and verify head is at 2:
	//         0
	//        /
	//       2 <- head
	if err := f.ProcessBlock(context.Background(), 0, indexToHash(2), params.BeaconConfig().ZeroHash, [32]byte{}, 1, 1); err != nil {
		t.Fatal(err)
	}
	r, err = f.Head(context.Background(), 1, params.BeaconConfig().ZeroHash, balances, 1)
	if err != nil {
		t.Fatal(err)
	}
	if r != indexToHash(2) {
		t.Error("Incorrect head for with justified epoch at 1")
	}

	// Insert block 1 into the tree and verify head is still at 2:
	//            0
	//           / \
	//  head -> 2  1
	if err := f.ProcessBlock(context.Background(), 0, indexToHash(1), params.BeaconConfig().ZeroHash, [32]byte{}, 1, 1); err != nil {
		t.Fatal(err)
	}
	r, err = f.Head(context.Background(), 1, params.BeaconConfig().ZeroHash, balances, 1)
	if err != nil {
		t.Fatal(err)
	}
	if r != indexToHash(2) {
		t.Error("Incorrect head for with justified epoch at 1")
	}

	// Insert block 3 into the tree and verify head is still at 2:
	//            0
	//           / \
	//  head -> 2  1
	//             |
	//             3
	if err := f.ProcessBlock(context.Background(), 0, indexToHash(3), indexToHash(1), [32]byte{}, 1, 1); err != nil {
		t.Fatal(err)
	}
	r, err = f.Head(context.Background(), 1, params.BeaconConfig().ZeroHash, balances, 1)
	if err != nil {
		t.Fatal(err)
	}
	if r != indexToHash(2) {
		t.Error("Incorrect head for with justified epoch at 1")
	}

	// Insert block 4 into the tree and verify head is at 4:
	//            0
	//           / \
	//          2  1
	//          |  |
	//  head -> 4  3
	if err := f.ProcessBlock(context.Background(), 0, indexToHash(4), indexToHash(2), [32]byte{}, 1, 1); err != nil {
		t.Fatal(err)
	}
	r, err = f.Head(context.Background(), 1, params.BeaconConfig().ZeroHash, balances, 1)
	if err != nil {
		t.Fatal(err)
	}
	if r != indexToHash(4) {
		t.Error("Incorrect head for with justified epoch at 1")
	}

	// Insert block 5 with justified epoch of 2, verify head is still at 4.
	//            0
	//           / \
	//          2  1
	//          |  |
	//  head -> 4  3
	//          |
	//          5 <- justified epoch = 2
	if err := f.ProcessBlock(context.Background(), 0, indexToHash(5), indexToHash(4), [32]byte{}, 2, 1); err != nil {
		t.Fatal(err)
	}
	r, err = f.Head(context.Background(), 1, params.BeaconConfig().ZeroHash, balances, 1)
	if err != nil {
		t.Fatal(err)
	}
	if r != indexToHash(4) {
		t.Error("Incorrect head for with justified epoch at 1")
	}

	// Verify there's an error when starting from a block with wrong justified epoch.
	//            0
	//           / \
	//          2  1
	//          |  |
	//  head -> 4  3
	//          |
	//          5 <- starting from 5 with justified epoch 0 should error
	r, err = f.Head(context.Background(), 1, indexToHash(5), balances, 1)
	wanted := "head at slot 0 with weight 0 is not eligible, finalizedEpoch 1 != 1, justifiedEpoch 2 != 1"
	if !strings.Contains(err.Error(), wanted) {
		t.Fatal("Did not get wanted error")
	}

	// Set the justified epoch to 2 and start block to 5 to verify head is 5.
	//            0
	//           / \
	//          2  1
	//          |  |
	//          4  3
	//          |
	//          5 <- head
	r, err = f.Head(context.Background(), 2, indexToHash(5), balances, 1)
	if err != nil {
		t.Fatal(err)
	}
	if r != indexToHash(5) {
		t.Error("Incorrect head for with justified epoch at 2")
	}

	// Insert block 6 with justified epoch of 2, verify head is at 6.
	//            0
	//           / \
	//          2  1
	//          |  |
	//          4  3
	//          |
	//          5
	//          |
	//          6 <- head
	if err := f.ProcessBlock(context.Background(), 0, indexToHash(6), indexToHash(5), [32]byte{}, 2, 1); err != nil {
		t.Fatal(err)
	}
	r, err = f.Head(context.Background(), 2, indexToHash(5), balances, 1)
	if err != nil {
		t.Fatal(err)
	}
	if r != indexToHash(6) {
		t.Error("Incorrect head for with justified epoch at 2")
	}
}
