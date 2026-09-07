// improvised min heap implementation lol

package heap


type Node[T any] struct {
	data T
	val  int
}

type MinHeap []Node[any]

func (h *MinHeap) Parent(idx int) int {
	return (idx - 1) / 2
}

func (h *MinHeap) LeftChild(idx int) int {
	return 2*idx + 1
}

func (h *MinHeap) RightChild(idx int) int {
	return 2*idx + 2
}

func (h *MinHeap) Length() int {
	return len(*h)
}

func (h *MinHeap) Node(idx int) Node[any] {
	return (*h)[idx]
}

func (h *MinHeap) Peek() Node[any] {
	return h.Node(0)
}

func (h *MinHeap) Insert(data any) {
	(*h)[h.Length()] = Node[any]{
		data: data,
		val:  len(*h),
	}
	h.HeapifyUp(h.Length())
}

func (h *MinHeap) HeapifyUp(idx int) {
	if idx == 0 {
		return
	}

	parent := h.Parent(idx)
	child := h.Node(idx)

	parentV := h.Node(parent).val
	childV := child.val

	if childV < parentV {
		temp := h.Node(parent)
		(*h)[parent] = child
		(*h)[idx] = temp

		h.HeapifyUp(parent)
	}
}

func (h *MinHeap) HeapifyDown(idx int) {
	leftChild := h.LeftChild(idx)
	rightChild := h.RightChild(idx)

	if idx >= h.Length() || leftChild >= h.Length() {
		return
	}

	lV := h.Node(leftChild).val
	rV := h.Node(rightChild).val

	v := h.Node(idx).val

	if lV > rV && v > rV {
		temp := h.Node(idx)
		(*h)[idx] = h.Node(rightChild)
		(*h)[rightChild] = temp
		h.HeapifyDown(rightChild)
	} else if rV > lV && v > lV {
		temp := h.Node(idx)
		(*h)[idx] = h.Node(leftChild)
		(*h)[leftChild] = temp
		h.HeapifyDown(leftChild)
	}
}
